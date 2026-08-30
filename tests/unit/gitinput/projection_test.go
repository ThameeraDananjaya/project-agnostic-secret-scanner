package gitinput_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/gitinput"
)

func TestExactRangeEdgesTreesMergesRenamesAndDeterminism(t *testing.T) {
	gitPath, digest := gitBinding(t)
	root := t.TempDir()
	source := filepath.Join(root, "source")
	if err := os.Mkdir(source, 0o700); err != nil {
		t.Fatal(err)
	}
	initFixtureRepo(t, gitPath, source)
	textCanary := `synthetic_api_key = "pscan_fixture_11111111111111111111111111111111"`
	binaryCanary := `synthetic_api_key = "pscan_fixture_22222222222222222222222222222222"`
	write(t, filepath.Join(source, "out-of-range.txt"), `synthetic_api_key = "pscan_fixture_00000000000000000000000000000000"`)
	runGit(t, gitPath, source, "add", "--", "out-of-range.txt")
	runGit(t, gitPath, source, "commit", "--quiet", "-m", "pre-range marker")
	if err := os.Remove(filepath.Join(source, "out-of-range.txt")); err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(source, "rename-old.txt"), "clean rename payload")
	write(t, filepath.Join(source, ".gitattributes"), "*.dat binary\n*.txt diff=hostile\n")
	runGit(t, gitPath, source, "add", "-A", "--", ".")
	runGit(t, gitPath, source, "commit", "--quiet", "-m", "range base")
	base := runGit(t, gitPath, source, "rev-parse", "HEAD")
	runGit(t, gitPath, source, "branch", "feature")
	runGit(t, gitPath, source, "mv", "--", "rename-old.txt", "rename-new.txt")
	runGit(t, gitPath, source, "commit", "--quiet", "-m", "rename in range")
	runGit(t, gitPath, source, "checkout", "--quiet", "feature")
	write(t, filepath.Join(source, "feature.txt"), "clean feature")
	runGit(t, gitPath, source, "add", "--", "feature.txt")
	runGit(t, gitPath, source, "commit", "--quiet", "-m", "feature commit")
	runGit(t, gitPath, source, "checkout", "--quiet", "main")
	runGit(t, gitPath, source, "merge", "--quiet", "--no-ff", "feature", "-m", "merge feature")
	write(t, filepath.Join(source, "deleted-text.txt"), textCanary)
	binary := append(append([]byte{0, 1, 2, 0}, []byte(binaryCanary)...), 0, 3, 4)
	if err := os.WriteFile(filepath.Join(source, "deleted-binary.dat"), binary, 0o600); err != nil {
		t.Fatal(err)
	}
	runGit(t, gitPath, source, "add", "--", "deleted-text.txt", "deleted-binary.dat")
	runGit(t, gitPath, source, "commit", "--quiet", "-m", "add deleted canaries")
	runGit(t, gitPath, source, "rm", "--quiet", "--", "deleted-text.txt", "deleted-binary.dat")
	runGit(t, gitPath, source, "commit", "--quiet", "-m", "delete canaries")
	head := runGit(t, gitPath, source, "rev-parse", "HEAD")

	private := filepath.Join(root, "private")
	bare := filepath.Join(root, "bare.git")
	g := gitinput.Git{Executable: gitPath, Digest: digest, Timeout: 15 * time.Second}
	if err := g.PrepareBare(context.Background(), source, bare, private); err != nil {
		t.Fatal(err)
	}
	limits := gitinput.ProjectionLimits{MaxBlobBytes: 1 << 20, MaxBlobCount: 1000, MaxTotalBytes: 32 << 20}
	planA, err := g.PlanRange(context.Background(), bare, base, head, false, filepath.Join(private, "home"), limits)
	if err != nil {
		t.Fatal(err)
	}
	planB, err := g.PlanRange(context.Background(), bare, base, head, false, filepath.Join(private, "home"), limits)
	if err != nil {
		t.Fatal(err)
	}
	if planA.PlanDigest != planB.PlanDigest || !reflect.DeepEqual(planA.Entries, planB.Entries) {
		t.Fatal("exact Git plan was not deterministic")
	}
	var sawText, sawBinary, sawRenameOld, sawRenameNew, sawMergeParentOne, sawMergeParentTwo, sawHeadTree bool
	for _, entry := range planA.Entries {
		if entry.Path == "out-of-range.txt" {
			t.Fatal("out-of-range-only blob entered the exact range or head tree")
		}
		switch entry.Path {
		case "deleted-text.txt":
			sawText = true
		case "deleted-binary.dat":
			sawBinary = true
		case "rename-old.txt":
			sawRenameOld = true
		case "rename-new.txt":
			sawRenameNew = true
		}
		if entry.ParentIndex == 1 && entry.Class != gitinput.ProjectionTrackedTree {
			sawMergeParentOne = sawMergeParentOne || strings.Contains(entry.RelativePath, "/p001/")
		}
		if entry.ParentIndex == 2 {
			sawMergeParentTwo = true
		}
		if entry.Class == gitinput.ProjectionTrackedTree && entry.Path == "rename-new.txt" {
			sawHeadTree = true
		}
	}
	if !sawText || !sawBinary || !sawRenameOld || !sawRenameNew || !sawMergeParentOne || !sawMergeParentTwo || !sawHeadTree {
		t.Fatalf("missing required exact coverage: text=%t binary=%t old=%t new=%t p1=%t p2=%t tree=%t", sawText, sawBinary, sawRenameOld, sawRenameNew, sawMergeParentOne, sawMergeParentTwo, sawHeadTree)
	}
	projectionA, err := g.Materialize(context.Background(), bare, filepath.Join(private, "home"), filepath.Join(root, "raw-a"), filepath.Join(root, "framed-a"), planA)
	if err != nil {
		t.Fatal(err)
	}
	projectionB, err := g.Materialize(context.Background(), bare, filepath.Join(private, "home"), filepath.Join(root, "raw-b"), filepath.Join(root, "framed-b"), planB)
	if err != nil {
		t.Fatal(err)
	}
	if projectionA.ContentDigest != projectionB.ContentDigest || projectionA.ProbeDigest != projectionB.ProbeDigest {
		t.Fatal("materialized coverage digests were not deterministic")
	}
	for i, entry := range planA.Entries {
		if entry.Path != "deleted-text.txt" && entry.Path != "deleted-binary.dat" {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(projectionA.ScanRoot, filepath.FromSlash(projectionA.Files[i].RelativePath)))
		if err != nil || entry.Path == "deleted-text.txt" && !bytes.Contains(raw, []byte(textCanary)) || entry.Path == "deleted-binary.dat" && !bytes.Contains(raw, []byte(binaryCanary)) {
			t.Fatalf("deleted canary bytes missing from %s", entry.Path)
		}
	}
	extra := filepath.Join(projectionA.ProbeRoot, "unbound.txt")
	if err := os.WriteFile(extra, []byte("unbound"), 0o600); err != nil {
		t.Fatal(err)
	}
	if gitinput.VerifyMaterializedProjection(projectionA) == nil {
		t.Fatal("unbound projected file did not fail completeness")
	}
	if err := os.Remove(extra); err != nil {
		t.Fatal(err)
	}
	missing := filepath.Join(projectionA.ScanRoot, filepath.FromSlash(projectionA.Files[0].RelativePath))
	if err := os.Remove(missing); err != nil {
		t.Fatal(err)
	}
	if gitinput.VerifyMaterializedProjection(projectionA) == nil {
		t.Fatal("missing projected file did not fail completeness")
	}
	mutated := filepath.Join(projectionB.ProbeRoot, filepath.FromSlash(projectionB.Files[0].RelativePath))
	if err := os.WriteFile(mutated, []byte("mutated"), 0o600); err != nil {
		t.Fatal(err)
	}
	if gitinput.VerifyMaterializedProjection(projectionB) == nil {
		t.Fatal("projection mutation did not fail integrity")
	}
}

func TestProjectionRejectsUnsafePathsModesAndAggregateLimits(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("unsafe Git path and mode fixture runs in pinned Linux Docker")
	}
	for _, test := range []struct {
		name   string
		mutate func(*testing.T, string, string)
		limits gitinput.ProjectionLimits
	}{
		{name: "portable-path", mutate: func(t *testing.T, gitPath, repository string) {
			write(t, filepath.Join(repository, "bad:name.txt"), "unsafe")
			runGit(t, gitPath, repository, "add", "--", "bad:name.txt")
		}},
		{name: "case-collision", mutate: func(t *testing.T, gitPath, repository string) {
			write(t, filepath.Join(repository, "Case.txt"), "one")
			write(t, filepath.Join(repository, "case.txt"), "two")
			runGit(t, gitPath, repository, "add", "--", "Case.txt", "case.txt")
		}},
		{name: "symlink-mode", mutate: func(t *testing.T, gitPath, repository string) {
			if err := os.Symlink("clean.txt", filepath.Join(repository, "link")); err != nil {
				t.Fatal(err)
			}
			runGit(t, gitPath, repository, "add", "--", "link")
		}},
		{name: "gitlink-mode", mutate: func(t *testing.T, gitPath, repository string) {
			oid := runGit(t, gitPath, repository, "rev-parse", "HEAD")
			runGit(t, gitPath, repository, "update-index", "--add", "--cacheinfo", "160000,"+oid+",vendor/module")
		}},
		{name: "blob-count", mutate: func(t *testing.T, gitPath, repository string) {
			write(t, filepath.Join(repository, "two.txt"), "two")
			runGit(t, gitPath, repository, "add", "--", "two.txt")
		}, limits: gitinput.ProjectionLimits{MaxBlobBytes: 1 << 20, MaxBlobCount: 1, MaxTotalBytes: 1 << 20}},
		{name: "total-bytes", mutate: func(t *testing.T, gitPath, repository string) {
			write(t, filepath.Join(repository, "large.txt"), strings.Repeat("x", 128))
			runGit(t, gitPath, repository, "add", "--", "large.txt")
		}, limits: gitinput.ProjectionLimits{MaxBlobBytes: 1 << 20, MaxBlobCount: 100, MaxTotalBytes: 16}},
	} {
		t.Run(test.name, func(t *testing.T) {
			gitPath, digest := gitBinding(t)
			root, source := t.TempDir(), ""
			source = filepath.Join(root, "source")
			if err := os.Mkdir(source, 0o700); err != nil {
				t.Fatal(err)
			}
			initFixtureRepo(t, gitPath, source)
			write(t, filepath.Join(source, "clean.txt"), "clean")
			runGit(t, gitPath, source, "add", "--", "clean.txt")
			runGit(t, gitPath, source, "commit", "--quiet", "-m", "base")
			base := runGit(t, gitPath, source, "rev-parse", "HEAD")
			test.mutate(t, gitPath, source)
			runGit(t, gitPath, source, "commit", "--quiet", "-m", "unsafe candidate")
			head := runGit(t, gitPath, source, "rev-parse", "HEAD")
			private, bare := filepath.Join(root, "private"), filepath.Join(root, "bare.git")
			g := gitinput.Git{Executable: gitPath, Digest: digest, Timeout: 15 * time.Second}
			if err := g.PrepareBare(context.Background(), source, bare, private); err != nil {
				t.Fatal(err)
			}
			limits := test.limits
			if limits.MaxBlobBytes == 0 {
				limits = gitinput.ProjectionLimits{MaxBlobBytes: 1 << 20, MaxBlobCount: 100, MaxTotalBytes: 8 << 20}
			}
			if _, err := g.PlanRange(context.Background(), bare, base, head, false, filepath.Join(private, "home"), limits); err == nil {
				t.Fatal("unsafe or over-limit input was admitted")
			}
		})
	}
}

func TestHostileGitConfigurationCannotChangeEnumeration(t *testing.T) {
	gitPath, digest := gitBinding(t)
	root := t.TempDir()
	source := filepath.Join(root, "source")
	if err := os.Mkdir(source, 0o700); err != nil {
		t.Fatal(err)
	}
	initFixtureRepo(t, gitPath, source)
	write(t, filepath.Join(source, "clean.txt"), "clean")
	runGit(t, gitPath, source, "add", "--", "clean.txt")
	runGit(t, gitPath, source, "commit", "--quiet", "-m", "base")
	base := runGit(t, gitPath, source, "rev-parse", "HEAD")
	write(t, filepath.Join(source, "changed.txt"), "changed")
	runGit(t, gitPath, source, "add", "--", "changed.txt")
	runGit(t, gitPath, source, "commit", "--quiet", "-m", "head")
	head := runGit(t, gitPath, source, "rev-parse", "HEAD")
	marker := filepath.Join(root, "must-not-exist")
	hostile := filepath.Join(root, "hostile.sh")
	if err := os.WriteFile(hostile, []byte("#!/bin/sh\ntouch '"+marker+"'\nexit 1\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	runGit(t, gitPath, source, "config", "core.hooksPath", filepath.Dir(hostile))
	runGit(t, gitPath, source, "config", "diff.external", hostile)
	runGit(t, gitPath, source, "config", "diff.hostile.command", hostile)
	globalConfig := filepath.Join(root, "global.gitconfig")
	if err := os.WriteFile(globalConfig, []byte("[diff]\n\texternal = "+hostile+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GIT_CONFIG_GLOBAL", globalConfig)
	private, bare := filepath.Join(root, "private"), filepath.Join(root, "bare.git")
	g := gitinput.Git{Executable: gitPath, Digest: digest, Timeout: 15 * time.Second}
	if err := g.PrepareBare(context.Background(), source, bare, private); err != nil {
		t.Fatal(err)
	}
	if _, err := g.PlanRange(context.Background(), bare, base, head, false, filepath.Join(private, "home"), gitinput.ProjectionLimits{MaxBlobBytes: 1 << 20, MaxBlobCount: 100, MaxTotalBytes: 8 << 20}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatal("hostile Git helper or hook executed")
	}
}

func initFixtureRepo(t *testing.T, gitPath, repository string) {
	t.Helper()
	runGit(t, gitPath, repository, "init", "--quiet", "-b", "main")
	runGit(t, gitPath, repository, "config", "user.email", "synthetic@example.invalid")
	runGit(t, gitPath, repository, "config", "user.name", "Synthetic Fixture")
}
