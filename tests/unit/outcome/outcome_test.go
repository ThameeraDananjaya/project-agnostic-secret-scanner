package outcome_test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/outcome"
)

var fixedTime = time.Date(2026, 8, 30, 20, 0, 0, 0, time.UTC)

func TestCompleteReasonStateExitMapping(t *testing.T) {
	if len(outcome.AllStates()) != 8 {
		t.Fatalf("got %d states", len(outcome.AllStates()))
	}
	expected := map[outcome.State]int{outcome.StatePass: 0, outcome.StateFail: 10, outcome.StateIndeterminate: 20, outcome.StateUnavailable: 30, outcome.StateTerminal: 40}
	reasons := outcome.AllReasonCodes()
	if len(reasons) != 24 {
		t.Fatalf("got %d reason codes", len(reasons))
	}
	counts := map[outcome.State]int{}
	for _, reason := range reasons {
		state, ok := outcome.StateForReason(reason)
		if !ok {
			t.Fatalf("unmapped %s", reason)
		}
		code, ok := outcome.ExitCode(state)
		if !ok || code != expected[state] {
			t.Fatalf("wrong exit for %s", reason)
		}
		counts[state]++
	}
	if !reflect.DeepEqual(counts, map[outcome.State]int{outcome.StatePass: 1, outcome.StateFail: 8, outcome.StateIndeterminate: 7, outcome.StateUnavailable: 4, outcome.StateTerminal: 4}) {
		t.Fatalf("unexpected groups: %#v", counts)
	}
	for _, state := range []outcome.State{outcome.StatePending, outcome.StateRunning, outcome.StateRetrying} {
		if _, ok := outcome.ExitCode(state); ok {
			t.Errorf("transient state %s has exit", state)
		}
		value, err := outcome.NewTransient("123e4567-e89b-42d3-a456-426614174000", "pr", state, 1, fixedTime)
		if err != nil || value.Validate() != nil || value.PermittedAction != outcome.ActionWait || value.ReasonCode != "" {
			t.Errorf("invalid transient state %s: %#v %v", state, value, err)
		}
	}
}

func TestSerializerEmitsOneContentFreeObject(t *testing.T) {
	value, err := outcome.NewTerminal("123e4567-e89b-42d3-a456-426614174000", "pr", outcome.ReasonUnavailableEngine, 1, fixedTime, fixedTime)
	if err != nil {
		t.Fatal(err)
	}
	value.Bindings = &outcome.Bindings{ScannerReleaseDigest: strings.Repeat("a", 64), EngineName: "gitleaks", EngineVersion: "1.0.0", EngineBinaryDigest: strings.Repeat("b", 64), AdapterVersion: "1.0.0", RulePackDigest: strings.Repeat("c", 64), PolicyDigest: strings.Repeat("d", 64), AllowlistDigest: strings.Repeat("e", 64), RequestSchemaVersion: "1.0", SourceBaseCommit: strings.Repeat("f", 64), SourceHeadCommit: strings.Repeat("a", 64), SourceMergeBase: strings.Repeat("b", 64), HistoryRangeDigest: strings.Repeat("c", 64), TrackedTreeDigest: strings.Repeat("d", 64), TrackedSourceManifestDigest: strings.Repeat("e", 64)}
	var buffer bytes.Buffer
	if err := outcome.Serialize(&buffer, value); err != nil {
		t.Fatal(err)
	}
	if bytes.Count(buffer.Bytes(), []byte("\n")) != 1 || !bytes.HasSuffix(buffer.Bytes(), []byte("\n")) {
		t.Fatalf("not exactly one JSON line: %q", buffer.Bytes())
	}
	var object map[string]any
	if err := json.Unmarshal(buffer.Bytes(), &object); err != nil {
		t.Fatal(err)
	}
	for _, forbidden := range []string{"findings", "findingCount", "path", "location", "ruleName", "providerName", "author", "excerpt", "allowlistDetails"} {
		if _, exists := object[forbidden]; exists {
			t.Errorf("forbidden field %s", forbidden)
		}
	}
}

func TestOutcomeBindingsAcceptOnlyFullGitObjectIDs(t *testing.T) {
	for _, width := range []int{40, 64} {
		value, _ := outcome.NewTerminal("123e4567-e89b-42d3-a456-426614174000", "pr", outcome.ReasonUnavailableEngine, 1, fixedTime, fixedTime)
		d := strings.Repeat("a", 64)
		oid := strings.Repeat("b", width)
		value.Bindings = &outcome.Bindings{ScannerReleaseDigest: d, EngineName: "gitleaks", EngineVersion: "1.0.0", EngineBinaryDigest: d, AdapterVersion: "1.0.0", RulePackDigest: d, PolicyDigest: d, AllowlistDigest: d, RequestSchemaVersion: "1.0", SourceBaseCommit: oid, SourceHeadCommit: oid, SourceMergeBase: oid, HistoryRangeDigest: d, TrackedTreeDigest: d, TrackedSourceManifestDigest: d}
		if err := value.Validate(); err != nil {
			t.Fatalf("%d-hex full object ID rejected: %v", width, err)
		}
	}
	value, _ := outcome.NewTerminal("123e4567-e89b-42d3-a456-426614174000", "pr", outcome.ReasonUnavailableEngine, 1, fixedTime, fixedTime)
	d := strings.Repeat("a", 64)
	value.Bindings = &outcome.Bindings{ScannerReleaseDigest: d, EngineName: "gitleaks", EngineVersion: "1.0.0", EngineBinaryDigest: d, AdapterVersion: "1.0.0", RulePackDigest: d, PolicyDigest: d, AllowlistDigest: d, RequestSchemaVersion: "1.0", SourceHeadCommit: strings.Repeat("b", 12), HistoryRangeDigest: d, TrackedTreeDigest: d, TrackedSourceManifestDigest: d}
	if err := value.Validate(); err == nil {
		t.Fatal("abbreviated object ID was accepted")
	}
}

func TestSerializerRejectsStateReasonAndInjection(t *testing.T) {
	value, _ := outcome.NewTerminal("123e4567-e89b-42d3-a456-426614174000", "pr", outcome.ReasonFailInputIntegrity, 1, fixedTime, fixedTime)
	value.State = outcome.StatePass
	if err := outcome.Serialize(&bytes.Buffer{}, value); err == nil {
		t.Fatal("accepted mismatched state")
	}
	value, _ = outcome.NewTerminal("123e4567-e89b-42d3-a456-426614174000", "pr\n::error::injected", outcome.ReasonFailInputIntegrity, 1, fixedTime, fixedTime)
	if err := outcome.Serialize(&bytes.Buffer{}, value); err == nil {
		t.Fatal("accepted injected mode")
	}
}

func TestSemanticDeterminismExcludesAllowedVolatileFields(t *testing.T) {
	a, _ := outcome.NewTerminal("123e4567-e89b-42d3-a456-426614174000", "pr", outcome.ReasonUnavailableEngine, 1, fixedTime, fixedTime.Add(time.Second))
	b, _ := outcome.NewTerminal("223e4567-e89b-42d3-a456-426614174000", "pr", outcome.ReasonUnavailableEngine, 1, fixedTime.Add(time.Hour), fixedTime.Add(time.Hour+2*time.Second))
	normalize := func(value outcome.Outcome) []string {
		return []string{value.OutcomeSchemaVersion, value.Mode, string(value.State), string(value.ReasonCode), string(value.PermittedAction)}
	}
	if !reflect.DeepEqual(normalize(a), normalize(b)) {
		t.Fatal("logical outcome drift")
	}
	values := outcome.AllReasonCodes()
	sort.Slice(values, func(i, j int) bool { return values[i] < values[j] })
	if len(values) == 0 {
		t.Fatal("missing reasons")
	}
}
