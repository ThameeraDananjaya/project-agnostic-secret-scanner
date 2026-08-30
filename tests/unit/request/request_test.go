package request_test

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

var testNow = time.Date(2026, 8, 30, 20, 0, 0, 0, time.UTC)

func validRequest(t *testing.T, mode string) request.ScanRequest {
	t.Helper()
	manifest := filepath.Join(t.TempDir(), "tracked-manifest.json")
	content := []byte(`{"files":[]}`)
	if err := os.WriteFile(manifest, content, 0o600); err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(content)
	d := strings.Repeat("a", 64)
	value := request.ScanRequest{
		RequestSchemaVersion:  "1.0",
		ScanID:                "123e4567-e89b-42d3-a456-426614174000",
		Mode:                  mode,
		ScannerReleaseDigest:  d,
		EngineBinding:         request.EngineBinding{Name: "gitleaks", Version: "1.0.0", BinaryDigest: d, AdapterVersion: "1.0.0"},
		RulePackDigest:        d,
		PolicyDigest:          d,
		AllowlistDigest:       d,
		SourceBinding:         request.SourceBinding{HeadCommit: d, HistoryRangeDigest: d, TrackedTreeDigest: d},
		TrackedSourceManifest: request.FileBinding{Path: manifest, Digest: hex.EncodeToString(digest[:])},
		FallbackRequirement:   request.FallbackRequirement{Mode: "disabled"},
		Limits:                request.Limits{TimeoutSeconds: 900, MaxArchiveDepth: 5, MaxArchiveEntries: 100000, MaxExpandedBytes: 2 << 30, MaxFileBytes: 512 << 20, MaxCompressionRatio: 1000, MaxMemoryBytes: 1 << 30, MaxCPUPercent: 100},
		OfflineRequired:       true,
		RedactionMode:         "full",
		RequestedAt:           testNow.Format(time.RFC3339),
	}
	if mode == "pr" {
		value.SourceBinding.BaseCommit = d
		value.SourceBinding.MergeBase = d
	}
	return value
}

func reason(t *testing.T, err error) outcome.ReasonCode {
	t.Helper()
	if err == nil {
		t.Fatal("expected rejection")
	}
	return request.ReasonFor(err, "unexpected")
}

func TestVersionCompatibility(t *testing.T) {
	tests := []struct {
		value             string
		supported, future bool
	}{
		{"1.0", true, false}, {"1.7", true, true}, {"2.0", false, false},
	}
	for _, test := range tests {
		version, err := request.ParseVersion(test.value)
		if err != nil {
			t.Fatalf("%s: %v", test.value, err)
		}
		if version.Supported() != test.supported || version.FutureMinor() != test.future {
			t.Fatalf("unexpected compatibility for %s", test.value)
		}
	}
	for _, invalid := range []string{"", "1", "1.0.0", "01.0", "1.00", "-1.0", "a.b"} {
		if _, err := request.ParseVersion(invalid); err == nil {
			t.Errorf("accepted invalid version %q", invalid)
		}
	}
}

func TestLoadRejectsCorruptionAndDuplicateKeys(t *testing.T) {
	cases := [][]byte{
		{}, []byte(`{`), []byte(`{"requestSchemaVersion":"1.0"}{}`),
		[]byte(`{"requestSchemaVersion":"1.0","requestSchemaVersion":"2.0"}`),
	}
	for _, data := range cases {
		_, err := request.Load(bytes.NewReader(data))
		if got := reason(t, err); got != outcome.ReasonFailInputIntegrity {
			t.Fatalf("got %s", got)
		}
	}
}

func TestLoadRejectsMissingSchemaRequiredPresence(t *testing.T) {
	base := validRequest(t, "release")
	d := strings.Repeat("b", 64)
	base.SourceBinding.FirstRelease = true
	base.BuildContextManifest = &request.FileBinding{Path: base.TrackedSourceManifest.Path, Digest: base.TrackedSourceManifest.Digest}
	base.ArtifactManifest = &request.ArtifactManifest{Digest: d, Entries: []request.ArtifactEntry{}}
	for _, test := range []struct {
		name   string
		mutate func(map[string]any)
	}{
		{"offlineRequired", func(v map[string]any) { delete(v, "offlineRequired") }},
		{"firstRelease", func(v map[string]any) { delete(v["sourceBinding"].(map[string]any), "firstRelease") }},
		{"artifactEntries", func(v map[string]any) { delete(v["artifactManifest"].(map[string]any), "entries") }},
	} {
		t.Run(test.name, func(t *testing.T) {
			data, _ := json.Marshal(base)
			var object map[string]any
			_ = json.Unmarshal(data, &object)
			test.mutate(object)
			data, _ = json.Marshal(object)
			if _, err := request.Load(bytes.NewReader(data)); reason(t, err) != outcome.ReasonFailInputIntegrity {
				t.Fatal("schema-required omission was accepted")
			}
		})
	}
}

func TestLoadMajorMinorAndUnknownFieldRules(t *testing.T) {
	value := validRequest(t, "pr")
	data, _ := json.Marshal(value)
	var object map[string]any
	_ = json.Unmarshal(data, &object)
	object["futureOptionalField"] = "ignored"
	current, _ := json.Marshal(object)
	if _, err := request.Load(bytes.NewReader(current)); reason(t, err) != outcome.ReasonFailInputIntegrity {
		t.Fatal("current minor accepted unknown field")
	}
	object["requestSchemaVersion"] = "1.9"
	future, _ := json.Marshal(object)
	loaded, err := request.Load(bytes.NewReader(future))
	if err != nil || loaded.RequestSchemaVersion != "1.9" {
		t.Fatalf("compatible future minor rejected: %v", err)
	}
	object["requestSchemaVersion"] = "2.0"
	major, _ := json.Marshal(object)
	if _, err := request.Load(bytes.NewReader(major)); reason(t, err) != outcome.ReasonIndeterminateSchemaUnsupported {
		t.Fatal("unknown major did not fail closed")
	}
}

func TestRequiredFeatureFailsClosed(t *testing.T) {
	value := validRequest(t, "pr")
	value.RequestSchemaVersion = "1.1"
	value.RequiredFeatures = []string{"future-required-semantics"}
	data, _ := json.Marshal(value)
	loaded, err := request.Load(bytes.NewReader(data))
	if err == nil {
		err = request.ValidateAt(loaded, testNow)
	}
	if reason(t, err) != outcome.ReasonIndeterminateSchemaUnsupported {
		t.Fatal("required feature did not fail closed")
	}
}

func TestValidationRejectsUnsafeAndMismatchedBindings(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*request.ScanRequest)
		want   outcome.ReasonCode
	}{
		{"bad scan id", func(v *request.ScanRequest) { v.ScanID = "source-derived" }, outcome.ReasonFailInputIntegrity},
		{"same recovery id", func(v *request.ScanRequest) { v.SupersedesScanID = v.ScanID }, outcome.ReasonFailBindingMismatch},
		{"wrong engine", func(v *request.ScanRequest) { v.EngineBinding.Name = "other" }, outcome.ReasonFailBindingMismatch},
		{"bad digest", func(v *request.ScanRequest) { v.PolicyDigest = strings.Repeat("z", 64) }, outcome.ReasonFailBindingMismatch},
		{"fallback enabled", func(v *request.ScanRequest) { v.FallbackRequirement.Mode = "enabled" }, outcome.ReasonIndeterminateUnsupportedInput},
		{"offline false", func(v *request.ScanRequest) { v.OfflineRequired = false }, outcome.ReasonFailBindingMismatch},
		{"redaction weakened", func(v *request.ScanRequest) { v.RedactionMode = "partial" }, outcome.ReasonFailBindingMismatch},
		{"stale request", func(v *request.ScanRequest) { v.RequestedAt = testNow.Add(-25 * time.Hour).Format(time.RFC3339) }, outcome.ReasonFailInputIntegrity},
		{"future request", func(v *request.ScanRequest) { v.RequestedAt = testNow.Add(6 * time.Minute).Format(time.RFC3339) }, outcome.ReasonFailInputIntegrity},
		{"unclean path", func(v *request.ScanRequest) {
			v.TrackedSourceManifest.Path = filepath.Dir(v.TrackedSourceManifest.Path) + string(filepath.Separator) + "child" + string(filepath.Separator) + ".." + string(filepath.Separator) + filepath.Base(v.TrackedSourceManifest.Path)
		}, outcome.ReasonFailInputIntegrity},
		{"control path", func(v *request.ScanRequest) { v.TrackedSourceManifest.Path += "\nname" }, outcome.ReasonFailInputIntegrity},
		{"resource excess", func(v *request.ScanRequest) { v.Limits.MaxArchiveDepth = 6 }, outcome.ReasonIndeterminateResourceLimit},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value := validRequest(t, "pr")
			test.mutate(&value)
			if got := reason(t, request.ValidateAt(value, testNow)); got != test.want {
				t.Fatalf("got %s want %s", got, test.want)
			}
		})
	}
}

func TestFullGitObjectIDsAcceptSHA1AndSHA256Only(t *testing.T) {
	for _, width := range []int{40, 64} {
		value := validRequest(t, "pr")
		oid := strings.Repeat("a", width)
		value.SourceBinding.BaseCommit = oid
		value.SourceBinding.HeadCommit = oid
		value.SourceBinding.MergeBase = oid
		if err := request.ValidateAt(value, testNow); err != nil {
			t.Fatalf("%d-hex full object ID rejected: %v", width, err)
		}
	}
	value := validRequest(t, "pr")
	value.SourceBinding.HeadCommit = strings.Repeat("a", 12)
	if got := reason(t, request.ValidateAt(value, testNow)); got != outcome.ReasonFailBindingMismatch {
		t.Fatalf("abbreviated object ID got %s", got)
	}
}

func TestRequestInputRejectsLinksAndUnsafeDescriptorTypes(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "request.json")
	if err := os.WriteFile(target, []byte(`{}`), 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "request-link.json")
	if err := os.Symlink(target, link); err == nil {
		if file, err := request.OpenLocal(link); err == nil {
			_ = file.Close()
			t.Fatal("symlink request was accepted")
		}
	}
	directory, err := os.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if file, err := request.OpenDescriptor(directory.Fd(), time.Second); err == nil {
		_ = file.Close()
		t.Fatal("directory descriptor was accepted")
	}
}

func TestReleaseRequiresCompleteUniqueBindings(t *testing.T) {
	value := validRequest(t, "release")
	if got := reason(t, request.ValidateAt(value, testNow)); got != outcome.ReasonUnavailableRequiredInput {
		t.Fatalf("got %s", got)
	}
	d := strings.Repeat("b", 64)
	value.SourceBinding.FirstRelease = true
	value.BuildContextManifest = &request.FileBinding{Path: value.TrackedSourceManifest.Path, Digest: value.TrackedSourceManifest.Digest}
	value.ArtifactManifest = &request.ArtifactManifest{Digest: d, Entries: []request.ArtifactEntry{
		{Path: value.TrackedSourceManifest.Path, Type: "file", Size: 12, Digest: value.TrackedSourceManifest.Digest},
		{Path: value.TrackedSourceManifest.Path, Type: "file", Size: 12, Digest: value.TrackedSourceManifest.Digest},
	}}
	if got := reason(t, request.ValidateAt(value, testNow)); got != outcome.ReasonFailBindingMismatch {
		t.Fatalf("got %s", got)
	}
	value.ArtifactManifest.Entries = value.ArtifactManifest.Entries[:1]
	if got := reason(t, request.ValidateAt(value, testNow)); got != outcome.ReasonFailBindingMismatch {
		t.Fatalf("mismatched artifact-manifest digest got %s", got)
	}
	encoded, _ := json.Marshal(value.ArtifactManifest.Entries)
	digest := sha256.Sum256(encoded)
	value.ArtifactManifest.Digest = hex.EncodeToString(digest[:])
	if err := request.ValidateAt(value, testNow); err != nil {
		t.Fatalf("valid first-release binding rejected: %v", err)
	}
	value.SourceBinding.FirstRelease = false
	if got := reason(t, request.ValidateAt(value, testNow)); got != outcome.ReasonFailBindingMismatch {
		t.Fatalf("missing previous release got %s", got)
	}
}

func TestBoundFileVerification(t *testing.T) {
	value := validRequest(t, "pr")
	if err := request.ValidateAt(value, testNow); err != nil {
		t.Fatal(err)
	}
	if err := request.ValidateBoundFiles(value); err != nil {
		t.Fatal(err)
	}
	value.TrackedSourceManifest.Digest = strings.Repeat("0", 64)
	if got := reason(t, request.ValidateBoundFiles(value)); got != outcome.ReasonFailBindingMismatch {
		t.Fatalf("got %s", got)
	}
	value.TrackedSourceManifest.Path = filepath.Join(t.TempDir(), "missing")
	if got := reason(t, request.ValidateBoundFiles(value)); got != outcome.ReasonUnavailableRequiredInput {
		t.Fatalf("got %s", got)
	}
}

func TestNewScanIDsAreUUIDv4AndUnique(t *testing.T) {
	seen := map[string]bool{}
	for range 256 {
		id, err := request.NewScanID()
		if err != nil || !request.IsUUIDv4(id) || seen[id] {
			t.Fatalf("invalid or duplicate id %q: %v", id, err)
		}
		seen[id] = true
	}
}
