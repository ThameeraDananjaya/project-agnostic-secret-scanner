package request

import (
	"errors"
	"os"
	"time"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/outcome"
)

// OpenLocal binds request input to one verified regular-file identity.
func OpenLocal(path string) (*os.File, error) {
	if err := ValidateLocalPath(path); err != nil {
		return nil, err
	}
	before, err := os.Lstat(path)
	if err != nil {
		return nil, reject(outcome.ReasonUnavailableRequiredInput)
	}
	if !safeRequestFile(before) {
		return nil, reject(outcome.ReasonFailInputIntegrity)
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, reject(outcome.ReasonUnavailableRequiredInput)
	}
	after, err := file.Stat()
	if err != nil || !safeRequestFile(after) || !os.SameFile(before, after) {
		_ = file.Close()
		return nil, reject(outcome.ReasonFailInputIntegrity)
	}
	return file, nil
}

// OpenDescriptor accepts only a bounded regular file or a deadline-capable pipe.
func OpenDescriptor(fd uintptr, timeout time.Duration) (*os.File, error) {
	file := os.NewFile(fd, "request")
	if file == nil {
		return nil, errors.New("descriptor unavailable")
	}
	info, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return nil, errors.New("descriptor unavailable")
	}
	if info.Mode().IsRegular() {
		if info.Size() < 0 || info.Size() > MaxPayloadBytes {
			_ = file.Close()
			return nil, errors.New("descriptor input is not bounded")
		}
		return file, nil
	}
	if info.Mode()&os.ModeNamedPipe == 0 || timeout <= 0 {
		_ = file.Close()
		return nil, errors.New("unsafe descriptor type")
	}
	if err := file.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		_ = file.Close()
		return nil, errors.New("descriptor cannot be bounded")
	}
	return file, nil
}

func safeRequestFile(info os.FileInfo) bool {
	return info.Mode().IsRegular() && info.Mode()&os.ModeSymlink == 0 && info.Size() >= 0 && info.Size() <= MaxPayloadBytes
}
