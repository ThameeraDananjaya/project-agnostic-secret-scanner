//go:build !unix && !windows

package artifact

import "os"

func hasMultipleLinks(string, os.FileInfo) bool { return false }
