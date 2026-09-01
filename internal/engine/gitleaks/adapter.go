// Package gitleaks binds the product to one exact Gitleaks CLI contract.
package gitleaks

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/artifact"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/engine"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/gitinput"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/outcome"
)

const (
	EngineName        = "gitleaks"
	EngineVersion     = "8.30.1"
	AdapterVersion    = "2.0.0"
	OutputBinding     = "json-v8.30.1"
	FindingExitCode   = 11
	MaximumReportSize = 1 << 30
	CoverageRuleID    = "pscan-projection-coverage"
)

var oidPattern = regexp.MustCompile(`^(?:[0-9a-f]{40}|[0-9a-f]{64})$`)

type Binding struct {
	Executable       string
	ExecutableDigest string
	Config           string
	ConfigDigest     string
	IgnoreFile       string
	IgnoreFileDigest string
	Environment      []string
	PrivateHome      string
}

type Adapter struct{ Binding Binding }

func (a Adapter) Verify() error {
	return a.VerifyContext(context.Background())
}

func (a Adapter) VerifyContext(ctx context.Context) error {
	if err := engine.VerifyRegularFileContext(ctx, a.Binding.Executable, a.Binding.ExecutableDigest); err != nil {
		return err
	}
	if err := engine.VerifyRegularFileContext(ctx, a.Binding.Config, a.Binding.ConfigDigest); err != nil {
		return err
	}
	if _, err := ProveRuleSpansContext(ctx, a.Binding.Config, a.Binding.ConfigDigest); err != nil {
		return err
	}
	if err := engine.VerifyRegularFileContext(ctx, a.Binding.IgnoreFile, a.Binding.IgnoreFileDigest); err != nil {
		return err
	}
	if !filepath.IsAbs(a.Binding.PrivateHome) {
		return errors.New("private engine home is not absolute")
	}
	return nil
}

func (a Adapter) ScanDirectory(ctx context.Context, target string, timeout time.Duration, maxFileBytes int64) engine.Result {
	if !filepath.IsAbs(target) || timeout <= 0 || maxFileBytes <= 0 {
		return engine.Result{Reason: outcome.ReasonFailInputIntegrity, ExitCode: -1}
	}
	if probe := a.probe(ctx, timeout); probe.Reason != outcome.ReasonPassNoBlockingFindings {
		return probe
	}
	args := append(a.commonArgs(a.Binding.Config, timeout, maxFileBytes), target)
	return a.run(ctx, append([]string{"dir"}, args...), timeout)
}

// ScanProjection verifies the byte-exact and framed projections, then runs one
// pinned Gitleaks invocation whose configuration contains both product rules
// and the private coverage rule. A skip therefore removes the expected marker
// from the same run and cannot become pass.
func (a Adapter) ScanProjection(ctx context.Context, projection gitinput.MaterializedProjection, timeout time.Duration, maxFileBytes int64) engine.Result {
	// Unnamed legacy limits cannot support pass.
	return engine.Result{Reason: outcome.ReasonIndeterminateIncompleteCoverage, ExitCode: -1}
}

func (a Adapter) ScanProjectionProfile(ctx context.Context, projection gitinput.MaterializedProjection, profile gitinput.CoverageProfile) engine.Result {
	if err := gitinput.VerifyMaterializedProjection(projection); err != nil {
		return engine.Result{Reason: outcome.ReasonFailInputIntegrity, ExitCode: -1}
	}
	var total int64
	for _, object := range projection.Ledger {
		if profile.AdmitBlob(object.Size) != nil || object.Size > profile.MaxTotalBytes-total {
			return engine.Result{Reason: outcome.ReasonIndeterminateResourceLimit, ExitCode: -1}
		}
		total += object.Size
	}
	if profile.Admit(total, projection.ObjectCount, len(projection.ExpectedProbeFiles)) != nil {
		return engine.Result{Reason: outcome.ReasonIndeterminateResourceLimit, ExitCode: -1}
	}
	if profile.AdmitRuntime(profile.Timeout, profile.MaxMemoryBytes, 1, profile.MaxReportBytes, gitinput.MaximumDetectorFile-1) != nil {
		return engine.Result{Reason: outcome.ReasonIndeterminateResourceLimit, ExitCode: -1}
	}
	timeout := profile.Timeout
	if timeout <= 0 || projection.EntryCount <= 0 {
		return engine.Result{Reason: outcome.ReasonFailInputIntegrity, ExitCode: -1}
	}
	profileContext, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	probeTimeout := 5 * time.Second
	if timeout < probeTimeout {
		probeTimeout = timeout
	}
	if probe := a.probe(profileContext, probeTimeout); probe.Reason != outcome.ReasonPassNoBlockingFindings {
		return probe
	}
	expected, ok := normalizeExpectedPaths(projection.ExpectedProbeFiles)
	if !ok || len(expected) != len(projection.ExpectedProbeFiles) {
		return engine.Result{Reason: outcome.ReasonFailInputIntegrity, ExitCode: -1}
	}
	args := append(a.commonArgs(a.Binding.Config, timeout, gitinput.MaximumDetectorFile), ".")
	return a.runProfile(profileContext, append([]string{"dir"}, args...), projection.ProbeRoot, profile, decodeProjection(expected))
}

// ScanArtifactProjection is the sole pass-capable artifact entry point. It
// verifies the private ledger and every overlapping detector projection before
// invoking the exact pinned detector with archive recursion disabled.
func (a Adapter) ScanArtifactProjection(ctx context.Context, projection artifact.Result) engine.Result {
	if err := artifact.VerifyContext(ctx, projection); err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return engine.Result{Reason: outcome.ReasonIndeterminateTimeout, ExitCode: -1}
		}
		if errors.Is(ctx.Err(), context.Canceled) {
			return engine.Result{Reason: outcome.ReasonIndeterminateIncompleteCoverage, ExitCode: -1}
		}
		return engine.Result{Reason: outcome.ReasonFailInputIntegrity, ExitCode: -1}
	}
	if !projection.Profile.Valid() {
		return engine.Result{Reason: outcome.ReasonIndeterminateResourceLimit, ExitCode: -1}
	}
	environment := append([]string(nil), a.Binding.Environment...)
	environment = append(environment,
		"GOMEMLIMIT="+strconv.FormatInt(projection.Profile.MaxMemoryBytes, 10)+"B",
		"GOMAXPROCS=1",
	)
	if !engine.IsSealedEnvironment(environment) {
		return engine.Result{Reason: outcome.ReasonUnavailableRuntime, ExitCode: -1}
	}
	expected, ok := normalizeExpectedPaths(projection.ExpectedProbeFiles)
	if !ok || len(expected) != len(projection.ExpectedProbeFiles) {
		return engine.Result{Reason: outcome.ReasonFailInputIntegrity, ExitCode: -1}
	}
	probeTimeout := 5 * time.Second
	if probeTimeout > projection.Profile.Timeout {
		probeTimeout = projection.Profile.Timeout
	}
	if probe := a.probe(ctx, probeTimeout); probe.Reason != outcome.ReasonPassNoBlockingFindings {
		return probe
	}
	args := append(a.commonArgs(a.Binding.Config, projection.Profile.Timeout, artifact.MaximumDetectorFile), ".")
	return engine.RunPrivateContext(ctx, engine.Command{
		Executable: a.Binding.Executable, ExpectedDigest: a.Binding.ExecutableDigest,
		Args: append([]string{"dir"}, args...), Dir: projection.ProbeRoot,
		Environment: environment, Timeout: projection.Profile.Timeout,
		CaptureLimit: MaximumReportSize,
	}, decodeProjectionContext(expected))
}

func (a Adapter) runProfile(ctx context.Context, args []string, directory string, profile gitinput.CoverageProfile, decoder engine.Decoder) engine.Result {
	if err := a.VerifyContext(ctx); err != nil || !profile.Valid() {
		return engine.Result{Reason: outcome.ReasonFailScannerIntegrity, ExitCode: -1}
	}
	environment := append([]string(nil), a.Binding.Environment...)
	environment = append(environment,
		"GOMEMLIMIT="+strconv.FormatInt(profile.MaxMemoryBytes, 10)+"B",
		"GOMAXPROCS=2",
	)
	return engine.RunPrivate(ctx, engine.Command{
		Executable: a.Binding.Executable, ExpectedDigest: a.Binding.ExecutableDigest,
		Args: args, Dir: directory, Environment: environment,
		Timeout: profile.Timeout, CaptureLimit: profile.MaxReportBytes / 2,
	}, decoder)
}

func (a Adapter) ScanGitRange(ctx context.Context, bareRepository, base, head string, firstRelease bool, timeout time.Duration, maxFileBytes int64) engine.Result {
	if !filepath.IsAbs(bareRepository) || !oidPattern.MatchString(head) || timeout <= 0 || maxFileBytes <= 0 {
		return engine.Result{Reason: outcome.ReasonFailInputIntegrity, ExitCode: -1}
	}
	// Native patch input cannot prove binary blob coverage. It is retained only
	// as a fail-closed compatibility surface; callers must use ScanProjection.
	return engine.Result{Reason: outcome.ReasonIndeterminateIncompleteCoverage, ExitCode: -1}
}

func (a Adapter) probe(ctx context.Context, timeout time.Duration) engine.Result {
	if err := a.VerifyContext(ctx); err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return engine.Result{Reason: outcome.ReasonIndeterminateTimeout, ExitCode: -1}
		}
		if errors.Is(ctx.Err(), context.Canceled) {
			return engine.Result{Reason: outcome.ReasonIndeterminateIncompleteCoverage, ExitCode: -1}
		}
		return engine.Result{Reason: outcome.ReasonFailScannerIntegrity, ExitCode: -1}
	}
	return engine.RunPrivate(ctx, engine.Command{
		Executable:     a.Binding.Executable,
		ExpectedDigest: a.Binding.ExecutableDigest,
		Args:           []string{"version"},
		Dir:            a.Binding.PrivateHome,
		Environment:    append([]string(nil), a.Binding.Environment...),
		Timeout:        timeout,
		CaptureLimit:   4096,
	}, func(output engine.PrivateOutput) outcome.ReasonCode {
		if output.Exit != 0 || len(bytes.TrimSpace(output.Stderr)) != 0 {
			return outcome.ReasonUnavailableEngine
		}
		if string(bytes.TrimSpace(output.Stdout)) != EngineVersion {
			return outcome.ReasonFailBindingMismatch
		}
		return outcome.ReasonPassNoBlockingFindings
	})
}

func (a Adapter) commonArgs(config string, timeout time.Duration, maxFileBytes int64) []string {
	seconds := int(timeout.Round(time.Second) / time.Second)
	if seconds < 1 {
		seconds = 1
	}
	megabytes := maxFileBytes / (1 << 20)
	if maxFileBytes%(1<<20) != 0 {
		megabytes++
	}
	if megabytes < 1 {
		megabytes = 1
	}
	return []string{
		"--config", config,
		"--gitleaks-ignore-path", a.Binding.IgnoreFile,
		"--ignore-gitleaks-allow",
		"--report-format", "json",
		"--report-path", "-",
		"--redact=100",
		"--no-banner",
		"--no-color",
		"--log-level", "error",
		"--exit-code", strconv.Itoa(FindingExitCode),
		"--max-archive-depth", "0",
		"--max-target-megabytes", strconv.FormatInt(megabytes, 10),
		"--timeout", strconv.Itoa(seconds),
	}
}

func (a Adapter) run(ctx context.Context, args []string, timeout time.Duration) engine.Result {
	return a.runWithDecoder(ctx, args, "", timeout, decode)
}

func (a Adapter) runWithDecoder(ctx context.Context, args []string, directory string, timeout time.Duration, decoder engine.Decoder) engine.Result {
	if err := a.VerifyContext(ctx); err != nil {
		return engine.Result{Reason: outcome.ReasonFailScannerIntegrity, ExitCode: -1}
	}
	return engine.RunPrivate(ctx, engine.Command{
		Executable:     a.Binding.Executable,
		ExpectedDigest: a.Binding.ExecutableDigest,
		Args:           args,
		Dir:            directory,
		Environment:    append([]string(nil), a.Binding.Environment...),
		Timeout:        timeout + 5*time.Second,
		CaptureLimit:   MaximumReportSize,
	}, decoder)
}

func normalizeExpectedPaths(paths []string) (map[string]bool, bool) {
	expected := make(map[string]bool, len(paths))
	for _, path := range paths {
		normalized := filepath.ToSlash(filepath.Clean(filepath.FromSlash(path)))
		if normalized == "." || normalized != path || filepath.IsAbs(path) || strings.HasPrefix(path, "../") || strings.ContainsRune(path, '\x00') || expected[path] {
			return nil, false
		}
		expected[path] = true
	}
	return expected, true
}

func decodeProjection(expected map[string]bool) engine.Decoder {
	contextDecoder := decodeProjectionContext(expected)
	return func(output engine.PrivateOutput) outcome.ReasonCode {
		return contextDecoder(context.Background(), output)
	}
}

func decodeProjectionContext(expected map[string]bool) engine.ContextDecoder {
	return func(ctx context.Context, output engine.PrivateOutput) outcome.ReasonCode {
		if len(bytes.TrimSpace(output.Stderr)) != 0 || output.Exit != FindingExitCode {
			return outcome.ReasonIndeterminateIncompleteCoverage
		}
		decoder := json.NewDecoder(bytes.NewReader(bytes.TrimSpace(output.Stdout)))
		token, err := decoder.Token()
		if err != nil || token != json.Delim('[') {
			return outcome.ReasonIndeterminateIncompleteCoverage
		}
		seen := make(map[string]bool, len(expected))
		candidateFinding := false
		findingCount := 0
		for decoder.More() {
			if ctx.Err() != nil {
				return outcome.ReasonIndeterminateTimeout
			}
			var object map[string]json.RawMessage
			if decoder.Decode(&object) != nil || leaksUnredactedCandidate(object) {
				return outcome.ReasonIndeterminateRedactionUnproven
			}
			findingCount++
			var ruleID, path string
			if json.Unmarshal(object["RuleID"], &ruleID) != nil || json.Unmarshal(object["File"], &path) != nil {
				return outcome.ReasonIndeterminateIncompleteCoverage
			}
			path = strings.TrimPrefix(filepath.ToSlash(path), "./")
			if !expected[path] {
				return outcome.ReasonIndeterminateIncompleteCoverage
			}
			if ruleID == CoverageRuleID {
				if seen[path] {
					return outcome.ReasonIndeterminateIncompleteCoverage
				}
				seen[path] = true
			} else {
				candidateFinding = true
			}
		}
		if token, err := decoder.Token(); err != nil || token != json.Delim(']') || decoder.Decode(&struct{}{}) != io.EOF || findingCount < len(expected) || len(seen) != len(expected) {
			return outcome.ReasonIndeterminateIncompleteCoverage
		}
		if candidateFinding {
			return outcome.ReasonFailFindingDetected
		}
		return outcome.ReasonPassNoBlockingFindings
	}
}

func decode(output engine.PrivateOutput) outcome.ReasonCode {
	if len(bytes.TrimSpace(output.Stderr)) != 0 {
		return outcome.ReasonIndeterminateIncompleteCoverage
	}
	trimmed := bytes.TrimSpace(output.Stdout)
	var findings []json.RawMessage
	if len(trimmed) == 0 || json.Unmarshal(trimmed, &findings) != nil {
		return outcome.ReasonIndeterminateSchemaUnsupported
	}
	for _, raw := range findings {
		var object map[string]json.RawMessage
		if json.Unmarshal(raw, &object) != nil || len(object) == 0 {
			return outcome.ReasonIndeterminateSchemaUnsupported
		}
		if leaksUnredactedCandidate(object) {
			return outcome.ReasonIndeterminateRedactionUnproven
		}
	}
	switch output.Exit {
	case 0:
		if len(findings) != 0 {
			return outcome.ReasonIndeterminateConflictingResults
		}
		return outcome.ReasonPassNoBlockingFindings
	case FindingExitCode:
		if len(findings) == 0 {
			return outcome.ReasonIndeterminateConflictingResults
		}
		return outcome.ReasonFailFindingDetected
	default:
		return outcome.ReasonIndeterminateIncompleteCoverage
	}
}

func leaksUnredactedCandidate(object map[string]json.RawMessage) bool {
	raw, ok := object["Secret"]
	if !ok {
		return true
	}
	var secret string
	if json.Unmarshal(raw, &secret) != nil || secret == "" {
		return true
	}
	if secret == "REDACTED" {
		return false
	}
	for _, r := range secret {
		if r != '*' {
			return true
		}
	}
	return false
}

func ValidateBinding(name, version, adapter, output string) error {
	if name != EngineName || version != EngineVersion || adapter != AdapterVersion || output != OutputBinding {
		return errors.New("unsupported Gitleaks binding")
	}
	return nil
}

func RangeSpec(base, head string, firstRelease bool) (string, error) {
	if !oidPattern.MatchString(head) {
		return "", fmt.Errorf("invalid head commit")
	}
	if firstRelease {
		if strings.TrimSpace(base) != "" {
			return "", fmt.Errorf("first release must not bind a base")
		}
		return head, nil
	}
	if !oidPattern.MatchString(base) || base == head {
		return "", fmt.Errorf("invalid base commit")
	}
	return base + ".." + head, nil
}
