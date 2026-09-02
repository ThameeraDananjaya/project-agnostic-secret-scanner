package verify

import "time"

type ChainRecord struct {
	SchemaFamily   string            `json:"schemaFamily"`
	SchemaVersion  string            `json:"schemaVersion"`
	RecordID       string            `json:"recordId"`
	Sequence       int64             `json:"sequence"`
	PreviousDigest string            `json:"previousDigest,omitempty"`
	Kind           string            `json:"kind"`
	TargetType     string            `json:"targetType"`
	TargetValue    string            `json:"targetValue"`
	AuthorityID    string            `json:"authorityId"`
	ReasonCode     string            `json:"reasonCode"`
	IssuedAt       string            `json:"issuedAt"`
	EffectiveAt    string            `json:"effectiveAt"`
	RecordDigest   string            `json:"recordDigest"`
	Signature      DetachedSignature `json:"signature"`
}

type Checkpoint struct {
	Sequence int64
	Digest   string
}

type ChainResult struct {
	Head        Checkpoint
	Revocations map[string]time.Time
}

func RevocationKey(targetType, targetValue string) string { return targetType + "\x00" + targetValue }

func ChainRecordMessage(record ChainRecord) ([]byte, error) {
	if !IsToken(record.SchemaFamily) || !IsUUID(record.RecordID) || record.Sequence < 1 || record.Kind != "revocation" ||
		!validTargetType(record.TargetType) || !IsToken(record.TargetValue) || !IsToken(record.AuthorityID) ||
		!IsToken(record.ReasonCode) ||
		!IsToken(record.Signature.TrustDomain) || !IsToken(record.Signature.KeyID) {
		return nil, ErrInvalidReference
	}
	if _, _, err := parseVersion(record.SchemaVersion); err != nil {
		return nil, err
	}
	if record.Sequence == 1 {
		if record.PreviousDigest != "" {
			return nil, ErrEvidenceConflict
		}
	} else if !IsDigest(record.PreviousDigest) {
		return nil, ErrEvidenceConflict
	}
	issued, err := parseCanonicalTime(record.IssuedAt)
	if err != nil {
		return nil, err
	}
	effective, err := parseCanonicalTime(record.EffectiveAt)
	if err != nil {
		return nil, err
	}
	if effective.Before(issued) {
		return nil, ErrInvalidReference
	}
	p, _ := newPreimage("PSCAN-EVIDENCE-CHAIN-RECORD-V1")
	p.add(record.Signature.TrustDomain)
	p.add(record.SchemaFamily)
	p.add(record.SchemaVersion)
	p.add(record.RecordID)
	p.addInt(record.Sequence)
	p.add(record.PreviousDigest)
	p.add(record.Kind)
	p.add(record.TargetType)
	p.add(record.TargetValue)
	p.add(record.AuthorityID)
	p.add(record.ReasonCode)
	p.add(record.IssuedAt)
	p.add(record.EffectiveAt)
	p.add(record.Signature.KeyID)
	return p.bytes(), nil
}

func VerifyChain(records []ChainRecord, anchor Checkpoint, window FamilyWindow, verifier SignatureVerifier, now time.Time) (ChainResult, error) {
	if verifier == nil || anchor.Sequence < 0 || (anchor.Sequence == 0) != (anchor.Digest == "") || (anchor.Sequence > 0 && !IsDigest(anchor.Digest)) {
		return ChainResult{}, ErrInvalidReference
	}
	result := ChainResult{Head: anchor, Revocations: map[string]time.Time{}}
	if len(records) == 0 {
		return result, nil
	}
	expectedSequence := anchor.Sequence + 1
	expectedPrevious := anchor.Digest
	seenRecords := map[string]bool{}
	seenIDs := map[string]bool{}
	seenTargets := map[string]bool{}
	for _, record := range records {
		if record.SchemaFamily != window.Family || record.Sequence < expectedSequence {
			return ChainResult{}, ErrEvidenceRollback
		}
		if record.Sequence != expectedSequence || record.PreviousDigest != expectedPrevious {
			return ChainResult{}, ErrEvidenceConflict
		}
		if _, err := window.Check(record.SchemaVersion, nil, now, verifier); err != nil {
			return ChainResult{}, err
		}
		if err := validateScannerOwnedGlobalRecord(record); err != nil {
			return ChainResult{}, err
		}
		message, err := ChainRecordMessage(record)
		if err != nil {
			return ChainResult{}, err
		}
		digest := DigestBytes(message)
		if !IsDigest(record.RecordDigest) || record.RecordDigest != digest || seenRecords[digest] || seenIDs[record.RecordID] {
			return ChainResult{}, ErrEvidenceConflict
		}
		if err := verifier.Verify(message, record.Signature); err != nil {
			return ChainResult{}, ErrInvalidSignature
		}
		effective, _ := parseCanonicalTime(record.EffectiveAt)
		issued, _ := parseCanonicalTime(record.IssuedAt)
		if issued.After(now) {
			return ChainResult{}, ErrInvalidReference
		}
		key := RevocationKey(record.TargetType, record.TargetValue)
		if seenTargets[key] {
			return ChainResult{}, ErrEvidenceConflict
		}
		seenRecords[digest] = true
		seenIDs[record.RecordID] = true
		seenTargets[key] = true
		if !effective.After(now) {
			result.Revocations[key] = effective
		}
		expectedSequence++
		expectedPrevious = digest
		result.Head = Checkpoint{Sequence: record.Sequence, Digest: digest}
	}
	return result, nil
}

func validateScannerOwnedGlobalRecord(record ChainRecord) error {
	if record.SchemaFamily != "global-scanner-revocation" {
		return nil
	}
	major, minor, err := parseVersion(record.SchemaVersion)
	if err != nil || major != 1 || minor < 1 || record.Signature.TrustDomain != "global-scanner-revocation" ||
		!IsDigest(record.TargetValue) || !validGlobalTargetType(record.TargetType) || !validGlobalReasonCode(record.ReasonCode) {
		return ErrInvalidReference
	}
	return nil
}

func validGlobalTargetType(value string) bool {
	switch value {
	case "scanner-release", "scanner-asset", "rule-pack", "scanner-schema":
		return true
	default:
		return false
	}
}

func validGlobalReasonCode(value string) bool {
	switch value {
	case "COMPROMISE", "INTEGRITY_DEFECT", "INCOMPLETE_COVERAGE", "AUTHORITY_DEFECT", "RETIREMENT_BREACH":
		return true
	default:
		return false
	}
}

func validTargetType(value string) bool {
	switch value {
	case "scanner-release", "scanner-asset", "rule-pack", "scanner-schema", "receipt", "receipt-key", "policy", "allowlist", "exception":
		return true
	default:
		return false
	}
}
