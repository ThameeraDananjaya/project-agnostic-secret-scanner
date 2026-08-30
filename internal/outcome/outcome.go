// Package outcome defines the scanner's content-free public result contract.
package outcome

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"time"
)

const SchemaVersion = "1.0"

type State string

const (
	StatePending       State = "PENDING"
	StateRunning       State = "RUNNING"
	StatePass          State = "PASS"
	StateFail          State = "FAIL"
	StateIndeterminate State = "INDETERMINATE"
	StateUnavailable   State = "UNAVAILABLE"
	StateRetrying      State = "RETRYING"
	StateTerminal      State = "TERMINAL"
)

type ReasonCode string

const (
	ReasonPassNoBlockingFindings          ReasonCode = "PASS_NO_BLOCKING_FINDINGS"
	ReasonFailFindingDetected             ReasonCode = "FAIL_FINDING_DETECTED"
	ReasonFailPolicyDenied                ReasonCode = "FAIL_POLICY_DENIED"
	ReasonFailAllowlistInvalid            ReasonCode = "FAIL_ALLOWLIST_INVALID"
	ReasonFailBindingMismatch             ReasonCode = "FAIL_BINDING_MISMATCH"
	ReasonFailInputIntegrity              ReasonCode = "FAIL_INPUT_INTEGRITY"
	ReasonFailScannerIntegrity            ReasonCode = "FAIL_SCANNER_INTEGRITY"
	ReasonFailSignatureInvalid            ReasonCode = "FAIL_SIGNATURE_INVALID"
	ReasonFailRevokedDependency           ReasonCode = "FAIL_REVOKED_DEPENDENCY"
	ReasonIndeterminateUnsupportedInput   ReasonCode = "INDETERMINATE_UNSUPPORTED_INPUT"
	ReasonIndeterminateIncompleteCoverage ReasonCode = "INDETERMINATE_INCOMPLETE_COVERAGE"
	ReasonIndeterminateSchemaUnsupported  ReasonCode = "INDETERMINATE_SCHEMA_UNSUPPORTED"
	ReasonIndeterminateConflictingResults ReasonCode = "INDETERMINATE_CONFLICTING_RESULTS"
	ReasonIndeterminateRedactionUnproven  ReasonCode = "INDETERMINATE_REDACTION_UNPROVEN"
	ReasonIndeterminateResourceLimit      ReasonCode = "INDETERMINATE_RESOURCE_LIMIT"
	ReasonIndeterminateTimeout            ReasonCode = "INDETERMINATE_TIMEOUT"
	ReasonUnavailableEngine               ReasonCode = "UNAVAILABLE_ENGINE"
	ReasonUnavailableRuntime              ReasonCode = "UNAVAILABLE_RUNTIME"
	ReasonUnavailableRequiredInput        ReasonCode = "UNAVAILABLE_REQUIRED_INPUT"
	ReasonUnavailableWorkspace            ReasonCode = "UNAVAILABLE_WORKSPACE"
	ReasonTerminalRetryExhausted          ReasonCode = "TERMINAL_RETRY_EXHAUSTED"
	ReasonTerminalTimeoutExhausted        ReasonCode = "TERMINAL_TIMEOUT_EXHAUSTED"
	ReasonTerminalUnresolvedConflict      ReasonCode = "TERMINAL_UNRESOLVED_CONFLICT"
	ReasonTerminalInternalInvariant       ReasonCode = "TERMINAL_INTERNAL_INVARIANT"
)

type PermittedAction string

const (
	ActionVerifyOtherGates         PermittedAction = "VERIFY_OTHER_GATES"
	ActionCreateCorrectedCandidate PermittedAction = "CREATE_CORRECTED_CANDIDATE"
	ActionDiagnoseInput            PermittedAction = "DIAGNOSE_INPUT"
	ActionRetryIfEligible          PermittedAction = "RETRY_IF_ELIGIBLE"
	ActionWait                     PermittedAction = "WAIT"
	ActionStopAndEscalate          PermittedAction = "STOP_AND_ESCALATE"
)

type Bindings struct {
	ScannerReleaseDigest        string `json:"scannerReleaseDigest"`
	EngineName                  string `json:"engineName"`
	EngineVersion               string `json:"engineVersion"`
	EngineBinaryDigest          string `json:"engineBinaryDigest"`
	AdapterVersion              string `json:"adapterVersion"`
	RulePackDigest              string `json:"rulePackDigest"`
	PolicyDigest                string `json:"policyDigest"`
	AllowlistDigest             string `json:"allowlistDigest"`
	RequestSchemaVersion        string `json:"requestSchemaVersion"`
	SourceBaseCommit            string `json:"sourceBaseCommit,omitempty"`
	SourceHeadCommit            string `json:"sourceHeadCommit"`
	SourceMergeBase             string `json:"sourceMergeBase,omitempty"`
	HistoryRangeDigest          string `json:"historyRangeDigest"`
	TrackedTreeDigest           string `json:"trackedTreeDigest"`
	TrackedSourceManifestDigest string `json:"trackedSourceManifestDigest"`
	BuildContextManifestDigest  string `json:"buildContextManifestDigest,omitempty"`
	ArtifactManifestDigest      string `json:"artifactManifestDigest,omitempty"`
}

type Outcome struct {
	OutcomeSchemaVersion string          `json:"outcomeSchemaVersion"`
	ScanID               string          `json:"scanId"`
	Mode                 string          `json:"mode"`
	State                State           `json:"state"`
	ReasonCode           ReasonCode      `json:"reasonCode,omitempty"`
	Bindings             *Bindings       `json:"bindings,omitempty"`
	AttemptNumber        int             `json:"attemptNumber"`
	StartedAt            time.Time       `json:"startedAt"`
	EndedAt              time.Time       `json:"endedAt"`
	DurationMillis       int64           `json:"durationMillis"`
	CoverageClasses      []string        `json:"coverageClasses"`
	PermittedAction      PermittedAction `json:"permittedAction"`
	SupersedesScanID     string          `json:"supersedesScanId,omitempty"`
}

var stateExit = map[State]int{
	StatePass:          0,
	StateFail:          10,
	StateIndeterminate: 20,
	StateUnavailable:   30,
	StateTerminal:      40,
}

var uuidV4Pattern = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
var digestPattern = regexp.MustCompile(`^[0-9a-f]{64}$`)
var gitOIDPattern = regexp.MustCompile(`^(?:[0-9a-f]{40}|[0-9a-f]{64})$`)
var versionTokenPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._+-]{0,63}$`)

var reasonState = map[ReasonCode]State{
	ReasonPassNoBlockingFindings:          StatePass,
	ReasonFailFindingDetected:             StateFail,
	ReasonFailPolicyDenied:                StateFail,
	ReasonFailAllowlistInvalid:            StateFail,
	ReasonFailBindingMismatch:             StateFail,
	ReasonFailInputIntegrity:              StateFail,
	ReasonFailScannerIntegrity:            StateFail,
	ReasonFailSignatureInvalid:            StateFail,
	ReasonFailRevokedDependency:           StateFail,
	ReasonIndeterminateUnsupportedInput:   StateIndeterminate,
	ReasonIndeterminateIncompleteCoverage: StateIndeterminate,
	ReasonIndeterminateSchemaUnsupported:  StateIndeterminate,
	ReasonIndeterminateConflictingResults: StateIndeterminate,
	ReasonIndeterminateRedactionUnproven:  StateIndeterminate,
	ReasonIndeterminateResourceLimit:      StateIndeterminate,
	ReasonIndeterminateTimeout:            StateIndeterminate,
	ReasonUnavailableEngine:               StateUnavailable,
	ReasonUnavailableRuntime:              StateUnavailable,
	ReasonUnavailableRequiredInput:        StateUnavailable,
	ReasonUnavailableWorkspace:            StateUnavailable,
	ReasonTerminalRetryExhausted:          StateTerminal,
	ReasonTerminalTimeoutExhausted:        StateTerminal,
	ReasonTerminalUnresolvedConflict:      StateTerminal,
	ReasonTerminalInternalInvariant:       StateTerminal,
}

func StateForReason(reason ReasonCode) (State, bool) {
	state, ok := reasonState[reason]
	return state, ok
}

func ExitCode(state State) (int, bool) {
	code, ok := stateExit[state]
	return code, ok
}

func ActionForState(state State) (PermittedAction, bool) {
	actions := map[State]PermittedAction{
		StatePending:       ActionWait,
		StateRunning:       ActionWait,
		StatePass:          ActionVerifyOtherGates,
		StateFail:          ActionCreateCorrectedCandidate,
		StateIndeterminate: ActionDiagnoseInput,
		StateUnavailable:   ActionRetryIfEligible,
		StateRetrying:      ActionWait,
		StateTerminal:      ActionStopAndEscalate,
	}
	action, ok := actions[state]
	return action, ok
}

func NewTransient(scanID, mode string, state State, attempt int, at time.Time) (Outcome, error) {
	if state != StatePending && state != StateRunning && state != StateRetrying {
		return Outcome{}, errors.New("state is not transient")
	}
	return Outcome{OutcomeSchemaVersion: SchemaVersion, ScanID: scanID, Mode: mode, State: state, AttemptNumber: attempt, StartedAt: at.UTC(), EndedAt: at.UTC(), CoverageClasses: []string{}, PermittedAction: ActionWait}, nil
}

func NewTerminal(scanID, mode string, reason ReasonCode, attempt int, start, end time.Time) (Outcome, error) {
	state, ok := StateForReason(reason)
	if !ok {
		return Outcome{}, fmt.Errorf("unknown reason code")
	}
	action, ok := ActionForState(state)
	if !ok {
		return Outcome{}, fmt.Errorf("state is not terminal")
	}
	if end.Before(start) {
		return Outcome{}, errors.New("end precedes start")
	}
	return Outcome{
		OutcomeSchemaVersion: SchemaVersion,
		ScanID:               scanID,
		Mode:                 mode,
		State:                state,
		ReasonCode:           reason,
		AttemptNumber:        attempt,
		StartedAt:            start.UTC(),
		EndedAt:              end.UTC(),
		DurationMillis:       end.Sub(start).Milliseconds(),
		CoverageClasses:      []string{},
		PermittedAction:      action,
	}, nil
}

func (o Outcome) Validate() error {
	if o.OutcomeSchemaVersion != SchemaVersion || !uuidV4Pattern.MatchString(o.ScanID) {
		return errors.New("missing required outcome identity")
	}
	if o.SupersedesScanID != "" && (!uuidV4Pattern.MatchString(o.SupersedesScanID) || o.SupersedesScanID == o.ScanID) {
		return errors.New("invalid recovery identity")
	}
	if o.Mode != "local" && o.Mode != "pr" && o.Mode != "release" && o.Mode != "unknown" {
		return errors.New("invalid mode")
	}
	_, terminal := ExitCode(o.State)
	if terminal {
		expected, ok := StateForReason(o.ReasonCode)
		if !ok || expected != o.State {
			return errors.New("state and reason code mismatch")
		}
	} else if (o.State != StatePending && o.State != StateRunning && o.State != StateRetrying) || o.ReasonCode != "" {
		return errors.New("invalid transient state or reason")
	}
	action, ok := ActionForState(o.State)
	if !ok || action != o.PermittedAction {
		return errors.New("state and permitted action mismatch")
	}
	if o.AttemptNumber < 1 || o.AttemptNumber > 3 || o.EndedAt.Before(o.StartedAt) || o.DurationMillis < 0 {
		return errors.New("invalid attempt or timing")
	}
	if o.CoverageClasses == nil {
		return errors.New("coverageClasses must be an explicit array")
	}
	allowedCoverage := map[string]bool{"git-history": true, "tracked-source": true, "build-context": true, "artifacts": true, "container-layers": true, "nested-archives": true}
	seenCoverage := map[string]bool{}
	for _, class := range o.CoverageClasses {
		if !allowedCoverage[class] || seenCoverage[class] {
			return errors.New("invalid coverage class")
		}
		seenCoverage[class] = true
	}
	if o.Bindings != nil {
		b := o.Bindings
		if !digestPattern.MatchString(b.ScannerReleaseDigest) || b.EngineName != "gitleaks" || !versionTokenPattern.MatchString(b.EngineVersion) || !digestPattern.MatchString(b.EngineBinaryDigest) || !versionTokenPattern.MatchString(b.AdapterVersion) || !digestPattern.MatchString(b.RulePackDigest) || !digestPattern.MatchString(b.PolicyDigest) || !digestPattern.MatchString(b.AllowlistDigest) || !gitOIDPattern.MatchString(b.SourceHeadCommit) || !digestPattern.MatchString(b.HistoryRangeDigest) || !digestPattern.MatchString(b.TrackedTreeDigest) || !digestPattern.MatchString(b.TrackedSourceManifestDigest) || !optionalGitOID(b.SourceBaseCommit) || !optionalGitOID(b.SourceMergeBase) || !optionalDigest(b.BuildContextManifestDigest) || !optionalDigest(b.ArtifactManifestDigest) {
			return errors.New("invalid bindings")
		}
		version, err := parseSchemaVersion(b.RequestSchemaVersion)
		if err != nil || version != 1 {
			return errors.New("unsupported request schema binding")
		}
	}
	return nil
}

func optionalDigest(value string) bool { return value == "" || digestPattern.MatchString(value) }
func optionalGitOID(value string) bool { return value == "" || gitOIDPattern.MatchString(value) }

func parseSchemaVersion(value string) (int, error) {
	match := regexp.MustCompile(`^([0-9]+)\.[0-9]+$`).FindStringSubmatch(value)
	if len(match) != 2 || match[1] != "1" {
		return 0, errors.New("unsupported schema version")
	}
	return 1, nil
}

func AllReasonCodes() []ReasonCode {
	result := make([]ReasonCode, 0, len(reasonState))
	for reason := range reasonState {
		result = append(result, reason)
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result
}

func AllStates() []State {
	return []State{StatePending, StateRunning, StatePass, StateFail, StateIndeterminate, StateUnavailable, StateRetrying, StateTerminal}
}
