package verify

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	MaximumReleaseManifestBytes = 4 << 20
	MaximumReleaseAssets        = 256
)

var (
	gitOIDPattern           = regexp.MustCompile(`^[0-9a-f]{40}$`)
	releaseVersionPattern   = regexp.MustCompile(`^v[0-9]+\.[0-9]+\.[0-9]+$`)
	runnerVersionPattern    = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)
	toolchainVersionPattern = regexp.MustCompile(`^go[0-9]+\.[0-9]+\.[0-9]+$`)
	assetPathPattern        = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._+/-]{0,255}$`)
)

type ReleaseAsset struct {
	Path   string `json:"path"`
	Kind   string `json:"kind"`
	OS     string `json:"os"`
	Arch   string `json:"arch"`
	Size   int64  `json:"size"`
	SHA256 string `json:"sha256"`
}

type EngineBinding struct {
	Name           string `json:"name"`
	Version        string `json:"version"`
	SourceRevision string `json:"sourceRevision"`
	OS             string `json:"os"`
	Arch           string `json:"arch"`
	Path           string `json:"path"`
	SHA256         string `json:"sha256"`
}

type PlatformDigest struct {
	OS     string `json:"os"`
	Arch   string `json:"arch"`
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

type NamedDigest struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

type SchemaBinding struct {
	Family  string `json:"family"`
	Version string `json:"version"`
	Path    string `json:"path"`
	SHA256  string `json:"sha256"`
}

type ReleaseIdentity struct {
	Repository          string `json:"repository"`
	RepositoryOwnerID   int64  `json:"repositoryOwnerId"`
	Workflow            string `json:"workflow"`
	Ref                 string `json:"ref"`
	WorkflowSHA         string `json:"workflowSha,omitempty"`
	Trigger             string `json:"trigger,omitempty"`
	OIDCIssuer          string `json:"oidcIssuer"`
	CertificateIdentity string `json:"certificateIdentity"`
}

type ProductSourceIdentity struct {
	Tag    string `json:"tag"`
	Commit string `json:"commit"`
	Tree   string `json:"tree"`
}

type ReleaseToolingIdentity struct {
	Tag         string `json:"tag"`
	Commit      string `json:"commit"`
	Tree        string `json:"tree"`
	Workflow    string `json:"workflow"`
	WorkflowRef string `json:"workflowRef"`
	WorkflowSHA string `json:"workflowSha"`
	Trigger     string `json:"trigger"`
}

type ReleaseCompatibility struct {
	MinimumRunnerVersion      string   `json:"minimumRunnerVersion"`
	SupportedOperatingSystems []string `json:"supportedOperatingSystems"`
	SupportedArchitectures    []string `json:"supportedArchitectures"`
	TestSummaryAsset          string   `json:"testSummaryAsset"`
	LimitationsAsset          string   `json:"limitationsAsset"`
}

type ReleaseRevocation struct {
	DiscoveryLocation       string `json:"discoveryLocation"`
	SnapshotAsset           string `json:"snapshotAsset"`
	CheckpointAsset         string `json:"checkpointAsset"`
	SchemaVersion           string `json:"schemaVersion"`
	MaximumSnapshotAgeHours int    `json:"maximumSnapshotAgeHours"`
}

type ReleaseManifest struct {
	BuildIdentity         *BuildIdentity          `json:"buildIdentity,omitempty"`
	ReleaseState          string                  `json:"releaseState,omitempty"`
	SchemaFamily          string                  `json:"schemaFamily"`
	ManifestSchemaVersion string                  `json:"manifestSchemaVersion"`
	ReleaseVersion        string                  `json:"releaseVersion"`
	SourceRevision        string                  `json:"sourceRevision,omitempty"`
	SourceTree            string                  `json:"sourceTree,omitempty"`
	ProductSource         *ProductSourceIdentity  `json:"productSource,omitempty"`
	ReleaseTooling        *ReleaseToolingIdentity `json:"releaseTooling,omitempty"`
	RunnerVersion         string                  `json:"runnerVersion"`
	GoToolchainVersion    string                  `json:"goToolchainVersion"`
	RunnerBindings        []PlatformDigest        `json:"runnerBindings"`
	EngineBindings        []EngineBinding         `json:"engineBindings"`
	RulePack              NamedDigest             `json:"rulePack"`
	SchemaBindings        []SchemaBinding         `json:"schemaBindings"`
	ReleaseIdentity       ReleaseIdentity         `json:"releaseIdentity"`
	Compatibility         ReleaseCompatibility    `json:"compatibility"`
	Revocation            ReleaseRevocation       `json:"revocation"`
	Assets                []ReleaseAsset          `json:"assets"`
	CreatedAt             string                  `json:"createdAt"`
}

type ReleaseTrustPolicy struct {
	Repository            string
	RepositoryOwnerID     int64
	Workflow              string
	Ref                   string
	OIDCIssuer            string
	CertificateIdentity   string
	ReleaseVersion        string
	ManifestSchemaVersion string
	ProductSourceTag      string
	ProductSourceCommit   string
	ProductSourceTree     string
	ReleaseToolingTag     string
	ReleaseToolingCommit  string
	ReleaseToolingTree    string
	WorkflowSHA           string
	Trigger               string
}

type ReleaseSignatureVerifier interface {
	VerifyManifest(ctx context.Context, manifestPath, bundlePath string, identity ReleaseIdentity) error
}

type ReleaseVerificationRequest struct {
	Directory    string
	ManifestPath string
	BundlePath   string
	Policy       ReleaseTrustPolicy
	Signature    ReleaseSignatureVerifier
	Now          time.Time
}

type ReleaseVerificationResult struct {
	ManifestDigest string
	AssetCount     int
	ReleaseVersion string
}

func ParseReleaseManifest(raw []byte) (ReleaseManifest, error) {
	if len(raw) == 0 || len(raw) > MaximumReleaseManifestBytes {
		return ReleaseManifest{}, ErrInvalidReference
	}
	if err := rejectDuplicateJSONMembers(raw); err != nil {
		return ReleaseManifest{}, err
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	var manifest ReleaseManifest
	if err := decoder.Decode(&manifest); err != nil {
		return ReleaseManifest{}, ErrInvalidReference
	}
	if err := ensureJSONEOF(decoder); err != nil {
		return ReleaseManifest{}, err
	}
	var fields map[string]json.RawMessage
	if json.Unmarshal(raw, &fields) != nil {
		return ReleaseManifest{}, ErrInvalidReference
	}
	if manifest.ManifestSchemaVersion == "2.3" || manifest.ManifestSchemaVersion == "2.4" || manifest.ManifestSchemaVersion == "2.5" {
		if !exactObjectFields(raw, "schemaFamily", "manifestSchemaVersion", "releaseVersion", "productSource", "releaseTooling", "runnerVersion", "goToolchainVersion", "runnerBindings", "engineBindings", "rulePack", "schemaBindings", "releaseIdentity", "compatibility", "revocation", "assets", "createdAt", "buildIdentity", "releaseState") || !exactBuildFields(raw) {
			return ReleaseManifest{}, ErrInvalidReference
		}
		if manifest.ReleaseState == "signing-pending" && !exactIdentityFields(fields["releaseIdentity"], true) {
			return ReleaseManifest{}, ErrInvalidReference
		}
		identity, exists := fields["releaseIdentity"]
		if !exists || (manifest.ReleaseState == "unsigned-candidate" && !bytes.Equal(bytes.TrimSpace(identity), []byte("null"))) {
			return ReleaseManifest{}, ErrInvalidReference
		}
	} else {
		if _, ok := fields["buildIdentity"]; ok {
			return ReleaseManifest{}, ErrInvalidReference
		}
		if _, ok := fields["releaseState"]; ok {
			return ReleaseManifest{}, ErrInvalidReference
		}
	}
	if err := validateReleaseManifest(manifest); err != nil {
		return ReleaseManifest{}, err
	}
	return manifest, nil
}

func VerifyRelease(ctx context.Context, request ReleaseVerificationRequest) (ReleaseVerificationResult, error) {
	if ctx == nil || request.Signature == nil || request.Now.IsZero() {
		return ReleaseVerificationResult{}, ErrInvalidReference
	}
	root, err := filepath.Abs(request.Directory)
	if err != nil {
		return ReleaseVerificationResult{}, ErrInvalidReference
	}
	manifestPath, err := confinedRegularFile(root, request.ManifestPath)
	if err != nil {
		return ReleaseVerificationResult{}, err
	}
	bundlePath, err := confinedRegularFile(root, request.BundlePath)
	if err != nil {
		return ReleaseVerificationResult{}, err
	}
	raw, _, err := readRegularFileBounded(manifestPath, MaximumReleaseManifestBytes)
	if err != nil {
		return ReleaseVerificationResult{}, ErrInvalidReference
	}
	bundleBefore, _, err := readRegularFileBounded(bundlePath, MaximumReleaseManifestBytes)
	if err != nil {
		return ReleaseVerificationResult{}, ErrInvalidReference
	}
	manifest, err := ParseReleaseManifest(raw)
	if err != nil {
		return ReleaseVerificationResult{}, err
	}
	if err := checkReleasePolicy(manifest, request.Policy, request.Now); err != nil {
		return ReleaseVerificationResult{}, err
	}
	if err := request.Signature.VerifyManifest(ctx, manifestPath, bundlePath, manifest.ReleaseIdentity); err != nil {
		return ReleaseVerificationResult{}, ErrInvalidSignature
	}
	if err := verifyReleaseAssets(root, manifest, request.ManifestPath, request.BundlePath); err != nil {
		return ReleaseVerificationResult{}, err
	}
	if err := verifyInitialRevocationEvidence(root, manifest); err != nil {
		return ReleaseVerificationResult{}, err
	}
	manifestAfter, _, err := readRegularFileBounded(manifestPath, MaximumReleaseManifestBytes)
	if err != nil || !bytes.Equal(raw, manifestAfter) {
		return ReleaseVerificationResult{}, ErrBindingMismatch
	}
	bundleAfter, _, err := readRegularFileBounded(bundlePath, MaximumReleaseManifestBytes)
	if err != nil || !bytes.Equal(bundleBefore, bundleAfter) {
		return ReleaseVerificationResult{}, ErrBindingMismatch
	}
	return ReleaseVerificationResult{ManifestDigest: DigestBytes(raw), AssetCount: len(manifest.Assets), ReleaseVersion: manifest.ReleaseVersion}, nil
}

type initialRevocationSnapshot struct {
	SchemaFamily  string            `json:"schemaFamily"`
	SchemaVersion string            `json:"schemaVersion"`
	CapturedAt    string            `json:"capturedAt"`
	Records       []json.RawMessage `json:"records"`
}

type initialRevocationCheckpoint struct {
	SchemaFamily      string  `json:"schemaFamily"`
	SchemaVersion     string  `json:"schemaVersion"`
	Sequence          int64   `json:"sequence"`
	Digest            *string `json:"digest"`
	CapturedAt        string  `json:"capturedAt"`
	DiscoveryLocation string  `json:"discoveryLocation"`
}

func verifyInitialRevocationEvidence(root string, manifest ReleaseManifest) error {
	snapshotPath, err := confinedRegularFile(root, manifest.Revocation.SnapshotAsset)
	if err != nil {
		return err
	}
	checkpointPath, err := confinedRegularFile(root, manifest.Revocation.CheckpointAsset)
	if err != nil {
		return err
	}
	var snapshot initialRevocationSnapshot
	if err := decodeStrictFile(snapshotPath, &snapshot); err != nil || snapshot.SchemaFamily != "global-scanner-revocation-snapshot" ||
		snapshot.SchemaVersion != "1.0" || len(snapshot.Records) != 0 {
		return ErrInvalidReference
	}
	var checkpoint initialRevocationCheckpoint
	if err := decodeStrictFile(checkpointPath, &checkpoint); err != nil || checkpoint.SchemaFamily != "global-scanner-revocation-checkpoint" ||
		checkpoint.SchemaVersion != "1.0" || checkpoint.Sequence != 0 || checkpoint.Digest != nil || checkpoint.DiscoveryLocation != manifest.Revocation.DiscoveryLocation ||
		checkpoint.CapturedAt != snapshot.CapturedAt || checkpoint.CapturedAt != manifest.CreatedAt {
		return ErrInvalidReference
	}
	return nil
}

func decodeStrictFile(path string, destination any) error {
	raw, _, err := readRegularFileBounded(path, MaximumReleaseManifestBytes)
	if err != nil || len(raw) == 0 {
		return ErrInvalidReference
	}
	if err := rejectDuplicateJSONMembers(raw); err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return ErrInvalidReference
	}
	return ensureJSONEOF(decoder)
}

func validateReleaseManifest(m ReleaseManifest) error {
	if m.SchemaFamily != "scanner-release-manifest" || !releaseVersionPattern.MatchString(m.ReleaseVersion) ||
		!runnerVersionPattern.MatchString(m.RunnerVersion) ||
		!toolchainVersionPattern.MatchString(m.GoToolchainVersion) || len(m.RunnerBindings) != 2 || len(m.EngineBindings) != 2 ||
		len(m.SchemaBindings) == 0 || len(m.SchemaBindings) > 32 || len(m.Assets) < 10 || len(m.Assets) > MaximumReleaseAssets {
		return ErrInvalidReference
	}
	if err := validateReleaseIdentityVersion(m); err != nil {
		return err
	}
	if _, err := parseCanonicalTime(m.CreatedAt); err != nil {
		return err
	}
	if !validAssetReference(m.RulePack.Path, m.RulePack.SHA256) || !runnerVersionPattern.MatchString(m.Compatibility.MinimumRunnerVersion) ||
		!equalStrings(m.Compatibility.SupportedOperatingSystems, []string{"linux", "windows"}) ||
		!equalStrings(m.Compatibility.SupportedArchitectures, []string{"amd64"}) ||
		m.Revocation.SchemaVersion != "1.1" || m.Revocation.MaximumSnapshotAgeHours != 24 ||
		!strings.HasPrefix(m.Revocation.DiscoveryLocation, "https://api.github.com/repos/ThameeraDananjaya/project-agnostic-secret-scanner/") {
		return ErrInvalidReference
	}
	seenPlatforms := map[string]bool{}
	for _, runner := range m.RunnerBindings {
		key := runner.OS + "/" + runner.Arch
		if (runner.OS != "linux" && runner.OS != "windows") || runner.Arch != "amd64" || seenPlatforms[key] ||
			!validAssetReference(runner.Path, runner.SHA256) {
			return ErrInvalidReference
		}
		seenPlatforms[key] = true
	}
	seenPlatforms = map[string]bool{}
	for _, engine := range m.EngineBindings {
		key := engine.OS + "/" + engine.Arch
		if engine.Name != "gitleaks" || engine.Version != "8.30.1" || engine.SourceRevision != "83d9cd684c87d95d656c1458ef04895a7f1cbd8e" ||
			(engine.OS != "linux" && engine.OS != "windows") || engine.Arch != "amd64" || seenPlatforms[key] ||
			!validAssetReference(engine.Path, engine.SHA256) {
			return ErrInvalidReference
		}
		seenPlatforms[key] = true
	}
	seenSchemas := map[string]bool{}
	for _, schema := range m.SchemaBindings {
		key := schema.Family + "@" + schema.Version
		if !IsToken(schema.Family) || !validSchemaVersion(schema.Version) || seenSchemas[key] || !validAssetReference(schema.Path, schema.SHA256) {
			return ErrInvalidReference
		}
		seenSchemas[key] = true
	}
	seenAssets := map[string]bool{}
	assetDigests := map[string]string{}
	for _, asset := range m.Assets {
		if seenAssets[asset.Path] || !validAssetReference(asset.Path, asset.SHA256) || asset.Size < 1 || asset.Size > 1<<30 ||
			!validReleaseKind(asset.Kind) || !validPlatform(asset.OS, asset.Arch) {
			return ErrInvalidReference
		}
		seenAssets[asset.Path] = true
		assetDigests[asset.Path] = asset.SHA256
	}
	for _, path := range []string{m.RulePack.Path, m.Compatibility.TestSummaryAsset, m.Compatibility.LimitationsAsset, m.Revocation.SnapshotAsset, m.Revocation.CheckpointAsset} {
		if !seenAssets[path] {
			return ErrBindingMismatch
		}
	}
	for _, runner := range m.RunnerBindings {
		if !seenAssets[runner.Path] || assetDigests[runner.Path] != runner.SHA256 {
			return ErrBindingMismatch
		}
	}
	for _, engine := range m.EngineBindings {
		if !seenAssets[engine.Path] || assetDigests[engine.Path] != engine.SHA256 {
			return ErrBindingMismatch
		}
	}
	for _, schema := range m.SchemaBindings {
		if !seenAssets[schema.Path] || assetDigests[schema.Path] != schema.SHA256 {
			return ErrBindingMismatch
		}
	}
	if assetDigests[m.RulePack.Path] != m.RulePack.SHA256 {
		return ErrBindingMismatch
	}
	return nil
}

func validateReleaseIdentityVersion(m ReleaseManifest) error {
	if m.ManifestSchemaVersion == "2.3" || m.ManifestSchemaVersion == "2.4" || m.ManifestSchemaVersion == "2.5" {
		return validateUnsignedBuildIdentity(m)
	}
	if m.BuildIdentity != nil || m.ReleaseState != "" {
		return ErrInvalidReference
	}
	switch m.ManifestSchemaVersion {
	case "1.1":
		if !gitOIDPattern.MatchString(m.SourceRevision) || !gitOIDPattern.MatchString(m.SourceTree) ||
			m.ProductSource != nil || m.ReleaseTooling != nil ||
			m.ReleaseIdentity.WorkflowSHA != "" || m.ReleaseIdentity.Trigger != "" {
			return ErrInvalidReference
		}
	case "2.0", "2.1", "2.2":
		if m.ProductSource == nil || m.ReleaseTooling == nil {
			return ErrInvalidReference
		}
		product := m.ProductSource
		tooling := m.ReleaseTooling
		identity := m.ReleaseIdentity
		expectedToolingTag := "release-tooling-v1.0.0-c1"
		expectedWorkflow := ".github/workflows/release-recovery-v1.0.0.yml"
		if m.ManifestSchemaVersion == "2.1" {
			expectedToolingTag = "release-tooling-v1.0.0-c2"
		} else if m.ManifestSchemaVersion == "2.2" {
			expectedToolingTag = "release-tooling-v1.0.0-c2-r6"
			expectedWorkflow = ".github/workflows/release-recovery-v1.0.0-c2-r6.yml"
		}
		expectedToolingRef := "refs/tags/" + expectedToolingTag
		expectedCertificateIdentity := "https://github.com/ThameeraDananjaya/project-agnostic-secret-scanner/" + expectedWorkflow + "@" + expectedToolingRef
		if m.SourceRevision != "" || m.SourceTree != "" || m.ReleaseVersion != "v1.0.0" ||
			product.Tag != "v1.0.0" || product.Commit != "a13c28fe7273bc8dc6545f97966a02889524eb4c" ||
			product.Tree != "217b711ddea51fd0ea7e808edd2e27fdecef8427" ||
			tooling.Tag != expectedToolingTag || !gitOIDPattern.MatchString(tooling.Commit) ||
			!gitOIDPattern.MatchString(tooling.Tree) || tooling.Workflow != expectedWorkflow ||
			tooling.WorkflowRef != expectedToolingRef || tooling.WorkflowSHA != tooling.Commit ||
			tooling.Trigger != "workflow_dispatch" || identity.Workflow != tooling.Workflow || identity.Ref != tooling.WorkflowRef ||
			identity.WorkflowSHA != tooling.WorkflowSHA || identity.Trigger != tooling.Trigger ||
			identity.Repository != "ThameeraDananjaya/project-agnostic-secret-scanner" || identity.RepositoryOwnerID != 50274860 ||
			identity.OIDCIssuer != "https://token.actions.githubusercontent.com" || identity.CertificateIdentity != expectedCertificateIdentity {
			return ErrInvalidReference
		}
	default:
		return ErrInvalidReference
	}
	return nil
}

func checkReleasePolicy(m ReleaseManifest, p ReleaseTrustPolicy, now time.Time) error {
	if m.ManifestSchemaVersion == "2.3" || m.ManifestSchemaVersion == "2.4" || m.ManifestSchemaVersion == "2.5" {
		return checkSeparateSignerPolicy(m, p, now)
	}
	identity := m.ReleaseIdentity
	if p.Repository == "" || p.RepositoryOwnerID <= 0 || p.Workflow == "" || p.Ref == "" || p.OIDCIssuer == "" || p.CertificateIdentity == "" || p.ReleaseVersion == "" ||
		identity.Repository != p.Repository || identity.RepositoryOwnerID != p.RepositoryOwnerID || identity.Workflow != p.Workflow || identity.Ref != p.Ref ||
		identity.OIDCIssuer != p.OIDCIssuer || identity.CertificateIdentity != p.CertificateIdentity || m.ReleaseVersion != p.ReleaseVersion ||
		identity.Repository != "ThameeraDananjaya/project-agnostic-secret-scanner" || identity.RepositoryOwnerID != 50274860 ||
		identity.OIDCIssuer != "https://token.actions.githubusercontent.com" ||
		identity.CertificateIdentity != "https://github.com/"+identity.Repository+"/"+identity.Workflow+"@"+identity.Ref {
		return ErrBindingMismatch
	}
	switch m.ManifestSchemaVersion {
	case "1.1":
		if p.ManifestSchemaVersion != "" && p.ManifestSchemaVersion != "1.1" {
			return ErrBindingMismatch
		}
		if identity.Workflow != ".github/workflows/release.yml" || identity.Ref != "refs/tags/"+m.ReleaseVersion ||
			identity.WorkflowSHA != "" || identity.Trigger != "" {
			return ErrBindingMismatch
		}
	case "2.0", "2.1", "2.2":
		if m.ProductSource == nil || m.ReleaseTooling == nil {
			return ErrBindingMismatch
		}
		product := m.ProductSource
		tooling := m.ReleaseTooling
		if p.ManifestSchemaVersion != m.ManifestSchemaVersion || p.ProductSourceTag == "" || p.ProductSourceCommit == "" || p.ProductSourceTree == "" ||
			p.ReleaseToolingTag == "" || p.ReleaseToolingCommit == "" || p.ReleaseToolingTree == "" || p.WorkflowSHA == "" || p.Trigger == "" ||
			product.Tag != p.ProductSourceTag || product.Commit != p.ProductSourceCommit || product.Tree != p.ProductSourceTree ||
			tooling.Tag != p.ReleaseToolingTag || tooling.Commit != p.ReleaseToolingCommit || tooling.Tree != p.ReleaseToolingTree ||
			tooling.Workflow != p.Workflow || tooling.WorkflowRef != p.Ref || tooling.WorkflowSHA != p.WorkflowSHA || tooling.Trigger != p.Trigger ||
			identity.WorkflowSHA != p.WorkflowSHA || identity.Trigger != p.Trigger {
			return ErrBindingMismatch
		}
	default:
		return ErrBindingMismatch
	}
	created, _ := parseCanonicalTime(m.CreatedAt)
	if created.After(now) {
		return ErrInvalidReference
	}
	return nil
}

func verifyReleaseAssets(root string, manifest ReleaseManifest, manifestName, bundleName string) error {
	expected := map[string]bool{filepath.ToSlash(manifestName): true, filepath.ToSlash(bundleName): true}
	for _, asset := range manifest.Assets {
		path, err := confinedRegularFile(root, asset.Path)
		if err != nil {
			return err
		}
		digest, size, err := digestRegularFile(path)
		if err != nil || size != asset.Size || digest != asset.SHA256 {
			return ErrBindingMismatch
		}
		expected[asset.Path] = true
	}
	var discovered []string
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == root {
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return ErrInvalidReference
		}
		if entry.IsDir() {
			return nil
		}
		if !entry.Type().IsRegular() {
			return ErrInvalidReference
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return ErrInvalidReference
		}
		discovered = append(discovered, filepath.ToSlash(relative))
		return nil
	})
	if err != nil {
		return ErrInvalidReference
	}
	sort.Strings(discovered)
	if len(discovered) != len(expected) {
		return ErrBindingMismatch
	}
	for _, path := range discovered {
		if !expected[path] {
			return ErrBindingMismatch
		}
	}
	return nil
}

func confinedRegularFile(root, relative string) (string, error) {
	if !validReleasePath(filepath.ToSlash(relative)) {
		return "", ErrInvalidReference
	}
	path := filepath.Join(root, filepath.FromSlash(relative))
	cleanRoot := filepath.Clean(root) + string(os.PathSeparator)
	if !strings.HasPrefix(filepath.Clean(path)+string(os.PathSeparator), cleanRoot) {
		return "", ErrInvalidReference
	}
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
		return "", ErrInvalidReference
	}
	return path, nil
}

func digestReader(reader io.Reader) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func digestRegularFile(path string) (string, int64, error) {
	file, before, err := openStableRegularFile(path)
	if err != nil {
		return "", 0, err
	}
	digest, digestErr := digestReader(file)
	afterOpen, statErr := file.Stat()
	closeErr := file.Close()
	afterPath, lstatErr := os.Lstat(path)
	if digestErr != nil || statErr != nil || closeErr != nil || lstatErr != nil || !afterPath.Mode().IsRegular() ||
		afterPath.Mode()&os.ModeSymlink != 0 || !os.SameFile(before, afterOpen) || !os.SameFile(before, afterPath) || before.Size() != afterOpen.Size() {
		return "", 0, ErrInvalidReference
	}
	return digest, before.Size(), nil
}

func readRegularFileBounded(path string, maximum int64) ([]byte, int64, error) {
	file, before, err := openStableRegularFile(path)
	if err != nil || before.Size() < 1 || before.Size() > maximum {
		if file != nil {
			_ = file.Close()
		}
		return nil, 0, ErrInvalidReference
	}
	raw, readErr := io.ReadAll(io.LimitReader(file, maximum+1))
	afterOpen, statErr := file.Stat()
	closeErr := file.Close()
	afterPath, lstatErr := os.Lstat(path)
	if readErr != nil || statErr != nil || closeErr != nil || lstatErr != nil || int64(len(raw)) != before.Size() ||
		!afterPath.Mode().IsRegular() || afterPath.Mode()&os.ModeSymlink != 0 || !os.SameFile(before, afterOpen) || !os.SameFile(before, afterPath) {
		return nil, 0, ErrInvalidReference
	}
	return raw, before.Size(), nil
}

func openStableRegularFile(path string) (*os.File, os.FileInfo, error) {
	before, err := os.Lstat(path)
	if err != nil || !before.Mode().IsRegular() || before.Mode()&os.ModeSymlink != 0 {
		return nil, nil, ErrInvalidReference
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, ErrInvalidReference
	}
	opened, err := file.Stat()
	if err != nil || !opened.Mode().IsRegular() || !os.SameFile(before, opened) {
		_ = file.Close()
		return nil, nil, ErrInvalidReference
	}
	return file, before, nil
}

func validAssetReference(path, digest string) bool { return validReleasePath(path) && IsDigest(digest) }

func validReleasePath(path string) bool {
	return assetPathPattern.MatchString(path) && !strings.Contains(path, "//") && !strings.Contains(path, "../") &&
		!strings.HasPrefix(path, "/") && !strings.HasSuffix(path, "/") && !strings.Contains(path, `\`)
}

func validPlatform(osName, arch string) bool {
	return (osName == "none" && arch == "none") || ((osName == "linux" || osName == "windows") && arch == "amd64")
}

func validReleaseKind(kind string) bool {
	switch kind {
	case "platform-bundle", "runner", "engine", "verifier", "rules", "schema", "checksums", "licence", "licence-manifest", "sbom", "test-summary", "limitations", "compatibility", "revocation-snapshot", "revocation-checkpoint", "documentation":
		return true
	default:
		return false
	}
}

func equalStrings(actual, expected []string) bool {
	if len(actual) != len(expected) {
		return false
	}
	for i := range actual {
		if actual[i] != expected[i] {
			return false
		}
	}
	return true
}

func ensureJSONEOF(decoder *json.Decoder) error {
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return ErrInvalidReference
	}
	return nil
}

func rejectDuplicateJSONMembers(raw []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	if err := inspectJSONValue(decoder); err != nil {
		return err
	}
	return ensureJSONEOF(decoder)
}

func inspectJSONValue(decoder *json.Decoder) error {
	token, err := decoder.Token()
	if err != nil {
		return ErrInvalidReference
	}
	delimiter, ok := token.(json.Delim)
	if !ok {
		return nil
	}
	switch delimiter {
	case '{':
		seen := map[string]bool{}
		for decoder.More() {
			keyToken, err := decoder.Token()
			key, ok := keyToken.(string)
			if err != nil || !ok || seen[key] {
				return ErrInvalidReference
			}
			seen[key] = true
			if err := inspectJSONValue(decoder); err != nil {
				return err
			}
		}
		closing, err := decoder.Token()
		if err != nil || closing != json.Delim('}') {
			return ErrInvalidReference
		}
	case '[':
		for decoder.More() {
			if err := inspectJSONValue(decoder); err != nil {
				return err
			}
		}
		closing, err := decoder.Token()
		if err != nil || closing != json.Delim(']') {
			return ErrInvalidReference
		}
	default:
		return ErrInvalidReference
	}
	return nil
}
