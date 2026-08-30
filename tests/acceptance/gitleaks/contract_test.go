package gitleaks_test

import (
	"strings"
	"testing"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/engine/gitleaks"
)

func TestExactAdapterContract(t *testing.T) {
	if err := gitleaks.ValidateBinding("gitleaks", "8.30.1", "1.0.0", "json-v8.30.1"); err != nil {
		t.Fatal(err)
	}
	invalid := [][]string{
		{"other", "8.30.1", "1.0.0", "json-v8.30.1"},
		{"gitleaks", "latest", "1.0.0", "json-v8.30.1"},
		{"gitleaks", "8.30.1", "2.0.0", "json-v8.30.1"},
		{"gitleaks", "8.30.1", "1.0.0", "json"},
	}
	for _, binding := range invalid {
		if gitleaks.ValidateBinding(binding[0], binding[1], binding[2], binding[3]) == nil {
			t.Fatalf("unsupported binding accepted: %v", binding)
		}
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
