// Package workspace creates one deterministic, isolated workspace per scan attempt.
package workspace

import (
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
	root     *os.Root
	manager  *Manager
	identity os.FileInfo
	marker   [32]byte
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
	if err != nil || !identity.IsDir() || identity.Mode()&os.ModeSymlink != 0 {
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
			if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
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
		if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
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
	name := "scan-" + scanID + "-attempt-" + strconv.Itoa(attempt)
	path := filepath.Join(m.base, name)
	if err := confined(m.base, path); err != nil {
		return nil, err
	}
	if err := m.validateRoot(); err != nil {
		return nil, err
	}
	if err := m.root.Mkdir(name, 0o700); err != nil {
		return nil, errors.New("workspace collision or unavailable")
	}
	for _, child := range []string{"private", "input", "output"} {
		if err := m.root.Mkdir(filepath.Join(name, child), 0o700); err != nil {
			_ = m.root.RemoveAll(name)
			return nil, errors.New("workspace initialization failed")
		}
	}
	var marker [32]byte
	if _, err := rand.Read(marker[:]); err != nil {
		_ = m.root.RemoveAll(name)
		return nil, errors.New("workspace ownership marker unavailable")
	}
	if err := m.root.WriteFile(filepath.Join(name, ".pscan-owned"), marker[:], 0o600); err != nil {
		_ = m.root.RemoveAll(name)
		return nil, errors.New("workspace ownership marker unavailable")
	}
	identity, err := m.root.Lstat(name)
	if err != nil {
		_ = m.root.RemoveAll(name)
		return nil, errors.New("workspace identity unavailable")
	}
	if err := m.validateRoot(); err != nil {
		_ = m.root.RemoveAll(name)
		return nil, err
	}
	return &Workspace{path: path, name: name, root: m.root, manager: m, identity: identity, marker: marker}, nil
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
	if err != nil || current.Mode()&os.ModeSymlink != 0 || !os.SameFile(w.identity, current) {
		return errors.New("workspace identity changed")
	}
	markerPath := filepath.Join(w.name, ".pscan-owned")
	markerInfo, err := w.root.Lstat(markerPath)
	if err != nil || !markerInfo.Mode().IsRegular() || markerInfo.Mode()&os.ModeSymlink != 0 || markerInfo.Size() != int64(len(w.marker)) {
		return errors.New("workspace ownership marker changed")
	}
	marker, err := w.root.ReadFile(markerPath)
	if err != nil || subtle.ConstantTimeCompare(marker, w.marker[:]) != 1 {
		return errors.New("workspace ownership marker changed")
	}
	if err := w.root.RemoveAll(w.name); err != nil {
		return errors.New("workspace cleanup failed")
	}
	return nil
}

func (m *Manager) validateRoot() error {
	if err := rejectLinks(m.base); err != nil {
		return errors.New("workspace root identity changed")
	}
	current, err := os.Lstat(m.base)
	if err != nil || !current.IsDir() || current.Mode()&os.ModeSymlink != 0 || !os.SameFile(m.identity, current) {
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
		if info.Mode()&os.ModeSymlink != 0 {
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
