package policy

import (
	"time"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/outcome"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/verify"
)

type EvaluationContext struct {
	Now                  time.Time
	ScannerReleaseDigest string
	RulePackDigest       string
	PolicyDigest         string
	AllowlistDigest      string
	SourceDigest         string
	Revocations          map[string]bool
}

type Finding struct {
	RuleID      string
	Class       string
	Severity    Severity
	ScopeDigest string
	Credential  bool
}

type Decision struct {
	Blocking bool
	Reason   outcome.ReasonCode
}

func Evaluate(context EvaluationContext, project PolicyProjection, allowlist AllowlistProjection, findings []Finding) Decision {
	if !validContext(context) || project.Digest != context.PolicyDigest || allowlist.Digest != context.AllowlistDigest {
		return Decision{Blocking: true, Reason: outcome.ReasonFailBindingMismatch}
	}
	for _, target := range [][2]string{
		{"scanner-release", context.ScannerReleaseDigest},
		{"rule-pack", context.RulePackDigest},
		{"policy", context.PolicyDigest},
		{"allowlist", context.AllowlistDigest},
	} {
		if context.Revocations != nil && context.Revocations[verify.RevocationKey(target[0], target[1])] {
			return Decision{Blocking: true, Reason: outcome.ReasonFailRevokedDependency}
		}
	}
	for _, finding := range findings {
		if !verify.IsToken(finding.RuleID) || !verify.IsToken(finding.Class) || !verify.IsDigest(finding.ScopeDigest) || !validSeverity(finding.Severity) {
			return Decision{Blocking: true, Reason: outcome.ReasonFailInputIntegrity}
		}
		if finding.Credential || credentialClass(finding.Class) {
			return Decision{Blocking: true, Reason: outcome.ReasonFailFindingDetected}
		}
		if !project.blocks(finding) {
			continue
		}
		exception, found := allowlist.exceptions[exceptionKey(finding.RuleID, finding.Class, finding.ScopeDigest)]
		if !found {
			return Decision{Blocking: true, Reason: outcome.ReasonFailPolicyDenied}
		}
		if !project.allowedEvidence[exception.EvidenceType] {
			return Decision{Blocking: true, Reason: outcome.ReasonFailAllowlistInvalid}
		}
		created, createdErr := verify.ParseCanonicalTime(exception.CreatedAt)
		expires, expiresErr := verify.ParseCanonicalTime(exception.ExpiresAt)
		maximum := time.Duration(project.maxExceptionDays) * 24 * time.Hour
		if createdErr != nil || expiresErr != nil || created.After(context.Now) || !expires.After(context.Now) || expires.Sub(created) > maximum {
			return Decision{Blocking: true, Reason: outcome.ReasonFailAllowlistInvalid}
		}
	}
	return Decision{Reason: outcome.ReasonPassNoBlockingFindings}
}

func (p PolicyProjection) blocks(finding Finding) bool {
	return p.blockedClasses[finding.Class] || severityRank(finding.Severity) >= severityRank(p.minimumSeverity)
}

func validContext(context EvaluationContext) bool {
	return !context.Now.IsZero() && context.Revocations != nil && verify.IsDigest(context.ScannerReleaseDigest) && verify.IsDigest(context.RulePackDigest) &&
		verify.IsDigest(context.PolicyDigest) && verify.IsDigest(context.AllowlistDigest) && verify.IsDigest(context.SourceDigest)
}

func severityRank(value Severity) int {
	switch value {
	case SeverityCritical:
		return 5
	case SeverityHigh:
		return 4
	case SeverityMedium:
		return 3
	case SeverityLow:
		return 2
	case SeverityInformational:
		return 1
	default:
		return 0
	}
}
