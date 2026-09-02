// Package policy projects transient project-owned policy and allowlist inputs
// into a minimal scanner evaluation model.
package policy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/verify"
)

const (
	PolicyAdapterVersion    = "policy-projection-1.0"
	AllowlistAdapterVersion = "allowlist-projection-1.0"
)

var ErrInvalidProjection = errors.New("project projection is invalid")

type Severity string

const (
	SeverityCritical      Severity = "critical"
	SeverityHigh          Severity = "high"
	SeverityMedium        Severity = "medium"
	SeverityLow           Severity = "low"
	SeverityInformational Severity = "informational"
)

type policyDocument struct {
	RequiredFeatures         []string `json:"requiredFeatures,omitempty"`
	BlockedClasses           []string `json:"blockedClasses"`
	MinimumSeverity          Severity `json:"minimumSeverity"`
	MaxExceptionDays         int      `json:"maxExceptionDays"`
	AllowedExceptionEvidence []string `json:"allowedExceptionEvidence"`
}

type Exception struct {
	ExceptionID          string   `json:"exceptionId"`
	Owner                string   `json:"owner"`
	Approver             string   `json:"approver"`
	Rationale            string   `json:"rationale"`
	RuleID               string   `json:"ruleId"`
	Class                string   `json:"class"`
	ScopeDigest          string   `json:"scopeDigest"`
	EvidenceType         string   `json:"evidenceType"`
	EvidenceDigest       string   `json:"evidenceDigest"`
	CreatedAt            string   `json:"createdAt"`
	ExpiresAt            string   `json:"expiresAt"`
	ScannerReleaseDigest string   `json:"scannerReleaseDigest"`
	RulePackDigest       string   `json:"rulePackDigest"`
	PolicyDigest         string   `json:"policyDigest"`
	SourceDigest         string   `json:"sourceDigest"`
	InvalidatesOn        []string `json:"invalidatesOn"`
}

type allowlistDocument struct {
	RequiredFeatures []string    `json:"requiredFeatures,omitempty"`
	Exceptions       []Exception `json:"exceptions"`
}

type PolicyProjection struct {
	Digest           string
	SchemaFamily     string
	SchemaVersion    string
	blockedClasses   map[string]bool
	minimumSeverity  Severity
	maxExceptionDays int
	allowedEvidence  map[string]bool
}

type AllowlistProjection struct {
	Digest        string
	SchemaFamily  string
	SchemaVersion string
	exceptions    map[string]Exception
}

type ExceptionContext struct {
	ScannerReleaseDigest string
	RulePackDigest       string
	PolicyDigest         string
	SourceDigest         string
}

func LoadPolicy(raw []byte, binding verify.DocumentBinding, window verify.FamilyWindow, verifier verify.SignatureVerifier, now time.Time) (PolicyProjection, error) {
	if binding.SchemaFamily != window.Family || binding.AdapterVersion != PolicyAdapterVersion || binding.Signature.TrustDomain != "project-policy" {
		return PolicyProjection{}, ErrInvalidProjection
	}
	compatibility, err := window.Check(binding.SchemaVersion, nil, now, verifier)
	if err != nil || verify.VerifyDocument(raw, binding, verifier) != nil {
		return PolicyProjection{}, ErrInvalidProjection
	}
	var document policyDocument
	if err := decodeProjection(raw, &document, compatibility.AllowUnknownFields); err != nil || len(document.RequiredFeatures) != 0 {
		return PolicyProjection{}, ErrInvalidProjection
	}
	if !validSeverity(document.MinimumSeverity) || document.MaxExceptionDays < 1 || document.MaxExceptionDays > 30 {
		return PolicyProjection{}, ErrInvalidProjection
	}
	blocked := map[string]bool{}
	for _, class := range document.BlockedClasses {
		if !verify.IsToken(class) || blocked[class] {
			return PolicyProjection{}, ErrInvalidProjection
		}
		blocked[class] = true
	}
	allowed := map[string]bool{}
	for _, evidence := range document.AllowedExceptionEvidence {
		if !validEvidenceType(evidence) || allowed[evidence] {
			return PolicyProjection{}, ErrInvalidProjection
		}
		allowed[evidence] = true
	}
	return PolicyProjection{
		Digest: binding.Digest, SchemaFamily: binding.SchemaFamily, SchemaVersion: binding.SchemaVersion,
		blockedClasses: blocked, minimumSeverity: document.MinimumSeverity,
		maxExceptionDays: document.MaxExceptionDays, allowedEvidence: allowed,
	}, nil
}

func LoadAllowlist(raw []byte, binding verify.DocumentBinding, window verify.FamilyWindow, verifier verify.SignatureVerifier, now time.Time, context ExceptionContext) (AllowlistProjection, error) {
	if binding.SchemaFamily != window.Family || binding.AdapterVersion != AllowlistAdapterVersion || binding.Signature.TrustDomain != "project-allowlist" || !validExceptionContext(context) {
		return AllowlistProjection{}, ErrInvalidProjection
	}
	compatibility, err := window.Check(binding.SchemaVersion, nil, now, verifier)
	if err != nil || verify.VerifyDocument(raw, binding, verifier) != nil {
		return AllowlistProjection{}, ErrInvalidProjection
	}
	var document allowlistDocument
	if err := decodeProjection(raw, &document, compatibility.AllowUnknownFields); err != nil || len(document.RequiredFeatures) != 0 {
		return AllowlistProjection{}, ErrInvalidProjection
	}
	exceptions := map[string]Exception{}
	seenIDs := map[string]bool{}
	for _, exception := range document.Exceptions {
		if err := validateException(exception, now, context); err != nil {
			return AllowlistProjection{}, ErrInvalidProjection
		}
		if seenIDs[exception.ExceptionID] {
			return AllowlistProjection{}, ErrInvalidProjection
		}
		key := exceptionKey(exception.RuleID, exception.Class, exception.ScopeDigest)
		if _, duplicate := exceptions[key]; duplicate {
			return AllowlistProjection{}, ErrInvalidProjection
		}
		seenIDs[exception.ExceptionID] = true
		exceptions[key] = exception
	}
	return AllowlistProjection{Digest: binding.Digest, SchemaFamily: binding.SchemaFamily, SchemaVersion: binding.SchemaVersion, exceptions: exceptions}, nil
}

func validateException(exception Exception, now time.Time, context ExceptionContext) error {
	if !verify.IsToken(exception.ExceptionID) || !verify.IsToken(exception.Owner) || !verify.IsToken(exception.Approver) || exception.Owner == exception.Approver ||
		len(exception.Rationale) < 1 || len(exception.Rationale) > 256 || strings.TrimSpace(exception.Rationale) != exception.Rationale ||
		!verify.IsToken(exception.RuleID) || !verify.IsToken(exception.Class) || credentialClass(exception.Class) ||
		!verify.IsDigest(exception.ScopeDigest) || !validEvidenceType(exception.EvidenceType) || !verify.IsDigest(exception.EvidenceDigest) ||
		exception.ScannerReleaseDigest != context.ScannerReleaseDigest || exception.RulePackDigest != context.RulePackDigest ||
		exception.PolicyDigest != context.PolicyDigest || exception.SourceDigest != context.SourceDigest {
		return ErrInvalidProjection
	}
	created, createdErr := verify.ParseCanonicalTime(exception.CreatedAt)
	expires, expiresErr := verify.ParseCanonicalTime(exception.ExpiresAt)
	if createdErr != nil || expiresErr != nil || created.After(now) || !expires.After(now) || expires.Sub(created) > 30*24*time.Hour {
		return ErrInvalidProjection
	}
	required := map[string]bool{"scope": true, "policy": true, "rule": true, "scanner": true, "source": true, "evidence": true}
	if len(exception.InvalidatesOn) != len(required) {
		return ErrInvalidProjection
	}
	for _, value := range exception.InvalidatesOn {
		if !required[value] {
			return ErrInvalidProjection
		}
		delete(required, value)
	}
	return nil
}

func validExceptionContext(context ExceptionContext) bool {
	return verify.IsDigest(context.ScannerReleaseDigest) && verify.IsDigest(context.RulePackDigest) &&
		verify.IsDigest(context.PolicyDigest) && verify.IsDigest(context.SourceDigest)
}

func decodeProjection(raw []byte, target any, allowUnknown bool) error {
	if len(raw) == 0 || len(raw) > verify.MaxDocumentBytes || rejectDuplicateKeys(raw) != nil {
		return ErrInvalidProjection
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	if !allowUnknown {
		decoder.DisallowUnknownFields()
	}
	if err := decoder.Decode(target); err != nil || decoder.Decode(&struct{}{}) != io.EOF {
		return ErrInvalidProjection
	}
	return nil
}

func rejectDuplicateKeys(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	var walk func() error
	walk = func() error {
		token, err := decoder.Token()
		if err != nil {
			return err
		}
		delimiter, ok := token.(json.Delim)
		if !ok {
			return nil
		}
		switch delimiter {
		case '{':
			seen := map[string]bool{}
			for decoder.More() {
				keyToken, keyErr := decoder.Token()
				key, isString := keyToken.(string)
				if keyErr != nil || !isString || seen[key] {
					return ErrInvalidProjection
				}
				seen[key] = true
				if err := walk(); err != nil {
					return err
				}
			}
			_, err = decoder.Token()
			return err
		case '[':
			for decoder.More() {
				if err := walk(); err != nil {
					return err
				}
			}
			_, err = decoder.Token()
			return err
		default:
			return fmt.Errorf("unexpected delimiter")
		}
	}
	if err := walk(); err != nil {
		return err
	}
	if decoder.Decode(&struct{}{}) != io.EOF {
		return ErrInvalidProjection
	}
	return nil
}

func validSeverity(value Severity) bool {
	switch value {
	case SeverityCritical, SeverityHigh, SeverityMedium, SeverityLow, SeverityInformational:
		return true
	default:
		return false
	}
}

func validEvidenceType(value string) bool {
	return value == "synthetic" || value == "public" || value == "invalid-by-construction"
}

func credentialClass(value string) bool {
	switch value {
	case "private-key", "access-token", "password", "connection-credential", "signing-key", "live-authentication":
		return true
	default:
		return false
	}
}

func exceptionKey(ruleID, class, scopeDigest string) string {
	return ruleID + "\x00" + class + "\x00" + scopeDigest
}
