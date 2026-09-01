package artifact_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/artifact"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/cleanup"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/outcome"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/workspace"
)

func TestArtifactLifecycleLeavesNoTransientMaterial(t *testing.T) {
	input := buildZIP(t, []byte("synthetic lifecycle input"))
	for attempt, reason := range []outcome.ReasonCode{
		outcome.ReasonPassNoBlockingFindings,
		outcome.ReasonFailFindingDetected,
		outcome.ReasonIndeterminateIncompleteCoverage,
	} {
		base := filepath.Join(t.TempDir(), "workspaces")
		manager, err := workspace.NewManager(base)
		if err != nil {
			t.Fatal(err)
		}
		got, err := cleanup.Run(context.Background(), manager, "123e4567-e89b-42d3-a456-426614174000", attempt+1, func(ctx context.Context, root string) outcome.ReasonCode {
			_, normalizeErr := artifact.Normalize(ctx, input, filepath.Join(root, "private", "artifact"), artifact.TestLimits(5, 100, 8<<20, 4<<20, 1000, time.Minute))
			if normalizeErr != nil {
				t.Fatal(normalizeErr)
			}
			return reason
		})
		if err != nil || got != reason {
			t.Fatalf("reason=%s err=%v", got, err)
		}
		entries, readErr := os.ReadDir(base)
		if readErr != nil || len(entries) != 0 {
			t.Fatalf("transient residue: %v %v", entries, readErr)
		}
		_ = manager.Close()
	}
}
