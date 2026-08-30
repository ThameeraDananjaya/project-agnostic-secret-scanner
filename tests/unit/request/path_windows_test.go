//go:build windows

package request_test

import (
	"path/filepath"
	"testing"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/request"
)

func TestWindowsRejectsDeviceUNCAndAlternateStreamPaths(t *testing.T) {
	root := filepath.VolumeName(t.TempDir()) + string(filepath.Separator)
	for _, path := range []string{
		`\\server\share\manifest.json`, `\\?\C:\safe\manifest.json`,
		filepath.Join(root, "safe", "CON.json"),
		filepath.Join(root, "safe", "file.json:stream"),
		filepath.Join(root, "safe", "trailing. "),
	} {
		if err := request.ValidateLocalPath(path); err == nil {
			t.Errorf("accepted unsafe Windows path %q", path)
		}
	}
}
