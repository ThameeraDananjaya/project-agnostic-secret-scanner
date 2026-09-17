package supplychain_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/verify"
)

func minimumManifestV23() verify.ReleaseManifest {
	m := minimumManifestV22()
	m.ManifestSchemaVersion = "2.3"
	m.ReleaseTooling.Tag = verify.UnsignedBuildTag
	m.ReleaseTooling.Workflow = verify.UnsignedBuildWorkflow
	m.ReleaseTooling.WorkflowRef = "refs/tags/" + verify.UnsignedBuildTag
	m.BuildIdentity = &verify.BuildIdentity{Repository: m.ReleaseIdentity.Repository, RepositoryOwnerID: 50274860, Workflow: verify.UnsignedBuildWorkflow, Ref: m.ReleaseTooling.WorkflowRef, WorkflowSHA: m.ReleaseTooling.Commit, Trigger: "workflow_dispatch"}
	m.ReleaseState = "signing-pending"
	m.ReleaseIdentity.Workflow = verify.ReleaseSignerWorkflow
	m.ReleaseIdentity.Ref = verify.ReleaseSignerRef
	m.ReleaseIdentity.WorkflowSHA = strings.Repeat("a", 40)
	m.ReleaseIdentity.CertificateIdentity = "https://github.com/" + m.ReleaseIdentity.Repository + "/" + verify.ReleaseSignerWorkflow + "@" + verify.ReleaseSignerRef
	return m
}

func v23Raw(t *testing.T, m verify.ReleaseManifest) []byte {
	t.Helper()
	raw, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	if m.ReleaseState == "unsigned-candidate" {
		var fields map[string]json.RawMessage
		if err = json.Unmarshal(raw, &fields); err != nil {
			t.Fatal(err)
		}
		fields["releaseIdentity"] = json.RawMessage("null")
		raw, err = json.Marshal(fields)
		if err != nil {
			t.Fatal(err)
		}
	}
	return raw
}

func TestSeparateBuildSignerIdentityAndUnsignedStructure(t *testing.T) {
	m := minimumManifestV23()
	if _, err := verify.ParseReleaseManifest(v23Raw(t, m)); err != nil {
		t.Fatal(err)
	}
	m.ReleaseState = "unsigned-candidate"
	m.ReleaseIdentity = verify.ReleaseIdentity{}
	if _, err := verify.ParseReleaseManifest(v23Raw(t, m)); err != nil {
		t.Fatal(err)
	}
	parsed, err := verify.ParseReleaseManifest(v23Raw(t, m))
	if err != nil {
		t.Fatal(err)
	}
	roundtrip, err := json.Marshal(parsed)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = verify.ParseReleaseManifest(roundtrip); err != nil {
		t.Fatalf("unsigned structural roundtrip lost null signer: %v", err)
	}
	badNull := strings.Replace(string(v23Raw(t, m)), `"releaseIdentity":null`, `"releaseIdentity":{}`, 1)
	if _, err := verify.ParseReleaseManifest([]byte(badNull)); err == nil {
		t.Fatal("unsigned identity object admitted")
	}
	for name, mutate := range map[string]func(*verify.ReleaseManifest){
		"wrong build sha":             func(m *verify.ReleaseManifest) { m.BuildIdentity.WorkflowSHA = strings.Repeat("b", 40) },
		"old tooling tag":             func(m *verify.ReleaseManifest) { m.ReleaseTooling.Tag = "release-tooling-v1.0.0-c2-r6" },
		"signer is builder":           func(m *verify.ReleaseManifest) { m.ReleaseIdentity.Workflow = verify.UnsignedBuildWorkflow },
		"wrong signer ref":            func(m *verify.ReleaseManifest) { m.ReleaseIdentity.Ref = m.ReleaseTooling.WorkflowRef },
		"invented verified state":     func(m *verify.ReleaseManifest) { m.ReleaseState = "verified" },
		"missing build":               func(m *verify.ReleaseManifest) { m.BuildIdentity = nil },
		"old version with new fields": func(m *verify.ReleaseManifest) { m.ManifestSchemaVersion = "2.2" },
	} {
		t.Run(name, func(t *testing.T) {
			m := minimumManifestV23()
			mutate(&m)
			if _, err := verify.ParseReleaseManifest(v23Raw(t, m)); err == nil {
				t.Fatal("cross-identity accepted")
			}
		})
	}
}

func signerPolicyBytes(t *testing.T, m verify.ReleaseManifest) []byte {
	t.Helper()
	raw, err := json.Marshal(map[string]any{"schema": "pscan-separate-signer-policy-v1", "productSource": m.ProductSource, "releaseTooling": m.ReleaseTooling, "buildIdentity": m.BuildIdentity, "signer": m.ReleaseIdentity})
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func TestIndependentSignerPolicyCannotChangeCompiledBuildOrFixedSigner(t *testing.T) {
	m := minimumManifestV23()
	candidate := t.TempDir()
	path := filepath.Join(t.TempDir(), "owner-policy.json")
	good := signerPolicyBytes(t, m)
	load := func(raw []byte, digest, commit, tree string) (verify.ReleaseTrustPolicy, error) {
		t.Helper()
		if err := os.WriteFile(path, raw, 0600); err != nil {
			t.Fatal(err)
		}
		return verify.LoadSeparateSignerPolicy(path, digest, candidate, commit, tree)
	}
	p, err := load(good, verify.DigestBytes(good), m.ReleaseTooling.Commit, m.ReleaseTooling.Tree)
	if err != nil || p.WorkflowSHA != m.ReleaseIdentity.WorkflowSHA || p.ReleaseToolingCommit != m.ReleaseTooling.Commit {
		t.Fatalf("valid independent policy rejected: %v", err)
	}
	for _, raw := range [][]byte{nil, []byte("null"), append(append([]byte{}, good...), []byte("{}")...), []byte(strings.Replace(string(good), `"schema":`, `"unknown":0,"schema":`, 1)), []byte(strings.Replace(string(good), `"schema":`, `"schema":"duplicate","schema":`, 1))} {
		if _, err := load(raw, verify.DigestBytes(raw), m.ReleaseTooling.Commit, m.ReleaseTooling.Tree); err == nil {
			t.Fatal("malformed policy admitted")
		}
	}
	for _, field := range []string{"productSource", "releaseTooling", "buildIdentity", "signer", "schema"} {
		var record map[string]any
		json.Unmarshal(good, &record)
		delete(record, field)
		raw, _ := json.Marshal(record)
		if _, err := load(raw, verify.DigestBytes(raw), m.ReleaseTooling.Commit, m.ReleaseTooling.Tree); err == nil {
			t.Fatalf("missing %s admitted", field)
		}
	}
	for _, mutate := range []func(*verify.ReleaseManifest){func(m *verify.ReleaseManifest) { m.ReleaseIdentity.Workflow = ".github/workflows/other.yml" }, func(m *verify.ReleaseManifest) { m.ReleaseIdentity.RepositoryOwnerID = 1 }, func(m *verify.ReleaseManifest) { m.ReleaseIdentity.OIDCIssuer = "https://other.invalid" }, func(m *verify.ReleaseManifest) { m.ProductSource.Commit = strings.Repeat("b", 40) }} {
		changed := minimumManifestV23()
		mutate(&changed)
		raw := signerPolicyBytes(t, changed)
		if _, err := load(raw, verify.DigestBytes(raw), m.ReleaseTooling.Commit, m.ReleaseTooling.Tree); err == nil {
			t.Fatal("broadened policy admitted")
		}
	}
	for _, args := range [][3]string{{strings.Repeat("b", 64), m.ReleaseTooling.Commit, m.ReleaseTooling.Tree}, {verify.DigestBytes(good), strings.Repeat("b", 40), m.ReleaseTooling.Tree}, {verify.DigestBytes(good), m.ReleaseTooling.Commit, strings.Repeat("b", 40)}} {
		if _, err := load(good, args[0], args[1], args[2]); err == nil {
			t.Fatal("wrong independent digest/build admitted")
		}
	}
	local := filepath.Join(candidate, "policy.json")
	os.WriteFile(local, good, 0600)
	if _, err := verify.LoadSeparateSignerPolicy(local, verify.DigestBytes(good), candidate, m.ReleaseTooling.Commit, m.ReleaseTooling.Tree); err == nil {
		t.Fatal("candidate-local policy admitted")
	}
}

type separateSignerProbe struct{ calls int }

func (p *separateSignerProbe) VerifyManifest(context.Context, string, string, verify.ReleaseIdentity) error {
	p.calls++
	return errors.New("synthetic provider always rejects")
}

func TestUnsignedWithBundleNeverCallsSignatureAndSeparateSignerIsRequired(t *testing.T) {
	m := minimumManifestV23()
	root := t.TempDir()
	ownerPath := filepath.Join(t.TempDir(), "policy.json")
	raw := signerPolicyBytes(t, m)
	os.WriteFile(ownerPath, raw, 0600)
	policy, err := verify.LoadSeparateSignerPolicy(ownerPath, verify.DigestBytes(raw), root, m.ReleaseTooling.Commit, m.ReleaseTooling.Tree)
	if err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(root, "bundle.json"), []byte("synthetic"), 0600)
	for _, mode := range []string{"unsigned", "wrong signer sha", "correct pending"} {
		t.Run(mode, func(t *testing.T) {
			candidate := minimumManifestV23()
			p := policy
			expectedCalls := 0
			if mode == "unsigned" {
				candidate.ReleaseState = "unsigned-candidate"
				candidate.ReleaseIdentity = verify.ReleaseIdentity{}
			} else if mode == "wrong signer sha" {
				p.WorkflowSHA = candidate.ReleaseTooling.Commit
			} else {
				expectedCalls = 1
			}
			os.WriteFile(filepath.Join(root, "manifest.json"), v23Raw(t, candidate), 0600)
			probe := &separateSignerProbe{}
			_, err := verify.VerifyRelease(context.Background(), verify.ReleaseVerificationRequest{Directory: root, ManifestPath: "manifest.json", BundlePath: "bundle.json", Policy: p, Signature: probe, Now: time.Date(2026, 9, 17, 0, 0, 0, 0, time.UTC)})
			if err == nil || probe.calls != expectedCalls {
				t.Fatalf("unexpected trust or provider calls: %v, %d", err, probe.calls)
			}
		})
	}
}
