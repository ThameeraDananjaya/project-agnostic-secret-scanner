package policy_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/outcome"
	. "github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/policy"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/verify"
)

type acceptingVerifier struct{}

func (acceptingVerifier) Verify(_ []byte, signature verify.DetachedSignature) error {
	if signature.Algorithm != "ed25519" || signature.TrustDomain == "" || signature.KeyID == "" {
		return verify.ErrInvalidSignature
	}
	return nil
}

var testNow = time.Date(2026, 9, 2, 0, 0, 0, 0, time.UTC)

func TestPolicyAndAllowlistProjection(t *testing.T) {
	policyRaw := []byte(`{"blockedClasses":["api-secret"],"minimumSeverity":"high","maxExceptionDays":7,"allowedExceptionEvidence":["synthetic"]}`)
	policyBinding := binding(policyRaw, "synthetic-policy", PolicyAdapterVersion, "project-policy")
	project, err := LoadPolicy(policyRaw, policyBinding, family("synthetic-policy"), acceptingVerifier{}, testNow)
	if err != nil {
		t.Fatal(err)
	}
	context := ExceptionContext{ScannerReleaseDigest: digest("scanner"), RulePackDigest: digest("rules"), PolicyDigest: policyBinding.Digest, SourceDigest: digest("source")}
	allowlistRaw := exceptionDocument(t, validException(context))
	allowlistBinding := binding(allowlistRaw, "synthetic-allowlist", AllowlistAdapterVersion, "project-allowlist")
	allowlist, err := LoadAllowlist(allowlistRaw, allowlistBinding, family("synthetic-allowlist"), acceptingVerifier{}, testNow, context)
	if err != nil {
		t.Fatal(err)
	}
	evaluation := EvaluationContext{
		Now:                  testNow,
		ScannerReleaseDigest: context.ScannerReleaseDigest, RulePackDigest: context.RulePackDigest,
		PolicyDigest: context.PolicyDigest, AllowlistDigest: allowlistBinding.Digest, SourceDigest: context.SourceDigest,
		Revocations: map[string]bool{},
	}
	finding := Finding{RuleID: "synthetic-rule", Class: "api-secret", Severity: SeverityHigh, ScopeDigest: digest("scope")}
	if decision := Evaluate(evaluation, project, allowlist, []Finding{finding}); decision.Blocking || decision.Reason != outcome.ReasonPassNoBlockingFindings {
		t.Fatalf("exact valid exception was not honoured: %+v", decision)
	}
	missingRevocations := evaluation
	missingRevocations.Revocations = nil
	if decision := Evaluate(missingRevocations, project, allowlist, []Finding{finding}); !decision.Blocking || decision.Reason != outcome.ReasonFailBindingMismatch {
		t.Fatalf("missing revocation state did not fail closed: %+v", decision)
	}
	evaluation.Revocations = map[string]bool{verify.RevocationKey("policy", context.PolicyDigest): true}
	if decision := Evaluate(evaluation, project, allowlist, []Finding{finding}); !decision.Blocking || decision.Reason != outcome.ReasonFailRevokedDependency {
		t.Fatalf("global revocation did not precede the exception: %+v", decision)
	}
}

func TestCredentialClassesCannotBeAllowlisted(t *testing.T) {
	policyRaw := []byte(`{"blockedClasses":[],"minimumSeverity":"critical","maxExceptionDays":7,"allowedExceptionEvidence":["synthetic"]}`)
	policyBinding := binding(policyRaw, "synthetic-policy", PolicyAdapterVersion, "project-policy")
	project, err := LoadPolicy(policyRaw, policyBinding, family("synthetic-policy"), acceptingVerifier{}, testNow)
	if err != nil {
		t.Fatal(err)
	}
	context := ExceptionContext{ScannerReleaseDigest: digest("scanner"), RulePackDigest: digest("rules"), PolicyDigest: policyBinding.Digest, SourceDigest: digest("source")}
	exception := validException(context)
	exception.Class = "access-token"
	raw := exceptionDocument(t, exception)
	if _, err := LoadAllowlist(raw, binding(raw, "synthetic-allowlist", AllowlistAdapterVersion, "project-allowlist"), family("synthetic-allowlist"), acceptingVerifier{}, testNow, context); err == nil {
		t.Fatal("credential-class exception was admitted")
	}
	emptyRaw := exceptionDocument(t)
	emptyBinding := binding(emptyRaw, "synthetic-allowlist", AllowlistAdapterVersion, "project-allowlist")
	allowlist, err := LoadAllowlist(emptyRaw, emptyBinding, family("synthetic-allowlist"), acceptingVerifier{}, testNow, context)
	if err != nil {
		t.Fatal(err)
	}
	evaluation := EvaluationContext{Now: testNow, ScannerReleaseDigest: context.ScannerReleaseDigest, RulePackDigest: context.RulePackDigest, PolicyDigest: context.PolicyDigest, AllowlistDigest: emptyBinding.Digest, SourceDigest: context.SourceDigest, Revocations: map[string]bool{}}
	decision := Evaluate(evaluation, project, allowlist, []Finding{{RuleID: "synthetic-rule", Class: "access-token", Severity: SeverityLow, ScopeDigest: digest("scope")}})
	if !decision.Blocking || decision.Reason != outcome.ReasonFailFindingDetected {
		t.Fatalf("credential class was not mandatory-blocking: %+v", decision)
	}
	decision = Evaluate(evaluation, project, allowlist, []Finding{{RuleID: "synthetic-rule", Class: "api-secret", Severity: SeverityLow, ScopeDigest: digest("scope"), Credential: true}})
	if !decision.Blocking || decision.Reason != outcome.ReasonFailFindingDetected {
		t.Fatalf("credential classification was not mandatory-blocking: %+v", decision)
	}
}

func TestInvalidExceptionsFailClosed(t *testing.T) {
	context := ExceptionContext{ScannerReleaseDigest: digest("scanner"), RulePackDigest: digest("rules"), PolicyDigest: digest("policy"), SourceDigest: digest("source")}
	cases := map[string]func(*Exception){
		"same owner and approver": func(e *Exception) { e.Approver = e.Owner },
		"expired":                 func(e *Exception) { e.ExpiresAt = "2026-09-01T00:00:00Z" },
		"more than 30 days":       func(e *Exception) { e.ExpiresAt = "2026-10-10T00:00:00Z" },
		"missing invalidation":    func(e *Exception) { e.InvalidatesOn = e.InvalidatesOn[:5] },
		"wrong policy":            func(e *Exception) { e.PolicyDigest = digest("other-policy") },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			exception := validException(context)
			mutate(&exception)
			raw := exceptionDocument(t, exception)
			if _, err := LoadAllowlist(raw, binding(raw, "synthetic-allowlist", AllowlistAdapterVersion, "project-allowlist"), family("synthetic-allowlist"), acceptingVerifier{}, testNow, context); err == nil {
				t.Fatal("invalid exception was admitted")
			}
		})
	}
}

func TestDuplicateExceptionIdentityFailsClosedAcrossDifferentScopes(t *testing.T) {
	context := ExceptionContext{ScannerReleaseDigest: digest("scanner"), RulePackDigest: digest("rules"), PolicyDigest: digest("policy"), SourceDigest: digest("source")}
	first := validException(context)
	second := validException(context)
	second.RuleID = "different-synthetic-rule"
	second.ScopeDigest = digest("different-scope")
	raw := exceptionDocument(t, first, second)
	if _, err := LoadAllowlist(raw, binding(raw, "synthetic-allowlist", AllowlistAdapterVersion, "project-allowlist"), family("synthetic-allowlist"), acceptingVerifier{}, testNow, context); err == nil {
		t.Fatal("duplicate exception identity was admitted across different scopes")
	}
}

func TestPolicySpecificExceptionMaximumIsEnforced(t *testing.T) {
	policyRaw := []byte(`{"blockedClasses":["api-secret"],"minimumSeverity":"critical","maxExceptionDays":7,"allowedExceptionEvidence":["synthetic"]}`)
	policyBinding := binding(policyRaw, "synthetic-policy", PolicyAdapterVersion, "project-policy")
	project, err := LoadPolicy(policyRaw, policyBinding, family("synthetic-policy"), acceptingVerifier{}, testNow)
	if err != nil {
		t.Fatal(err)
	}
	context := ExceptionContext{ScannerReleaseDigest: digest("scanner"), RulePackDigest: digest("rules"), PolicyDigest: policyBinding.Digest, SourceDigest: digest("source")}
	exception := validException(context)
	exception.ExpiresAt = "2026-09-10T00:00:00Z"
	allowlistRaw := exceptionDocument(t, exception)
	allowlistBinding := binding(allowlistRaw, "synthetic-allowlist", AllowlistAdapterVersion, "project-allowlist")
	allowlist, err := LoadAllowlist(allowlistRaw, allowlistBinding, family("synthetic-allowlist"), acceptingVerifier{}, testNow, context)
	if err != nil {
		t.Fatal(err)
	}
	evaluation := EvaluationContext{Now: testNow, ScannerReleaseDigest: context.ScannerReleaseDigest, RulePackDigest: context.RulePackDigest, PolicyDigest: context.PolicyDigest, AllowlistDigest: allowlistBinding.Digest, SourceDigest: context.SourceDigest, Revocations: map[string]bool{}}
	decision := Evaluate(evaluation, project, allowlist, []Finding{{RuleID: exception.RuleID, Class: exception.Class, Severity: SeverityHigh, ScopeDigest: exception.ScopeDigest}})
	if !decision.Blocking || decision.Reason != outcome.ReasonFailAllowlistInvalid {
		t.Fatalf("project-specific seven-day cap was not enforced: %+v", decision)
	}
}

func TestExceptionExpiryIsRecheckedAtEvaluationTime(t *testing.T) {
	policyRaw := []byte(`{"blockedClasses":["api-secret"],"minimumSeverity":"critical","maxExceptionDays":30,"allowedExceptionEvidence":["synthetic"]}`)
	policyBinding := binding(policyRaw, "synthetic-policy", PolicyAdapterVersion, "project-policy")
	project, err := LoadPolicy(policyRaw, policyBinding, family("synthetic-policy"), acceptingVerifier{}, testNow)
	if err != nil {
		t.Fatal(err)
	}
	context := ExceptionContext{ScannerReleaseDigest: digest("scanner"), RulePackDigest: digest("rules"), PolicyDigest: policyBinding.Digest, SourceDigest: digest("source")}
	exception := validException(context)
	allowlistRaw := exceptionDocument(t, exception)
	allowlistBinding := binding(allowlistRaw, "synthetic-allowlist", AllowlistAdapterVersion, "project-allowlist")
	allowlist, err := LoadAllowlist(allowlistRaw, allowlistBinding, family("synthetic-allowlist"), acceptingVerifier{}, testNow, context)
	if err != nil {
		t.Fatal(err)
	}
	evaluation := EvaluationContext{Now: testNow.Add(7 * 24 * time.Hour), ScannerReleaseDigest: context.ScannerReleaseDigest, RulePackDigest: context.RulePackDigest, PolicyDigest: context.PolicyDigest, AllowlistDigest: allowlistBinding.Digest, SourceDigest: context.SourceDigest, Revocations: map[string]bool{}}
	decision := Evaluate(evaluation, project, allowlist, []Finding{{RuleID: exception.RuleID, Class: exception.Class, Severity: SeverityHigh, ScopeDigest: exception.ScopeDigest}})
	if !decision.Blocking || decision.Reason != outcome.ReasonFailAllowlistInvalid {
		t.Fatalf("expired retained projection was not rejected: %+v", decision)
	}
}

func TestProjectionRejectsDuplicateAndUnknownMembers(t *testing.T) {
	duplicate := []byte(`{"blockedClasses":[],"blockedClasses":[],"minimumSeverity":"high","maxExceptionDays":7,"allowedExceptionEvidence":[]}`)
	if _, err := LoadPolicy(duplicate, binding(duplicate, "synthetic-policy", PolicyAdapterVersion, "project-policy"), family("synthetic-policy"), acceptingVerifier{}, testNow); err == nil {
		t.Fatal("duplicate member was admitted")
	}
	unknown := []byte(`{"blockedClasses":[],"minimumSeverity":"high","maxExceptionDays":7,"allowedExceptionEvidence":[],"newAuthority":true}`)
	if _, err := LoadPolicy(unknown, binding(unknown, "synthetic-policy", PolicyAdapterVersion, "project-policy"), family("synthetic-policy"), acceptingVerifier{}, testNow); err == nil {
		t.Fatal("unknown current-minor member was admitted")
	}
	badBinding := binding([]byte(`{"blockedClasses":[],"minimumSeverity":"high","maxExceptionDays":7,"allowedExceptionEvidence":[]}`), "synthetic-policy", PolicyAdapterVersion, "project-policy")
	badBinding.SchemaVersion = "2.0"
	if _, err := LoadPolicy([]byte(`{"blockedClasses":[],"minimumSeverity":"high","maxExceptionDays":7,"allowedExceptionEvidence":[]}`), badBinding, family("synthetic-policy"), acceptingVerifier{}, testNow); err == nil {
		t.Fatal("unknown major was admitted")
	}
}

func binding(raw []byte, schemaFamily, adapter, trustDomain string) verify.DocumentBinding {
	return verify.DocumentBinding{
		SchemaFamily: schemaFamily, SchemaVersion: "1.0", AdapterVersion: adapter, Digest: verify.DigestBytes(raw),
		Signature: verify.DetachedSignature{TrustDomain: trustDomain, Algorithm: "ed25519", KeyID: "synthetic-key", Value: "unused-by-test-verifier"},
	}
}

func family(name string) verify.FamilyWindow {
	return verify.FamilyWindow{Family: name, CurrentMajor: 1, Minors: []verify.MinorCompatibility{{Minor: 0}}}
}

func validException(context ExceptionContext) Exception {
	return Exception{
		ExceptionID: "synthetic-exception", Owner: "synthetic-owner", Approver: "synthetic-approver", Rationale: "Synthetic false positive evidence.",
		RuleID: "synthetic-rule", Class: "api-secret", ScopeDigest: digest("scope"), EvidenceType: "synthetic", EvidenceDigest: digest("evidence"),
		CreatedAt: "2026-09-01T00:00:00Z", ExpiresAt: "2026-09-08T00:00:00Z", ScannerReleaseDigest: context.ScannerReleaseDigest,
		RulePackDigest: context.RulePackDigest, PolicyDigest: context.PolicyDigest, SourceDigest: context.SourceDigest,
		InvalidatesOn: []string{"scope", "policy", "rule", "scanner", "source", "evidence"},
	}
}

func exceptionDocument(t *testing.T, exceptions ...Exception) []byte {
	t.Helper()
	value, err := json.Marshal(struct {
		Exceptions []Exception `json:"exceptions"`
	}{Exceptions: exceptions})
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func digest(seed string) string { return verify.DigestBytes([]byte(seed)) }
