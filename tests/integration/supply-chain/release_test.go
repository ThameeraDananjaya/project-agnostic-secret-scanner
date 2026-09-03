package supplychain_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/verify"
)

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
	manifest.ManifestSchemaVersion = "2.0"
	raw, _ = json.Marshal(manifest)
	if _, err := verify.ParseReleaseManifest(raw); !errors.Is(err, verify.ErrInvalidReference) {
		t.Fatalf("unknown schema major was not rejected: %v", err)
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

func releaseIdentity() verify.ReleaseIdentity {
	return verify.ReleaseIdentity{Repository: "ThameeraDananjaya/project-agnostic-secret-scanner", RepositoryOwnerID: 50274860, Workflow: ".github/workflows/release.yml", Ref: "refs/tags/v1.0.0", OIDCIssuer: "https://token.actions.githubusercontent.com", CertificateIdentity: "https://github.com/ThameeraDananjaya/project-agnostic-secret-scanner/.github/workflows/release.yml@refs/tags/v1.0.0"}
}
