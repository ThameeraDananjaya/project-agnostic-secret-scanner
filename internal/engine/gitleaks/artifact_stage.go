package gitleaks

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/artifact"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/cleanup"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/engine"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/outcome"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/workspace"
)

// ScanArtifact runs the bounded PSCAN-04 internal artifact stage under one
// attempt deadline. It does not alter the public CLI or claim consumer wiring.
func (a Adapter) ScanArtifact(parent context.Context, manager *workspace.Manager, scanID string, attempt int, input string, limits artifact.Limits) engine.Result {
	if parent == nil || manager == nil || !limits.Valid() || !filepath.IsAbs(input) {
		return engine.Result{Reason: outcome.ReasonFailInputIntegrity, ExitCode: -1}
	}
	ctx, cancel := context.WithTimeout(parent, limits.Timeout)
	defer cancel()
	result := engine.Result{Reason: outcome.ReasonTerminalInternalInvariant, ExitCode: -1}
	reason, err := cleanup.Run(ctx, manager, scanID, attempt, func(attemptContext context.Context, root string) outcome.ReasonCode {
		privateHome := filepath.Join(root, "private", "engine-home")
		if err := os.Mkdir(privateHome, 0o700); err != nil {
			return outcome.ReasonUnavailableWorkspace
		}
		projection, err := artifact.Normalize(attemptContext, input, filepath.Join(root, "input", "normalized"), limits)
		if err != nil {
			return artifactReason(err)
		}
		bound := a
		bound.Binding.PrivateHome = privateHome
		bound.Binding.Environment = engine.SafeEnvironment(filepath.Dir(bound.Binding.Executable), privateHome)
		result = bound.ScanArtifactProjection(attemptContext, projection)
		return result.Reason
	})
	if err != nil {
		return engine.Result{Reason: outcome.ReasonTerminalInternalInvariant, ExitCode: -1}
	}
	if result.Reason == outcome.ReasonTerminalInternalInvariant {
		result.Reason = reason
	}
	return result
}

func artifactReason(err error) outcome.ReasonCode {
	var rejection artifact.Rejection
	if !errors.As(err, &rejection) {
		return outcome.ReasonTerminalInternalInvariant
	}
	switch rejection.Code {
	case artifact.RejectResource:
		return outcome.ReasonIndeterminateResourceLimit
	case artifact.RejectTimeout:
		return outcome.ReasonIndeterminateTimeout
	case artifact.RejectCancelled:
		return outcome.ReasonIndeterminateIncompleteCoverage
	case artifact.RejectUnsupported:
		return outcome.ReasonIndeterminateUnsupportedInput
	case artifact.RejectUnsafe, artifact.RejectMalformed:
		return outcome.ReasonFailInputIntegrity
	default:
		return outcome.ReasonTerminalInternalInvariant
	}
}
