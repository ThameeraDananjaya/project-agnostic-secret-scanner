// Package workspace creates one deterministic, isolated workspace per scan attempt.
package workspace

import (
	"bytes"
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/request"
)

type Manager struct {
	base     string
	root     *os.Root
	identity os.FileInfo
}

type Workspace struct {
	path     string
	name     string
	lease    string
	root     *os.Root
	manager  *Manager
	identity os.FileInfo
	marker   []byte
}

const markerPrefix = "PSCAN_WORKSPACE_V1\n"

func leaseName(scanID string, attempt int) string {
	return ".pscan-lease-" + scanID + "-attempt-" + strconv.Itoa(attempt)
}

func workspaceName(scanID string, attempt int) string {
	return "scan-" + scanID + "-attempt-" + strconv.Itoa(attempt)
}

func NewManager(base string) (*Manager, error) {
	if base == "" || !filepath.IsAbs(base) || filepath.Clean(base) != base || isRoot(base) || unsafeWindowsRoot(base) {
		return nil, errors.New("unsafe workspace root")
	}
	if err := createSafeDirectory(base); err != nil {
		return nil, errors.New("workspace root unavailable")
	}
	if err := rejectLinks(base); err != nil {
		return nil, err
	}
	identity, err := os.Lstat(base)
	if err != nil || !identity.IsDir() || identity.Mode()&os.ModeSymlink != 0 || isReparse(identity) {
		return nil, errors.New("workspace root unavailable")
	}
	root, err := os.OpenRoot(base)
	if err != nil {
		return nil, errors.New("workspace root unavailable")
	}
	opened, err := root.Stat(".")
	if err != nil || !os.SameFile(identity, opened) {
		_ = root.Close()
		return nil, errors.New("workspace root identity changed")
	}
	manager := &Manager{base: base, root: root, identity: identity}
	if err := manager.validateRoot(); err != nil {
		_ = root.Close()
		return nil, err
	}
	return manager, nil
}

func createSafeDirectory(path string) error {
	current := path
	missing := []string{}
	for {
		info, err := os.Lstat(current)
		if err == nil {
			if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 || isReparse(info) {
				return errors.New("workspace ancestor is unsafe")
			}
			break
		}
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		parent := filepath.Dir(current)
		if parent == current {
			return errors.New("workspace ancestor unavailable")
		}
		missing = append(missing, filepath.Base(current))
		current = parent
	}
	if err := rejectLinks(current); err != nil {
		return err
	}
	for i := len(missing) - 1; i >= 0; i-- {
		current = filepath.Join(current, missing[i])
		if err := os.Mkdir(current, 0o700); err != nil {
			return err
		}
		info, err := os.Lstat(current)
		if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 || isReparse(info) {
			return errors.New("workspace component is unsafe")
		}
	}
	return nil
}

func DefaultRoot() string { return filepath.Join(os.TempDir(), "project-agnostic-secret-scanner") }

func (m *Manager) Close() error { return m.root.Close() }

func (m *Manager) Create(scanID string, attempt int) (*Workspace, error) {
	if !request.IsUUIDv4(scanID) || attempt < 1 || attempt > 3 {
		return nil, errors.New("invalid workspace identity")
	}
	name := workspaceName(scanID, attempt)
	lease := leaseName(scanID, attempt)
	path := filepath.Join(m.base, name)
	if err := confined(m.base, path); err != nil {
		return nil, err
	}
	if err := m.validateRoot(); err != nil {
		return nil, err
	}
	marker := make([]byte, len(markerPrefix)+32)
	copy(marker, markerPrefix)
	if _, err := rand.Read(marker[len(markerPrefix):]); err != nil {
		return nil, errors.New("workspace ownership marker unavailable")
	}
	leaseFile, err := m.root.OpenFile(lease, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return nil, errors.New("workspace collision or unavailable")
	}
	if _, err := leaseFile.Write(marker); err != nil {
		_ = leaseFile.Close()
		_ = m.root.Remove(lease)
		return nil, errors.New("workspace ownership lease unavailable")
	}
	if err := leaseFile.Close(); err != nil {
		_ = m.root.Remove(lease)
		return nil, errors.New("workspace ownership lease unavailable")
	}
	if err := m.root.Mkdir(name, 0o700); err != nil {
		_ = m.root.Remove(lease)
		return nil, errors.New("workspace collision or unavailable")
	}
	for _, child := range []string{"private", "input", "output"} {
		if err := m.root.Mkdir(filepath.Join(name, child), 0o700); err != nil {
			_ = m.root.RemoveAll(name)
			_ = m.root.Remove(lease)
			return nil, errors.New("workspace initialization failed")
		}
	}
	if err := m.root.WriteFile(filepath.Join(name, ".pscan-owned"), marker, 0o600); err != nil {
		_ = m.root.RemoveAll(name)
		_ = m.root.Remove(lease)
		return nil, errors.New("workspace ownership marker unavailable")
	}
	identity, err := m.root.Lstat(name)
	if err != nil {
		_ = m.root.RemoveAll(name)
		_ = m.root.Remove(lease)
		return nil, errors.New("workspace identity unavailable")
	}
	if err := m.validateRoot(); err != nil {
		_ = m.root.RemoveAll(name)
		_ = m.root.Remove(lease)
		return nil, err
	}
	return &Workspace{path: path, name: name, lease: lease, root: m.root, manager: m, identity: identity, marker: marker}, nil
}

func (w *Workspace) Path() string { return w.path }

func (w *Workspace) Cleanup() error {
	if err := w.manager.validateRoot(); err != nil {
		return err
	}
	current, err := w.root.Lstat(w.name)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil || current.Mode()&os.ModeSymlink != 0 || isReparse(current) || !os.SameFile(w.identity, current) {
		return errors.New("workspace identity changed")
	}
	markerPath := filepath.Join(w.name, ".pscan-owned")
	markerInfo, err := w.root.Lstat(markerPath)
	if err != nil || !markerInfo.Mode().IsRegular() || markerInfo.Mode()&os.ModeSymlink != 0 || isReparse(markerInfo) || markerInfo.Size() != int64(len(w.marker)) {
		return errors.New("workspace ownership marker changed")
	}
	marker, err := w.root.ReadFile(markerPath)
	if err != nil || subtle.ConstantTimeCompare(marker, w.marker) != 1 {
		return errors.New("workspace ownership marker changed")
	}
	leaseInfo, err := w.root.Lstat(w.lease)
	if err != nil || !leaseInfo.Mode().IsRegular() || leaseInfo.Mode()&os.ModeSymlink != 0 || isReparse(leaseInfo) || leaseInfo.Size() != int64(len(w.marker)) {
		return errors.New("workspace ownership lease changed")
	}
	lease, err := w.root.ReadFile(w.lease)
	if err != nil || subtle.ConstantTimeCompare(lease, w.marker) != 1 {
		return errors.New("workspace ownership lease changed")
	}
	if err := w.root.RemoveAll(w.name); err != nil {
		return errors.New("workspace cleanup failed")
	}
	if err := w.root.Remove(w.lease); err != nil {
		return errors.New("workspace lease cleanup failed")
	}
	return nil
}

// Recover removes only abandoned workspaces carrying the scanner-owned marker.
// It is safe to call after a process crash and fails closed on any unexpected
// name, link, reparse point, special file or marker.
func (m *Manager) Recover() error {
	if err := m.validateRoot(); err != nil {
		return err
	}
	directory, err := os.Open(m.base)
	if err != nil {
		return errors.New("workspace recovery unavailable")
	}
	opened, err := directory.Stat()
	if err != nil || !os.SameFile(m.identity, opened) {
		_ = directory.Close()
		return errors.New("workspace recovery unavailable")
	}
	entries, err := directory.ReadDir(-1)
	closeErr := directory.Close()
	if err != nil || closeErr != nil {
		return errors.New("workspace recovery unavailable")
	}
	entryByName := make(map[string]os.DirEntry, len(entries))
	for _, entry := range entries {
		entryByName[entry.Name()] = entry
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".pscan-lease-") {
			if !entry.Type().IsRegular() {
				return errors.New("workspace recovery encountered unsafe lease")
			}
			continue
		}
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "scan-") {
			return errors.New("workspace recovery encountered unowned material")
		}
		lease := ".pscan-lease-" + strings.TrimPrefix(entry.Name(), "scan-")
		if _, ok := entryByName[lease]; !ok {
			return errors.New("workspace recovery encountered unowned material")
		}
		if err := m.recoverOne(entry.Name(), lease); err != nil {
			return err
		}
		delete(entryByName, lease)
		delete(entryByName, entry.Name())
	}
	for name, entry := range entryByName {
		if !strings.HasPrefix(name, ".pscan-lease-") || !entry.Type().IsRegular() {
			return errors.New("workspace recovery encountered unowned material")
		}
		marker, err := m.root.ReadFile(name)
		if err != nil || len(marker) != len(markerPrefix)+32 || !bytes.HasPrefix(marker, []byte(markerPrefix)) {
			return errors.New("workspace recovery marker mismatch")
		}
		if err := m.root.Remove(name); err != nil {
			return errors.New("workspace recovery cleanup failed")
		}
	}
	return nil
}

func (m *Manager) recoverOne(name, leaseName string) error {
	info, err := m.root.Lstat(name)
	if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 || isReparse(info) {
		return errors.New("workspace recovery identity mismatch")
	}
	markerPath := filepath.Join(name, ".pscan-owned")
	markerInfo, err := m.root.Lstat(markerPath)
	if err != nil || !markerInfo.Mode().IsRegular() || markerInfo.Mode()&os.ModeSymlink != 0 || isReparse(markerInfo) || markerInfo.Size() != int64(len(markerPrefix)+32) {
		return errors.New("workspace recovery marker mismatch")
	}
	marker, err := m.root.ReadFile(markerPath)
	if err != nil || !bytes.HasPrefix(marker, []byte(markerPrefix)) {
		return errors.New("workspace recovery marker mismatch")
	}
	leaseInfo, err := m.root.Lstat(leaseName)
	if err != nil || !leaseInfo.Mode().IsRegular() || leaseInfo.Mode()&os.ModeSymlink != 0 || isReparse(leaseInfo) || leaseInfo.Size() != int64(len(markerPrefix)+32) {
		return errors.New("workspace recovery lease mismatch")
	}
	lease, err := m.root.ReadFile(leaseName)
	if err != nil || subtle.ConstantTimeCompare(marker, lease) != 1 {
		return errors.New("workspace recovery lease mismatch")
	}
	if err := verifyRecoverableTree(m.base, filepath.Join(m.base, name)); err != nil {
		return err
	}
	if err := m.root.RemoveAll(name); err != nil {
		return errors.New("workspace recovery cleanup failed")
	}
	if err := m.root.Remove(leaseName); err != nil {
		return errors.New("workspace recovery cleanup failed")
	}
	return nil
}

func verifyRecoverableTree(base, root string) error {
	if err := confined(base, root); err != nil {
		return err
	}
	return filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return errors.New("workspace recovery traversal failed")
		}
		info, err := os.Lstat(path)
		if err != nil || info.Mode()&os.ModeSymlink != 0 || isReparse(info) || (!info.IsDir() && !info.Mode().IsRegular()) {
			return errors.New("workspace recovery encountered unsafe material")
		}
		return nil
	})
}

func (m *Manager) validateRoot() error {
	if err := rejectLinks(m.base); err != nil {
		return errors.New("workspace root identity changed")
	}
	current, err := os.Lstat(m.base)
	if err != nil || !current.IsDir() || current.Mode()&os.ModeSymlink != 0 || isReparse(current) || !os.SameFile(m.identity, current) {
		return errors.New("workspace root identity changed")
	}
	opened, err := m.root.Stat(".")
	if err != nil || !os.SameFile(m.identity, opened) {
		return errors.New("workspace root identity changed")
	}
	return nil
}

func confined(base, path string) error {
	relative, err := filepath.Rel(base, path)
	if err != nil || relative == "." || relative == "" || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) || filepath.IsAbs(relative) {
		return errors.New("workspace escaped root")
	}
	return nil
}

func isRoot(path string) bool {
	volume := filepath.VolumeName(path)
	return filepath.Clean(path) == filepath.Clean(volume+string(filepath.Separator)) || filepath.Clean(path) == string(filepath.Separator)
}

func unsafeWindowsRoot(path string) bool {
	lower := strings.ToLower(path)
	return strings.HasPrefix(lower, `\\`) || strings.HasPrefix(lower, `\\?\`) || strings.HasPrefix(lower, `\\.\`)
}

func rejectLinks(path string) error {
	current := filepath.Clean(path)
	for {
		info, err := os.Lstat(current)
		if err != nil {
			return errors.New("workspace component unavailable")
		}
		if info.Mode()&os.ModeSymlink != 0 || isReparse(info) {
			return errors.New("workspace root contains a link")
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	return nil
}

func PlatformDescription() string { return fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH) }
