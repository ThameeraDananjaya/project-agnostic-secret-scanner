package verify

import "time"

type MinorCompatibility struct {
	Minor              int
	AllowUnknownFields bool
}

type RetirementNotice struct {
	SchemaFamily     string            `json:"schemaFamily"`
	RetiredMajor     int               `json:"retiredMajor"`
	Replacement      string            `json:"replacement"`
	AnnouncedAt      string            `json:"announcedAt"`
	LastSupportedAt  string            `json:"lastSupportedAt"`
	SuccessfulCycles int               `json:"successfulCycles"`
	Signature        DetachedSignature `json:"signature"`
}

type FamilyWindow struct {
	Family       string
	CurrentMajor int
	Minors       []MinorCompatibility
	Retirements  []RetirementNotice
}

type CompatibilityResult struct {
	AllowUnknownFields bool
	RetiredMajor       bool
}

func (w FamilyWindow) Check(version string, requiredFeatures []string, now time.Time, verifier SignatureVerifier) (CompatibilityResult, error) {
	if !IsToken(w.Family) || w.CurrentMajor < 1 || len(w.Minors) == 0 || len(requiredFeatures) != 0 {
		return CompatibilityResult{}, ErrUnsupportedSchema
	}
	seenMinors := map[int]bool{}
	for _, supported := range w.Minors {
		if supported.Minor < 0 || seenMinors[supported.Minor] {
			return CompatibilityResult{}, ErrInvalidReference
		}
		seenMinors[supported.Minor] = true
	}
	major, minor, err := parseVersion(version)
	if err != nil {
		return CompatibilityResult{}, err
	}
	if major == w.CurrentMajor {
		for _, supported := range w.Minors {
			if minor == supported.Minor {
				return CompatibilityResult{AllowUnknownFields: supported.AllowUnknownFields}, nil
			}
		}
		return CompatibilityResult{}, ErrUnsupportedSchema
	}
	var selected *RetirementNotice
	for index := range w.Retirements {
		notice := &w.Retirements[index]
		if notice.RetiredMajor != major || notice.SchemaFamily != w.Family {
			continue
		}
		if selected != nil {
			return CompatibilityResult{}, ErrInvalidReference
		}
		selected = notice
	}
	if selected != nil {
		notice := *selected
		if verifier == nil || !IsToken(notice.Replacement) || notice.SuccessfulCycles < 2 {
			return CompatibilityResult{}, ErrUnsupportedSchema
		}
		replacementMajor, _, err := parseVersion(notice.Replacement)
		if err != nil || replacementMajor != w.CurrentMajor {
			return CompatibilityResult{}, ErrUnsupportedSchema
		}
		announced, announcedErr := parseCanonicalTime(notice.AnnouncedAt)
		last, lastErr := parseCanonicalTime(notice.LastSupportedAt)
		if announcedErr != nil || lastErr != nil || announced.After(now) || !last.After(announced) ||
			last.Sub(announced) < 90*24*time.Hour || last.Sub(announced) > 180*24*time.Hour || now.After(last) {
			return CompatibilityResult{}, ErrUnsupportedSchema
		}
		message, messageErr := retirementMessage(notice)
		if messageErr != nil || verifier.Verify(message, notice.Signature) != nil {
			return CompatibilityResult{}, ErrInvalidSignature
		}
		return CompatibilityResult{RetiredMajor: true}, nil
	}
	return CompatibilityResult{}, ErrUnsupportedSchema
}

func retirementMessage(notice RetirementNotice) ([]byte, error) {
	if !IsToken(notice.SchemaFamily) || !IsToken(notice.Replacement) || notice.RetiredMajor < 1 || notice.SuccessfulCycles < 0 ||
		!IsToken(notice.Signature.TrustDomain) || !IsToken(notice.Signature.KeyID) {
		return nil, ErrInvalidReference
	}
	p, _ := newPreimage("PSCAN-SCHEMA-RETIREMENT-V1")
	p.add(notice.Signature.TrustDomain)
	p.add(notice.SchemaFamily)
	p.addInt(int64(notice.RetiredMajor))
	p.add(notice.Replacement)
	p.add(notice.AnnouncedAt)
	p.add(notice.LastSupportedAt)
	p.addInt(int64(notice.SuccessfulCycles))
	p.add(notice.Signature.KeyID)
	return p.bytes(), nil
}
