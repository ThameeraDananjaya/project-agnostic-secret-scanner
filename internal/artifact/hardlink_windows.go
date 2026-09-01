//go:build windows

package artifact

import (
	"os"
	"syscall"
)

func hasMultipleLinks(path string, _ os.FileInfo) bool {
	pointer, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return true
	}
	handle, err := syscall.CreateFile(pointer, syscall.GENERIC_READ,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE|syscall.FILE_SHARE_DELETE,
		nil, syscall.OPEN_EXISTING, syscall.FILE_FLAG_OPEN_REPARSE_POINT, 0)
	if err != nil {
		return true
	}
	defer syscall.CloseHandle(handle)
	var information syscall.ByHandleFileInformation
	if syscall.GetFileInformationByHandle(handle, &information) != nil {
		return true
	}
	return information.NumberOfLinks > 1
}
