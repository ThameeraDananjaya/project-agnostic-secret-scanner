// Package cleanup guarantees removal of transient material for every returned
// lifecycle result and supports marker-bound recovery after process death.
package cleanup

import (
	"context"
	"errors"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/outcome"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/workspace"
)

type Operation func(context.Context, string) outcome.ReasonCode

func Run(ctx context.Context, manager *workspace.Manager, scanID string, attempt int, operation Operation) (reason outcome.ReasonCode, err error) {
	if ctx == nil || manager == nil || operation == nil {
		return outcome.ReasonTerminalInternalInvariant, errors.New("lifecycle unavailable")
	}
	space, err := manager.Create(scanID, attempt)
	if err != nil {
		return outcome.ReasonUnavailableWorkspace, errors.New("lifecycle unavailable")
	}
	reason = outcome.ReasonTerminalInternalInvariant
	defer func() {
		if recover() != nil {
			reason = outcome.ReasonTerminalInternalInvariant
			err = errors.New("lifecycle unavailable")
		}
		if cleanupErr := space.Cleanup(); cleanupErr != nil {
			reason = outcome.ReasonTerminalInternalInvariant
			err = errors.New("lifecycle unavailable")
		}
	}()
	reason = operation(ctx, space.Path())
	if _, known := outcome.StateForReason(reason); !known {
		return outcome.ReasonTerminalInternalInvariant, errors.New("lifecycle unavailable")
	}
	return reason, nil
}

func Recover(manager *workspace.Manager) error {
	if manager == nil {
		return errors.New("lifecycle unavailable")
	}
	if err := manager.Recover(); err != nil {
		return errors.New("lifecycle unavailable")
	}
	return nil
}
