//go:build !windows

package workspace

import "os"

func isReparse(os.FileInfo) bool { return false }
