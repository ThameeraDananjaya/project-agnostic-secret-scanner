package gitleaks_test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/engine/gitleaks"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/outcome"
)

func TestExactAdapterContract(t *testing.T) {
	if err := gitleaks.ValidateBinding("gitleaks", "8.30.1", "2.0.0", "json-v8.30.1"); err != nil {
		t.Fatal(err)
	}
	invalid := [][]string{
		{"other", "8.30.1", "2.0.0", "json-v8.30.1"},
		{"gitleaks", "latest", "2.0.0", "json-v8.30.1"},
		{"gitleaks", "8.30.1", "1.1.0", "json-v8.30.1"},
		{"gitleaks", "8.30.1", "2.0.0", "json"},
	}
	for _, binding := range invalid {
		if gitleaks.ValidateBinding(binding[0], binding[1], binding[2], binding[3]) == nil {
			t.Fatalf("unsupported binding accepted: %v", binding)
		}
	}
}

func TestPinnedRuleAndCoverageIntegrityBindings(t *testing.T) {
	root := filepath.Clean(filepath.Join("..", "..", ".."))
	raw, err := os.ReadFile(filepath.Join(root, "build", "gitleaks", "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest struct {
		Engine         string `json:"engine"`
		Version        string `json:"version"`
		AdapterVersion string `json:"adapterVersion"`
		OutputBinding  string `json:"outputBinding"`
		Rules          struct {
			ProductConfig       string `json:"productConfig"`
			ProductConfigDigest string `json:"productConfigSha256"`
			CoverageRuleID      string `json:"coverageRuleId"`
			IgnoreFile          string `json:"ignoreFile"`
			IgnoreFileDigest    string `json:"ignoreFileSha256"`
		} `json:"rules"`
	}
	if err := json.Unmarshal(raw, &manifest); err != nil {
		t.Fatal(err)
	}
	if err := gitleaks.ValidateBinding(manifest.Engine, manifest.Version, manifest.AdapterVersion, manifest.OutputBinding); err != nil {
		t.Fatal(err)
	}
	if manifest.Rules.CoverageRuleID != gitleaks.CoverageRuleID {
		t.Fatal("coverage rule binding mismatch")
	}
	for path, expected := range map[string]string{
		manifest.Rules.ProductConfig: manifest.Rules.ProductConfigDigest,
		manifest.Rules.IgnoreFile:    manifest.Rules.IgnoreFileDigest,
	} {
		content, readErr := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
		if readErr != nil {
			t.Fatal(readErr)
		}
		digest := sha256.Sum256(content)
		if hex.EncodeToString(digest[:]) != expected {
			t.Fatalf("pinned rule asset digest mismatch: %s", path)
		}
	}
	config, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(manifest.Rules.ProductConfig)))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(config), `id = "`+gitleaks.CoverageRuleID+`"`) != 1 {
		t.Fatal("coverage rule is missing or ambiguous in pinned product config")
	}
}

func TestNativePatchGitModeCannotPass(t *testing.T) {
	base, head := strings.Repeat("a", 40), strings.Repeat("b", 40)
	result := (gitleaks.Adapter{}).ScanGitRange(t.Context(), t.TempDir(), base, head, false, time.Second, 1<<20)
	if result.Reason != outcome.ReasonIndeterminateIncompleteCoverage {
		t.Fatalf("native patch mode did not fail closed: %#v", result)
	}
}

func TestRangeContractDoesNotCreateOptionSyntax(t *testing.T) {
	base := strings.Repeat("a", 40)
	head := strings.Repeat("b", 40)
	rangeSpec, err := gitleaks.RangeSpec(base, head, false)
	if err != nil || rangeSpec != base+".."+head {
		t.Fatalf("wrong exact range: %q (%v)", rangeSpec, err)
	}
	if _, err := gitleaks.RangeSpec("--all", head, false); err == nil {
		t.Fatal("candidate option syntax accepted as commit identity")
	}
	if _, err := gitleaks.RangeSpec("", head, true); err != nil {
		t.Fatalf("first-release history rejected: %v", err)
	}
}
