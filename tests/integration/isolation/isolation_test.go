package isolation_test

import (
	"sync"
	"testing"
	"time"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/policy"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/verify"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/workspace"
)

type acceptingVerifier struct{}

func (acceptingVerifier) Verify(_ []byte, _ verify.DetachedSignature) error { return nil }

func TestParallelAndSequentialAttemptsDoNotShareWorkspaceState(t *testing.T) {
	manager, err := workspace.NewManager(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer manager.Close()

	ids := []string{
		"00000000-0000-4000-8000-000000000011",
		"00000000-0000-4000-8000-000000000022",
	}
	paths := make(chan string, len(ids))
	errors := make(chan error, len(ids))
	var group sync.WaitGroup
	for _, id := range ids {
		group.Add(1)
		go func(scanID string) {
			defer group.Done()
			created, createErr := manager.Create(scanID, 1)
			if createErr != nil {
				errors <- createErr
				return
			}
			paths <- created.Path()
			errors <- created.Cleanup()
		}(id)
	}
	group.Wait()
	close(errors)
	for err := range errors {
		if err != nil {
			t.Fatal(err)
		}
	}
	first, second := <-paths, <-paths
	if first == second {
		t.Fatal("parallel scans shared a workspace path")
	}

	for attempt := 1; attempt <= 2; attempt++ {
		created, err := manager.Create("00000000-0000-4000-8000-000000000033", attempt)
		if err != nil {
			t.Fatal(err)
		}
		if err := created.Cleanup(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestParallelPolicyProjectionsDoNotCrossScopes(t *testing.T) {
	now := time.Date(2026, 9, 2, 0, 0, 0, 0, time.UTC)
	raw := [][]byte{
		[]byte(`{"blockedClasses":["synthetic-class-a"],"minimumSeverity":"critical","maxExceptionDays":7,"allowedExceptionEvidence":[]}`),
		[]byte(`{"blockedClasses":["synthetic-class-b"],"minimumSeverity":"critical","maxExceptionDays":7,"allowedExceptionEvidence":[]}`),
	}
	results := make(chan policy.PolicyProjection, 2)
	errors := make(chan error, 2)
	var group sync.WaitGroup
	for _, document := range raw {
		group.Add(1)
		go func(value []byte) {
			defer group.Done()
			binding := verify.DocumentBinding{
				SchemaFamily: "synthetic-policy", SchemaVersion: "1.0", AdapterVersion: policy.PolicyAdapterVersion,
				Digest: verify.DigestBytes(value), Signature: verify.DetachedSignature{TrustDomain: "project-policy", Algorithm: "ed25519", KeyID: "synthetic-key"},
			}
			projection, err := policy.LoadPolicy(value, binding, verify.FamilyWindow{Family: "synthetic-policy", CurrentMajor: 1, Minors: []verify.MinorCompatibility{{Minor: 0}}}, acceptingVerifier{}, now)
			results <- projection
			errors <- err
		}(document)
	}
	group.Wait()
	close(errors)
	for err := range errors {
		if err != nil {
			t.Fatal(err)
		}
	}
	first, second := <-results, <-results
	if first.Digest == second.Digest {
		t.Fatal("distinct transient policy inputs shared a projection identity")
	}
}
