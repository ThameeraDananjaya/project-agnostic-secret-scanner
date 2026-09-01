package gitinput

import "time"

type CoverageProfile struct {
	Name             string
	MaxBlobBytes     int64
	MaxBlobCount     int
	MaxTotalBytes    int64
	MaxDetectorFiles int
	MaxReportBytes   int64
	Timeout          time.Duration
	MaxMemoryBytes   int64
	MaxProcesses     int
}

var (
	PRProfile = CoverageProfile{
		Name: "pr", MaxBlobBytes: 512 << 20, MaxBlobCount: 100_000,
		MaxTotalBytes: 2 << 30, MaxDetectorFiles: 125_000,
		MaxReportBytes: 1 << 30, Timeout: 15 * time.Minute,
		MaxMemoryBytes: 4 << 30, MaxProcesses: 1,
	}
	ReleaseProfile = CoverageProfile{
		Name: "release", MaxBlobBytes: 512 << 20, MaxBlobCount: 100_000,
		MaxTotalBytes: 10 << 30, MaxDetectorFiles: 125_000,
		MaxReportBytes: 1 << 30, Timeout: 60 * time.Minute,
		MaxMemoryBytes: 8 << 30, MaxProcesses: 1,
	}
)

func (p CoverageProfile) Valid() bool {
	return (p == PRProfile || p == ReleaseProfile) && p.MaxBlobBytes > 0 && p.MaxBlobCount > 0 && p.MaxTotalBytes > 0 && p.MaxDetectorFiles > 0 && p.MaxReportBytes > 0 && p.Timeout > 0 && p.MaxMemoryBytes > 0 && p.MaxProcesses == 1
}

func (p CoverageProfile) ProjectionLimits() ProjectionLimits {
	return ProjectionLimits{MaxBlobBytes: p.MaxBlobBytes, MaxBlobCount: p.MaxBlobCount, MaxTotalBytes: p.MaxTotalBytes}
}

func (p CoverageProfile) Admit(blobBytes int64, blobCount, detectorFiles int) error {
	if !p.Valid() || blobBytes < 0 || blobBytes > p.MaxTotalBytes || blobCount < 1 || blobCount > p.MaxBlobCount || detectorFiles < 1 || detectorFiles > p.MaxDetectorFiles {
		return ErrResourceLimit
	}
	return nil
}

func (p CoverageProfile) AdmitBlob(size int64) error {
	if !p.Valid() || size < 0 || size > p.MaxBlobBytes {
		return ErrResourceLimit
	}
	return nil
}

func (p CoverageProfile) AdmitRuntime(timeout time.Duration, memoryBytes int64, processes int, reportBytes, detectorFileBytes int64) error {
	if !p.Valid() || timeout <= 0 || timeout > p.Timeout || memoryBytes <= 0 || memoryBytes > p.MaxMemoryBytes || processes < 1 || processes > p.MaxProcesses || reportBytes <= 0 || reportBytes > p.MaxReportBytes || detectorFileBytes < 0 || detectorFileBytes >= MaximumDetectorFile {
		return ErrResourceLimit
	}
	return nil
}
