package gitinput_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/gitinput"
)

func TestSafeBareCloneAndExactRangeBinding(t *testing.T) {
	gitPath, digest := gitBinding(t)
	root := t.TempDir()
	source := filepath.Join(root, "source;candidate-not-shell")
	if err := os.Mkdir(source, 0o700); err != nil {
		t.Fatal(err)
	}
	runGit(t, gitPath, source, "init", "--quiet")
	runGit(t, gitPath, source, "config", "user.email", "synthetic@example.invalid")
	runGit(t, gitPath, source, "config", "user.name", "Synthetic Fixture")
	write(t, filepath.Join(source, "clean.txt"), "clean")
	runGit(t, gitPath, source, "add", "--", "clean.txt")
	runGit(t, gitPath, source, "commit", "--quiet", "-m", "synthetic base")
	base := runGit(t, gitPath, source, "rev-parse", "HEAD")

	binaryCanary := "PSCAN_SYNTHETIC_SECRET_NOT_A_CREDENTIAL_4f9be9"
	if err := os.WriteFile(filepath.Join(source, "deleted.bin"), append([]byte{0, 1, 2, 0}, []byte(binaryCanary)...), 0o600); err != nil {
		t.Fatal(err)
	}
	runGit(t, gitPath, source, "add", "--", "deleted.bin")
	runGit(t, gitPath, source, "commit", "--quiet", "-m", "synthetic in-range finding")
	runGit(t, gitPath, source, "rm", "--quiet", "--", "deleted.bin")
	runGit(t, gitPath, source, "commit", "--quiet", "-m", "delete synthetic finding")
	head := runGit(t, gitPath, source, "rev-parse", "HEAD")
	patch := runGit(t, gitPath, source, "log", "-p", "-U0", base+".."+head)
	if strings.Contains(patch, binaryCanary) || !strings.Contains(patch, "Binary files") {
		t.Fatal("fixture did not prove Git binary-patch omission")
	}

	private := filepath.Join(root, "private")
	bare := filepath.Join(root, "bare.git")
	g := gitinput.Git{Executable: gitPath, Digest: digest, Timeout: 15 * time.Second}
	if err := g.PrepareBare(context.Background(), source, bare, private); err != nil {
		t.Fatalf("safe clone failed: %v", err)
	}
	binding, err := g.BindRange(context.Background(), bare, base, head, false, filepath.Join(private, "home"))
	if err != nil {
		t.Fatalf("range binding failed: %v", err)
	}
	if binding.Base != base || binding.Head != head || binding.MergeBase != base || binding.CommitCount != 2 {
		t.Fatalf("wrong range binding: %#v", binding)
	}
	if len(binding.HistoryRangeDigest) != 64 || len(binding.TrackedTreeDigest) != 64 || len(binding.HeadTreeOID) < 40 {
		t.Fatalf("missing digests: %#v", binding)
	}
	plan, err := g.PlanRange(context.Background(), bare, base, head, false, filepath.Join(private, "home"), gitinput.ProjectionLimits{MaxBlobBytes: 1 << 20, MaxBlobCount: 100, MaxTotalBytes: 8 << 20})
	if err != nil {
		t.Fatalf("history plan failed: %v", err)
	}
	projection, err := g.Materialize(context.Background(), bare, filepath.Join(private, "home"), filepath.Join(root, "raw-projection"), filepath.Join(root, "framed-projection"), plan)
	if err != nil {
		t.Fatalf("history projection failed: %v", err)
	}
	if projection.EntryCount != 3 || plan.Range.HistoryRangeDigest != binding.HistoryRangeDigest || len(plan.PlanDigest) != 64 || len(projection.ContentDigest) != 64 || len(projection.ProbeDigest) != 64 {
		t.Fatalf("wrong projection binding: plan=%#v projection=%#v", plan, projection)
	}
	var foundDeletedBinary bool
	for i, entry := range plan.Entries {
		if entry.Path != "deleted.bin" {
			continue
		}
		raw, rawErr := os.ReadFile(filepath.Join(projection.ScanRoot, filepath.FromSlash(projection.Files[i].RelativePath)))
		framed, framedErr := os.ReadFile(filepath.Join(projection.ProbeRoot, filepath.FromSlash(projection.Files[i].Chunks[0].RelativePath)))
		if rawErr != nil || framedErr != nil || !bytes.Contains(raw, []byte(binaryCanary)) || !bytes.Contains(framed, []byte(binaryCanary)) || !bytes.Contains(framed, []byte("PSCAN_COVERAGE_MARKER_")) {
			t.Fatal("deleted binary bytes did not reach the bound scan projection")
		}
		foundDeletedBinary = true
	}
	if !foundDeletedBinary {
		t.Fatal("deleted binary blob was absent from exact history projection")
	}
	if err := gitinput.VerifyMaterializedProjection(projection); err != nil {
		t.Fatalf("materialized projection did not verify: %v", err)
	}

	if _, err := g.BindRange(context.Background(), bare, head, base, false, filepath.Join(private, "home")); err == nil {
		t.Fatal("non-ancestor range was accepted")
	}
	if _, err := os.Stat(filepath.Join(root, "candidate-not-shell")); !os.IsNotExist(err) {
		t.Fatal("candidate path became executable syntax")
	}
}

func TestHistoryProjectionRejectsResourceLimit(t *testing.T) {
	gitPath, digest := gitBinding(t)
	root := t.TempDir()
	source := filepath.Join(root, "source")
	if err := os.Mkdir(source, 0o700); err != nil {
		t.Fatal(err)
	}
	runGit(t, gitPath, source, "init", "--quiet")
	runGit(t, gitPath, source, "config", "user.email", "synthetic@example.invalid")
	runGit(t, gitPath, source, "config", "user.name", "Synthetic Fixture")
	write(t, filepath.Join(source, "large.txt"), strings.Repeat("x", 1024))
	runGit(t, gitPath, source, "add", "--", "large.txt")
	runGit(t, gitPath, source, "commit", "--quiet", "-m", "bounded")
	head := runGit(t, gitPath, source, "rev-parse", "HEAD")
	private := filepath.Join(root, "private")
	bare := filepath.Join(root, "bare.git")
	g := gitinput.Git{Executable: gitPath, Digest: digest, Timeout: 15 * time.Second}
	if err := g.PrepareBare(context.Background(), source, bare, private); err != nil {
		t.Fatal(err)
	}
	if _, err := g.PlanRange(context.Background(), bare, "", head, true, filepath.Join(private, "home"), gitinput.ProjectionLimits{MaxBlobBytes: 10, MaxBlobCount: 10, MaxTotalBytes: 100}); err == nil {
		t.Fatal("oversized historical blob was silently admitted")
	}
}

func gitBinding(t *testing.T) (string, string) {
	t.Helper()
	path, err := exec.LookPath("git")
	if err != nil {
		t.Fatal("Git is required for PSCAN-03 tests")
	}
	path, err = filepath.Abs(path)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(raw)
	return path, hex.EncodeToString(sum[:])
}

func runGit(t *testing.T, gitPath, directory string, args ...string) string {
	t.Helper()
	cmd := exec.Command(gitPath, args...)
	cmd.Dir = directory
	if runtime.GOOS == "windows" {
		cmd.Env = append(os.Environ(), "GIT_CONFIG_NOSYSTEM=1", "GIT_TERMINAL_PROMPT=0")
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("fixture Git command failed: %v", err)
	}
	return strings.TrimSpace(string(out))
}
