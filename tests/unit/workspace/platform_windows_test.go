//go:build windows

package workspace_test

import (
	"path/filepath"
	"testing"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/workspace"
)

func TestWindowsRejectsUNCWorkspaceRoot(t *testing.T) {
	if _, err := workspace.NewManager(`\\server\share\scanner`); err == nil {
		t.Fatal("accepted UNC root")
	}
	if filepath.Separator != '\\' {
		t.Fatal("unexpected Windows separator")
	}
}
