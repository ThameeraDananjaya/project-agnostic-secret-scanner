//go:build linux

package request_test

import (
	"path/filepath"
	"testing"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/request"
)

func TestLinuxAcceptsCleanBoundedAbsoluteInputPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "manifest.json")
	if err := request.ValidateLocalPath(path); err != nil {
		t.Fatal(err)
	}
}
