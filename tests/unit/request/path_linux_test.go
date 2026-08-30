//go:build linux

package request_test

import (
	"path/filepath"
	"syscall"
	"testing"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/request"
)

func TestLinuxAcceptsCleanBoundedAbsoluteInputPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "manifest.json")
	if err := request.ValidateLocalPath(path); err != nil {
		t.Fatal(err)
	}
}

func TestLinuxRejectsFIFORequestWithoutOpening(t *testing.T) {
	path := filepath.Join(t.TempDir(), "request.pipe")
	if err := syscall.Mkfifo(path, 0o600); err != nil {
		t.Fatal(err)
	}
	if file, err := request.OpenLocal(path); err == nil {
		_ = file.Close()
		t.Fatal("FIFO request was accepted")
	}
}
