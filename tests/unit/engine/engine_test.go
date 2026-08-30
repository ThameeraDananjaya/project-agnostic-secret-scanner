package engine_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/engine"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/outcome"
)

const syntheticCanary = "PSCAN_SYNTHETIC_CANARY_NOT_A_CREDENTIAL_83e08b"

func TestEngineHelper(t *testing.T) {
	if os.Getenv("PSCAN_ENGINE_HELPER") != "1" {
		return
	}
	args := os.Args
	mode := ""
	for i := range args {
		if args[i] == "--" && i+1 < len(args) {
			mode = args[i+1]
		}
	}
	switch mode {
	case "canary":
		fmt.Print(syntheticCanary)
	case "overflow":
		fmt.Print(strings.Repeat("x", 8192))
	case "touch":
		_ = os.WriteFile(os.Getenv("PSCAN_TOUCH_PATH"), []byte("executed"), 0o600)
	}
	os.Exit(0)
}

func TestPrivateOutputDoesNotCrossResultBoundary(t *testing.T) {
	executable, digest := selfBinding(t)
	result := engine.RunPrivate(context.Background(), engine.Command{
		Executable: executable, ExpectedDigest: digest,
		Args:        []string{"-test.run=TestEngineHelper", "--", "canary"},
		Environment: append(os.Environ(), "PSCAN_ENGINE_HELPER=1"),
		Timeout:     5 * time.Second,
	}, func(private engine.PrivateOutput) outcome.ReasonCode {
		if string(private.Stdout) != syntheticCanary {
			t.Fatal("private decoder did not receive exact synthetic output")
		}
		return outcome.ReasonFailFindingDetected
	})
	if result.Reason != outcome.ReasonFailFindingDetected {
		t.Fatalf("unexpected result: %#v", result)
	}
	if strings.Contains(fmt.Sprintf("%#v", result), syntheticCanary) {
		t.Fatal("candidate content crossed the public result boundary")
	}
}

func TestDigestMismatchFailsBeforeExecution(t *testing.T) {
	executable, digest := selfBinding(t)
	marker := filepath.Join(t.TempDir(), "must-not-exist")
	digest = strings.Repeat("0", 64)
	result := engine.RunPrivate(context.Background(), engine.Command{
		Executable: executable, ExpectedDigest: digest,
		Args:        []string{"-test.run=TestEngineHelper", "--", "touch"},
		Environment: append(os.Environ(), "PSCAN_ENGINE_HELPER=1", "PSCAN_TOUCH_PATH="+marker),
		Timeout:     time.Second,
	}, func(engine.PrivateOutput) outcome.ReasonCode { return outcome.ReasonPassNoBlockingFindings })
	if result.Reason != outcome.ReasonFailScannerIntegrity {
		t.Fatalf("mismatch did not fail closed: %#v", result)
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatal("engine executed after its digest failed")
	}
	_, _ = executable, digest
}

func TestCaptureLimitIsExplicitNonPass(t *testing.T) {
	executable, digest := selfBinding(t)
	result := engine.RunPrivate(context.Background(), engine.Command{
		Executable: executable, ExpectedDigest: digest,
		Args:        []string{"-test.run=TestEngineHelper", "--", "overflow"},
		Environment: append(os.Environ(), "PSCAN_ENGINE_HELPER=1"),
		Timeout:     5 * time.Second, CaptureLimit: 128,
	}, func(engine.PrivateOutput) outcome.ReasonCode { return outcome.ReasonPassNoBlockingFindings })
	if result.Reason != outcome.ReasonIndeterminateResourceLimit {
		t.Fatalf("overflow was not explicit: %#v", result)
	}
}

func selfBinding(t *testing.T) (string, string) {
	t.Helper()
	path, err := filepath.Abs(os.Args[0])
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(raw)
	return path, hex.EncodeToString(sum[:])
}
