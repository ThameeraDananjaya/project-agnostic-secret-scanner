package redaction_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/engine"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/outcome"
)

func TestPrivateDecoderPanicIsContentFree(t *testing.T) {
	if os.Getenv("PSCAN_ENGINE_PANIC_HELPER") == "1" {
		os.Exit(0)
	}
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(executable)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(raw)
	result := engine.RunPrivate(context.Background(), engine.Command{
		Executable: executable, ExpectedDigest: hex.EncodeToString(sum[:]),
		Args: []string{"-test.run", "^TestPrivateDecoderPanicIsContentFree$"},
		Dir:  filepath.Dir(executable), Environment: append(os.Environ(), "PSCAN_ENGINE_PANIC_HELPER=1"),
		Timeout: 10 * time.Second, CaptureLimit: 4096,
	}, func(engine.PrivateOutput) outcome.ReasonCode {
		panic("PSCAN_PRIVATE_PANIC_CANARY")
	})
	if result.Reason != outcome.ReasonTerminalInternalInvariant || result.ExitCode != -1 {
		t.Fatalf("panic did not fail content-free: %#v", result)
	}
}
