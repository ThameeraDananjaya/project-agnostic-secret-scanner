package redaction_test

import (
	"bytes"
	"encoding/base64"
	"errors"
	"testing"
	"time"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/outcome"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/redaction"
)

func TestFirewallEmitsOnlyContentFreeSchema(t *testing.T) {
	start := time.Unix(0, 0).UTC()
	value, err := outcome.NewTerminal("123e4567-e89b-42d3-a456-426614174000", "release", outcome.ReasonFailFindingDetected, 1, start, start)
	if err != nil {
		t.Fatal(err)
	}
	diagnostic, err := redaction.DiagnosticFromOutcome(value)
	if err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := redaction.WriteDiagnostic(&output, diagnostic); err != nil {
		t.Fatal(err)
	}
	if err := redaction.ValidateDiagnostic(output.Bytes()); err != nil {
		t.Fatal(err)
	}
	if err := redaction.OracleMetadata(output.Bytes()); err != nil {
		t.Fatal(err)
	}
}

func TestCanaryOracleCatchesEveryDerivedForm(t *testing.T) {
	canary := []byte("PSCAN-CANARY-7fd284")
	forms := [][]byte{canary, canary[:8], []byte(base64.StdEncoding.EncodeToString(canary))}
	for _, form := range forms {
		if redaction.Oracle(form, canary) == nil {
			t.Fatalf("oracle missed %q", form)
		}
	}
	if redaction.Oracle([]byte(`{"state":"FAIL"}`), canary) != nil {
		t.Fatal("oracle rejected content-free output")
	}
}

func TestOutputInjectionCannotAddFieldsOrFrames(t *testing.T) {
	attacks := [][]byte{
		[]byte("{}\n{}\n"),
		[]byte(`{"schemaVersion":"pscan.diagnostic.v1","state":"FAIL","reasonCode":"FAIL_FINDING_DETECTED","permittedAction":"CREATE_CORRECTED_CANDIDATE","path":"injected"}`),
		[]byte(`{"schemaVersion":"pscan.diagnostic.v1","state":"FAIL\nPASS","reasonCode":"FAIL_FINDING_DETECTED","permittedAction":"CREATE_CORRECTED_CANDIDATE"}`),
		[]byte(`{"schemaVersion":"pscan.diagnostic.v1","state":"FAIL","state":"PASS","reasonCode":"FAIL_FINDING_DETECTED","permittedAction":"CREATE_CORRECTED_CANDIDATE"}`),
	}
	for _, attack := range attacks {
		if redaction.ValidateDiagnostic(attack) == nil {
			t.Fatalf("accepted injection %q", attack)
		}
	}
}

func TestPanicPayloadNeverEscapes(t *testing.T) {
	reason := outcome.ReasonPassNoBlockingFindings
	func() {
		defer redaction.Recover(&reason)
		panic("PSCAN-PANIC-CANARY")
	}()
	if reason != outcome.ReasonTerminalInternalInvariant {
		t.Fatal("panic did not fail closed")
	}
}

type hostileWriter struct{}

func (hostileWriter) Write([]byte) (int, error) {
	return 0, errors.New("zq7w9x3v5b1n8m4")
}

func TestWriterErrorIsFixedAndContentFree(t *testing.T) {
	value := redaction.Diagnostic{SchemaVersion: redaction.DiagnosticSchemaVersion,
		State: outcome.StateFail, ReasonCode: outcome.ReasonFailFindingDetected,
		PermittedAction: outcome.ActionCreateCorrectedCandidate}
	err := redaction.WriteDiagnostic(hostileWriter{}, value)
	if err == nil || err.Error() != "diagnostic unavailable" {
		t.Fatalf("writer error escaped: %v", err)
	}
}
