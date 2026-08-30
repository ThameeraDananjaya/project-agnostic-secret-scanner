//go:build linux

package workspace_test

import (
	"path/filepath"
	"testing"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/workspace"
)

func TestLinuxAcceptsBoundedAbsoluteWorkspaceRoot(t *testing.T) {
	root := filepath.Join(t.TempDir(), "scanner")
	manager, err := workspace.NewManager(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := manager.Close(); err != nil {
		t.Fatal(err)
	}
}
