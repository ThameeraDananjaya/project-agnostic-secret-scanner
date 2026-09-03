package supplychain_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/verify"
)

type exactSyntheticBundle struct {
	manifestDigest string
	bundle         []byte
}

func (v exactSyntheticBundle) VerifyManifest(_ context.Context, manifestPath, bundlePath string, _ verify.ReleaseIdentity) error {
	manifest, manifestErr := os.ReadFile(manifestPath)
	bundle, bundleErr := os.ReadFile(bundlePath)
	if manifestErr != nil || bundleErr != nil || verify.DigestBytes(manifest) != v.manifestDigest || string(bundle) != string(v.bundle) {
		return verify.ErrInvalidSignature
	}
	return nil
}

func TestReleaseVerificationRejectsEveryRequiredMutationBeforeScan(t *testing.T) {
	directory, request := writeRelease(t)
	if _, err := verify.VerifyRelease(context.Background(), request); err != nil {
		t.Fatalf("valid synthetic release rejected: %v", err)
	}

	cases := map[string]func(string, *verify.ReleaseVerificationRequest){
		"one-byte runner": func(root string, _ *verify.ReleaseVerificationRequest) {
			mutate(t, filepath.Join(root, "runner-linux"))
		},
		"licence": func(root string, _ *verify.ReleaseVerificationRequest) { mutate(t, filepath.Join(root, "LICENSE.txt")) },
		"sbom": func(root string, _ *verify.ReleaseVerificationRequest) {
			mutate(t, filepath.Join(root, "sbom.spdx.json"))
		},
		"revocation": func(root string, _ *verify.ReleaseVerificationRequest) {
			mutate(t, filepath.Join(root, "revocations.json"))
		},
		"bundle": func(root string, _ *verify.ReleaseVerificationRequest) {
			mutate(t, filepath.Join(root, "release-manifest.sigstore.json"))
		},
		"identity": func(_ string, request *verify.ReleaseVerificationRequest) {
			request.Policy.CertificateIdentity += "-wrong"
		},
		"issuer": func(_ string, request *verify.ReleaseVerificationRequest) {
			request.Policy.OIDCIssuer = "https://example.invalid"
		},
		"tag": func(_ string, request *verify.ReleaseVerificationRequest) { request.Policy.Ref = "refs/tags/v1.0.1" },
	}
	for name, change := range cases {
		t.Run(name, func(t *testing.T) {
			copyRoot := t.TempDir()
			copyFiles(t, directory, copyRoot)
			candidate := request
			candidate.Directory = copyRoot
			change(copyRoot, &candidate)
			if _, err := verify.VerifyRelease(context.Background(), candidate); err == nil {
				t.Fatal("mutation was admitted")
			}
		})
	}
}

func writeRelease(t *testing.T) (string, verify.ReleaseVerificationRequest) {
	t.Helper()
	root := t.TempDir()
	contents := map[string][]byte{
		"runner-linux": []byte("L"), "runner-windows.exe": []byte("W"), "engine-linux": []byte("E"), "engine-windows.exe": []byte("F"),
		"rules.toml": []byte("R"), "schema.json": []byte("S"), "tests.json": []byte("T"), "limitations.md": []byte("M"),
		"LICENSE.txt": []byte("I"), "sbom.spdx.json": []byte("B"),
		"revocations.json": []byte(`{"schemaFamily":"global-scanner-revocation-snapshot","schemaVersion":"1.0","capturedAt":"2026-09-03T00:00:00Z","records":[]}`),
		"checkpoint.json":  []byte(`{"schemaFamily":"global-scanner-revocation-checkpoint","schemaVersion":"1.0","sequence":0,"digest":null,"capturedAt":"2026-09-03T00:00:00Z","discoveryLocation":"https://api.github.com/repos/ThameeraDananjaya/project-agnostic-secret-scanner/releases?per_page=100"}`),
	}
	assets := make([]verify.ReleaseAsset, 0, len(contents))
	for path, raw := range contents {
		kind := "documentation"
		switch path {
		case "runner-linux", "runner-windows.exe":
			kind = "runner"
		case "engine-linux", "engine-windows.exe":
			kind = "engine"
		case "rules.toml":
			kind = "rules"
		case "schema.json":
			kind = "schema"
		case "tests.json":
			kind = "test-summary"
		case "limitations.md":
			kind = "limitations"
		case "LICENSE.txt":
			kind = "licence"
		case "sbom.spdx.json":
			kind = "sbom"
		case "revocations.json":
			kind = "revocation-snapshot"
		case "checkpoint.json":
			kind = "revocation-checkpoint"
		}
		if err := os.WriteFile(filepath.Join(root, path), raw, 0o600); err != nil {
			t.Fatal(err)
		}
		assets = append(assets, verify.ReleaseAsset{Path: path, Kind: kind, OS: "none", Arch: "none", Size: int64(len(raw)), SHA256: verify.DigestBytes(raw)})
	}
	lookup := func(path string) string {
		for _, asset := range assets {
			if asset.Path == path {
				return asset.SHA256
			}
		}
		t.Fatal(path)
		return ""
	}
	manifest := verify.ReleaseManifest{
		SchemaFamily: "scanner-release-manifest", ManifestSchemaVersion: "1.1", ReleaseVersion: "v1.0.0", SourceRevision: "0123456789abcdef0123456789abcdef01234567", SourceTree: "89abcdef0123456789abcdef0123456789abcdef", RunnerVersion: "1.0.0", GoToolchainVersion: "go1.27.1",
		RunnerBindings: []verify.PlatformDigest{{OS: "linux", Arch: "amd64", Path: "runner-linux", SHA256: lookup("runner-linux")}, {OS: "windows", Arch: "amd64", Path: "runner-windows.exe", SHA256: lookup("runner-windows.exe")}},
		EngineBindings: []verify.EngineBinding{{Name: "gitleaks", Version: "8.30.1", SourceRevision: "83d9cd684c87d95d656c1458ef04895a7f1cbd8e", OS: "linux", Arch: "amd64", Path: "engine-linux", SHA256: lookup("engine-linux")}, {Name: "gitleaks", Version: "8.30.1", SourceRevision: "83d9cd684c87d95d656c1458ef04895a7f1cbd8e", OS: "windows", Arch: "amd64", Path: "engine-windows.exe", SHA256: lookup("engine-windows.exe")}},
		RulePack:       verify.NamedDigest{Path: "rules.toml", SHA256: lookup("rules.toml")}, SchemaBindings: []verify.SchemaBinding{{Family: "scan-request", Version: "1.1", Path: "schema.json", SHA256: lookup("schema.json")}},
		ReleaseIdentity: verify.ReleaseIdentity{Repository: "ThameeraDananjaya/project-agnostic-secret-scanner", RepositoryOwnerID: 50274860, Workflow: ".github/workflows/release.yml", Ref: "refs/tags/v1.0.0", OIDCIssuer: "https://token.actions.githubusercontent.com", CertificateIdentity: "https://github.com/ThameeraDananjaya/project-agnostic-secret-scanner/.github/workflows/release.yml@refs/tags/v1.0.0"},
		Compatibility:   verify.ReleaseCompatibility{MinimumRunnerVersion: "1.0.0", SupportedOperatingSystems: []string{"linux", "windows"}, SupportedArchitectures: []string{"amd64"}, TestSummaryAsset: "tests.json", LimitationsAsset: "limitations.md"},
		Revocation:      verify.ReleaseRevocation{DiscoveryLocation: "https://api.github.com/repos/ThameeraDananjaya/project-agnostic-secret-scanner/releases?per_page=100", SnapshotAsset: "revocations.json", CheckpointAsset: "checkpoint.json", SchemaVersion: "1.1", MaximumSnapshotAgeHours: 24}, Assets: assets, CreatedAt: "2026-09-03T00:00:00Z",
	}
	raw, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "release-manifest.json"), raw, 0o600); err != nil {
		t.Fatal(err)
	}
	bundle := []byte("synthetic-bundle")
	if err := os.WriteFile(filepath.Join(root, "release-manifest.sigstore.json"), bundle, 0o600); err != nil {
		t.Fatal(err)
	}
	request := verify.ReleaseVerificationRequest{Directory: root, ManifestPath: "release-manifest.json", BundlePath: "release-manifest.sigstore.json", Policy: verify.ReleaseTrustPolicy{Repository: "ThameeraDananjaya/project-agnostic-secret-scanner", RepositoryOwnerID: 50274860, Workflow: ".github/workflows/release.yml", Ref: "refs/tags/v1.0.0", OIDCIssuer: "https://token.actions.githubusercontent.com", CertificateIdentity: "https://github.com/ThameeraDananjaya/project-agnostic-secret-scanner/.github/workflows/release.yml@refs/tags/v1.0.0", ReleaseVersion: "v1.0.0"}, Signature: exactSyntheticBundle{manifestDigest: verify.DigestBytes(raw), bundle: bundle}, Now: time.Date(2026, 9, 3, 1, 0, 0, 0, time.UTC)}
	return root, request
}

func mutate(t *testing.T, path string) {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	raw[0] ^= 1
	if err := os.WriteFile(path, raw, 0o600); err != nil {
		t.Fatal(err)
	}
}
func copyFiles(t *testing.T, source, destination string) {
	t.Helper()
	entries, err := os.ReadDir(source)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		raw, err := os.ReadFile(filepath.Join(source, entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(destination, entry.Name()), raw, 0o600); err != nil {
			t.Fatal(err)
		}
	}
}
