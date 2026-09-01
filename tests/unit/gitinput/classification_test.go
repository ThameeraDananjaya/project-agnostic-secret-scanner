package gitinput_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/gitinput"
)

func TestRawClassificationPrecedesPreparation(t *testing.T) {
	for _, test := range []struct {
		name, logical string
		data          []byte
		class         gitinput.RawClass
		unsupported   bool
	}{
		{name: "text", logical: "sample.txt", data: []byte("plain UTF-8\n"), class: gitinput.RawText},
		{name: "binary", logical: "sample.bin", data: []byte{0, 1, 2, 3, 0xff}, class: gitinput.RawBinary},
		{name: "zip-magic", logical: "disguised.bin", data: []byte{'P', 'K', 3, 4, 0, 0}, class: gitinput.RawArchive, unsupported: true},
		{name: "extension-only", logical: "named.tar", data: []byte("not an archive"), class: gitinput.RawArchive, unsupported: true},
		{name: "container", logical: "document.bin", data: []byte("%PDF-1.7"), class: gitinput.RawContainer, unsupported: true},
		{name: "ambiguous", logical: "double.zip", data: []byte{0x1f, 0x8b, 8, 0}, class: gitinput.RawAmbiguous, unsupported: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "raw")
			if err := os.WriteFile(path, test.data, 0o600); err != nil {
				t.Fatal(err)
			}
			class, err := gitinput.ClassifyRawFile(path, test.logical)
			if class != test.class || errors.Is(err, gitinput.ErrUnsupportedRawClass) != test.unsupported {
				t.Fatalf("class=%s err=%v", class, err)
			}
		})
	}
}

func TestNamedProfilesAtBelowAndAboveEveryResourceBound(t *testing.T) {
	for _, profile := range []gitinput.CoverageProfile{gitinput.PRProfile, gitinput.ReleaseProfile} {
		if !profile.Valid() {
			t.Fatalf("invalid named profile: %#v", profile)
		}
		for _, size := range []int64{profile.MaxBlobBytes - 1, profile.MaxBlobBytes} {
			if err := profile.AdmitBlob(size); err != nil {
				t.Fatalf("%s rejected blob boundary %d", profile.Name, size)
			}
		}
		if profile.AdmitBlob(profile.MaxBlobBytes+1) == nil {
			t.Fatalf("%s accepted blob above maximum", profile.Name)
		}
		if profile.Admit(profile.MaxTotalBytes-1, profile.MaxBlobCount-1, profile.MaxDetectorFiles-1) != nil ||
			profile.Admit(profile.MaxTotalBytes, profile.MaxBlobCount, profile.MaxDetectorFiles) != nil {
			t.Fatalf("%s rejected aggregate at/below boundary", profile.Name)
		}
		for _, values := range [][3]int64{
			{profile.MaxTotalBytes + 1, int64(profile.MaxBlobCount), int64(profile.MaxDetectorFiles)},
			{profile.MaxTotalBytes, int64(profile.MaxBlobCount + 1), int64(profile.MaxDetectorFiles)},
			{profile.MaxTotalBytes, int64(profile.MaxBlobCount), int64(profile.MaxDetectorFiles + 1)},
		} {
			if profile.Admit(values[0], int(values[1]), int(values[2])) == nil {
				t.Fatalf("%s accepted aggregate above maximum", profile.Name)
			}
		}
		if profile.AdmitRuntime(profile.Timeout-time.Second, profile.MaxMemoryBytes-1, 1, profile.MaxReportBytes-1, gitinput.MaximumDetectorFile-2) != nil ||
			profile.AdmitRuntime(profile.Timeout, profile.MaxMemoryBytes, profile.MaxProcesses, profile.MaxReportBytes, gitinput.MaximumDetectorFile-1) != nil {
			t.Fatalf("%s rejected runtime at/below boundary", profile.Name)
		}
		if profile.AdmitRuntime(profile.Timeout+time.Nanosecond, profile.MaxMemoryBytes, 1, profile.MaxReportBytes, gitinput.MaximumDetectorFile-1) == nil ||
			profile.AdmitRuntime(profile.Timeout, profile.MaxMemoryBytes+1, 1, profile.MaxReportBytes, gitinput.MaximumDetectorFile-1) == nil ||
			profile.AdmitRuntime(profile.Timeout, profile.MaxMemoryBytes, profile.MaxProcesses+1, profile.MaxReportBytes, gitinput.MaximumDetectorFile-1) == nil ||
			profile.AdmitRuntime(profile.Timeout, profile.MaxMemoryBytes, 1, profile.MaxReportBytes+1, gitinput.MaximumDetectorFile-1) == nil ||
			profile.AdmitRuntime(profile.Timeout, profile.MaxMemoryBytes, 1, profile.MaxReportBytes, gitinput.MaximumDetectorFile) == nil {
			t.Fatalf("%s accepted runtime above maximum", profile.Name)
		}
	}
}
