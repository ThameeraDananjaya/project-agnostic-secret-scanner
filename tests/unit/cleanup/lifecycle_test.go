package cleanup_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/cleanup"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/outcome"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/workspace"
)

const scanID = "123e4567-e89b-42d3-a456-426614174000"

func TestEveryReturnedLifecycleCleans(t *testing.T) {
	reasons := []outcome.ReasonCode{
		outcome.ReasonPassNoBlockingFindings, outcome.ReasonFailFindingDetected,
		outcome.ReasonIndeterminateIncompleteCoverage, outcome.ReasonIndeterminateTimeout,
		outcome.ReasonTerminalInternalInvariant,
	}
	for index, expected := range reasons {
		base := filepath.Join(t.TempDir(), "workspaces")
		manager, err := workspace.NewManager(base)
		if err != nil {
			t.Fatal(err)
		}
		reason, runErr := cleanup.Run(context.Background(), manager, scanID, index%3+1, func(_ context.Context, root string) outcome.ReasonCode {
			if err := os.WriteFile(filepath.Join(root, "private", "transient"), []byte("synthetic"), 0o600); err != nil {
				t.Fatal(err)
			}
			return expected
		})
		if runErr != nil || reason != expected {
			t.Fatalf("reason=%s err=%v", reason, runErr)
		}
		entries, err := os.ReadDir(base)
		if err != nil || len(entries) != 0 {
			t.Fatalf("lifecycle residue: %v %v", entries, err)
		}
		_ = manager.Close()
	}
}

func TestCancellationAndPanicClean(t *testing.T) {
	base := filepath.Join(t.TempDir(), "workspaces")
	manager, err := workspace.NewManager(base)
	if err != nil {
		t.Fatal(err)
	}
	defer manager.Close()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, _ = cleanup.Run(ctx, manager, scanID, 1, func(ctx context.Context, root string) outcome.ReasonCode {
		return outcome.ReasonIndeterminateTimeout
	})
	_, err = cleanup.Run(context.Background(), manager, scanID, 2, func(context.Context, string) outcome.ReasonCode {
		panic("private crash payload")
	})
	if err == nil {
		t.Fatal("panic was not converted")
	}
	entries, _ := os.ReadDir(base)
	if len(entries) != 0 {
		t.Fatal("cancellation or panic residue remains")
	}
}

func TestCrashRecoveryRemovesOnlyOwnedWorkspace(t *testing.T) {
	base := filepath.Join(t.TempDir(), "workspaces")
	if os.Getenv("PSCAN_CRASH_HELPER") == "1" {
		manager, err := workspace.NewManager(os.Getenv("PSCAN_CRASH_ROOT"))
		if err != nil {
			os.Exit(92)
		}
		space, err := manager.Create(scanID, 1)
		if err != nil {
			os.Exit(93)
		}
		_ = os.WriteFile(filepath.Join(space.Path(), "private", "crash-material"), []byte("synthetic"), 0o600)
		os.Exit(91)
	}
	command := exec.Command(os.Args[0], "-test.run", "^TestCrashRecoveryRemovesOnlyOwnedWorkspace$")
	command.Env = append(os.Environ(), "PSCAN_CRASH_HELPER=1", "PSCAN_CRASH_ROOT="+base)
	if err := command.Run(); err == nil || command.ProcessState.ExitCode() != 91 {
		t.Fatalf("crash helper did not terminate as expected: %v", err)
	}
	restarted, err := workspace.NewManager(base)
	if err != nil {
		t.Fatal(err)
	}
	defer restarted.Close()
	if err := cleanup.Recover(restarted); err != nil {
		t.Fatal(err)
	}
	entries, _ := os.ReadDir(base)
	if len(entries) != 0 {
		t.Fatal("crash residue remains")
	}
}

func TestRestartRecoveryWithoutProcessDeath(t *testing.T) {
	base := filepath.Join(t.TempDir(), "workspaces")
	manager, err := workspace.NewManager(base)
	if err != nil {
		t.Fatal(err)
	}
	space, err := manager.Create(scanID, 1)
	if err != nil {
		t.Fatal(err)
	}
	write := filepath.Join(space.Path(), "private", "crash-material")
	if err := os.WriteFile(write, []byte("synthetic"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := manager.Close(); err != nil {
		t.Fatal(err)
	}
	restarted, err := workspace.NewManager(base)
	if err != nil {
		t.Fatal(err)
	}
	defer restarted.Close()
	if err := cleanup.Recover(restarted); err != nil {
		t.Fatal(err)
	}
	entries, _ := os.ReadDir(base)
	if len(entries) != 0 {
		t.Fatal("crash residue remains")
	}
}

func TestRecoveryPreservesForgedOrUnleasedWorkspace(t *testing.T) {
	base := filepath.Join(t.TempDir(), "workspaces")
	manager, err := workspace.NewManager(base)
	if err != nil {
		t.Fatal(err)
	}
	defer manager.Close()
	forged := filepath.Join(base, "scan-123e4567-e89b-42d3-a456-426614174001-attempt-1")
	if err := os.Mkdir(forged, 0o700); err != nil {
		t.Fatal(err)
	}
	marker := append([]byte("PSCAN_WORKSPACE_V1\n"), make([]byte, 32)...)
	if err := os.WriteFile(filepath.Join(forged, ".pscan-owned"), marker, 0o600); err != nil {
		t.Fatal(err)
	}
	payload := filepath.Join(forged, "must-survive")
	if err := os.WriteFile(payload, []byte("user material"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := cleanup.Recover(manager); err == nil {
		t.Fatal("unleased forged workspace was accepted")
	}
	if _, err := os.Stat(payload); err != nil {
		t.Fatal("forged workspace was deleted")
	}
}
