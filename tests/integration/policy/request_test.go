package policy_test

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/outcome"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/request"
)

func TestRequest11CarriesExactPolicyAndProofBindings(t *testing.T) {
	now := time.Date(2026, 9, 2, 0, 0, 0, 0, time.UTC)
	value := validRequest11(t, now)
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	loaded, err := request.Load(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if err := request.ValidateAt(loaded, now); err != nil {
		t.Fatalf("valid 1.1 request rejected: %v", err)
	}

	for name, mutate := range map[string]func(*request.ScanRequest){
		"missing policy binding": func(v *request.ScanRequest) { v.PolicyBinding = nil },
		"policy digest mismatch": func(v *request.ScanRequest) { v.PolicyBinding.Digest = strings.Repeat("b", 64) },
		"wrong trust domain":     func(v *request.ScanRequest) { v.AllowlistBinding.Signature.TrustDomain = "project-policy" },
		"missing proof bindings": func(v *request.ScanRequest) { v.ProofBindings = nil },
		"ambiguous proof schema": func(v *request.ScanRequest) { v.ProofBindings.AdmissionLedgerSchemaVersion = "01.0" },
		"bad inspection digest":  func(v *request.ScanRequest) { v.ProofBindings.InspectionProofDigest = "not-a-digest" },
	} {
		t.Run(name, func(t *testing.T) {
			candidate := validRequest11(t, now)
			mutate(&candidate)
			if err := request.ValidateAt(candidate, now); request.ReasonFor(err, "unexpected") != outcome.ReasonFailBindingMismatch {
				t.Fatalf("invalid binding did not fail closed: %v", err)
			}
		})
	}
}

func TestRequest11RequiredFieldsCannotBeOmittedFromJSON(t *testing.T) {
	value := validRequest11(t, time.Date(2026, 9, 2, 0, 0, 0, 0, time.UTC))
	data, _ := json.Marshal(value)
	var object map[string]any
	_ = json.Unmarshal(data, &object)
	delete(object, "proofBindings")
	data, _ = json.Marshal(object)
	if _, err := request.Load(bytes.NewReader(data)); request.ReasonFor(err, "unexpected") != outcome.ReasonFailInputIntegrity {
		t.Fatalf("omitted proof bindings did not fail closed: %v", err)
	}
}

func validRequest11(t *testing.T, now time.Time) request.ScanRequest {
	t.Helper()
	digest := strings.Repeat("a", 64)
	signature := request.SignatureReference{TrustDomain: "project-policy", Algorithm: "ed25519", KeyID: "synthetic-key", Value: strings.Repeat("A", 86)}
	return request.ScanRequest{
		RequestSchemaVersion: "1.1", ScanID: "123e4567-e89b-42d3-a456-426614174000", Mode: "pr",
		ScannerReleaseDigest: digest,
		EngineBinding:        request.EngineBinding{Name: "gitleaks", Version: "8.30.1", BinaryDigest: digest, AdapterVersion: "gitleaks-adapter-1.0"},
		RulePackDigest:       digest, PolicyDigest: digest, AllowlistDigest: digest,
		PolicyBinding:    &request.ProjectDocumentBinding{SchemaFamily: "synthetic-policy", SchemaVersion: "1.0", AdapterVersion: "policy-projection-1.0", Digest: digest, Signature: signature},
		AllowlistBinding: &request.ProjectDocumentBinding{SchemaFamily: "synthetic-allowlist", SchemaVersion: "1.0", AdapterVersion: "allowlist-projection-1.0", Digest: digest, Signature: request.SignatureReference{TrustDomain: "project-allowlist", Algorithm: "ed25519", KeyID: "synthetic-key", Value: strings.Repeat("A", 86)}},
		ProofBindings: &request.ProofBindings{
			AdmissionLedgerSchemaVersion: "1.0", AdmissionLedgerDigest: digest, RawClassifierVersion: "classifier-1.0", RawClassifierDigest: digest,
			PreparationVersion: "preparation-1.0", PreparationDigest: digest, ResourceProfileID: "pr-bounded-1",
			InspectionProofFormat: "inspection-proof-1.0", InspectionProofDigest: digest,
		},
		SourceBinding:         request.SourceBinding{BaseCommit: digest, HeadCommit: digest, MergeBase: digest, HistoryRangeDigest: digest, TrackedTreeDigest: digest},
		TrackedSourceManifest: request.FileBinding{Path: filepath.Join(t.TempDir(), "synthetic-manifest.json"), Digest: digest},
		FallbackRequirement:   request.FallbackRequirement{Mode: "disabled"},
		Limits:                request.Limits{TimeoutSeconds: 900, MaxArchiveDepth: 5, MaxArchiveEntries: 100000, MaxExpandedBytes: 2 << 30, MaxFileBytes: 512 << 20, MaxCompressionRatio: 1000, MaxMemoryBytes: 1 << 30, MaxCPUPercent: 100},
		OfflineRequired:       true, RedactionMode: "full", RequestedAt: now.Format(time.RFC3339),
	}
}
