// Package artifact normalizes untrusted release artifacts without executing
// candidate-controlled content. All detailed identities stay private.
package artifact

import (
	"errors"
	"time"
)

const (
	NormalizationVersion = "pscan.artifact.v1"
	MaximumRuleSpan      = int64(4020)
	DetectorPayloadSize  = int64(90_000)
	DetectorOverlap      = MaximumRuleSpan - 1
	MaximumDetectorFile  = int64(100_000)
)

type Class string

const (
	ClassDirectory  Class = "directory"
	ClassFile       Class = "file"
	ClassZIP        Class = "zip"
	ClassTAR        Class = "tar"
	ClassGZIP       Class = "gzip"
	ClassTGZ        Class = "tgz"
	ClassOCILayout  Class = "oci-layout"
	ClassDockerSave Class = "docker-save"
)

type RejectCode string

const (
	RejectUnsupported RejectCode = "UNSUPPORTED_INPUT"
	RejectUnsafe      RejectCode = "UNSAFE_INPUT"
	RejectMalformed   RejectCode = "MALFORMED_INPUT"
	RejectResource    RejectCode = "RESOURCE_LIMIT"
	RejectTimeout     RejectCode = "TIMEOUT"
	RejectCancelled   RejectCode = "CANCELLED"
	RejectInvariant   RejectCode = "INTERNAL_INVARIANT"
)

// Rejection intentionally exposes no candidate-controlled detail.
type Rejection struct{ Code RejectCode }

func (r Rejection) Error() string { return "artifact rejected" }

func IsCode(err error, code RejectCode) bool {
	var rejection Rejection
	return errors.As(err, &rejection) && rejection.Code == code
}

type Limits struct {
	MaxDepth            int
	MaxEntries          int
	MaxExpandedBytes    int64
	MaxFileBytes        int64
	MaxCompressionRatio int64
	Timeout             time.Duration
	MaxMemoryBytes      int64
	MaxProcesses        int
}

type Usage struct {
	Depth         int
	Entries       int
	ExpandedBytes int64
	FileBytes     int64
	Expanded      int64
	Compressed    int64
	Elapsed       time.Duration
	MemoryBytes   int64
	Processes     int
}

// CheckUsage is the single arithmetic boundary for exact at/below/above
// profile tests and runtime admission. It allocates no candidate-sized data.
func (l Limits) CheckUsage(u Usage) error {
	if !l.internallyValid() || u.Depth < 0 || u.Entries < 0 || u.ExpandedBytes < 0 || u.FileBytes < 0 ||
		u.Expanded < 0 || u.Compressed < 0 || u.Elapsed < 0 || u.MemoryBytes < 0 || u.Processes < 0 {
		return Rejection{Code: RejectInvariant}
	}
	if u.Depth > l.MaxDepth || u.Entries > l.MaxEntries || u.ExpandedBytes > l.MaxExpandedBytes ||
		u.FileBytes > l.MaxFileBytes || u.Elapsed > l.Timeout || u.MemoryBytes > l.MaxMemoryBytes || u.Processes > l.MaxProcesses ||
		exceedsRatio(u.Expanded, u.Compressed, l.MaxCompressionRatio) {
		return Rejection{Code: RejectResource}
	}
	return nil
}

func PRLimits() Limits {
	return Limits{MaxDepth: 5, MaxEntries: 100_000, MaxExpandedBytes: 2 << 30,
		MaxFileBytes: 512 << 20, MaxCompressionRatio: 1000,
		Timeout: 15 * time.Minute, MaxMemoryBytes: 4 << 30, MaxProcesses: 1}
}

func ReleaseLimits() Limits {
	return Limits{MaxDepth: 5, MaxEntries: 100_000, MaxExpandedBytes: 10 << 30,
		MaxFileBytes: 512 << 20, MaxCompressionRatio: 1000,
		Timeout: time.Hour, MaxMemoryBytes: 8 << 30, MaxProcesses: 1}
}

func (l Limits) Valid() bool {
	return l.MaxDepth == 5 && l.MaxEntries == 100_000 &&
		(l.MaxExpandedBytes == 2<<30 || l.MaxExpandedBytes == 10<<30) &&
		l.MaxFileBytes == 512<<20 && l.MaxCompressionRatio == 1000 &&
		(l.Timeout == 15*time.Minute || l.Timeout == time.Hour) &&
		(l.MaxMemoryBytes == 4<<30 || l.MaxMemoryBytes == 8<<30) && l.MaxProcesses == 1
}

// TestLimits returns a stricter profile for synthetic boundary tests. Product
// code accepts it only when explicitly passed; it cannot support a product pass.
func TestLimits(depth, entries int, expanded, file, ratio int64, timeout time.Duration) Limits {
	return Limits{MaxDepth: depth, MaxEntries: entries, MaxExpandedBytes: expanded,
		MaxFileBytes: file, MaxCompressionRatio: ratio, Timeout: timeout,
		MaxMemoryBytes: 64 << 20, MaxProcesses: 1}
}

func (l Limits) internallyValid() bool {
	return l.MaxDepth > 0 && l.MaxDepth <= 5 && l.MaxEntries > 0 && l.MaxEntries <= 100_000 &&
		l.MaxExpandedBytes > 0 && l.MaxExpandedBytes <= 10<<30 &&
		l.MaxFileBytes > 0 && l.MaxFileBytes <= 512<<20 &&
		l.MaxCompressionRatio > 0 && l.MaxCompressionRatio <= 1000 &&
		l.Timeout > 0 && l.Timeout <= time.Hour && l.MaxMemoryBytes > 0 &&
		l.MaxMemoryBytes <= 8<<30 && l.MaxProcesses == 1
}

type Chunk struct {
	RelativePath string
	Start        int64
	End          int64
	Size         int64
	Digest       string
	Marker       string
}

type Entry struct {
	Class      Class
	Depth      int
	Size       int64
	Digest     string
	Projection string
	Chunks     []Chunk
}

type Result struct {
	Class              Class
	RawRoot            string
	ProbeRoot          string
	Entries            []Entry
	ExpectedProbeFiles []string
	ExpandedBytes      int64
	LedgerDigest       string
	ProjectionDigest   string
	Profile            Limits
}
