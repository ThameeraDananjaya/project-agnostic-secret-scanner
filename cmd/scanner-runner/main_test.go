package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/outcome"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/request"
)

var cliTestTime = time.Date(2026, 8, 30, 20, 0, 0, 0, time.UTC)

func cliRequest(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	manifest := filepath.Join(root, "tracked.json")
	content := []byte(`{"files":[]}`)
	if err := os.WriteFile(manifest, content, 0o600); err != nil {
		t.Fatal(err)
	}
	hash := sha256.Sum256(content)
	d := strings.Repeat("a", 64)
	value := request.ScanRequest{
		RequestSchemaVersion: "1.0", ScanID: "123e4567-e89b-42d3-a456-426614174000", Mode: "pr", ScannerReleaseDigest: d,
		EngineBinding:  request.EngineBinding{Name: "gitleaks", Version: "1.0.0", BinaryDigest: d, AdapterVersion: "1.0.0"},
		RulePackDigest: d, PolicyDigest: d, AllowlistDigest: d,
		SourceBinding:         request.SourceBinding{BaseCommit: d, HeadCommit: d, MergeBase: d, HistoryRangeDigest: d, TrackedTreeDigest: d},
		TrackedSourceManifest: request.FileBinding{Path: manifest, Digest: hex.EncodeToString(hash[:])},
		FallbackRequirement:   request.FallbackRequirement{Mode: "disabled"},
		Limits:                request.Limits{TimeoutSeconds: 900, MaxArchiveDepth: 5, MaxArchiveEntries: 100000, MaxExpandedBytes: 2 << 30, MaxFileBytes: 512 << 20, MaxCompressionRatio: 1000, MaxMemoryBytes: 1 << 30, MaxCPUPercent: 100},
		OfflineRequired:       true, RedactionMode: "full", RequestedAt: cliTestTime.Format(time.RFC3339),
	}
	payload, _ := json.Marshal(value)
	path := filepath.Join(root, "request.json")
	if err := os.WriteFile(path, payload, 0o600); err != nil {
		t.Fatal(err)
	}
	return path, filepath.Join(root, "workspaces")
}

func TestRunValidRequestIsContentFreeAndStopsBeforeEngine(t *testing.T) {
	path, root := cliRequest(t)
	var stdout, stderr bytes.Buffer
	code := run([]string{"--request", path, "--workspace-root", root}, &stdout, &stderr, func() time.Time { return cliTestTime })
	if code != 30 || stderr.String() != "scanner-runner: engine unavailable\n" {
		t.Fatalf("code=%d stderr=%q", code, stderr.String())
	}
	var value outcome.Outcome
	if err := json.Unmarshal(stdout.Bytes(), &value); err != nil {
		t.Fatal(err)
	}
	if value.State != outcome.StateUnavailable || value.ReasonCode != outcome.ReasonUnavailableEngine || value.Bindings == nil {
		t.Fatalf("unexpected outcome: %#v", value)
	}
	if bytes.Contains(stdout.Bytes(), []byte(path)) || bytes.Contains(stdout.Bytes(), []byte(root)) {
		t.Fatal("path crossed output boundary")
	}
}

func TestRunCorruptionUnknownMajorAndArgumentInjectionFailClosed(t *testing.T) {
	root := t.TempDir()
	for _, test := range []struct {
		name, payload string
		code          int
		reason        outcome.ReasonCode
	}{
		{"corrupt", `{"requestSchemaVersion":"1.0",`, 10, outcome.ReasonFailInputIntegrity},
		{"major", `{"requestSchemaVersion":"2.0"}`, 20, outcome.ReasonIndeterminateSchemaUnsupported},
	} {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(root, test.name+".json")
			if err := os.WriteFile(path, []byte(test.payload), 0o600); err != nil {
				t.Fatal(err)
			}
			var stdout, stderr bytes.Buffer
			code := run([]string{"--request", path}, &stdout, &stderr, func() time.Time { return cliTestTime })
			var value outcome.Outcome
			if err := json.Unmarshal(stdout.Bytes(), &value); err != nil {
				t.Fatal(err)
			}
			if code != test.code || value.ReasonCode != test.reason || bytes.Contains(stdout.Bytes(), []byte(path)) {
				t.Fatalf("code=%d outcome=%#v", code, value)
			}
		})
	}
	injection := "\n::error file=synthetic::candidate"
	var stdout, stderr bytes.Buffer
	code := run([]string{"--bad", injection}, &stdout, &stderr, func() time.Time { return cliTestTime })
	if code != 10 || bytes.Contains(stdout.Bytes(), []byte(injection)) || bytes.Contains(stderr.Bytes(), []byte(injection)) || stderr.String() != "scanner-runner: invalid invocation\n" {
		t.Fatal("argument injection escaped")
	}
}

func TestRunInvalidRequestIdentityStillEmitsOneSafeOutcome(t *testing.T) {
	for _, test := range []struct {
		name   string
		mutate func(map[string]any)
	}{
		{"scanId", func(v map[string]any) { v["scanId"] = "invalid" }},
		{"mode", func(v map[string]any) { v["mode"] = "invalid" }},
		{"supersedes", func(v map[string]any) { v["supersedesScanId"] = "invalid" }},
	} {
		t.Run(test.name, func(t *testing.T) {
			path, _ := cliRequest(t)
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			var object map[string]any
			_ = json.Unmarshal(data, &object)
			test.mutate(object)
			data, _ = json.Marshal(object)
			if err := os.WriteFile(path, data, 0o600); err != nil {
				t.Fatal(err)
			}
			var stdout, stderr bytes.Buffer
			code := run([]string{"--request", path}, &stdout, &stderr, func() time.Time { return cliTestTime })
			decoder := json.NewDecoder(&stdout)
			var value outcome.Outcome
			if err := decoder.Decode(&value); err != nil {
				t.Fatal(err)
			}
			if code != 10 || value.ScanID == "invalid" || value.Mode != "unknown" || value.SupersedesScanID != "" || value.ReasonCode != outcome.ReasonFailInputIntegrity {
				t.Fatalf("unsafe rejection outcome: code=%d value=%#v", code, value)
			}
		})
	}
}

func TestRunInvalidAttemptEmitsRepresentableInputFailure(t *testing.T) {
	path, _ := cliRequest(t)
	var stdout, stderr bytes.Buffer
	code := run([]string{"--request", path, "--attempt", "0"}, &stdout, &stderr, func() time.Time { return cliTestTime })
	var value outcome.Outcome
	if err := json.Unmarshal(stdout.Bytes(), &value); err != nil {
		t.Fatal(err)
	}
	if code != 10 || value.AttemptNumber != 1 || value.ReasonCode != outcome.ReasonFailInputIntegrity || value.State != outcome.StateFail {
		t.Fatalf("invalid invocation was not safely representable: code=%d value=%#v", code, value)
	}
}
