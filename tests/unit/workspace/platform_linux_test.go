//go:build linux

package workspace_test

import (
	"path/filepath"
	"testing"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/workspace"
)

func TestLinuxAcceptsBoundedAbsoluteWorkspaceRoot(t *testing.T) {
	root := filepath.Join(t.TempDir(), "scanner")
	if _, err := workspace.NewManager(root); err != nil {
		t.Fatal(err)
	}
}
