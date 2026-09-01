package redaction_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/outcome"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/redaction"
)

func TestEveryTerminalSurfaceIsOneContentFreeDocument(t *testing.T) {
	canaries := [][]byte{
		[]byte("zq7w9x3v5b1n8m4"), []byte("/k6j2h9g4/f8d3s7"),
		[]byte("v4c8x2z6q9w3e7"), []byte("n5m1b7v3c9x6z2"),
		[]byte("r8t4y1u7i3o9p5"), []byte("a6s2d8f4g1h7j3"),
	}
	at := time.Unix(0, 0).UTC()
	for _, reason := range outcome.AllReasonCodes() {
		value, err := outcome.NewTerminal("123e4567-e89b-42d3-a456-426614174000", "unknown", reason, 1, at, at)
		if err != nil {
			t.Fatal(err)
		}
		diagnostic, err := redaction.DiagnosticFromOutcome(value)
		if err != nil {
			t.Fatal(err)
		}
		var surface bytes.Buffer
		if err := redaction.WriteDiagnostic(&surface, diagnostic); err != nil {
			t.Fatal(err)
		}
		if bytes.Count(surface.Bytes(), []byte("\n")) != 1 || redaction.ValidateDiagnostic(surface.Bytes()) != nil || redaction.OracleMetadata(surface.Bytes()) != nil {
			t.Fatalf("invalid public document for %s", reason)
		}
		if err := redaction.Oracle(surface.Bytes(), canaries...); err != nil {
			t.Fatalf("canary escaped for %s", reason)
		}
		if bytes.Contains(surface.Bytes(), []byte(fmt.Sprintf("%x", canaries[0]))) {
			t.Fatal("encoded canary escaped")
		}
	}
}

func TestCandidateControlledOutputCannotReachWriter(t *testing.T) {
	for _, raw := range [][]byte{
		[]byte(`{"schemaVersion":"pscan.diagnostic.v1","state":"FAIL","reasonCode":"FAIL_FINDING_DETECTED","permittedAction":"CREATE_CORRECTED_CANDIDATE","message":"injected"}`),
		[]byte("{}\n{}\n"),
		[]byte("panic: PSCAN_SECRET_VALUE_d54e78\n"),
	} {
		if redaction.ValidateDiagnostic(raw) == nil {
			t.Fatalf("candidate output admitted: %q", raw)
		}
	}
}
