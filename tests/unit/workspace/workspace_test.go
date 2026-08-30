package workspace_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/workspace"
)

const scanID = "123e4567-e89b-42d3-a456-426614174000"

func manager(t *testing.T) (*workspace.Manager, string) {
	t.Helper()
	base := filepath.Join(t.TempDir(), "workspaces")
	value, err := workspace.NewManager(base)
	if err != nil {
		t.Fatal(err)
	}
	return value, base
}

func TestDeterministicWorkspaceAndCollision(t *testing.T) {
	m, base := manager(t)
	space, err := m.Create(scanID, 1)
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(base, "scan-"+scanID+"-attempt-1")
	if space.Path() != want {
		t.Fatalf("got %q want %q", space.Path(), want)
	}
	if _, err := m.Create(scanID, 1); err == nil {
		t.Fatal("workspace collision was overwritten")
	}
	for _, child := range []string{"private", "input", "output"} {
		if info, err := os.Stat(filepath.Join(space.Path(), child)); err != nil || !info.IsDir() {
			t.Fatalf("missing %s", child)
		}
	}
	if err := space.Cleanup(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(space.Path()); !os.IsNotExist(err) {
		t.Fatal("workspace residue remains")
	}
}

func TestSequentialAttemptsAreIsolated(t *testing.T) {
	m, _ := manager(t)
	one, err := m.Create(scanID, 1)
	if err != nil {
		t.Fatal(err)
	}
	two, err := m.Create(scanID, 2)
	if err != nil {
		t.Fatal(err)
	}
	if one.Path() == two.Path() {
		t.Fatal("attempts collided")
	}
	if err := one.Cleanup(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(two.Path()); err != nil {
		t.Fatal("cleanup crossed attempt boundary")
	}
	if err := two.Cleanup(); err != nil {
		t.Fatal(err)
	}
}

func TestParallelScansDoNotCollide(t *testing.T) {
	m, base := manager(t)
	var wait sync.WaitGroup
	errors := make(chan error, 64)
	paths := make(chan string, 64)
	for i := range 64 {
		wait.Add(1)
		go func(i int) {
			defer wait.Done()
			id := fmt.Sprintf("123e4567-e89b-4%03x-a456-%012x", i, i)
			space, err := m.Create(id, 1)
			if err != nil {
				errors <- err
				return
			}
			paths <- space.Path()
		}(i)
	}
	wait.Wait()
	close(errors)
	close(paths)
	for err := range errors {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for path := range paths {
		if seen[path] || !strings.HasPrefix(path, base+string(filepath.Separator)) {
			t.Fatalf("collision or escape: %q", path)
		}
		seen[path] = true
	}
	if len(seen) != 64 {
		t.Fatalf("got %d workspaces", len(seen))
	}
}

func TestUnsafeRootsAndIdentityReplacementFailClosed(t *testing.T) {
	volumeRoot := filepath.VolumeName(t.TempDir()) + string(filepath.Separator)
	if _, err := workspace.NewManager(volumeRoot); err == nil {
		t.Fatal("accepted filesystem root")
	}
	m, _ := manager(t)
	space, err := m.Create(scanID, 1)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(space.Path()); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(space.Path(), 0o700); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(space.Path(), "must-survive")
	if err := os.WriteFile(marker, []byte("replacement"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := space.Cleanup(); err == nil {
		t.Fatal("identity replacement was deleted")
	}
	if _, err := os.Stat(marker); err != nil {
		t.Fatal("replacement content was removed")
	}
}

func TestWorkspaceRootRejectsSymlinkWhenSupported(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "target")
	if err := os.Mkdir(target, 0o700); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "link")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	if _, err := workspace.NewManager(link); err == nil {
		t.Fatal("accepted symlink workspace root")
	}
}

func TestInvalidWorkspaceIdentity(t *testing.T) {
	m, _ := manager(t)
	for _, id := range []string{"", "source-derived", "123e4567-e89b-12d3-a456-426614174000"} {
		if _, err := m.Create(id, 1); err == nil {
			t.Errorf("accepted %q", id)
		}
	}
	if _, err := m.Create(scanID, 4); err == nil {
		t.Fatal("accepted attempt 4")
	}
	if workspace.PlatformDescription() == "/" {
		t.Fatal("missing platform description")
	}
}
