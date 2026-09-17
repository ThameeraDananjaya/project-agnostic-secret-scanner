package verify

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"time"
)

const (
	UnsignedBuildTag      = "release-tooling-v1.0.0-c2-linux-boundary"
	UnsignedBuildTagV24   = "release-tooling-v1.0.0-c2-linux-build-v2"
	UnsignedBuildWorkflow = ".github/workflows/release-build-unsigned.yml"
	ReleaseSignerWorkflow = ".github/workflows/release-sign.yml"
	ReleaseSignerRef      = "refs/tags/release-signing-v1.0.0-c2-linux-boundary"
	releaseRepository     = "ThameeraDananjaya/project-agnostic-secret-scanner"
	releaseIssuer         = "https://token.actions.githubusercontent.com"
)

// BuildIdentity is provenance. It is never a claim that a signature exists.
type BuildIdentity struct {
	Repository        string `json:"repository"`
	RepositoryOwnerID int64  `json:"repositoryOwnerId"`
	Workflow          string `json:"workflow"`
	Ref               string `json:"ref"`
	WorkflowSHA       string `json:"workflowSha"`
	Trigger           string `json:"trigger"`
}

// Preserve the explicit absence of a signer when structural inspection is
// serialized again. Nonzero conflicting signer fields are never normalized away.
func (m ReleaseManifest) MarshalJSON() ([]byte, error) {
	type plainManifest ReleaseManifest
	if (m.ManifestSchemaVersion == "2.3" || m.ManifestSchemaVersion == "2.4") && m.ReleaseState == "unsigned-candidate" && m.ReleaseIdentity == (ReleaseIdentity{}) {
		return json.Marshal(struct {
			plainManifest
			Signer any `json:"releaseIdentity"`
		}{plainManifest: plainManifest(m), Signer: nil})
	}
	return json.Marshal(plainManifest(m))
}

func validateUnsignedBuildIdentity(m ReleaseManifest) error {
	if m.ProductSource == nil || m.ReleaseTooling == nil || m.BuildIdentity == nil || m.SourceRevision != "" || m.SourceTree != "" || m.ReleaseVersion != "v1.0.0" {
		return ErrInvalidReference
	}
	p, t, b := m.ProductSource, m.ReleaseTooling, m.BuildIdentity
	buildTag := UnsignedBuildTag
	switch m.ManifestSchemaVersion {
	case "2.3":
	case "2.4":
		buildTag = UnsignedBuildTagV24
	default:
		return ErrInvalidReference
	}
	if *p != (ProductSourceIdentity{Tag: "v1.0.0", Commit: "a13c28fe7273bc8dc6545f97966a02889524eb4c", Tree: "217b711ddea51fd0ea7e808edd2e27fdecef8427"}) ||
		t.Tag != buildTag || !gitOIDPattern.MatchString(t.Commit) || !gitOIDPattern.MatchString(t.Tree) ||
		t.Workflow != UnsignedBuildWorkflow || t.WorkflowRef != "refs/tags/"+buildTag || t.WorkflowSHA != t.Commit || t.Trigger != "workflow_dispatch" ||
		b.Repository != releaseRepository || b.RepositoryOwnerID != 50274860 || b.Workflow != t.Workflow || b.Ref != t.WorkflowRef || b.WorkflowSHA != t.Commit || b.Trigger != t.Trigger {
		return ErrInvalidReference
	}
	switch m.ReleaseState {
	case "unsigned-candidate":
		if m.ReleaseIdentity != (ReleaseIdentity{}) {
			return ErrInvalidReference
		}
	case "signing-pending":
		if !validSeparateSigner(m.ReleaseIdentity) {
			return ErrInvalidReference
		}
	default:
		return ErrInvalidReference
	}
	return nil
}

func validSeparateSigner(s ReleaseIdentity) bool {
	return s.Repository == releaseRepository && s.RepositoryOwnerID == 50274860 && s.Workflow == ReleaseSignerWorkflow && s.Ref == ReleaseSignerRef &&
		gitOIDPattern.MatchString(s.WorkflowSHA) && s.WorkflowSHA != strings.Repeat("0", 40) && s.Trigger == "workflow_dispatch" && s.OIDCIssuer == releaseIssuer &&
		s.CertificateIdentity == "https://github.com/"+releaseRepository+"/"+ReleaseSignerWorkflow+"@"+ReleaseSignerRef
}

func checkSeparateSignerPolicy(m ReleaseManifest, p ReleaseTrustPolicy, now time.Time) error {
	// Even a supplied signature bundle cannot turn the unsigned state into trust.
	if m.ReleaseState != "signing-pending" || validateUnsignedBuildIdentity(m) != nil {
		return ErrInvalidReference
	}
	t, product, s := m.ReleaseTooling, m.ProductSource, m.ReleaseIdentity
	if p.ManifestSchemaVersion != m.ManifestSchemaVersion || p.ReleaseVersion != m.ReleaseVersion || p.ProductSourceTag != product.Tag || p.ProductSourceCommit != product.Commit || p.ProductSourceTree != product.Tree ||
		p.ReleaseToolingTag != t.Tag || p.ReleaseToolingCommit != t.Commit || p.ReleaseToolingTree != t.Tree ||
		p.Repository != s.Repository || p.RepositoryOwnerID != s.RepositoryOwnerID || p.Workflow != s.Workflow || p.Ref != s.Ref || p.WorkflowSHA != s.WorkflowSHA || p.Trigger != s.Trigger || p.OIDCIssuer != s.OIDCIssuer || p.CertificateIdentity != s.CertificateIdentity {
		return ErrBindingMismatch
	}
	created, _ := parseCanonicalTime(m.CreatedAt)
	if created.After(now) {
		return ErrInvalidReference
	}
	return nil
}

func exactObjectFields(raw []byte, names ...string) bool {
	var fields map[string]json.RawMessage
	if json.Unmarshal(raw, &fields) != nil || len(fields) != len(names) {
		return false
	}
	for _, name := range names {
		if _, ok := fields[name]; !ok {
			return false
		}
	}
	return true
}

func exactIdentityFields(raw []byte, signer bool) bool {
	names := []string{"repository", "repositoryOwnerId", "workflow", "ref", "workflowSha", "trigger"}
	if signer {
		names = append(names, "oidcIssuer", "certificateIdentity")
	}
	return exactObjectFields(raw, names...)
}

func exactBuildFields(raw []byte) bool {
	var fields map[string]json.RawMessage
	if json.Unmarshal(raw, &fields) != nil {
		return false
	}
	return exactObjectFields(fields["productSource"], "tag", "commit", "tree") &&
		exactObjectFields(fields["releaseTooling"], "tag", "commit", "tree", "workflow", "workflowRef", "workflowSha", "trigger") && exactIdentityFields(fields["buildIdentity"], false)
}

// LoadSeparateSignerPolicy admits independently obtained owner trust data, never
// a candidate-adjacent or discovered policy. The digest is a separate trust input.
func LoadSeparateSignerPolicy(path, expectedDigest, candidateDirectory, buildCommit, buildTree string) (ReleaseTrustPolicy, error) {
	var empty ReleaseTrustPolicy
	if !filepath.IsAbs(path) || !digestPattern.MatchString(expectedDigest) || expectedDigest == strings.Repeat("0", 64) || !gitOIDPattern.MatchString(buildCommit) || !gitOIDPattern.MatchString(buildTree) {
		return empty, ErrInvalidReference
	}
	root, err := filepath.Abs(candidateDirectory)
	if err != nil {
		return empty, ErrInvalidReference
	}
	rel, err := filepath.Rel(root, path)
	if err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return empty, ErrInvalidReference
	}
	// Resolve the parent too, so a link cannot disguise a candidate-local policy.
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil || filepath.Clean(resolved) != filepath.Clean(path) {
		return empty, ErrInvalidReference
	}
	raw, _, err := readRegularFileBounded(path, 8192)
	if err != nil || len(raw) == 0 || DigestBytes(raw) != expectedDigest || rejectDuplicateJSONMembers(raw) != nil {
		return empty, ErrInvalidReference
	}
	if !exactObjectFields(raw, "schema", "productSource", "releaseTooling", "buildIdentity", "signer") || !exactBuildFields(raw) {
		return empty, ErrInvalidReference
	}
	var fields map[string]json.RawMessage
	if json.Unmarshal(raw, &fields) != nil || !exactIdentityFields(fields["signer"], true) {
		return empty, ErrInvalidReference
	}
	var record struct {
		Schema  string                 `json:"schema"`
		Product ProductSourceIdentity  `json:"productSource"`
		Tooling ReleaseToolingIdentity `json:"releaseTooling"`
		Build   BuildIdentity          `json:"buildIdentity"`
		Signer  ReleaseIdentity        `json:"signer"`
	}
	d := json.NewDecoder(bytes.NewReader(raw))
	d.DisallowUnknownFields()
	if d.Decode(&record) != nil || ensureJSONEOF(d) != nil {
		return empty, ErrInvalidReference
	}
	version := ""
	switch record.Schema {
	case "pscan-separate-signer-policy-v1":
		version = "2.3"
	case "pscan-separate-signer-policy-v1.1":
		version = "2.4"
	default:
		return empty, ErrInvalidReference
	}
	m := ReleaseManifest{ManifestSchemaVersion: version, ReleaseVersion: "v1.0.0", ProductSource: &record.Product, ReleaseTooling: &record.Tooling, BuildIdentity: &record.Build, ReleaseState: "signing-pending", ReleaseIdentity: record.Signer}
	if validateUnsignedBuildIdentity(m) != nil || record.Tooling.Commit != buildCommit || record.Tooling.Tree != buildTree {
		return empty, ErrBindingMismatch
	}
	s := record.Signer
	return ReleaseTrustPolicy{Repository: s.Repository, RepositoryOwnerID: s.RepositoryOwnerID, Workflow: s.Workflow, Ref: s.Ref, WorkflowSHA: s.WorkflowSHA, Trigger: s.Trigger, OIDCIssuer: s.OIDCIssuer, CertificateIdentity: s.CertificateIdentity,
		ReleaseVersion: "v1.0.0", ManifestSchemaVersion: version, ProductSourceTag: record.Product.Tag, ProductSourceCommit: record.Product.Commit, ProductSourceTree: record.Product.Tree, ReleaseToolingTag: record.Tooling.Tag, ReleaseToolingCommit: buildCommit, ReleaseToolingTree: buildTree}, nil
}
