package request

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
	"unicode"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/outcome"
)

const (
	maxRequestAge                 = 24 * time.Hour
	maxFutureSkew                 = 5 * time.Minute
	maxPRExpandedBytes      int64 = 2 << 30
	maxReleaseExpandedBytes int64 = 10 << 30
	maxFileBytes            int64 = 512 << 20
)

var digestPattern = regexp.MustCompile(`^[0-9a-f]{64}$`)
var gitOIDPattern = regexp.MustCompile(`^(?:[0-9a-f]{40}|[0-9a-f]{64})$`)
var versionPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._+-]{0,63}$`)
var uuidV4Pattern = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

func IsUUIDv4(value string) bool { return uuidV4Pattern.MatchString(value) }

func ValidateLocalPath(path string) error { return validateSafeAbsolutePath(path) }

func ValidateAt(value ScanRequest, now time.Time) error {
	version, err := ParseVersion(value.RequestSchemaVersion)
	if err != nil || !version.Supported() || len(value.RequiredFeatures) > 0 {
		return reject(outcome.ReasonIndeterminateSchemaUnsupported)
	}
	if !IsUUIDv4(value.ScanID) || (value.SupersedesScanID != "" && !IsUUIDv4(value.SupersedesScanID)) {
		return reject(outcome.ReasonFailInputIntegrity)
	}
	if value.SupersedesScanID == value.ScanID {
		return reject(outcome.ReasonFailBindingMismatch)
	}
	if value.Mode != "local" && value.Mode != "pr" && value.Mode != "release" {
		return reject(outcome.ReasonFailInputIntegrity)
	}
	if !allDigests(value.ScannerReleaseDigest, value.EngineBinding.BinaryDigest, value.RulePackDigest, value.PolicyDigest, value.AllowlistDigest, value.SourceBinding.HistoryRangeDigest, value.SourceBinding.TrackedTreeDigest, value.TrackedSourceManifest.Digest) || !gitOIDPattern.MatchString(value.SourceBinding.HeadCommit) || !optionalGitOID(value.SourceBinding.BaseCommit) || !optionalGitOID(value.SourceBinding.MergeBase) {
		return reject(outcome.ReasonFailBindingMismatch)
	}
	if value.EngineBinding.Name != "gitleaks" || !versionPattern.MatchString(value.EngineBinding.Version) || !versionPattern.MatchString(value.EngineBinding.AdapterVersion) {
		return reject(outcome.ReasonFailBindingMismatch)
	}
	if value.FallbackRequirement.Mode != "disabled" {
		return reject(outcome.ReasonIndeterminateUnsupportedInput)
	}
	if value.RedactionMode != "full" || ((value.Mode == "pr" || value.Mode == "release") && !value.OfflineRequired) {
		return reject(outcome.ReasonFailBindingMismatch)
	}
	requestedAt, err := time.Parse(time.RFC3339, value.RequestedAt)
	if err != nil || requestedAt.Before(now.Add(-maxRequestAge)) || requestedAt.After(now.Add(maxFutureSkew)) {
		return reject(outcome.ReasonFailInputIntegrity)
	}
	if value.Mode == "pr" {
		if value.SourceBinding.FirstRelease || !allGitOIDs(value.SourceBinding.BaseCommit, value.SourceBinding.MergeBase) {
			return reject(outcome.ReasonFailBindingMismatch)
		}
	}
	if value.Mode == "release" {
		if value.BuildContextManifest == nil || value.ArtifactManifest == nil || !digestPattern.MatchString(value.ArtifactManifest.Digest) {
			return reject(outcome.ReasonUnavailableRequiredInput)
		}
		if !digestPattern.MatchString(value.BuildContextManifest.Digest) {
			return reject(outcome.ReasonFailBindingMismatch)
		}
		if value.SourceBinding.FirstRelease {
			if value.SourceBinding.BaseCommit != "" {
				return reject(outcome.ReasonFailBindingMismatch)
			}
		} else if !gitOIDPattern.MatchString(value.SourceBinding.BaseCommit) {
			return reject(outcome.ReasonFailBindingMismatch)
		}
	}
	if err := validateLimits(value.Mode, value.Limits); err != nil {
		return err
	}
	if err := validateSafeAbsolutePath(value.TrackedSourceManifest.Path); err != nil {
		return err
	}
	if value.BuildContextManifest != nil {
		if err := validateSafeAbsolutePath(value.BuildContextManifest.Path); err != nil {
			return err
		}
	}
	seen := map[string]struct{}{}
	if value.ArtifactManifest != nil {
		for _, entry := range value.ArtifactManifest.Entries {
			if err := validateSafeAbsolutePath(entry.Path); err != nil {
				return err
			}
			if entry.Size < 0 || entry.Size > value.Limits.MaxFileBytes || !digestPattern.MatchString(entry.Digest) || !versionPattern.MatchString(entry.Type) {
				return reject(outcome.ReasonFailInputIntegrity)
			}
			key := normalizedPathKey(entry.Path)
			if _, exists := seen[key]; exists {
				return reject(outcome.ReasonFailBindingMismatch)
			}
			seen[key] = struct{}{}
		}
		encoded, err := json.Marshal(value.ArtifactManifest.Entries)
		if err != nil {
			return reject(outcome.ReasonFailInputIntegrity)
		}
		digest := sha256.Sum256(encoded)
		if hex.EncodeToString(digest[:]) != value.ArtifactManifest.Digest {
			return reject(outcome.ReasonFailBindingMismatch)
		}
	}
	return nil
}

func ValidateBoundFiles(value ScanRequest) error {
	if err := verifyFile(value.TrackedSourceManifest.Path, value.TrackedSourceManifest.Digest, -1); err != nil {
		return err
	}
	if value.BuildContextManifest != nil {
		if err := verifyFile(value.BuildContextManifest.Path, value.BuildContextManifest.Digest, -1); err != nil {
			return err
		}
	}
	if value.ArtifactManifest != nil {
		for _, entry := range value.ArtifactManifest.Entries {
			if err := verifyFile(entry.Path, entry.Digest, entry.Size); err != nil {
				return err
			}
		}
	}
	return nil
}

func allDigests(values ...string) bool {
	for _, value := range values {
		if !digestPattern.MatchString(value) {
			return false
		}
	}
	return true
}

func allGitOIDs(values ...string) bool {
	for _, value := range values {
		if !gitOIDPattern.MatchString(value) {
			return false
		}
	}
	return true
}

func optionalGitOID(value string) bool { return value == "" || gitOIDPattern.MatchString(value) }

func validateLimits(mode string, limits Limits) error {
	maxExpanded := maxPRExpandedBytes
	maxTimeout := 900
	if mode == "release" {
		maxExpanded = maxReleaseExpandedBytes
		maxTimeout = 3600
	}
	if limits.TimeoutSeconds < 1 || limits.TimeoutSeconds > maxTimeout || limits.MaxArchiveDepth < 1 || limits.MaxArchiveDepth > 5 || limits.MaxArchiveEntries < 1 || limits.MaxArchiveEntries > 100000 || limits.MaxExpandedBytes < 1 || limits.MaxExpandedBytes > maxExpanded || limits.MaxFileBytes < 1 || limits.MaxFileBytes > maxFileBytes || limits.MaxCompressionRatio < 1 || limits.MaxCompressionRatio > 1000 || limits.MaxMemoryBytes < 1 || limits.MaxCPUPercent < 1 || limits.MaxCPUPercent > 100 {
		return reject(outcome.ReasonIndeterminateResourceLimit)
	}
	return nil
}

func validateSafeAbsolutePath(path string) error {
	if path == "" || !filepath.IsAbs(path) || filepath.Clean(path) != path || hasControl(path) {
		return reject(outcome.ReasonFailInputIntegrity)
	}
	volume := filepath.VolumeName(path)
	remainder := strings.TrimPrefix(path, volume)
	if remainder == string(filepath.Separator) || path == string(filepath.Separator) {
		return reject(outcome.ReasonFailInputIntegrity)
	}
	lower := strings.ToLower(path)
	if strings.HasPrefix(lower, `\\?\`) || strings.HasPrefix(lower, `\\.\`) || strings.HasPrefix(lower, `\\`) {
		return reject(outcome.ReasonFailInputIntegrity)
	}
	for _, part := range strings.Split(strings.TrimPrefix(remainder, string(filepath.Separator)), string(filepath.Separator)) {
		if part == "" || part == "." || part == ".." {
			return reject(outcome.ReasonFailInputIntegrity)
		}
		if runtime.GOOS == "windows" && strings.Contains(part, ":") {
			return reject(outcome.ReasonFailInputIntegrity)
		}
		if runtime.GOOS == "windows" && unsafeWindowsComponent(part) {
			return reject(outcome.ReasonFailInputIntegrity)
		}
	}
	return nil
}

func unsafeWindowsComponent(part string) bool {
	if strings.TrimRight(part, " .") != part {
		return true
	}
	stem := strings.ToUpper(strings.SplitN(part, ".", 2)[0])
	if stem == "CON" || stem == "PRN" || stem == "AUX" || stem == "NUL" {
		return true
	}
	if len(stem) == 4 && (strings.HasPrefix(stem, "COM") || strings.HasPrefix(stem, "LPT")) && stem[3] >= '1' && stem[3] <= '9' {
		return true
	}
	return false
}

func hasControl(value string) bool {
	for _, r := range value {
		if unicode.IsControl(r) {
			return true
		}
	}
	return false
}

func normalizedPathKey(path string) string {
	clean := filepath.Clean(path)
	if runtime.GOOS == "windows" {
		return strings.ToLower(clean)
	}
	return clean
}

func verifyFile(path, expected string, expectedSize int64) error {
	info, err := os.Lstat(path)
	if err != nil {
		return reject(outcome.ReasonUnavailableRequiredInput)
	}
	if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 || (expectedSize >= 0 && info.Size() != expectedSize) {
		return reject(outcome.ReasonFailInputIntegrity)
	}
	if expectedSize < 0 && info.Size() > 64<<20 {
		return reject(outcome.ReasonIndeterminateResourceLimit)
	}
	file, err := os.Open(path)
	if err != nil {
		return reject(outcome.ReasonUnavailableRequiredInput)
	}
	defer file.Close()
	opened, err := file.Stat()
	if err != nil || !opened.Mode().IsRegular() || !os.SameFile(info, opened) {
		return reject(outcome.ReasonFailInputIntegrity)
	}
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return reject(outcome.ReasonUnavailableRequiredInput)
	}
	if hex.EncodeToString(hash.Sum(nil)) != expected {
		return reject(outcome.ReasonFailBindingMismatch)
	}
	return nil
}
