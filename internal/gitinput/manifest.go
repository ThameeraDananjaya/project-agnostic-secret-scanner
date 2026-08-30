package gitinput

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const ManifestVersion = "1.0"

type Manifest struct {
	Version string          `json:"version"`
	Files   []ManifestEntry `json:"files"`
}

type ManifestEntry struct {
	Path   string `json:"path"`
	Size   int64  `json:"size"`
	Digest string `json:"digest"`
}

// VerifyManifest proves that a directory contains exactly the regular files in
// the manifest. It rejects links, special files, duplicates, unsorted paths,
// added files, missing files, size changes, and digest changes.
func VerifyManifest(root, manifestPath, expectedManifestDigest string) error {
	if !safeAbsolute(root) || !safeAbsolute(manifestPath) || len(expectedManifestDigest) != 64 {
		return errors.New("invalid manifest binding")
	}
	if err := engineVerify(manifestPath, expectedManifestDigest); err != nil {
		return err
	}
	raw, err := os.ReadFile(manifestPath)
	if err != nil || len(raw) > 16<<20 {
		return errors.New("manifest cannot be read")
	}
	decoder := json.NewDecoder(strings.NewReader(string(raw)))
	decoder.DisallowUnknownFields()
	var manifest Manifest
	if err := decoder.Decode(&manifest); err != nil || manifest.Version != ManifestVersion || len(manifest.Files) == 0 {
		return errors.New("invalid manifest")
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		return errors.New("manifest has trailing payload")
	}
	expected := make(map[string]ManifestEntry, len(manifest.Files))
	previous := ""
	for _, entry := range manifest.Files {
		if !safeRelative(entry.Path) || entry.Path <= previous || entry.Size < 0 || !lowerHexDigest(entry.Digest) {
			return errors.New("invalid or unsorted manifest entry")
		}
		previous = entry.Path
		expected[entry.Path] = entry
	}
	seen := make(map[string]bool, len(expected))
	err = filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return errors.New("tracked source cannot be traversed")
		}
		if path == root {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return errors.New("tracked source path cannot be normalized")
		}
		rel = filepath.ToSlash(rel)
		if entry.IsDir() {
			return nil
		}
		bound, ok := expected[rel]
		if !ok {
			return errors.New("unbound tracked-source file")
		}
		info, err := entry.Info()
		if err != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 || info.Size() != bound.Size {
			return errors.New("tracked-source file identity mismatch")
		}
		digest, err := fileDigest(path)
		if err != nil || digest != bound.Digest {
			return errors.New("tracked-source file digest mismatch")
		}
		seen[rel] = true
		return nil
	})
	if err != nil {
		return err
	}
	if len(seen) != len(expected) {
		return errors.New("manifest contains missing tracked-source file")
	}
	return nil
}

func BuildManifest(root string) (Manifest, error) {
	if !safeAbsolute(root) {
		return Manifest{}, errors.New("invalid root")
	}
	manifest := Manifest{Version: ManifestVersion, Files: []ManifestEntry{}}
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == root || entry.IsDir() {
			return nil
		}
		info, err := entry.Info()
		if err != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
			return errors.New("source contains a non-regular file")
		}
		rel, err := filepath.Rel(root, path)
		if err != nil || !safeRelative(filepath.ToSlash(rel)) {
			return errors.New("source contains an unsafe path")
		}
		digest, err := fileDigest(path)
		if err != nil {
			return err
		}
		manifest.Files = append(manifest.Files, ManifestEntry{Path: filepath.ToSlash(rel), Size: info.Size(), Digest: digest})
		return nil
	})
	if err != nil || len(manifest.Files) == 0 {
		return Manifest{}, errors.New("cannot build complete manifest")
	}
	sort.Slice(manifest.Files, func(i, j int) bool { return manifest.Files[i].Path < manifest.Files[j].Path })
	return manifest, nil
}

func safeRelative(path string) bool {
	return path != "" && path == filepath.ToSlash(filepath.Clean(filepath.FromSlash(path))) && path != "." && !strings.HasPrefix(path, "../") && !strings.HasPrefix(path, "/") && !strings.Contains(path, "\\") && !strings.ContainsRune(path, '\x00')
}

func fileDigest(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func lowerHexDigest(value string) bool {
	if len(value) != 64 {
		return false
	}
	for _, char := range value {
		if !(char >= '0' && char <= '9') && !(char >= 'a' && char <= 'f') {
			return false
		}
	}
	return true
}

// Isolated to keep manifest.go independent of process orchestration details.
func engineVerify(path, digest string) error {
	if !lowerHexDigest(digest) {
		return fmt.Errorf("invalid manifest digest binding")
	}
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
		return errors.New("manifest is not a regular file")
	}
	actual, err := fileDigest(path)
	if err != nil || actual != digest {
		return errors.New("manifest digest mismatch")
	}
	return nil
}
