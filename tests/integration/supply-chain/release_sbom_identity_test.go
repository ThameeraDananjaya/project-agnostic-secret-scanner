package supplychain_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/verify"
)

func minimumManifestV25() verify.ReleaseManifest {
	m := minimumManifestV24()
	m.ManifestSchemaVersion = "2.5"
	m.ReleaseTooling.Tag = verify.UnsignedBuildTagV25
	m.ReleaseTooling.WorkflowRef = "refs/tags/" + verify.UnsignedBuildTagV25
	m.BuildIdentity.Ref = m.ReleaseTooling.WorkflowRef
	return m
}
func signerPolicyV12Bytes(t *testing.T, m verify.ReleaseManifest) []byte {
	t.Helper()
	return bytes.Replace(signerPolicyBytes(t, m), []byte(`"pscan-separate-signer-policy-v1"`), []byte(`"pscan-separate-signer-policy-v1.2"`), 1)
}
func TestV25ClosedAdditiveSchemaProjections(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	for _, pair := range [][2]string{{"contracts/release-manifest/schema-2.4.json", "contracts/release-manifest/schema-2.5.json"}, {"contracts/release-signer-policy/schema-1.1.json", "contracts/release-signer-policy/schema-1.2.json"}} {
		old, e := os.ReadFile(filepath.Join(root, pair[0]))
		if e != nil {
			t.Fatal(e)
		}
		current, e := os.ReadFile(filepath.Join(root, pair[1]))
		if e != nil {
			t.Fatal(e)
		}
		current = bytes.ReplaceAll(current, []byte(verify.UnsignedBuildTagV25), []byte(verify.UnsignedBuildTagV24))
		if strings.Contains(pair[0], "release-manifest") {
			current = bytes.ReplaceAll(current, []byte("2.5"), []byte("2.4"))
		} else {
			current = bytes.ReplaceAll(current, []byte("release-signer-policy:1.2"), []byte("release-signer-policy:1.1"))
			current = bytes.ReplaceAll(current, []byte("pscan-separate-signer-policy-v1.2"), []byte("pscan-separate-signer-policy-v1.1"))
		}
		if !bytes.Equal(old, current) {
			t.Fatal("new schema changed beyond selected identity/version", pair[1])
		}
	}
}
func TestV25RejectsManifestTagSourceAndUnsignedMixtures(t *testing.T) {
	for name, mutate := range map[string]func(*verify.ReleaseManifest){
		"old build-v2 tag": func(m *verify.ReleaseManifest) { m.ReleaseTooling.Tag = verify.UnsignedBuildTagV24 },
		"old tooling ref": func(m *verify.ReleaseManifest) {
			m.ReleaseTooling.WorkflowRef = "refs/tags/" + verify.UnsignedBuildTagV24
		},
		"old build ref":                 func(m *verify.ReleaseManifest) { m.BuildIdentity.Ref = "refs/tags/" + verify.UnsignedBuildTagV24 },
		"old manifest 2.4":              func(m *verify.ReleaseManifest) { m.ManifestSchemaVersion = "2.4" },
		"old manifest 2.3":              func(m *verify.ReleaseManifest) { m.ManifestSchemaVersion = "2.3" },
		"unknown manifest":              func(m *verify.ReleaseManifest) { m.ManifestSchemaVersion = "2.6" },
		"product revision":              func(m *verify.ReleaseManifest) { m.ProductSource.Commit = strings.Repeat("b", 40) },
		"product tree":                  func(m *verify.ReleaseManifest) { m.ProductSource.Tree = strings.Repeat("b", 40) },
		"build source disagreement":     func(m *verify.ReleaseManifest) { m.BuildIdentity.WorkflowSHA = strings.Repeat("b", 40) },
		"tooling workflow disagreement": func(m *verify.ReleaseManifest) { m.ReleaseTooling.WorkflowSHA = strings.Repeat("b", 40) },
		"missing build":                 func(m *verify.ReleaseManifest) { m.BuildIdentity = nil },
		"builder masquerades signer":    func(m *verify.ReleaseManifest) { m.ReleaseIdentity.Workflow = verify.UnsignedBuildWorkflow },
		"invented trusted state":        func(m *verify.ReleaseManifest) { m.ReleaseState = "verified" },
	} {
		t.Run(name, func(t *testing.T) {
			m := minimumManifestV25()
			mutate(&m)
			raw, e := json.Marshal(m)
			if e != nil {
				t.Fatal(e)
			}
			if _, e = verify.ParseReleaseManifest(raw); e == nil {
				t.Fatal("mixed or unbound identity admitted")
			}
		})
	}
}
func TestV25PolicyVersionMatrixAndIndependentTrust(t *testing.T) {
	cases := []struct {
		schema   string
		manifest func() verify.ReleaseManifest
		policy   func(*testing.T, verify.ReleaseManifest) []byte
	}{{"2.3", minimumManifestV23, signerPolicyBytes}, {"2.4", minimumManifestV24, signerPolicyV11Bytes}, {"2.5", minimumManifestV25, signerPolicyV12Bytes}}
	candidate := t.TempDir()
	path := filepath.Join(t.TempDir(), "INERT-policy.json")
	load := func(raw []byte, digest, commit, tree string) (verify.ReleaseTrustPolicy, error) {
		t.Helper()
		if e := os.WriteFile(path, raw, 0600); e != nil {
			t.Fatal(e)
		}
		return verify.LoadSeparateSignerPolicy(path, digest, candidate, commit, tree)
	}
	for i, c := range cases {
		for j, p := range cases {
			t.Run(c.schema+"-policy-for-"+p.schema, func(t *testing.T) {
				m := c.manifest()
				raw := p.policy(t, m)
				got, e := load(raw, verify.DigestBytes(raw), m.ReleaseTooling.Commit, m.ReleaseTooling.Tree)
				if i == j {
					if e != nil || got.ManifestSchemaVersion != c.schema {
						t.Fatalf("matching pair rejected: %v", e)
					}
				} else if e == nil {
					t.Fatal("mixed policy version admitted")
				}
			})
		}
	}
	m := minimumManifestV25()
	good := signerPolicyV12Bytes(t, m)
	for name, raw := range map[string][]byte{
		"unknown field":    bytes.Replace(good, []byte(`"schema":`), []byte(`"unknown":0,"schema":`), 1),
		"duplicate field":  bytes.Replace(good, []byte(`"schema":`), []byte(`"schema":"duplicate","schema":`), 1),
		"unknown policy":   bytes.Replace(good, []byte("policy-v1.2"), []byte("policy-v1.3"), 1),
		"broadened signer": bytes.ReplaceAll(good, []byte(verify.ReleaseSignerWorkflow), []byte(".github/workflows/other.yml")),
		"old builder tag":  bytes.ReplaceAll(good, []byte(verify.UnsignedBuildTagV25), []byte(verify.UnsignedBuildTagV24)),
	} {
		t.Run(name, func(t *testing.T) {
			if _, e := load(raw, verify.DigestBytes(raw), m.ReleaseTooling.Commit, m.ReleaseTooling.Tree); e == nil {
				t.Fatal("untrusted policy admitted")
			}
		})
	}
	for name, args := range map[string][3]string{
		"digest":          {strings.Repeat("b", 64), m.ReleaseTooling.Commit, m.ReleaseTooling.Tree},
		"compiled commit": {verify.DigestBytes(good), strings.Repeat("b", 40), m.ReleaseTooling.Tree},
		"compiled tree":   {verify.DigestBytes(good), m.ReleaseTooling.Commit, strings.Repeat("b", 40)},
	} {
		t.Run(name, func(t *testing.T) {
			if _, e := load(good, args[0], args[1], args[2]); e == nil {
				t.Fatal("unbound policy admitted")
			}
		})
	}
	local := filepath.Join(candidate, "INERT-policy.json")
	if e := os.WriteFile(local, good, 0600); e != nil {
		t.Fatal(e)
	}
	if _, e := verify.LoadSeparateSignerPolicy(local, verify.DigestBytes(good), candidate, m.ReleaseTooling.Commit, m.ReleaseTooling.Tree); e == nil {
		t.Fatal("candidate-local policy admitted")
	}
}
func TestV25OuterManifestFinalizationPreservesArchivesAndVerificationGates(t *testing.T) {
	exerciseOuterManifestFinalization(t, minimumManifestV25(), signerPolicyV12Bytes)
}
