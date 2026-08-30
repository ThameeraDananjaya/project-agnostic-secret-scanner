//go:build windows

package contract_test

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWindowsRequestPathWithSpacesAndMetacharactersIsData(t *testing.T) {
	requestPath, workspaceRoot, _ := validRequestFile(t)
	directory := filepath.Join(filepath.Dir(requestPath), "safe & literal name")
	if err := os.Mkdir(directory, 0o700); err != nil {
		t.Fatal(err)
	}
	moved := filepath.Join(directory, "request [literal].json")
	if err := os.Rename(requestPath, moved); err != nil {
		t.Fatal(err)
	}
	code, _, _ := runCLI(t, "--request", moved, "--workspace-root", workspaceRoot)
	if code != 30 {
		t.Fatalf("literal Windows path was not handled as data: %d", code)
	}
}
