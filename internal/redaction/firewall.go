// Package redaction provides a content-free diagnostic boundary for logs,
// errors and lifecycle notices. The authoritative scanner outcome remains the
// scanner-owned outcome schema and is not replaced by this package.
package redaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/outcome"
)

const DiagnosticSchemaVersion = "pscan.diagnostic.v1"

type Diagnostic struct {
	SchemaVersion   string                  `json:"schemaVersion"`
	State           outcome.State           `json:"state"`
	ReasonCode      outcome.ReasonCode      `json:"reasonCode"`
	PermittedAction outcome.PermittedAction `json:"permittedAction"`
}

var allowedReasons = func() map[outcome.ReasonCode]bool {
	result := map[outcome.ReasonCode]bool{}
	for _, reason := range outcome.AllReasonCodes() {
		result[reason] = true
	}
	return result
}()

func DiagnosticFromOutcome(value outcome.Outcome) (Diagnostic, error) {
	if err := value.Validate(); err != nil {
		return Diagnostic{}, errors.New("diagnostic unavailable")
	}
	if _, terminal := outcome.ExitCode(value.State); !terminal || !allowedReasons[value.ReasonCode] {
		return Diagnostic{}, errors.New("diagnostic unavailable")
	}
	return Diagnostic{SchemaVersion: DiagnosticSchemaVersion, State: value.State,
		ReasonCode: value.ReasonCode, PermittedAction: value.PermittedAction}, nil
}

// WriteDiagnostic is deliberately typed: callers cannot pass logs, errors,
// findings, paths, panic payloads or candidate-controlled strings to a
// diagnostic writer.
func WriteDiagnostic(writer io.Writer, value Diagnostic) error {
	if writer == nil || value.SchemaVersion != DiagnosticSchemaVersion || !allowedReasons[value.ReasonCode] {
		return errors.New("diagnostic unavailable")
	}
	expected, ok := outcome.StateForReason(value.ReasonCode)
	action, actionOK := outcome.ActionForState(value.State)
	if !ok || !actionOK || expected != value.State || action != value.PermittedAction {
		return errors.New("diagnostic unavailable")
	}
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(true)
	if err := encoder.Encode(value); err != nil {
		return errors.New("diagnostic unavailable")
	}
	if err := ValidateDiagnostic(buffer.Bytes()); err != nil {
		return err
	}
	if _, err := writer.Write(buffer.Bytes()); err != nil {
		return errors.New("diagnostic unavailable")
	}
	return nil
}

func ValidateDiagnostic(raw []byte) error {
	if hasDuplicateJSONKey(raw) {
		return errors.New("diagnostic unavailable")
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	var value Diagnostic
	if decoder.Decode(&value) != nil || decoder.Decode(&struct{}{}) != io.EOF {
		return errors.New("diagnostic unavailable")
	}
	if value.SchemaVersion != DiagnosticSchemaVersion || !allowedReasons[value.ReasonCode] {
		return errors.New("diagnostic unavailable")
	}
	expected, ok := outcome.StateForReason(value.ReasonCode)
	action, actionOK := outcome.ActionForState(value.State)
	if !ok || !actionOK || value.State != expected || value.PermittedAction != action {
		return errors.New("diagnostic unavailable")
	}
	return nil
}

// Oracle rejects raw, partial, encoded, hashed and fingerprinted canary forms.
// It is an acceptance oracle and never emits the matched material.
func Oracle(surface []byte, canaries ...[]byte) error {
	lower := bytes.ToLower(surface)
	for _, canary := range canaries {
		if len(canary) == 0 {
			continue
		}
		variants := [][]byte{canary, bytes.ToLower(canary),
			[]byte(base64.StdEncoding.EncodeToString(canary)),
			[]byte(base64.RawStdEncoding.EncodeToString(canary)),
			[]byte(base64.URLEncoding.EncodeToString(canary)),
			[]byte(hex.EncodeToString(canary))}
		sum := sha256.Sum256(canary)
		variants = append(variants, []byte(hex.EncodeToString(sum[:])))
		for _, variant := range variants {
			if len(variant) > 0 && bytes.Contains(lower, bytes.ToLower(variant)) {
				return errors.New("redaction proof failed")
			}
		}
		// Four-byte substrings catch prefix/suffix and partial-value leakage
		// without making one-character JSON tokens unusable.
		if len(canary) >= 4 {
			for i := 0; i+4 <= len(canary); i++ {
				if bytes.Contains(lower, bytes.ToLower(canary[i:i+4])) {
					return errors.New("redaction proof failed")
				}
			}
		}
	}
	return nil
}

func OracleMetadata(surface []byte) error {
	if ValidateDiagnostic(surface) != nil {
		return errors.New("redaction proof failed")
	}
	return nil
}

func hasDuplicateJSONKey(raw []byte) bool {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	if duplicateJSONValue(decoder) || decoder.Decode(&struct{}{}) != io.EOF {
		return true
	}
	return false
}

func duplicateJSONValue(decoder *json.Decoder) bool {
	token, err := decoder.Token()
	if err != nil {
		return true
	}
	delimiter, ok := token.(json.Delim)
	if !ok {
		return false
	}
	switch delimiter {
	case '{':
		seen := map[string]bool{}
		for decoder.More() {
			keyToken, err := decoder.Token()
			key, keyOK := keyToken.(string)
			if err != nil || !keyOK || seen[key] {
				return true
			}
			seen[key] = true
			if duplicateJSONValue(decoder) {
				return true
			}
		}
		end, err := decoder.Token()
		return err != nil || end != json.Delim('}')
	case '[':
		for decoder.More() {
			if duplicateJSONValue(decoder) {
				return true
			}
		}
		end, err := decoder.Token()
		return err != nil || end != json.Delim(']')
	default:
		return true
	}
}
