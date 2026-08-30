package gitinput_test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/gitinput"
)

func TestManifestRequiresExactCompleteRegularTree(t *testing.T) {
	root := filepath.Join(t.TempDir(), "tracked")
	if err := os.Mkdir(root, 0o700); err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(root, "alpha.txt"), "clean synthetic text")
	write(t, filepath.Join(root, "name;not-command.txt"), "candidate text")
	manifest, err := gitinput.BuildManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := json.Marshal(manifest)
	manifestPath := filepath.Join(t.TempDir(), "manifest.json")
	if err := os.WriteFile(manifestPath, raw, 0o600); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(raw)
	digest := hex.EncodeToString(sum[:])
	if err := gitinput.VerifyManifest(root, manifestPath, digest); err != nil {
		t.Fatalf("exact manifest failed: %v", err)
	}

	write(t, filepath.Join(root, "unbound.txt"), "must not be skipped")
	if err := gitinput.VerifyManifest(root, manifestPath, digest); err == nil {
		t.Fatal("unbound file did not fail completeness")
	}
	if err := os.Remove(filepath.Join(root, "unbound.txt")); err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(root, "alpha.txt"), "changed")
	if err := gitinput.VerifyManifest(root, manifestPath, digest); err == nil {
		t.Fatal("changed file did not fail integrity")
	}
}

func TestManifestRejectsUnsortedAndTraversalEntries(t *testing.T) {
	root := filepath.Join(t.TempDir(), "tracked")
	if err := os.Mkdir(root, 0o700); err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(root, "a"), "a")
	for _, path := range []string{"../outside", "b\\escape", "/absolute"} {
		manifest := gitinput.Manifest{Version: gitinput.ManifestVersion, Files: []gitinput.ManifestEntry{{Path: path, Size: 1, Digest: string(make([]byte, 64))}}}
		raw, _ := json.Marshal(manifest)
		manifestPath := filepath.Join(t.TempDir(), "manifest.json")
		if err := os.WriteFile(manifestPath, raw, 0o600); err != nil {
			t.Fatal(err)
		}
		sum := sha256.Sum256(raw)
		if err := gitinput.VerifyManifest(root, manifestPath, hex.EncodeToString(sum[:])); err == nil {
			t.Fatalf("unsafe manifest path accepted: %q", path)
		}
	}
}

func write(t *testing.T, path, value string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(value), 0o600); err != nil {
		t.Fatal(err)
	}
}
