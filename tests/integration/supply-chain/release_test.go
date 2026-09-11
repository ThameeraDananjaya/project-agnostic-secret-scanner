package supplychain_test

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
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

func TestReleaseManifestV21BindsCorrectionC2AndPreservesC1Parsing(t *testing.T) {
	c1 := minimumManifestV2()
	raw, err := json.Marshal(c1)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := verify.ParseReleaseManifest(raw); err != nil {
		t.Fatalf("historical schema 2.0 C1 manifest rejected: %v", err)
	}

	c2 := minimumManifestV21()
	raw, err = json.Marshal(c2)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := verify.ParseReleaseManifest(raw); err != nil {
		t.Fatalf("schema 2.1 C2 manifest rejected: %v", err)
	}

	for name, mutate := range map[string]func(*verify.ReleaseManifest){
		"C1 tag under 2.1": func(m *verify.ReleaseManifest) { m.ReleaseTooling.Tag = "release-tooling-v1.0.0-c1" },
		"C1 ref under 2.1": func(m *verify.ReleaseManifest) { m.ReleaseTooling.WorkflowRef = "refs/tags/release-tooling-v1.0.0-c1" },
		"role swap": func(m *verify.ReleaseManifest) { m.ProductSource.Tag, m.ReleaseTooling.Tag = m.ReleaseTooling.Tag, m.ProductSource.Tag },
	} {
		t.Run(name, func(t *testing.T) {
			candidate := minimumManifestV21()
			mutate(&candidate)
			raw, _ := json.Marshal(candidate)
			if _, err := verify.ParseReleaseManifest(raw); !errors.Is(err, verify.ErrInvalidReference) {
				t.Fatalf("invalid C2 identity was not rejected: %v", err)
			}
		})
	}
}

func TestReleaseManifestV22BindsR6AndRejectsCrossVersionMixtures(t *testing.T) {
	for name, manifest := range map[string]verify.ReleaseManifest{
		"historical 2.0 C1": minimumManifestV2(),
		"historical 2.1 C2": minimumManifestV21(),
		"current 2.2 R6":    minimumManifestV22(),
	} {
		t.Run(name, func(t *testing.T) {
			raw, err := json.Marshal(manifest)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := verify.ParseReleaseManifest(raw); err != nil {
				t.Fatalf("valid versioned manifest rejected: %v", err)
			}
		})
	}

	for name, mutate := range map[string]func(*verify.ReleaseManifest){
		"old tag under 2.2":      func(m *verify.ReleaseManifest) { m.ReleaseTooling.Tag = "release-tooling-v1.0.0-c2" },
		"old ref under 2.2":      func(m *verify.ReleaseManifest) { m.ReleaseTooling.WorkflowRef = "refs/tags/release-tooling-v1.0.0-c2" },
		"old path under 2.2":     func(m *verify.ReleaseManifest) { m.ReleaseTooling.Workflow = ".github/workflows/release-recovery-v1.0.0.yml" },
		"old identity under 2.2": func(m *verify.ReleaseManifest) { m.ReleaseIdentity.CertificateIdentity = "https://github.com/ThameeraDananjaya/project-agnostic-secret-scanner/.github/workflows/release-recovery-v1.0.0.yml@refs/tags/release-tooling-v1.0.0-c2" },
	} {
		t.Run(name, func(t *testing.T) {
			candidate := minimumManifestV22()
			mutate(&candidate)
			raw, _ := json.Marshal(candidate)
			if _, err := verify.ParseReleaseManifest(raw); !errors.Is(err, verify.ErrInvalidReference) {
				t.Fatalf("2.2 cross-version identity was not rejected: %v", err)
			}
		})
	}

	for name, mutate := range map[string]func(*verify.ReleaseManifest){
		"new tag under 2.1":      func(m *verify.ReleaseManifest) { m.ReleaseTooling.Tag = "release-tooling-v1.0.0-c2-r6" },
		"new ref under 2.1":      func(m *verify.ReleaseManifest) { m.ReleaseTooling.WorkflowRef = "refs/tags/release-tooling-v1.0.0-c2-r6" },
		"new path under 2.1":     func(m *verify.ReleaseManifest) { m.ReleaseTooling.Workflow = ".github/workflows/release-recovery-v1.0.0-c2-r6.yml" },
		"new identity under 2.1": func(m *verify.ReleaseManifest) { m.ReleaseIdentity.CertificateIdentity = "https://github.com/ThameeraDananjaya/project-agnostic-secret-scanner/.github/workflows/release-recovery-v1.0.0-c2-r6.yml@refs/tags/release-tooling-v1.0.0-c2-r6" },
	} {
		t.Run(name, func(t *testing.T) {
			candidate := minimumManifestV21()
			mutate(&candidate)
			raw, _ := json.Marshal(candidate)
			if _, err := verify.ParseReleaseManifest(raw); !errors.Is(err, verify.ErrInvalidReference) {
				t.Fatalf("2.1 cross-version identity was not rejected: %v", err)
			}
		})
	}
}

func TestIteration007RepositoryIdentityAgreementAndPreservation(t *testing.T) {
	root := filepath.Clean(filepath.Join("..", "..", ".."))
	read := func(path string) string {
		t.Helper()
		raw, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
		if err != nil {
			t.Fatal(err)
		}
		return string(raw)
	}

	oldSchema := read("contracts/release-manifest/schema-2.1.json")
	if verify.DigestBytes([]byte(oldSchema)) != "caa9cd26665cc3a3550affea7490a0f1f277a69b10b1f1537787616e8ab973ce" {
		t.Fatal("historical schema 2.1 bytes changed")
	}
	newSchema := read("contracts/release-manifest/schema-2.2.json")
	normalizedSchema := strings.ReplaceAll(newSchema, "release-tooling-v1.0.0-c2-r6", "release-tooling-v1.0.0-c2")
	normalizedSchema = strings.ReplaceAll(normalizedSchema, ".github/workflows/release-recovery-v1.0.0-c2-r6.yml", ".github/workflows/release-recovery-v1.0.0.yml")
	normalizedSchema = strings.ReplaceAll(normalizedSchema, "2.2", "2.1")
	if !reflect.DeepEqual([]byte(normalizedSchema), []byte(oldSchema)) {
		t.Fatal("schema 2.2 differs from schema 2.1 outside the selected identity/version fields")
	}

	oldWorkflow := read(".github/workflows/release-recovery-v1.0.0.yml")
	if verify.DigestBytes([]byte(oldWorkflow)) != "c5f40f1b32e87c005fe33ee607af7e3d19ee4f0173c21619e21158c31aa0bdb4" {
		t.Fatal("historical C2 workflow bytes changed")
	}
	newWorkflow := read(".github/workflows/release-recovery-v1.0.0-c2-r6.yml")
	normalizedWorkflow := strings.ReplaceAll(newWorkflow, "gated-v1.0.0-c2-r6-recovery", "gated-v1.0.0-c2-recovery")
	normalizedWorkflow = strings.ReplaceAll(normalizedWorkflow, "release-v1.0.0-c2-r6-recovery", "release-v1.0.0-c2-recovery")
	normalizedWorkflow = strings.ReplaceAll(normalizedWorkflow, "release-tooling-v1.0.0-c2-r6", "release-tooling-v1.0.0-c2")
	normalizedWorkflow = strings.ReplaceAll(normalizedWorkflow, ".github/workflows/release-recovery-v1.0.0-c2-r6.yml", ".github/workflows/release-recovery-v1.0.0.yml")
	normalizedWorkflow = strings.ReplaceAll(normalizedWorkflow, "2.2", "2.1")
	if normalizedWorkflow != oldWorkflow {
		t.Fatal("R6 workflow differs from the C2 workflow outside the selected identity/version fields")
	}

	for path, required := range map[string][]string{
		"build/release/build.ps1": {"release-tooling-v1.0.0-c2-r6", ".github/workflows/release-recovery-v1.0.0-c2-r6.yml", "manifestSchemaVersion='2.2'", "schema-release-manifest-2.1.json", "schema-release-manifest-2.2.json"},
		"build/release/cmd/release-verifier/main.go": {"release-tooling-v1.0.0-c2-r6", ".github/workflows/release-recovery-v1.0.0-c2-r6.yml", `ManifestSchemaVersion: "2.2"`},
		"internal/verify/release.go": {"release-tooling-v1.0.0-c2", "release-tooling-v1.0.0-c2-r6", ".github/workflows/release-recovery-v1.0.0.yml", ".github/workflows/release-recovery-v1.0.0-c2-r6.yml"},
	} {
		source := read(path)
		for _, token := range required {
			if !strings.Contains(source, token) {
				t.Fatalf("%s does not contain required identity token %q", path, token)
			}
		}
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

func minimumManifestV21() verify.ReleaseManifest {
	manifest := minimumManifestV2()
	manifest.ManifestSchemaVersion = "2.1"
	manifest.ReleaseTooling.Tag = "release-tooling-v1.0.0-c2"
	manifest.ReleaseTooling.WorkflowRef = "refs/tags/release-tooling-v1.0.0-c2"
	manifest.ReleaseIdentity.Ref = manifest.ReleaseTooling.WorkflowRef
	manifest.ReleaseIdentity.CertificateIdentity = "https://github.com/ThameeraDananjaya/project-agnostic-secret-scanner/.github/workflows/release-recovery-v1.0.0.yml@refs/tags/release-tooling-v1.0.0-c2"
	return manifest
}

func minimumManifestV22() verify.ReleaseManifest {
	manifest := minimumManifestV21()
	manifest.ManifestSchemaVersion = "2.2"
	manifest.ReleaseTooling.Tag = "release-tooling-v1.0.0-c2-r6"
	manifest.ReleaseTooling.Workflow = ".github/workflows/release-recovery-v1.0.0-c2-r6.yml"
	manifest.ReleaseTooling.WorkflowRef = "refs/tags/release-tooling-v1.0.0-c2-r6"
	manifest.ReleaseIdentity.Workflow = manifest.ReleaseTooling.Workflow
	manifest.ReleaseIdentity.Ref = manifest.ReleaseTooling.WorkflowRef
	manifest.ReleaseIdentity.CertificateIdentity = "https://github.com/ThameeraDananjaya/project-agnostic-secret-scanner/.github/workflows/release-recovery-v1.0.0-c2-r6.yml@refs/tags/release-tooling-v1.0.0-c2-r6"
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
