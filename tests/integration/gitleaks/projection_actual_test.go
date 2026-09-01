package gitleaks_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/engine"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/engine/gitleaks"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/gitinput"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/outcome"
)

const actualSyntheticCanary = "PSCAN_SYNTHETIC_SECRET_83e08b4f9be9467ba5c614c51d7f2a9083e08b4f9be9467ba5c614c51d7f2a90"

func TestPinnedGitleaksProjectionCoverage(t *testing.T) {
	binary := requireAbsoluteEnvironment(t, "PSCAN_GITLEAKS_BINARY")
	gitBinary := requireAbsoluteEnvironment(t, "PSCAN_GIT_BINARY")
	config := requireAbsoluteEnvironment(t, "PSCAN_GITLEAKS_CONFIG")
	ignore := requireAbsoluteEnvironment(t, "PSCAN_GITLEAKS_IGNORE")
	if got := strings.TrimSpace(runActual(t, binary, "version")); got != gitleaks.EngineVersion {
		t.Fatalf("wrong actual Gitleaks version: %q", got)
	}

	t.Run("text-head", func(t *testing.T) {
		root, repository := newActualRepository(t, gitBinary)
		writeActual(t, filepath.Join(repository, "clean.txt"), []byte("clean\n"))
		commitActual(t, gitBinary, repository, "base", "clean.txt")
		base := strings.TrimSpace(runActualDir(t, repository, gitBinary, "rev-parse", "HEAD"))
		writeActual(t, filepath.Join(repository, "text.txt"), []byte(actualSyntheticCanary))
		commitActual(t, gitBinary, repository, "text", "text.txt")
		head := strings.TrimSpace(runActualDir(t, repository, gitBinary, "rev-parse", "HEAD"))
		assertActualFinding(t, root, repository, gitBinary, binary, config, ignore, base, head)
	})

	t.Run("binary-head", func(t *testing.T) {
		root, repository := newActualRepository(t, gitBinary)
		writeActual(t, filepath.Join(repository, "clean.txt"), []byte("clean\n"))
		commitActual(t, gitBinary, repository, "base", "clean.txt")
		base := strings.TrimSpace(runActualDir(t, repository, gitBinary, "rev-parse", "HEAD"))
		writeActual(t, filepath.Join(repository, "binary.bin"), append(append([]byte{0, 1, 2}, actualSyntheticCanary...), 0xff))
		commitActual(t, gitBinary, repository, "binary", "binary.bin")
		head := strings.TrimSpace(runActualDir(t, repository, gitBinary, "rev-parse", "HEAD"))
		assertActualFinding(t, root, repository, gitBinary, binary, config, ignore, base, head)
	})

	t.Run("deleted-binary-history", func(t *testing.T) {
		root, repository := newActualRepository(t, gitBinary)
		writeActual(t, filepath.Join(repository, "clean.txt"), []byte("clean\n"))
		commitActual(t, gitBinary, repository, "base", "clean.txt")
		base := strings.TrimSpace(runActualDir(t, repository, gitBinary, "rev-parse", "HEAD"))
		writeActual(t, filepath.Join(repository, "deleted.bin"), append(append([]byte{0, 1}, actualSyntheticCanary...), 0xff))
		commitActual(t, gitBinary, repository, "add binary", "deleted.bin")
		if err := os.Remove(filepath.Join(repository, "deleted.bin")); err != nil {
			t.Fatal(err)
		}
		runActualDir(t, repository, gitBinary, "add", "-u", "--", "deleted.bin")
		runActualDir(t, repository, gitBinary, "commit", "--quiet", "-m", "delete binary")
		head := strings.TrimSpace(runActualDir(t, repository, gitBinary, "rev-parse", "HEAD"))
		assertActualFinding(t, root, repository, gitBinary, binary, config, ignore, base, head)
	})

	t.Run("head-tree-unchanged", func(t *testing.T) {
		root, repository := newActualRepository(t, gitBinary)
		writeActual(t, filepath.Join(repository, "unchanged.txt"), []byte(actualSyntheticCanary))
		commitActual(t, gitBinary, repository, "base finding", "unchanged.txt")
		base := strings.TrimSpace(runActualDir(t, repository, gitBinary, "rev-parse", "HEAD"))
		writeActual(t, filepath.Join(repository, "clean.txt"), []byte("clean\n"))
		commitActual(t, gitBinary, repository, "clean head", "clean.txt")
		head := strings.TrimSpace(runActualDir(t, repository, gitBinary, "rev-parse", "HEAD"))
		assertActualFinding(t, root, repository, gitBinary, binary, config, ignore, base, head)
	})

	t.Run("rename", func(t *testing.T) {
		root, repository := newActualRepository(t, gitBinary)
		writeActual(t, filepath.Join(repository, "old.txt"), []byte(actualSyntheticCanary))
		commitActual(t, gitBinary, repository, "base", "old.txt")
		base := strings.TrimSpace(runActualDir(t, repository, gitBinary, "rev-parse", "HEAD"))
		runActualDir(t, repository, gitBinary, "mv", "--", "old.txt", "new.txt")
		runActualDir(t, repository, gitBinary, "commit", "--quiet", "-m", "rename")
		head := strings.TrimSpace(runActualDir(t, repository, gitBinary, "rev-parse", "HEAD"))
		assertActualFinding(t, root, repository, gitBinary, binary, config, ignore, base, head)
	})

	t.Run("merge-parent-edge", func(t *testing.T) {
		root, repository := newActualRepository(t, gitBinary)
		writeActual(t, filepath.Join(repository, "clean.txt"), []byte("clean\n"))
		commitActual(t, gitBinary, repository, "base", "clean.txt")
		base := strings.TrimSpace(runActualDir(t, repository, gitBinary, "rev-parse", "HEAD"))
		runActualDir(t, repository, gitBinary, "checkout", "--quiet", "-b", "feature")
		writeActual(t, filepath.Join(repository, "feature.txt"), []byte(actualSyntheticCanary))
		commitActual(t, gitBinary, repository, "feature", "feature.txt")
		runActualDir(t, repository, gitBinary, "checkout", "--quiet", "main")
		writeActual(t, filepath.Join(repository, "main.txt"), []byte("clean main\n"))
		commitActual(t, gitBinary, repository, "main", "main.txt")
		runActualDir(t, repository, gitBinary, "merge", "--quiet", "--no-ff", "feature", "-m", "merge")
		head := strings.TrimSpace(runActualDir(t, repository, gitBinary, "rev-parse", "HEAD"))
		assertActualFinding(t, root, repository, gitBinary, binary, config, ignore, base, head)
	})

	t.Run("deleted-text-and-binary", func(t *testing.T) {
		root, repository := newActualRepository(t, gitBinary)
		writeActual(t, filepath.Join(repository, "clean.txt"), []byte("clean\n"))
		writeActual(t, filepath.Join(repository, ".gitattributes"), []byte("*.dat binary\n"))
		commitActual(t, gitBinary, repository, "base", ".")
		base := strings.TrimSpace(runActualDir(t, repository, gitBinary, "rev-parse", "HEAD"))
		writeActual(t, filepath.Join(repository, "deleted-text.txt"), []byte(actualSyntheticCanary+"\n"))
		binaryBytes := append(append([]byte{0, 1, 2, 0}, []byte(actualSyntheticCanary)...), 0, 3, 4)
		writeActual(t, filepath.Join(repository, "deleted-binary.dat"), binaryBytes)
		commitActual(t, gitBinary, repository, "add canaries", "deleted-text.txt", "deleted-binary.dat")
		if err := os.Remove(filepath.Join(repository, "deleted-text.txt")); err != nil {
			t.Fatal(err)
		}
		if err := os.Remove(filepath.Join(repository, "deleted-binary.dat")); err != nil {
			t.Fatal(err)
		}
		runActualDir(t, repository, gitBinary, "add", "-u", "--", "deleted-text.txt", "deleted-binary.dat")
		runActualDir(t, repository, gitBinary, "commit", "--quiet", "-m", "delete canaries")
		head := strings.TrimSpace(runActualDir(t, repository, gitBinary, "rev-parse", "HEAD"))
		adapter, projection := actualProjection(t, root, repository, gitBinary, binary, config, ignore, base, head)
		t.Logf("bindings plan=%s content=%s coverage=%s entries=%d", projection.PlanDigest, projection.ContentDigest, projection.ProbeDigest, projection.EntryCount)
		result := adapter.ScanProjectionProfile(context.Background(), projection, gitinput.PRProfile)
		if result.Reason != outcome.ReasonFailFindingDetected {
			t.Fatalf("deleted text/binary canaries were not detected: %#v", result)
		}
		if strings.Contains(fmt.Sprintf("%#v", result), "PSCAN_SYNTHETIC_SECRET_") {
			t.Fatal("candidate content crossed the private result boundary")
		}
	})

	t.Run("clean-and-out-of-range", func(t *testing.T) {
		root, repository := newActualRepository(t, gitBinary)
		writeActual(t, filepath.Join(repository, "out-of-range.txt"), []byte(actualSyntheticCanary+"\n"))
		commitActual(t, gitBinary, repository, "pre-range", "out-of-range.txt")
		if err := os.Remove(filepath.Join(repository, "out-of-range.txt")); err != nil {
			t.Fatal(err)
		}
		writeActual(t, filepath.Join(repository, "clean.txt"), []byte("clean base\n"))
		runActualDir(t, repository, gitBinary, "add", "-A", "--", ".")
		runActualDir(t, repository, gitBinary, "commit", "--quiet", "-m", "range base")
		base := strings.TrimSpace(runActualDir(t, repository, gitBinary, "rev-parse", "HEAD"))
		writeActual(t, filepath.Join(repository, "clean.txt"), []byte("clean head\n"))
		commitActual(t, gitBinary, repository, "clean head", "clean.txt")
		head := strings.TrimSpace(runActualDir(t, repository, gitBinary, "rev-parse", "HEAD"))
		adapter, projection := actualProjection(t, root, repository, gitBinary, binary, config, ignore, base, head)
		t.Logf("bindings plan=%s content=%s coverage=%s entries=%d", projection.PlanDigest, projection.ContentDigest, projection.ProbeDigest, projection.EntryCount)
		result := adapter.ScanProjectionProfile(context.Background(), projection, gitinput.PRProfile)
		if result.Reason != outcome.ReasonPassNoBlockingFindings {
			t.Fatalf("clean exact range did not pass: %#v", result)
		}
	})

	for _, boundary := range []struct {
		name     string
		position int
	}{
		{name: "first-byte", position: 0},
		{name: "internal-fragment-boundary-no-whitespace", position: 89990},
		{name: "last-byte", position: 200000 - len(actualSyntheticCanary)},
	} {
		t.Run("span-"+boundary.name, func(t *testing.T) {
			root, repository := newActualRepository(t, gitBinary)
			writeActual(t, filepath.Join(repository, "clean.txt"), []byte("clean\n"))
			commitActual(t, gitBinary, repository, "base", "clean.txt")
			base := strings.TrimSpace(runActualDir(t, repository, gitBinary, "rev-parse", "HEAD"))
			payload := []byte(strings.Repeat("x", 200000))
			copy(payload[boundary.position:], actualSyntheticCanary)
			writeActual(t, filepath.Join(repository, "boundary.bin"), payload)
			commitActual(t, gitBinary, repository, "boundary", "boundary.bin")
			head := strings.TrimSpace(runActualDir(t, repository, gitBinary, "rev-parse", "HEAD"))
			adapter, projection := actualProjection(t, root, repository, gitBinary, binary, config, ignore, base, head)
			if result := adapter.ScanProjectionProfile(context.Background(), projection, gitinput.PRProfile); result.Reason != outcome.ReasonFailFindingDetected {
				t.Fatalf("boundary finding was not inspected: %#v", result)
			}
		})
	}

	t.Run("exact-maximum-span-crosses-fragment-boundary", func(t *testing.T) {
		root, repository := newActualRepository(t, gitBinary)
		writeActual(t, filepath.Join(repository, "clean.txt"), []byte("clean\n"))
		commitActual(t, gitBinary, repository, "base", "clean.txt")
		base := strings.TrimSpace(runActualDir(t, repository, gitBinary, "rev-parse", "HEAD"))
		payload := []byte(strings.Repeat("x", 200000))
		canary := "PSCAN_MAX_SPAN_BEGIN_" + strings.Repeat("A", 3995) + "_END"
		if len(canary) != int(gitinput.MaximumRuleSpan) {
			t.Fatal("maximum-span fixture has wrong length")
		}
		copy(payload[87990:], canary)
		writeActual(t, filepath.Join(repository, "maximum-span.bin"), payload)
		commitActual(t, gitBinary, repository, "maximum span", "maximum-span.bin")
		head := strings.TrimSpace(runActualDir(t, repository, gitBinary, "rev-parse", "HEAD"))
		adapter, projection := actualProjection(t, root, repository, gitBinary, binary, config, ignore, base, head)
		if result := adapter.ScanProjectionProfile(context.Background(), projection, gitinput.PRProfile); result.Reason != outcome.ReasonFailFindingDetected {
			t.Fatalf("maximum-span finding was not inspected: %#v", result)
		}
	})

	t.Run("declared-maximum-blob", func(t *testing.T) {
		if os.Getenv("PSCAN_RUN_MAXIMUM_SIZE") != "1" {
			t.Skip("set PSCAN_RUN_MAXIMUM_SIZE=1 for the bounded 512 MiB fixture")
		}
		root, repository := newActualRepository(t, gitBinary)
		writeActual(t, filepath.Join(repository, "clean.txt"), []byte("clean\n"))
		commitActual(t, gitBinary, repository, "base", "clean.txt")
		base := strings.TrimSpace(runActualDir(t, repository, gitBinary, "rev-parse", "HEAD"))
		path := filepath.Join(repository, "maximum.bin")
		f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_EXCL, 0o600)
		if err != nil {
			t.Fatal(err)
		}
		if err := f.Truncate(gitinput.PRProfile.MaxBlobBytes); err == nil {
			_, err = f.WriteAt([]byte(actualSyntheticCanary), gitinput.PRProfile.MaxBlobBytes-int64(len(actualSyntheticCanary)))
		}
		if closeErr := f.Close(); err == nil {
			err = closeErr
		}
		if err != nil {
			t.Fatal(err)
		}
		commitActual(t, gitBinary, repository, "maximum blob", "maximum.bin")
		head := strings.TrimSpace(runActualDir(t, repository, gitBinary, "rev-parse", "HEAD"))
		assertActualFinding(t, root, repository, gitBinary, binary, config, ignore, base, head)
	})

	t.Run("path-skip-is-non-pass", func(t *testing.T) {
		root, repository := newActualRepository(t, gitBinary)
		writeActual(t, filepath.Join(repository, "clean.txt"), []byte("clean\n"))
		commitActual(t, gitBinary, repository, "base", "clean.txt")
		base := strings.TrimSpace(runActualDir(t, repository, gitBinary, "rev-parse", "HEAD"))
		writeActual(t, filepath.Join(repository, "skipped.bin"), []byte("clean binary-class path\n"))
		commitActual(t, gitBinary, repository, "binary extension", "skipped.bin")
		head := strings.TrimSpace(runActualDir(t, repository, gitBinary, "rev-parse", "HEAD"))
		adapter, projection := actualProjection(t, root, repository, gitBinary, binary, config, ignore, base, head)
		t.Logf("bindings plan=%s content=%s coverage=%s entries=%d", projection.PlanDigest, projection.ContentDigest, projection.ProbeDigest, projection.EntryCount)
		result := adapter.ScanProjectionProfile(context.Background(), projection, gitinput.PRProfile)
		if result.Reason != outcome.ReasonPassNoBlockingFindings {
			t.Fatalf("binary-extension path was not inspected: %#v", result)
		}
	})
}

func actualProjection(t *testing.T, root, repository, gitBinary, gitleaksBinary, config, ignore, base, head string) (gitleaks.Adapter, gitinput.MaterializedProjection) {
	t.Helper()
	private := filepath.Join(root, "private")
	bare := filepath.Join(root, "bare.git")
	git := gitinput.Git{Executable: gitBinary, Digest: digestFile(t, gitBinary), Timeout: 15 * time.Second}
	if err := git.PrepareBare(context.Background(), repository, bare, private); err != nil {
		t.Fatal(err)
	}
	plan, err := git.PlanRange(context.Background(), bare, base, head, false, filepath.Join(private, "home"), gitinput.PRProfile.ProjectionLimits())
	if err != nil {
		t.Fatal(err)
	}
	projection, err := git.Materialize(context.Background(), bare, filepath.Join(private, "home"), filepath.Join(root, "raw"), filepath.Join(root, "framed"), plan)
	if err != nil {
		t.Fatal(err)
	}
	home := filepath.Join(root, "engine-home")
	if err := os.Mkdir(home, 0o700); err != nil {
		t.Fatal(err)
	}
	adapter := gitleaks.Adapter{Binding: gitleaks.Binding{
		Executable: gitleaksBinary, ExecutableDigest: digestFile(t, gitleaksBinary),
		Config: config, ConfigDigest: digestFile(t, config),
		IgnoreFile: ignore, IgnoreFileDigest: digestFile(t, ignore),
		PrivateHome: home, Environment: engine.SafeEnvironment(filepath.Dir(gitleaksBinary), home),
	}}
	return adapter, projection
}

func assertActualFinding(t *testing.T, root, repository, gitBinary, binary, config, ignore, base, head string) {
	t.Helper()
	adapter, projection := actualProjection(t, root, repository, gitBinary, binary, config, ignore, base, head)
	if result := adapter.ScanProjectionProfile(context.Background(), projection, gitinput.PRProfile); result.Reason != outcome.ReasonFailFindingDetected {
		t.Fatalf("isolated required class was not detected: %#v", result)
	}
}

func newActualRepository(t *testing.T, gitBinary string) (string, string) {
	t.Helper()
	root := t.TempDir()
	repository := filepath.Join(root, "source")
	if err := os.Mkdir(repository, 0o700); err != nil {
		t.Fatal(err)
	}
	runActualDir(t, repository, gitBinary, "init", "--quiet", "-b", "main")
	runActualDir(t, repository, gitBinary, "config", "user.email", "synthetic@example.invalid")
	runActualDir(t, repository, gitBinary, "config", "user.name", "Synthetic Fixture")
	return root, repository
}

func commitActual(t *testing.T, gitBinary, repository, message string, paths ...string) {
	t.Helper()
	args := append([]string{"add", "--"}, paths...)
	runActualDir(t, repository, gitBinary, args...)
	runActualDir(t, repository, gitBinary, "commit", "--quiet", "-m", message)
}

func writeActual(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
}

func requireAbsoluteEnvironment(t *testing.T, name string) string {
	t.Helper()
	value := os.Getenv(name)
	if value == "" {
		t.Skip(name + " is required for pinned Gitleaks integration")
	}
	if !filepath.IsAbs(value) {
		t.Fatalf("%s must be absolute", name)
	}
	return value
}

func runActual(t *testing.T, executable string, args ...string) string {
	t.Helper()
	return runActualDir(t, "", executable, args...)
}

func runActualDir(t *testing.T, directory, executable string, args ...string) string {
	t.Helper()
	cmd := exec.Command(executable, args...)
	cmd.Dir = directory
	cmd.Env = append(os.Environ(), "GIT_AUTHOR_DATE=2000-01-01T00:00:00Z", "GIT_COMMITTER_DATE=2000-01-01T00:00:00Z")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("synthetic command failed: %v", err)
	}
	return string(output)
}
