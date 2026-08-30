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
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/outcome"
)

const (
	EngineName        = "gitleaks"
	EngineVersion     = "8.30.1"
	AdapterVersion    = "1.0.0"
	OutputBinding     = "json-v8.30.1"
	FindingExitCode   = 11
	MaximumReportSize = 16 << 20
)

var oidPattern = regexp.MustCompile(`^(?:[0-9a-f]{40}|[0-9a-f]{64})$`)

type Binding struct {
	Executable       string
	ExecutableDigest string
	Config           string
	ConfigDigest     string
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
	return nil
}

func (a Adapter) ScanDirectory(ctx context.Context, target string, timeout time.Duration, maxFileBytes int64) engine.Result {
	if !filepath.IsAbs(target) || timeout <= 0 || maxFileBytes <= 0 {
		return engine.Result{Reason: outcome.ReasonFailInputIntegrity, ExitCode: -1}
	}
	if probe := a.probe(ctx, timeout); probe.Reason != outcome.ReasonPassNoBlockingFindings {
		return probe
	}
	args := append(a.commonArgs(timeout, maxFileBytes), target)
	return a.run(ctx, append([]string{"dir"}, args...), timeout)
}

func (a Adapter) ScanGitRange(ctx context.Context, bareRepository, base, head string, firstRelease bool, timeout time.Duration, maxFileBytes int64) engine.Result {
	if !filepath.IsAbs(bareRepository) || !oidPattern.MatchString(head) || timeout <= 0 || maxFileBytes <= 0 {
		return engine.Result{Reason: outcome.ReasonFailInputIntegrity, ExitCode: -1}
	}
	if probe := a.probe(ctx, timeout); probe.Reason != outcome.ReasonPassNoBlockingFindings {
		return probe
	}
	var rangeSpec string
	if firstRelease {
		if base != "" {
			return engine.Result{Reason: outcome.ReasonFailBindingMismatch, ExitCode: -1}
		}
		rangeSpec = head
	} else {
		if !oidPattern.MatchString(base) || base == head {
			return engine.Result{Reason: outcome.ReasonFailBindingMismatch, ExitCode: -1}
		}
		rangeSpec = base + ".." + head
	}
	args := []string{"git", "--log-opts=" + rangeSpec}
	args = append(args, a.commonArgs(timeout, maxFileBytes)...)
	args = append(args, bareRepository)
	return a.run(ctx, args, timeout)
}

func (a Adapter) probe(ctx context.Context, timeout time.Duration) engine.Result {
	if err := a.Verify(); err != nil {
		return engine.Result{Reason: outcome.ReasonFailScannerIntegrity, ExitCode: -1}
	}
	return engine.RunPrivate(ctx, engine.Command{
		Executable:     a.Binding.Executable,
		ExpectedDigest: a.Binding.ExecutableDigest,
		Args:           []string{"version"},
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

func (a Adapter) commonArgs(timeout time.Duration, maxFileBytes int64) []string {
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
		"--config", a.Binding.Config,
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
	if err := a.Verify(); err != nil {
		return engine.Result{Reason: outcome.ReasonFailScannerIntegrity, ExitCode: -1}
	}
	return engine.RunPrivate(ctx, engine.Command{
		Executable:     a.Binding.Executable,
		ExpectedDigest: a.Binding.ExecutableDigest,
		Args:           args,
		Environment:    append([]string(nil), a.Binding.Environment...),
		Timeout:        timeout + 5*time.Second,
		CaptureLimit:   MaximumReportSize,
	}, decode)
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
