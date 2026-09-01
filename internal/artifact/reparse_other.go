//go:build !windows

package artifact

import "os"

func isReparse(os.FileInfo) bool { return false }
