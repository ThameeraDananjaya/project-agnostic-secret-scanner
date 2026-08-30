package contract_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/outcome"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/request"
)

var runnerPath string
var repoRoot string

func TestMain(m *testing.M) {
	_, file, _, _ := runtime.Caller(0)
	repoRoot = filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", ".."))
	temporary, err := os.MkdirTemp("", "pscan-contract-test-")
	if err != nil {
		panic(err)
	}
	extension := ""
	if runtime.GOOS == "windows" {
		extension = ".exe"
	}
	runnerPath = filepath.Join(temporary, "scanner-runner"+extension)
	goCommand := filepath.Join(runtime.GOROOT(), "bin", "go")
	command := exec.Command(goCommand, "build", "-o", runnerPath, "./cmd/scanner-runner")
	command.Dir = repoRoot
	command.Env = append(os.Environ(), "GOTOOLCHAIN=local")
	if output, buildErr := command.CombinedOutput(); buildErr != nil {
		panic(string(output))
	}
	code := m.Run()
	_ = os.RemoveAll(temporary)
	os.Exit(code)
}

func validRequestFile(t *testing.T) (string, string, request.ScanRequest) {
	t.Helper()
	root := t.TempDir()
	manifest := filepath.Join(root, "tracked-manifest.json")
	data := []byte(`{"files":[]}`)
	if err := os.WriteFile(manifest, data, 0o600); err != nil {
		t.Fatal(err)
	}
	hash := sha256.Sum256(data)
	d := strings.Repeat("a", 64)
	value := request.ScanRequest{
		RequestSchemaVersion: "1.0", ScanID: "123e4567-e89b-42d3-a456-426614174000", Mode: "pr",
		ScannerReleaseDigest: d, EngineBinding: request.EngineBinding{Name: "gitleaks", Version: "1.0.0", BinaryDigest: d, AdapterVersion: "1.0.0"},
		RulePackDigest: d, PolicyDigest: d, AllowlistDigest: d,
		SourceBinding:         request.SourceBinding{BaseCommit: d, HeadCommit: d, MergeBase: d, HistoryRangeDigest: d, TrackedTreeDigest: d},
		TrackedSourceManifest: request.FileBinding{Path: manifest, Digest: hex.EncodeToString(hash[:])},
		FallbackRequirement:   request.FallbackRequirement{Mode: "disabled"},
		Limits:                request.Limits{TimeoutSeconds: 900, MaxArchiveDepth: 5, MaxArchiveEntries: 100000, MaxExpandedBytes: 2 << 30, MaxFileBytes: 512 << 20, MaxCompressionRatio: 1000, MaxMemoryBytes: 1 << 30, MaxCPUPercent: 100},
		OfflineRequired:       true, RedactionMode: "full", RequestedAt: time.Now().UTC().Format(time.RFC3339),
	}
	requestPath := filepath.Join(root, "request.json")
	payload, _ := json.Marshal(value)
	if err := os.WriteFile(requestPath, payload, 0o600); err != nil {
		t.Fatal(err)
	}
	return requestPath, filepath.Join(root, "workspaces"), value
}

func runCLI(t *testing.T, args ...string) (int, []byte, []byte) {
	t.Helper()
	command := exec.Command(runnerPath, args...)
	var stdout, stderr bytes.Buffer
	command.Stdout, command.Stderr = &stdout, &stderr
	err := command.Run()
	if err == nil {
		return 0, stdout.Bytes(), stderr.Bytes()
	}
	var exit *exec.ExitError
	if !errors.As(err, &exit) {
		t.Fatal(err)
	}
	return exit.ExitCode(), stdout.Bytes(), stderr.Bytes()
}

func decodeOne(t *testing.T, data []byte) map[string]any {
	t.Helper()
	decoder := json.NewDecoder(bytes.NewReader(data))
	var object map[string]any
	if err := decoder.Decode(&object); err != nil {
		t.Fatal(err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		t.Fatalf("stdout contained more than one object: %q", data)
	}
	return object
}

func TestValidRequestStopsBeforeEngineWithOneContentFreeOutcome(t *testing.T) {
	requestPath, workspaceRoot, _ := validRequestFile(t)
	code, stdout, stderr := runCLI(t, "--request", requestPath, "--workspace-root", workspaceRoot, "--attempt", "1")
	if code != 30 {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	object := decodeOne(t, stdout)
	if object["state"] != string(outcome.StateUnavailable) || object["reasonCode"] != string(outcome.ReasonUnavailableEngine) {
		t.Fatalf("unexpected outcome: %#v", object)
	}
	if string(stderr) != "scanner-runner: engine unavailable\n" {
		t.Fatalf("unsafe stderr: %q", stderr)
	}
	if bytes.Contains(stdout, []byte(requestPath)) || bytes.Contains(stdout, []byte(workspaceRoot)) || bytes.Contains(stderr, []byte(requestPath)) {
		t.Fatal("local path crossed output boundary")
	}
	if _, err := os.Stat(filepath.Join(workspaceRoot, "scan-123e4567-e89b-42d3-a456-426614174000-attempt-1")); !os.IsNotExist(err) {
		t.Fatal("workspace residue remains")
	}
}

func TestCorruptUnknownMajorAndInvocationInjectionFailClosed(t *testing.T) {
	root := t.TempDir()
	tests := []struct {
		name, payload string
		args          []string
		code          int
		reason        outcome.ReasonCode
		diagnostic    string
	}{
		{"corrupt", `{"requestSchemaVersion":"1.0",`, nil, 10, outcome.ReasonFailInputIntegrity, "scanner-runner: request rejected\n"},
		{"unknown-major", `{"requestSchemaVersion":"2.0"}`, nil, 20, outcome.ReasonIndeterminateSchemaUnsupported, "scanner-runner: request rejected\n"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(root, test.name+".json")
			if err := os.WriteFile(path, []byte(test.payload), 0o600); err != nil {
				t.Fatal(err)
			}
			code, stdout, stderr := runCLI(t, "--request", path)
			if code != test.code || decodeOne(t, stdout)["reasonCode"] != string(test.reason) || string(stderr) != test.diagnostic {
				t.Fatalf("code=%d stdout=%s stderr=%s", code, stdout, stderr)
			}
		})
	}
	injection := "\n::error file=synthetic::candidate text"
	code, stdout, stderr := runCLI(t, "--unknown-flag", injection)
	if code != 10 || bytes.Contains(stdout, []byte(injection)) || bytes.Contains(stderr, []byte(injection)) || string(stderr) != "scanner-runner: invalid invocation\n" {
		t.Fatalf("injection escaped: %d %q %q", code, stdout, stderr)
	}
}

func TestInvalidRequestIdentityEmitsExactlyOneSafeOutcome(t *testing.T) {
	for _, test := range []struct {
		name   string
		mutate func(map[string]any)
	}{
		{"scan-id", func(v map[string]any) { v["scanId"] = "invalid" }},
		{"mode", func(v map[string]any) { v["mode"] = "invalid" }},
		{"supersedes", func(v map[string]any) { v["supersedesScanId"] = "invalid" }},
	} {
		t.Run(test.name, func(t *testing.T) {
			requestPath, _, _ := validRequestFile(t)
			data, err := os.ReadFile(requestPath)
			if err != nil {
				t.Fatal(err)
			}
			var requestObject map[string]any
			_ = json.Unmarshal(data, &requestObject)
			test.mutate(requestObject)
			data, _ = json.Marshal(requestObject)
			if err := os.WriteFile(requestPath, data, 0o600); err != nil {
				t.Fatal(err)
			}
			code, stdout, _ := runCLI(t, "--request", requestPath)
			object := decodeOne(t, stdout)
			if code != 10 || object["scanId"] == "invalid" || object["mode"] != "unknown" || object["supersedesScanId"] != nil || object["reasonCode"] != string(outcome.ReasonFailInputIntegrity) {
				t.Fatalf("unsafe rejection: code=%d outcome=%#v", code, object)
			}
		})
	}
}

func TestDeterministicLogicalOutcomeAndWorkspaceCollision(t *testing.T) {
	requestPath, workspaceRoot, value := validRequestFile(t)
	_, first, _ := runCLI(t, "--request", requestPath, "--workspace-root", workspaceRoot)
	_, second, _ := runCLI(t, "--request", requestPath, "--workspace-root", workspaceRoot)
	a, b := decodeOne(t, first), decodeOne(t, second)
	for _, volatile := range []string{"startedAt", "endedAt", "durationMillis"} {
		delete(a, volatile)
		delete(b, volatile)
	}
	if !equalJSON(a, b) {
		t.Fatalf("logical outcomes differ:\n%#v\n%#v", a, b)
	}
	collision := filepath.Join(workspaceRoot, "scan-"+value.ScanID+"-attempt-1")
	if err := os.MkdirAll(collision, 0o700); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(collision, "must-survive")
	if err := os.WriteFile(marker, []byte("synthetic"), 0o600); err != nil {
		t.Fatal(err)
	}
	code, stdout, _ := runCLI(t, "--request", requestPath, "--workspace-root", workspaceRoot)
	if code != 30 || decodeOne(t, stdout)["reasonCode"] != string(outcome.ReasonUnavailableWorkspace) {
		t.Fatal("collision did not fail closed")
	}
	if _, err := os.Stat(marker); err != nil {
		t.Fatal("collision content was overwritten")
	}
}

func equalJSON(a, b map[string]any) bool {
	left, _ := json.Marshal(a)
	right, _ := json.Marshal(b)
	return bytes.Equal(left, right)
}

func TestScannerOwnedSchemasAreVersionedAndOutcomeReasonsExact(t *testing.T) {
	paths := []string{"contracts/scan-request/schema-1.0.json", "contracts/scan-outcome/schema-1.0.json", "contracts/release-manifest/schema-1.0.json", "contracts/global-revocation/schema-1.0.json", "contracts/rule-pack/schema-1.0.json"}
	ids := map[string]bool{}
	for _, relative := range paths {
		data, err := os.ReadFile(filepath.Join(repoRoot, filepath.FromSlash(relative)))
		if err != nil {
			t.Fatal(err)
		}
		for _, forbidden := range [][]byte{[]byte(`"projectPolicySchema"`), []byte(`"projectReceiptSchema"`), []byte(`"projectAllowlistSchema"`)} {
			if bytes.Contains(data, forbidden) {
				t.Fatalf("project-owned schema authority in %s", relative)
			}
		}
		var schema map[string]any
		if err := json.Unmarshal(data, &schema); err != nil {
			t.Fatalf("%s: %v", relative, err)
		}
		id, ok := schema["$id"].(string)
		if !ok || ids[id] || !strings.HasSuffix(id, "1.0") {
			t.Fatalf("invalid schema id %q", id)
		}
		ids[id] = true
	}
	data, _ := os.ReadFile(filepath.Join(repoRoot, "contracts", "scan-outcome", "schema-1.0.json"))
	var schema map[string]any
	_ = json.Unmarshal(data, &schema)
	properties := schema["properties"].(map[string]any)
	reasonSchema := properties["reasonCode"].(map[string]any)
	values := reasonSchema["enum"].([]any)
	want := map[string]bool{}
	for _, reason := range outcome.AllReasonCodes() {
		want[string(reason)] = true
	}
	if len(values) != len(want) {
		t.Fatalf("schema reasons=%d Go reasons=%d", len(values), len(want))
	}
	for _, value := range values {
		if !want[value.(string)] {
			t.Fatalf("schema-only reason %s", value)
		}
	}
}

func TestGitObjectIDSchemaDefinitionsAreFullLength(t *testing.T) {
	for _, relative := range []string{"contracts/scan-request/schema-1.0.json", "contracts/scan-outcome/schema-1.0.json", "contracts/release-manifest/schema-1.0.json"} {
		data, err := os.ReadFile(filepath.Join(repoRoot, filepath.FromSlash(relative)))
		if err != nil {
			t.Fatal(err)
		}
		var schema map[string]any
		if err := json.Unmarshal(data, &schema); err != nil {
			t.Fatal(err)
		}
		definitions := schema["$defs"].(map[string]any)
		pattern := definitions["gitOid"].(map[string]any)["pattern"].(string)
		compiled, err := regexp.Compile(pattern)
		if err != nil {
			t.Fatalf("%s: %v", relative, err)
		}
		if !compiled.MatchString(strings.Repeat("a", 40)) || !compiled.MatchString(strings.Repeat("b", 64)) || compiled.MatchString(strings.Repeat("c", 12)) || compiled.MatchString(strings.Repeat("d", 63)) {
			t.Fatalf("%s has unsafe Git object ID semantics", relative)
		}
	}
}
