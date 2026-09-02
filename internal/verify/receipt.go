package verify

import "time"

const MaximumReceiptReuse = 30 * 24 * time.Hour

type ReceiptBindings struct {
	ReleaseCommit                string
	HistoryRangeDigest           string
	TrackedTreeDigest            string
	BuildContextDigest           string
	ArtifactSetDigest            string
	BuildProvenanceDigest        string
	ScannerReleaseDigest         string
	RunnerDigest                 string
	EngineDigest                 string
	RulePackDigest               string
	PolicyDigest                 string
	AllowlistDigest              string
	PolicyBindingDigest          string
	AllowlistBindingDigest       string
	ExceptionSetDigest           string
	RequestSchemaVersion         string
	RequestSchemaDigest          string
	OutcomeSchemaVersion         string
	OutcomeSchemaDigest          string
	AdmissionLedgerSchemaVersion string
	AdmissionLedgerDigest        string
	RawClassifierVersion         string
	RawClassifierDigest          string
	PreparationVersion           string
	PreparationDigest            string
	ResourceProfileID            string
	InspectionProofFormat        string
	InspectionProofDigest        string
	OutcomeDigest                string
	OutcomeTimestamp             string
}

type ReceiptReference struct {
	ReceiptID         string            `json:"receiptId"`
	SchemaFamily      string            `json:"schemaFamily"`
	SchemaVersion     string            `json:"schemaVersion"`
	IssuerID          string            `json:"issuerId"`
	ScanID            string            `json:"scanId"`
	Mode              string            `json:"mode"`
	PassReason        string            `json:"passReason"`
	CreatedAt         string            `json:"createdAt"`
	PromotionDeadline string            `json:"promotionDeadline"`
	Sequence          int64             `json:"sequence"`
	PreviousDigest    string            `json:"previousDigest,omitempty"`
	Bindings          ReceiptBindings   `json:"bindings"`
	RecordDigest      string            `json:"recordDigest"`
	Signature         DetachedSignature `json:"signature"`
}

type ReferenceRequest struct {
	Receipt             ReceiptReference
	ExpectedBindings    ReceiptBindings
	Now                 time.Time
	ReceiptWindow       FamilyWindow
	GlobalWindow        FamilyWindow
	ProjectWindow       FamilyWindow
	GlobalAnchor        Checkpoint
	ProjectAnchor       Checkpoint
	ExpectedGlobalHead  Checkpoint
	ExpectedProjectHead Checkpoint
	GlobalRecords       []ChainRecord
	ProjectRecords      []ChainRecord
	ReceiptAnchor       Checkpoint
	ReceiptHead         Checkpoint
	AppliedExceptions   []string
}

type TrustSet struct {
	Receipt SignatureVerifier
	Global  SignatureVerifier
	Project SignatureVerifier
}

type ReferenceDecision struct {
	ReferenceValid bool
	ReceiptHead    Checkpoint
	GlobalHead     Checkpoint
	ProjectHead    Checkpoint
}

func ReceiptMessage(receipt ReceiptReference) ([]byte, error) {
	if !IsUUID(receipt.ReceiptID) || !IsUUID(receipt.ScanID) || !IsToken(receipt.SchemaFamily) ||
		!IsToken(receipt.IssuerID) || receipt.Mode != "release" || receipt.PassReason != "PASS_NO_BLOCKING_FINDINGS" ||
		receipt.Sequence < 1 || !IsToken(receipt.Signature.TrustDomain) || !IsToken(receipt.Signature.KeyID) {
		return nil, ErrInvalidReference
	}
	if _, _, err := parseVersion(receipt.SchemaVersion); err != nil {
		return nil, err
	}
	if receipt.Sequence == 1 {
		if receipt.PreviousDigest != "" {
			return nil, ErrEvidenceConflict
		}
	} else if !IsDigest(receipt.PreviousDigest) {
		return nil, ErrEvidenceConflict
	}
	bindings, err := receipt.Bindings.message()
	if err != nil {
		return nil, err
	}
	p, _ := newPreimage("PSCAN-PROJECT-RECEIPT-REFERENCE-V1")
	p.add(receipt.Signature.TrustDomain)
	p.add(receipt.ReceiptID)
	p.add(receipt.SchemaFamily)
	p.add(receipt.SchemaVersion)
	p.add(receipt.IssuerID)
	p.add(receipt.Signature.KeyID)
	p.add(receipt.ScanID)
	p.add(receipt.Mode)
	p.add(receipt.PassReason)
	p.add(receipt.CreatedAt)
	p.add(receipt.PromotionDeadline)
	p.addInt(receipt.Sequence)
	p.add(receipt.PreviousDigest)
	p.add(DigestBytes(bindings))
	return p.bytes(), nil
}

func VerifyReference(request ReferenceRequest, trust TrustSet) (ReferenceDecision, error) {
	if trust.Receipt == nil || trust.Global == nil || trust.Project == nil || request.Now.IsZero() {
		return ReferenceDecision{}, ErrInvalidReference
	}
	if request.GlobalWindow.Family != "global-scanner-revocation" {
		return ReferenceDecision{}, ErrInvalidReference
	}
	if request.GlobalAnchor != (Checkpoint{}) || request.ProjectAnchor != (Checkpoint{}) {
		return ReferenceDecision{}, ErrInvalidReference
	}
	receipt := request.Receipt
	if receipt.SchemaFamily != request.ReceiptWindow.Family {
		return ReferenceDecision{}, ErrUnsupportedSchema
	}
	if _, err := request.ReceiptWindow.Check(receipt.SchemaVersion, nil, request.Now, trust.Receipt); err != nil {
		return ReferenceDecision{}, err
	}
	created, createdErr := parseCanonicalTime(receipt.CreatedAt)
	deadline, deadlineErr := parseCanonicalTime(receipt.PromotionDeadline)
	outcomeAt, outcomeErr := parseCanonicalTime(receipt.Bindings.OutcomeTimestamp)
	if createdErr != nil || deadlineErr != nil || outcomeErr != nil || outcomeAt.After(created) || created.After(request.Now) || deadline.Before(created) || deadline.Sub(created) > MaximumReceiptReuse || request.Now.After(deadline) {
		return ReferenceDecision{}, ErrInvalidReference
	}
	if receipt.Bindings != request.ExpectedBindings {
		return ReferenceDecision{}, ErrBindingMismatch
	}
	exceptionSetDigest, err := ExceptionSetDigest(request.AppliedExceptions)
	if err != nil || exceptionSetDigest != receipt.Bindings.ExceptionSetDigest {
		return ReferenceDecision{}, ErrBindingMismatch
	}
	message, err := ReceiptMessage(receipt)
	if err != nil {
		return ReferenceDecision{}, err
	}
	digest := DigestBytes(message)
	if !IsDigest(receipt.RecordDigest) || receipt.RecordDigest != digest {
		return ReferenceDecision{}, ErrBindingMismatch
	}
	if err := trust.Receipt.Verify(message, receipt.Signature); err != nil {
		return ReferenceDecision{}, ErrInvalidSignature
	}
	if request.ReceiptAnchor.Sequence < 0 || (request.ReceiptAnchor.Sequence == 0) != (request.ReceiptAnchor.Digest == "") ||
		(request.ReceiptAnchor.Sequence > 0 && !IsDigest(request.ReceiptAnchor.Digest)) {
		return ReferenceDecision{}, ErrInvalidReference
	}
	if receipt.Sequence < request.ReceiptAnchor.Sequence+1 {
		return ReferenceDecision{}, ErrEvidenceRollback
	}
	if receipt.Sequence != request.ReceiptAnchor.Sequence+1 || receipt.PreviousDigest != request.ReceiptAnchor.Digest {
		return ReferenceDecision{}, ErrEvidenceConflict
	}
	if request.ReceiptHead.Sequence < 1 || request.ReceiptHead.Sequence != receipt.Sequence || request.ReceiptHead.Digest != digest {
		return ReferenceDecision{}, ErrEvidenceRollback
	}
	global, err := VerifyChain(request.GlobalRecords, request.GlobalAnchor, request.GlobalWindow, trust.Global, request.Now)
	if err != nil {
		return ReferenceDecision{}, err
	}
	project, err := VerifyChain(request.ProjectRecords, request.ProjectAnchor, request.ProjectWindow, trust.Project, request.Now)
	if err != nil {
		return ReferenceDecision{}, err
	}
	if global.Head != request.ExpectedGlobalHead || project.Head != request.ExpectedProjectHead {
		return ReferenceDecision{}, ErrEvidenceRollback
	}
	if receiptRevoked(receipt, digest, request.AppliedExceptions, global.Revocations) || receiptRevoked(receipt, digest, request.AppliedExceptions, project.Revocations) {
		return ReferenceDecision{}, ErrRevoked
	}
	return ReferenceDecision{ReferenceValid: true, ReceiptHead: request.ReceiptHead, GlobalHead: global.Head, ProjectHead: project.Head}, nil
}

func (b ReceiptBindings) message() ([]byte, error) {
	digests := []string{
		b.HistoryRangeDigest, b.TrackedTreeDigest, b.BuildContextDigest, b.ArtifactSetDigest,
		b.BuildProvenanceDigest, b.ScannerReleaseDigest, b.RunnerDigest, b.EngineDigest,
		b.RulePackDigest, b.PolicyDigest, b.AllowlistDigest, b.PolicyBindingDigest, b.AllowlistBindingDigest,
		b.ExceptionSetDigest, b.RequestSchemaDigest, b.OutcomeSchemaDigest, b.AdmissionLedgerDigest,
		b.RawClassifierDigest, b.PreparationDigest, b.InspectionProofDigest, b.OutcomeDigest,
	}
	for _, digest := range digests {
		if !IsDigest(digest) {
			return nil, ErrInvalidReference
		}
	}
	if len(b.ReleaseCommit) != 40 && len(b.ReleaseCommit) != 64 || !isLowerHex(b.ReleaseCommit) ||
		!validSchemaVersion(b.RequestSchemaVersion) || !validSchemaVersion(b.OutcomeSchemaVersion) ||
		!validSchemaVersion(b.AdmissionLedgerSchemaVersion) || !validVersionOrToken(b.RawClassifierVersion) ||
		!validVersionOrToken(b.PreparationVersion) || !IsToken(b.ResourceProfileID) || !IsToken(b.InspectionProofFormat) {
		return nil, ErrInvalidReference
	}
	p, _ := newPreimage("PSCAN-RECEIPT-BINDINGS-V1")
	p.add(b.ReleaseCommit)
	for _, value := range []string{
		b.HistoryRangeDigest, b.TrackedTreeDigest, b.BuildContextDigest, b.ArtifactSetDigest,
		b.BuildProvenanceDigest, b.ScannerReleaseDigest, b.RunnerDigest, b.EngineDigest,
		b.RulePackDigest, b.PolicyDigest, b.AllowlistDigest, b.PolicyBindingDigest,
		b.AllowlistBindingDigest, b.ExceptionSetDigest, b.RequestSchemaVersion, b.RequestSchemaDigest,
		b.OutcomeSchemaVersion, b.OutcomeSchemaDigest, b.AdmissionLedgerSchemaVersion, b.AdmissionLedgerDigest,
		b.RawClassifierVersion, b.RawClassifierDigest, b.PreparationVersion,
		b.PreparationDigest, b.ResourceProfileID, b.InspectionProofFormat,
		b.InspectionProofDigest, b.OutcomeDigest, b.OutcomeTimestamp,
	} {
		p.add(value)
	}
	return p.bytes(), nil
}

func ExceptionSetDigest(exceptionIDs []string) (string, error) {
	p, _ := newPreimage("PSCAN-APPLIED-EXCEPTION-SET-V1")
	p.addInt(int64(len(exceptionIDs)))
	previous := ""
	for _, exceptionID := range exceptionIDs {
		if !IsToken(exceptionID) || (previous != "" && exceptionID <= previous) {
			return "", ErrInvalidReference
		}
		p.add(exceptionID)
		previous = exceptionID
	}
	return DigestBytes(p.bytes()), nil
}

func receiptRevoked(receipt ReceiptReference, receiptDigest string, exceptionIDs []string, revoked map[string]time.Time) bool {
	targets := [][2]string{
		{"receipt", receiptDigest},
		{"receipt", receipt.ReceiptID},
		{"receipt-key", receipt.Signature.KeyID},
		{"scanner-release", receipt.Bindings.ScannerReleaseDigest},
		{"scanner-asset", receipt.Bindings.RunnerDigest},
		{"scanner-asset", receipt.Bindings.EngineDigest},
		{"rule-pack", receipt.Bindings.RulePackDigest},
		{"policy", receipt.Bindings.PolicyDigest},
		{"allowlist", receipt.Bindings.AllowlistDigest},
	}
	for _, exceptionID := range exceptionIDs {
		targets = append(targets, [2]string{"exception", exceptionID})
	}
	for _, target := range targets {
		if _, found := revoked[RevocationKey(target[0], target[1])]; found {
			return true
		}
	}
	return false
}

func isLowerHex(value string) bool {
	for _, c := range value {
		if c < '0' || c > '9' {
			if c < 'a' || c > 'f' {
				return false
			}
		}
	}
	return true
}

func validVersionOrToken(value string) bool {
	if _, _, err := parseVersion(value); err == nil {
		return true
	}
	return IsToken(value)
}

func validSchemaVersion(value string) bool {
	_, _, err := parseVersion(value)
	return err == nil
}
