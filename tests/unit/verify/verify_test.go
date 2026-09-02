package verify_test

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"testing"
	"time"

	. "github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/verify"
)

type acceptingVerifier struct{}

func (acceptingVerifier) Verify(_ []byte, signature DetachedSignature) error {
	if signature.Algorithm != "ed25519" || !IsToken(signature.TrustDomain) || !IsToken(signature.KeyID) {
		return ErrInvalidSignature
	}
	return nil
}

func TestEd25519VerifierRFC8032Vector(t *testing.T) {
	publicKey, _ := hex.DecodeString("d75a980182b10ab7d54bfed3c964073a0ee172f3daa62325af021a68f707511a")
	signatureBytes, _ := hex.DecodeString("e5564300c360ac729086e2cc806e828a84877f1eb8e5d974d873e065224901555fb8821590a33bacc61e39701cf9b46bd25bf5f0595bbe24655141438e7a100b")
	verifier, err := NewEd25519Verifier("rfc8032-test", []PublicKey{{KeyID: "vector-1", Bytes: publicKey}})
	if err != nil {
		t.Fatal(err)
	}
	reference := DetachedSignature{
		TrustDomain: "rfc8032-test", Algorithm: "ed25519", KeyID: "vector-1",
		Value: base64.RawURLEncoding.EncodeToString(signatureBytes),
	}
	if err := verifier.Verify(nil, reference); err != nil {
		t.Fatalf("valid RFC 8032 vector rejected: %v", err)
	}
	reference.Value = base64.RawURLEncoding.EncodeToString(signatureBytes[:63])
	if !errors.Is(verifier.Verify(nil, reference), ErrInvalidSignature) {
		t.Fatal("short signature did not fail closed")
	}
	if _, err := NewEd25519Verifier("rfc8032-test", []PublicKey{{KeyID: "short", Bytes: publicKey[:31]}}); !errors.Is(err, ErrInvalidReference) {
		t.Fatal("short public key did not fail before ed25519 verification")
	}
}

func TestDocumentBindingUsesExactBytes(t *testing.T) {
	raw := []byte(`{"blockedClasses":[]}`)
	binding := DocumentBinding{
		SchemaFamily: "synthetic-policy", SchemaVersion: "1.0", AdapterVersion: "policy-projection-1.0",
		Digest: DigestBytes(raw), Signature: testSignature("project-policy"),
	}
	if err := VerifyDocument(raw, binding, acceptingVerifier{}); err != nil {
		t.Fatal(err)
	}
	if !errors.Is(VerifyDocument(append(raw, '\n'), binding, acceptingVerifier{}), ErrBindingMismatch) {
		t.Fatal("formatting change did not invalidate the exact-byte binding")
	}
}

func TestFamilyWindowRejectsUnknownAndInvalidRetirement(t *testing.T) {
	now := time.Date(2026, 9, 2, 0, 0, 0, 0, time.UTC)
	window := FamilyWindow{Family: "synthetic-policy", CurrentMajor: 2, Minors: []MinorCompatibility{{Minor: 0}}}
	if _, err := window.Check("2.0", nil, now, acceptingVerifier{}); err != nil {
		t.Fatal(err)
	}
	if _, err := window.Check("3.0", nil, now, acceptingVerifier{}); !errors.Is(err, ErrUnsupportedSchema) {
		t.Fatal("unknown major was not rejected")
	}
	ambiguous := FamilyWindow{Family: "synthetic-policy", CurrentMajor: 2, Minors: []MinorCompatibility{{Minor: 0}, {Minor: 0, AllowUnknownFields: true}}}
	if _, err := ambiguous.Check("2.0", nil, now, acceptingVerifier{}); !errors.Is(err, ErrInvalidReference) {
		t.Fatal("ambiguous minor compatibility was admitted")
	}
	window.Retirements = []RetirementNotice{{
		SchemaFamily: "synthetic-policy", RetiredMajor: 1, Replacement: "2.0",
		AnnouncedAt: "2026-06-04T00:00:00Z", LastSupportedAt: "2026-09-02T00:00:00Z",
		SuccessfulCycles: 2, Signature: testSignature("schema-authority"),
	}}
	result, err := window.Check("1.0", nil, now, acceptingVerifier{})
	if err != nil || !result.RetiredMajor {
		t.Fatalf("valid 90-day retirement window rejected: result=%+v err=%v", result, err)
	}
	window.Retirements[0].SuccessfulCycles = 1
	if _, err := window.Check("1.0", nil, now, acceptingVerifier{}); !errors.Is(err, ErrUnsupportedSchema) {
		t.Fatal("retirement with fewer than two cycles was admitted")
	}
	window.Retirements = append(window.Retirements, window.Retirements[0])
	if _, err := window.Check("1.0", nil, now, acceptingVerifier{}); !errors.Is(err, ErrInvalidReference) {
		t.Fatal("conflicting retirement notices were admitted")
	}
}

func TestChainRejectsRollbackDivergenceAndDuplicateTarget(t *testing.T) {
	now := time.Date(2026, 9, 2, 0, 0, 0, 0, time.UTC)
	window := FamilyWindow{Family: "synthetic-revocations", CurrentMajor: 1, Minors: []MinorCompatibility{{Minor: 0}}}
	first := makeRecord(t, 1, "", "policy", digest("policy-a"))
	valid, err := VerifyChain([]ChainRecord{first}, Checkpoint{}, window, acceptingVerifier{}, now)
	if err != nil || valid.Head.Sequence != 1 {
		t.Fatalf("valid chain rejected: %+v %v", valid, err)
	}
	if _, err := VerifyChain([]ChainRecord{first}, valid.Head, window, acceptingVerifier{}, now); !errors.Is(err, ErrEvidenceRollback) {
		t.Fatalf("rollback was not rejected: %v", err)
	}
	second := makeRecord(t, 2, digest("wrong-head"), "allowlist", digest("allowlist-a"))
	if _, err := VerifyChain([]ChainRecord{second}, valid.Head, window, acceptingVerifier{}, now); !errors.Is(err, ErrEvidenceConflict) {
		t.Fatalf("divergence was not rejected: %v", err)
	}
	duplicate := makeRecord(t, 2, first.RecordDigest, "policy", digest("policy-a"))
	if _, err := VerifyChain([]ChainRecord{first, duplicate}, Checkpoint{}, window, acceptingVerifier{}, now); !errors.Is(err, ErrEvidenceConflict) {
		t.Fatalf("duplicate target was not rejected: %v", err)
	}
	duplicateID := makeRecord(t, 2, first.RecordDigest, "allowlist", digest("allowlist-b"))
	duplicateID.RecordID = first.RecordID
	duplicateID.RecordDigest = ""
	duplicateID = finalizeRecord(t, duplicateID)
	if _, err := VerifyChain([]ChainRecord{first, duplicateID}, Checkpoint{}, window, acceptingVerifier{}, now); !errors.Is(err, ErrEvidenceConflict) {
		t.Fatalf("duplicate record identity was not rejected: %v", err)
	}
}

func TestReceiptReferenceRequiresEveryBindingAndHonoursRevocation(t *testing.T) {
	now := time.Date(2026, 9, 2, 0, 0, 0, 0, time.UTC)
	bindings := testBindings()
	receipt := ReceiptReference{
		ReceiptID: "00000000-0000-4000-8000-000000000001", SchemaFamily: "synthetic-receipt", SchemaVersion: "1.0",
		IssuerID: "synthetic-issuer", ScanID: "00000000-0000-4000-8000-000000000002", Mode: "release",
		PassReason: "PASS_NO_BLOCKING_FINDINGS", CreatedAt: "2026-08-31T00:00:00Z", PromotionDeadline: "2026-09-30T00:00:00Z",
		Sequence: 1, Bindings: bindings, Signature: testSignature("receipt-authority"),
	}
	message, err := ReceiptMessage(receipt)
	if err != nil {
		t.Fatal(err)
	}
	receipt.RecordDigest = DigestBytes(message)
	request := ReferenceRequest{
		Receipt: receipt, ExpectedBindings: bindings, Now: now,
		ReceiptWindow: window("synthetic-receipt"), GlobalWindow: window("synthetic-global-revocations"), ProjectWindow: window("synthetic-project-revocations"),
		ReceiptHead: Checkpoint{Sequence: receipt.Sequence, Digest: receipt.RecordDigest},
	}
	trust := TrustSet{Receipt: acceptingVerifier{}, Global: acceptingVerifier{}, Project: acceptingVerifier{}}
	decision, err := VerifyReference(request, trust)
	if err != nil || !decision.ReferenceValid {
		t.Fatalf("valid reference rejected: %+v %v", decision, err)
	}
	rolledBack := request
	rolledBack.ReceiptHead = Checkpoint{}
	if _, err := VerifyReference(rolledBack, trust); !errors.Is(err, ErrEvidenceRollback) {
		t.Fatalf("receipt head rollback was not rejected: %v", err)
	}
	truncated := request
	truncated.ExpectedGlobalHead = Checkpoint{Sequence: 1, Digest: digest("missing-global-head")}
	if _, err := VerifyReference(truncated, trust); !errors.Is(err, ErrEvidenceRollback) {
		t.Fatalf("truncated global chain was not rejected: %v", err)
	}
	changed := request
	changed.ExpectedBindings.OutcomeDigest = digest("different-outcome")
	if _, err := VerifyReference(changed, trust); !errors.Is(err, ErrBindingMismatch) {
		t.Fatalf("changed binding was not rejected: %v", err)
	}
	future := request
	future.Receipt.CreatedAt = "2026-09-03T00:00:00Z"
	future.Receipt.PromotionDeadline = "2026-09-30T00:00:00Z"
	futureMessage, err := ReceiptMessage(future.Receipt)
	if err != nil {
		t.Fatal(err)
	}
	future.Receipt.RecordDigest = DigestBytes(futureMessage)
	future.ReceiptHead = Checkpoint{Sequence: future.Receipt.Sequence, Digest: future.Receipt.RecordDigest}
	if _, err := VerifyReference(future, trust); !errors.Is(err, ErrInvalidReference) {
		t.Fatalf("future receipt was not rejected: %v", err)
	}
	revocation := makeRecord(t, 1, "", "receipt", receipt.RecordDigest)
	revocation.SchemaFamily = "synthetic-global-revocations"
	revocation.RecordDigest = ""
	revocation = finalizeRecord(t, revocation)
	request.GlobalRecords = []ChainRecord{revocation}
	request.ExpectedGlobalHead = Checkpoint{Sequence: revocation.Sequence, Digest: revocation.RecordDigest}
	if _, err := VerifyReference(request, trust); !errors.Is(err, ErrRevoked) {
		t.Fatalf("revocation did not win: %v", err)
	}
	exceptionSet, _ := ExceptionSetDigest([]string{"synthetic-exception"})
	exceptionRequest := request
	exceptionRequest.Receipt.Bindings.ExceptionSetDigest = exceptionSet
	exceptionRequest.ExpectedBindings.ExceptionSetDigest = exceptionSet
	exceptionRequest.AppliedExceptions = []string{"synthetic-exception"}
	exceptionMessage, err := ReceiptMessage(exceptionRequest.Receipt)
	if err != nil {
		t.Fatal(err)
	}
	exceptionRequest.Receipt.RecordDigest = DigestBytes(exceptionMessage)
	exceptionRequest.ReceiptHead = Checkpoint{Sequence: 1, Digest: exceptionRequest.Receipt.RecordDigest}
	exceptionRevocation := makeRecord(t, 1, "", "exception", "synthetic-exception")
	exceptionRevocation.SchemaFamily = "synthetic-project-revocations"
	exceptionRevocation.RecordDigest = ""
	exceptionRevocation = finalizeRecord(t, exceptionRevocation)
	exceptionRequest.GlobalRecords = nil
	exceptionRequest.ExpectedGlobalHead = Checkpoint{}
	exceptionRequest.ProjectRecords = []ChainRecord{exceptionRevocation}
	exceptionRequest.ExpectedProjectHead = Checkpoint{Sequence: 1, Digest: exceptionRevocation.RecordDigest}
	if _, err := VerifyReference(exceptionRequest, trust); !errors.Is(err, ErrRevoked) {
		t.Fatalf("applied exception revocation did not win: %v", err)
	}
}

func testSignature(domain string) DetachedSignature {
	return DetachedSignature{TrustDomain: domain, Algorithm: "ed25519", KeyID: "synthetic-key", Value: "unused-by-test-verifier"}
}

func window(family string) FamilyWindow {
	return FamilyWindow{Family: family, CurrentMajor: 1, Minors: []MinorCompatibility{{Minor: 0}}}
}

func makeRecord(t *testing.T, sequence int64, previous, targetType, targetValue string) ChainRecord {
	t.Helper()
	return finalizeRecord(t, ChainRecord{
		SchemaFamily: "synthetic-revocations", SchemaVersion: "1.0", RecordID: recordID(sequence), Sequence: sequence, PreviousDigest: previous,
		Kind: "revocation", TargetType: targetType, TargetValue: targetValue, AuthorityID: "synthetic-authority",
		ReasonCode: "synthetic-integrity-defect",
		IssuedAt:   "2026-09-01T00:00:00Z", EffectiveAt: "2026-09-01T00:00:00Z", Signature: testSignature("revocation-authority"),
	})
}

func recordID(sequence int64) string {
	if sequence == 1 {
		return "00000000-0000-4000-8000-000000000101"
	}
	return "00000000-0000-4000-8000-000000000102"
}

func finalizeRecord(t *testing.T, record ChainRecord) ChainRecord {
	t.Helper()
	message, err := ChainRecordMessage(record)
	if err != nil {
		t.Fatal(err)
	}
	record.RecordDigest = DigestBytes(message)
	return record
}

func digest(seed string) string { return DigestBytes([]byte(seed)) }

func testBindings() ReceiptBindings {
	d := func(name string) string { return digest(name) }
	emptyExceptions, _ := ExceptionSetDigest(nil)
	return ReceiptBindings{
		ReleaseCommit: "0123456789abcdef0123456789abcdef01234567", HistoryRangeDigest: d("history"), TrackedTreeDigest: d("tree"),
		BuildContextDigest: d("context"), ArtifactSetDigest: d("artifacts"), BuildProvenanceDigest: d("provenance"),
		ScannerReleaseDigest: d("scanner"), RunnerDigest: d("runner"), EngineDigest: d("engine"), RulePackDigest: d("rules"),
		PolicyDigest: d("policy"), AllowlistDigest: d("allowlist"), ExceptionSetDigest: emptyExceptions, RequestSchemaVersion: "1.1", OutcomeSchemaVersion: "1.0",
		AdmissionLedgerSchemaVersion: "1.0", AdmissionLedgerDigest: d("ledger"), RawClassifierVersion: "classifier-1.0",
		RawClassifierDigest: d("classifier"), PreparationVersion: "preparation-1.0", PreparationDigest: d("preparation"),
		ResourceProfileID: "release-bounded-1", InspectionProofFormat: "inspection-proof-1.0", InspectionProofDigest: d("inspection"),
		OutcomeDigest: d("outcome"),
	}
}
