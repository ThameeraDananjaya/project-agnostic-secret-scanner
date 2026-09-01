package artifact

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func safeArchivePath(name string) bool {
	if name == "" || strings.ContainsRune(name, 0) || strings.Contains(name, "\\") ||
		strings.HasPrefix(name, "/") || filepath.IsAbs(name) || filepath.VolumeName(name) != "" {
		return false
	}
	clean := filepath.ToSlash(filepath.Clean(filepath.FromSlash(name)))
	return clean != "." && clean == strings.TrimSuffix(name, "/") && clean != ".." &&
		!strings.HasPrefix(clean, "../") && !strings.Contains(clean, ":")
}

func safeAbsolute(path string) bool {
	return path != "" && filepath.IsAbs(path) && filepath.Clean(path) == path &&
		filepath.Dir(path) != path && !strings.HasPrefix(strings.ToLower(path), `\\`) &&
		!strings.HasPrefix(strings.ToLower(path), `\\?\`) && !strings.HasPrefix(strings.ToLower(path), `\\.\`)
}

func confined(root, path string) bool {
	relative, err := filepath.Rel(root, path)
	return err == nil && relative != "." && relative != "" && relative != ".." &&
		!filepath.IsAbs(relative) && !strings.HasPrefix(relative, ".."+string(filepath.Separator))
}

func regularNoLink(path string) (os.FileInfo, error) {
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 || isReparse(info) || hasMultipleLinks(path, info) {
		return nil, Rejection{Code: RejectUnsafe}
	}
	return info, nil
}

func rejectTreeHazards(root string) error {
	return filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return Rejection{Code: RejectUnsafe}
		}
		info, err := os.Lstat(path)
		if err != nil || info.Mode()&os.ModeSymlink != 0 || isReparse(info) || info.Mode().IsRegular() && hasMultipleLinks(path, info) {
			return Rejection{Code: RejectUnsafe}
		}
		if path != root && !info.IsDir() && !info.Mode().IsRegular() {
			return Rejection{Code: RejectUnsafe}
		}
		return nil
	})
}

func createPrivateRoot(path string) error {
	if !safeAbsolute(path) {
		return Rejection{Code: RejectInvariant}
	}
	if _, err := os.Lstat(path); !errors.Is(err, os.ErrNotExist) {
		return Rejection{Code: RejectInvariant}
	}
	for current := filepath.Dir(path); ; current = filepath.Dir(current) {
		info, err := os.Lstat(current)
		if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 || isReparse(info) {
			return Rejection{Code: RejectInvariant}
		}
		if parent := filepath.Dir(current); parent == current {
			break
		}
	}
	if err := os.Mkdir(path, 0o700); err != nil {
		return Rejection{Code: RejectInvariant}
	}
	return nil
}
