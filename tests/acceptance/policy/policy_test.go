package policy_test

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/verify"
)

type acceptingVerifier struct{}

func (acceptingVerifier) Verify(_ []byte, _ verify.DetachedSignature) error { return nil }

func TestDivergentRevocationFixtureFailsClosed(t *testing.T) {
	fixture := readFixture(t, "evidence", "divergent-chain.json")
	var document struct {
		Anchor  verify.Checkpoint    `json:"anchor"`
		Records []verify.ChainRecord `json:"records"`
	}
	if err := json.Unmarshal(fixture, &document); err != nil {
		t.Fatal(err)
	}
	window := verify.FamilyWindow{Family: "synthetic-global-revocations", CurrentMajor: 1, Minors: []verify.MinorCompatibility{{Minor: 0}}}
	_, err := verify.VerifyChain(document.Records, document.Anchor, window, acceptingVerifier{}, time.Date(2026, 9, 2, 0, 0, 0, 0, time.UTC))
	if !errors.Is(err, verify.ErrEvidenceConflict) {
		t.Fatalf("divergent chain did not fail closed: %v", err)
	}
}

func TestAdversarialCorpusContainsNoPrivateKeyMaterial(t *testing.T) {
	root := readFixtureRoot(t)
	paths := []string{
		filepath.Join(root, "policy", "invalid-unknown-major.json"),
		filepath.Join(root, "allowlist", "expired-exception.json"),
		filepath.Join(root, "evidence", "divergent-chain.json"),
		filepath.Join(root, "isolation", "README.md"),
	}
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		for _, forbidden := range [][]byte{[]byte("BEGIN PRIVATE KEY"), []byte("BEGIN OPENSSH PRIVATE KEY"), []byte("PRIVATE KEY-----")} {
			if contains(data, forbidden) {
				t.Fatalf("private-key marker found in synthetic fixture %s", path)
			}
		}
	}
}

func readFixture(t *testing.T, group, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(readFixtureRoot(t), group, name))
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func readFixtureRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("caller path unavailable")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", "..", "fixtures", "adversarial"))
}

func contains(data, target []byte) bool {
	if len(target) == 0 || len(target) > len(data) {
		return false
	}
	for index := 0; index <= len(data)-len(target); index++ {
		if string(data[index:index+len(target)]) == string(target) {
			return true
		}
	}
	return false
}
