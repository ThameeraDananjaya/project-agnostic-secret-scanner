// Package gitleaks binds the product to one exact Gitleaks CLI contract.
package gitleaks

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/engine"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/gitinput"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/outcome"
)

const (
	EngineName        = "gitleaks"
	EngineVersion     = "8.30.1"
	AdapterVersion    = "1.1.0"
	OutputBinding     = "json-v8.30.1"
	FindingExitCode   = 11
	MaximumReportSize = 16 << 20
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
	if err := engine.VerifyRegularFile(a.Binding.Executable, a.Binding.ExecutableDigest); err != nil {
		return err
	}
	if err := engine.VerifyRegularFile(a.Binding.Config, a.Binding.ConfigDigest); err != nil {
		return err
	}
	if err := engine.VerifyRegularFile(a.Binding.IgnoreFile, a.Binding.IgnoreFileDigest); err != nil {
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
	if timeout <= 0 || maxFileBytes <= 0 || projection.EntryCount <= 0 {
		return engine.Result{Reason: outcome.ReasonFailInputIntegrity, ExitCode: -1}
	}
	if err := gitinput.VerifyMaterializedProjection(projection); err != nil {
		return engine.Result{Reason: outcome.ReasonFailInputIntegrity, ExitCode: -1}
	}
	if probe := a.probe(ctx, timeout); probe.Reason != outcome.ReasonPassNoBlockingFindings {
		return probe
	}
	expected, ok := normalizeExpectedPaths(projection.ExpectedProbeFiles)
	if !ok || len(expected) != projection.EntryCount {
		return engine.Result{Reason: outcome.ReasonFailInputIntegrity, ExitCode: -1}
	}
	args := append(a.commonArgs(a.Binding.Config, timeout, maxFileBytes+256), ".")
	return a.runWithDecoder(ctx, append([]string{"dir"}, args...), projection.ProbeRoot, timeout, decodeProjection(expected))
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
	if err := a.Verify(); err != nil {
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
	if err := a.Verify(); err != nil {
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
	return func(output engine.PrivateOutput) outcome.ReasonCode {
		if len(bytes.TrimSpace(output.Stderr)) != 0 || output.Exit != FindingExitCode {
			return outcome.ReasonIndeterminateIncompleteCoverage
		}
		var findings []json.RawMessage
		if json.Unmarshal(bytes.TrimSpace(output.Stdout), &findings) != nil || len(findings) < len(expected) {
			return outcome.ReasonIndeterminateIncompleteCoverage
		}
		seen := make(map[string]bool, len(findings))
		candidateFinding := false
		for _, raw := range findings {
			var object map[string]json.RawMessage
			if json.Unmarshal(raw, &object) != nil || leaksUnredactedCandidate(object) {
				return outcome.ReasonIndeterminateRedactionUnproven
			}
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
		if len(seen) != len(expected) {
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
