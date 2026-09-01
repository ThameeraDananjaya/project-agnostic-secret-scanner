package engine_test

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/engine/gitleaks"
)

func TestPinnedRulePackHasFiniteExactMaximumSpan(t *testing.T) {
	root := filepath.Clean(filepath.Join("..", "..", ".."))
	path := filepath.Join(root, "rules", "generic", "gitleaks-v8.30.1.toml")
	path, err := filepath.Abs(path)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(raw)
	proof, err := gitleaks.ProveRuleSpans(path, hex.EncodeToString(digest[:]))
	if err != nil {
		t.Fatal(err)
	}
	if proof.RuleCount != gitleaks.RequiredRuleCount || proof.Maximum != 4020 {
		t.Fatalf("unexpected span proof: %#v", proof)
	}
}
