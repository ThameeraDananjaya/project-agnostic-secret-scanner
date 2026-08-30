package request

import (
	"errors"
	"strconv"
	"strings"
)

const (
	SupportedMajor = 1
	SupportedMinor = 0
)

type Version struct {
	Major int
	Minor int
}

func ParseVersion(value string) (Version, error) {
	parts := strings.Split(value, ".")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return Version{}, errors.New("schema version must be major.minor")
	}
	if (len(parts[0]) > 1 && parts[0][0] == '0') || (len(parts[1]) > 1 && parts[1][0] == '0') {
		return Version{}, errors.New("schema version has a leading zero")
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil || major < 0 {
		return Version{}, errors.New("invalid schema major")
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil || minor < 0 {
		return Version{}, errors.New("invalid schema minor")
	}
	return Version{Major: major, Minor: minor}, nil
}

func (v Version) Supported() bool   { return v.Major == SupportedMajor }
func (v Version) FutureMinor() bool { return v.Major == SupportedMajor && v.Minor > SupportedMinor }
