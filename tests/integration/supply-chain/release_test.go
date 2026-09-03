package supplychain_test

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/verify"
)

func TestCosignCommandVerifierBindsExecutableAndTrustedRootBytes(t *testing.T) {
	root := t.TempDir()
	cosignPath := filepath.Join(root, "cosign")
	trustedRootPath := filepath.Join(root, "trusted-root.json")
	cosignBytes := []byte("synthetic-cosign-executable")
	trustedRootBytes := []byte(`{"synthetic":"trusted-root"}`)
	if err := os.WriteFile(cosignPath, cosignBytes, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(trustedRootPath, trustedRootBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := verify.NewCosignCommandVerifier(cosignPath, verify.DigestBytes(cosignBytes), trustedRootPath, verify.DigestBytes(trustedRootBytes)); err != nil {
		t.Fatalf("exact Cosign and trusted-root bytes rejected: %v", err)
	}
	trustedRootBytes[0] ^= 1
	if err := os.WriteFile(trustedRootPath, trustedRootBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := verify.NewCosignCommandVerifier(cosignPath, verify.DigestBytes(cosignBytes), trustedRootPath, verify.DigestBytes([]byte(`{"synthetic":"trusted-root"}`))); !errors.Is(err, verify.ErrBindingMismatch) {
		t.Fatalf("mutated trusted-root bytes did not fail closed: %v", err)
	}
}

func TestCosignCommandVerifierBindsEveryGitHubWorkflowClaim(t *testing.T) {
	root := t.TempDir()
	argumentsPath := filepath.Join(root, "arguments.txt")
	cosignPath := filepath.Join(root, "cosign")
	trustedRootPath := filepath.Join(root, "trusted-root.json")
	manifestPath := filepath.Join(root, "release-manifest.json")
	bundlePath := filepath.Join(root, "release-manifest.sigstore.json")
	script := []byte("#!/bin/sh\nprintf '%s\\n' \"$@\" > '" + argumentsPath + "'\n")
	if err := os.WriteFile(cosignPath, script, 0o700); err != nil {
		t.Fatal(err)
	}
	trustedRoot := []byte(`{"synthetic":"trusted-root"}`)
	if err := os.WriteFile(trustedRootPath, trustedRoot, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(manifestPath, []byte("manifest"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(bundlePath, []byte("bundle"), 0o600); err != nil {
		t.Fatal(err)
	}
	verifier, err := verify.NewCosignCommandVerifier(cosignPath, verify.DigestBytes(script), trustedRootPath, verify.DigestBytes(trustedRoot))
	if err != nil {
		t.Fatal(err)
	}
	identity := releaseIdentityV2()
	if err := verifier.VerifyManifest(t.Context(), manifestPath, bundlePath, identity); err != nil {
		t.Fatalf("synthetic Cosign invocation rejected: %v", err)
	}
	raw, err := os.ReadFile(argumentsPath)
	if err != nil {
		t.Fatal(err)
	}
	arguments := string(raw)
	for _, required := range []string{
		"--certificate-identity\n" + identity.CertificateIdentity,
		"--certificate-oidc-issuer\n" + identity.OIDCIssuer,
		"--certificate-github-workflow-repository\n" + identity.Repository,
		"--certificate-github-workflow-ref\n" + identity.Ref,
		"--certificate-github-workflow-sha\n" + identity.WorkflowSHA,
		"--certificate-github-workflow-trigger\n" + identity.Trigger,
	} {
		if !strings.Contains(arguments, required) {
			t.Fatalf("required Cosign claim argument absent: %q in %q", required, arguments)
		}
	}
}

func TestReleaseManifestParserRejectsUnknownDuplicateAndUnsupportedMembers(t *testing.T) {
	if _, err := verify.ParseReleaseManifest([]byte(`{"schemaFamily":"scanner-release-manifest","schemaFamily":"scanner-release-manifest"}`)); !errors.Is(err, verify.ErrInvalidReference) {
		t.Fatalf("duplicate member was not rejected: %v", err)
	}
	manifest := minimumManifest()
	raw, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	raw = append(raw[:len(raw)-1], []byte(`,"unexpected":true}`)...)
	if _, err := verify.ParseReleaseManifest(raw); !errors.Is(err, verify.ErrInvalidReference) {
		t.Fatalf("unknown member was not rejected: %v", err)
	}
	manifest.ManifestSchemaVersion = "3.0"
	raw, _ = json.Marshal(manifest)
	if _, err := verify.ParseReleaseManifest(raw); !errors.Is(err, verify.ErrInvalidReference) {
		t.Fatalf("unknown schema major was not rejected: %v", err)
	}
}

func TestReleaseManifestV2RequiresDistinctCompleteIdentities(t *testing.T) {
	manifest := minimumManifestV2()
	raw, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := verify.ParseReleaseManifest(raw); err != nil {
		t.Fatalf("exact schema 2.0 manifest rejected: %v", err)
	}

	tests := map[string]func(*verify.ReleaseManifest){
		"old schema masquerade": func(m *verify.ReleaseManifest) { m.SourceRevision = m.ProductSource.Commit },
		"missing product role":  func(m *verify.ReleaseManifest) { m.ProductSource = nil },
		"missing tooling role":  func(m *verify.ReleaseManifest) { m.ReleaseTooling = nil },
		"swapped tags": func(m *verify.ReleaseManifest) {
			m.ProductSource.Tag, m.ReleaseTooling.Tag = m.ReleaseTooling.Tag, m.ProductSource.Tag
		},
		"workflow sha differs from tooling commit": func(m *verify.ReleaseManifest) {
			m.ReleaseTooling.WorkflowSHA = "1111111111111111111111111111111111111111"
		},
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			candidate := minimumManifestV2()
			mutate(&candidate)
			raw, err := json.Marshal(candidate)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := verify.ParseReleaseManifest(raw); !errors.Is(err, verify.ErrInvalidReference) {
				t.Fatalf("invalid identity arrangement was not rejected: %v", err)
			}
		})
	}
}

func minimumManifest() verify.ReleaseManifest {
	digest := verify.DigestBytes([]byte("x"))
	assets := []verify.ReleaseAsset{}
	for _, path := range []string{"runner-linux", "runner-windows.exe", "engine-linux", "engine-windows.exe", "rules.toml", "schema.json", "tests.json", "limitations.md", "revocations.json", "checkpoint.json"} {
		assets = append(assets, verify.ReleaseAsset{Path: path, Kind: "documentation", OS: "none", Arch: "none", Size: 1, SHA256: digest})
	}
	return verify.ReleaseManifest{
		SchemaFamily: "scanner-release-manifest", ManifestSchemaVersion: "1.1", ReleaseVersion: "v1.0.0",
		SourceRevision: "0123456789abcdef0123456789abcdef01234567", SourceTree: "89abcdef0123456789abcdef0123456789abcdef",
		RunnerVersion: "1.0.0", GoToolchainVersion: "go1.27.1",
		RunnerBindings: []verify.PlatformDigest{{OS: "linux", Arch: "amd64", Path: "runner-linux", SHA256: digest}, {OS: "windows", Arch: "amd64", Path: "runner-windows.exe", SHA256: digest}},
		EngineBindings: []verify.EngineBinding{
			{Name: "gitleaks", Version: "8.30.1", SourceRevision: "83d9cd684c87d95d656c1458ef04895a7f1cbd8e", OS: "linux", Arch: "amd64", Path: "engine-linux", SHA256: digest},
			{Name: "gitleaks", Version: "8.30.1", SourceRevision: "83d9cd684c87d95d656c1458ef04895a7f1cbd8e", OS: "windows", Arch: "amd64", Path: "engine-windows.exe", SHA256: digest},
		},
		RulePack:        verify.NamedDigest{Path: "rules.toml", SHA256: digest},
		SchemaBindings:  []verify.SchemaBinding{{Family: "scan-request", Version: "1.1", Path: "schema.json", SHA256: digest}},
		ReleaseIdentity: releaseIdentity(),
		Compatibility:   verify.ReleaseCompatibility{MinimumRunnerVersion: "1.0.0", SupportedOperatingSystems: []string{"linux", "windows"}, SupportedArchitectures: []string{"amd64"}, TestSummaryAsset: "tests.json", LimitationsAsset: "limitations.md"},
		Revocation:      verify.ReleaseRevocation{DiscoveryLocation: "https://api.github.com/repos/ThameeraDananjaya/project-agnostic-secret-scanner/releases?per_page=100", SnapshotAsset: "revocations.json", CheckpointAsset: "checkpoint.json", SchemaVersion: "1.1", MaximumSnapshotAgeHours: 24},
		Assets:          assets, CreatedAt: "2026-09-03T00:00:00Z",
	}
}

func minimumManifestV2() verify.ReleaseManifest {
	manifest := minimumManifest()
	manifest.ManifestSchemaVersion = "2.0"
	manifest.SourceRevision = ""
	manifest.SourceTree = ""
	manifest.ProductSource = &verify.ProductSourceIdentity{Tag: "v1.0.0", Commit: "a13c28fe7273bc8dc6545f97966a02889524eb4c", Tree: "217b711ddea51fd0ea7e808edd2e27fdecef8427"}
	manifest.ReleaseTooling = &verify.ReleaseToolingIdentity{
		Tag: "release-tooling-v1.0.0-c1", Commit: "3fb1b0a55dc4f48dd35464c63c768f497efbc89b", Tree: "544867396910969d20cfd2acd454ac0d69c1d7e4",
		Workflow: ".github/workflows/release-recovery-v1.0.0.yml", WorkflowRef: "refs/tags/release-tooling-v1.0.0-c1", WorkflowSHA: "3fb1b0a55dc4f48dd35464c63c768f497efbc89b", Trigger: "workflow_dispatch",
	}
	manifest.ReleaseIdentity = releaseIdentityV2()
	return manifest
}

func releaseIdentity() verify.ReleaseIdentity {
	return verify.ReleaseIdentity{Repository: "ThameeraDananjaya/project-agnostic-secret-scanner", RepositoryOwnerID: 50274860, Workflow: ".github/workflows/release.yml", Ref: "refs/tags/v1.0.0", OIDCIssuer: "https://token.actions.githubusercontent.com", CertificateIdentity: "https://github.com/ThameeraDananjaya/project-agnostic-secret-scanner/.github/workflows/release.yml@refs/tags/v1.0.0"}
}

func releaseIdentityV2() verify.ReleaseIdentity {
	return verify.ReleaseIdentity{
		Repository: "ThameeraDananjaya/project-agnostic-secret-scanner", RepositoryOwnerID: 50274860,
		Workflow: ".github/workflows/release-recovery-v1.0.0.yml", Ref: "refs/tags/release-tooling-v1.0.0-c1",
		WorkflowSHA: "3fb1b0a55dc4f48dd35464c63c768f497efbc89b", Trigger: "workflow_dispatch",
		OIDCIssuer:          "https://token.actions.githubusercontent.com",
		CertificateIdentity: "https://github.com/ThameeraDananjaya/project-agnostic-secret-scanner/.github/workflows/release-recovery-v1.0.0.yml@refs/tags/release-tooling-v1.0.0-c1",
	}
}
