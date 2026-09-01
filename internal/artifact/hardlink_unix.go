//go:build unix

package artifact

import (
	"os"
	"syscall"
)

func hasMultipleLinks(_ string, info os.FileInfo) bool {
	stat, ok := info.Sys().(*syscall.Stat_t)
	return ok && stat.Nlink > 1
}
