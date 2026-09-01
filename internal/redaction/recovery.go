package redaction

import (
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/outcome"
)

// Recover converts a panic or internal diagnostic into a fixed content-free
// reason. The panic value is never formatted, logged or returned.
func Recover(reason *outcome.ReasonCode) {
	if recover() != nil {
		*reason = outcome.ReasonTerminalInternalInvariant
	}
}
