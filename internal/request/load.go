package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/outcome"
)

const MaxPayloadBytes int64 = 1 << 20

type ValidationError struct {
	Reason outcome.ReasonCode
}

func (e *ValidationError) Error() string { return "request rejected" }

func reject(reason outcome.ReasonCode) error { return &ValidationError{Reason: reason} }

func ReasonFor(err error, fallback outcome.ReasonCode) outcome.ReasonCode {
	var validation *ValidationError
	if errors.As(err, &validation) {
		return validation.Reason
	}
	return fallback
}

func Load(reader io.Reader) (ScanRequest, error) {
	data, err := io.ReadAll(io.LimitReader(reader, MaxPayloadBytes+1))
	if err != nil || int64(len(data)) > MaxPayloadBytes || len(bytes.TrimSpace(data)) == 0 {
		return ScanRequest{}, reject(outcome.ReasonFailInputIntegrity)
	}
	if err := rejectDuplicateKeys(data); err != nil {
		return ScanRequest{}, reject(outcome.ReasonFailInputIntegrity)
	}
	var header struct {
		RequestSchemaVersion string `json:"requestSchemaVersion"`
	}
	if err := json.Unmarshal(data, &header); err != nil {
		return ScanRequest{}, reject(outcome.ReasonFailInputIntegrity)
	}
	version, err := ParseVersion(header.RequestSchemaVersion)
	if err != nil {
		return ScanRequest{}, reject(outcome.ReasonIndeterminateSchemaUnsupported)
	}
	if !version.Supported() {
		return ScanRequest{}, reject(outcome.ReasonIndeterminateSchemaUnsupported)
	}
	var value ScanRequest
	decoder := json.NewDecoder(bytes.NewReader(data))
	if !version.FutureMinor() {
		decoder.DisallowUnknownFields()
	}
	if err := decoder.Decode(&value); err != nil {
		return ScanRequest{}, reject(outcome.ReasonFailInputIntegrity)
	}
	if decoder.Decode(&struct{}{}) != io.EOF {
		return ScanRequest{}, reject(outcome.ReasonFailInputIntegrity)
	}
	if len(value.RequiredFeatures) > 0 {
		return ScanRequest{}, reject(outcome.ReasonIndeterminateSchemaUnsupported)
	}
	return value, nil
}

func rejectDuplicateKeys(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	var walk func() error
	walk = func() error {
		token, err := decoder.Token()
		if err != nil {
			return err
		}
		delim, ok := token.(json.Delim)
		if !ok {
			return nil
		}
		switch delim {
		case '{':
			seen := map[string]struct{}{}
			for decoder.More() {
				keyToken, err := decoder.Token()
				if err != nil {
					return err
				}
				key, ok := keyToken.(string)
				if !ok {
					return fmt.Errorf("object key is not a string")
				}
				if _, exists := seen[key]; exists {
					return fmt.Errorf("duplicate key")
				}
				seen[key] = struct{}{}
				if err := walk(); err != nil {
					return err
				}
			}
			_, err = decoder.Token()
			return err
		case '[':
			for decoder.More() {
				if err := walk(); err != nil {
					return err
				}
			}
			_, err = decoder.Token()
			return err
		default:
			return fmt.Errorf("unexpected delimiter")
		}
	}
	if err := walk(); err != nil {
		return err
	}
	if decoder.Decode(&struct{}{}) != io.EOF {
		return fmt.Errorf("trailing data")
	}
	return nil
}
