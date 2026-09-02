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
		RequestSchemaVersion string   `json:"requestSchemaVersion"`
		RequiredFeatures     []string `json:"requiredFeatures"`
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
	if len(header.RequiredFeatures) > 0 {
		return ScanRequest{}, reject(outcome.ReasonIndeterminateSchemaUnsupported)
	}
	if err := validateRequiredPresence(data, version); err != nil {
		return ScanRequest{}, reject(outcome.ReasonFailInputIntegrity)
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

func validateRequiredPresence(data []byte, version Version) error {
	top, err := requiredObject(data,
		"requestSchemaVersion", "scanId", "mode", "scannerReleaseDigest", "engineBinding",
		"rulePackDigest", "policyDigest", "allowlistDigest", "sourceBinding",
		"trackedSourceManifest", "fallbackRequirement", "limits", "offlineRequired",
		"redactionMode", "requestedAt")
	if err != nil {
		return err
	}
	if raw, ok := top["requiredFeatures"]; ok && !isJSONArray(raw) {
		return fmt.Errorf("requiredFeatures must be an array")
	}
	if !isJSONBoolean(top["offlineRequired"]) {
		return fmt.Errorf("offlineRequired must be a boolean")
	}
	if _, err := requiredObject(top["engineBinding"], "name", "version", "binaryDigest", "adapterVersion"); err != nil {
		return err
	}
	source, err := requiredObject(top["sourceBinding"], "headCommit", "historyRangeDigest", "trackedTreeDigest")
	if err != nil {
		return err
	}
	if _, err := requiredObject(top["trackedSourceManifest"], "path", "digest"); err != nil {
		return err
	}
	if _, err := requiredObject(top["fallbackRequirement"], "mode"); err != nil {
		return err
	}
	if _, err := requiredObject(top["limits"], "timeoutSeconds", "maxArchiveDepth", "maxArchiveEntries", "maxExpandedBytes", "maxFileBytes", "maxCompressionRatio", "maxMemoryBytes", "maxCpuPercent"); err != nil {
		return err
	}
	if version.Major == 1 && version.Minor >= 1 && !version.FutureMinor() {
		if err := requireFields(top, "policyBinding", "allowlistBinding", "proofBindings"); err != nil {
			return err
		}
		for _, field := range []string{"policyBinding", "allowlistBinding"} {
			binding, err := requiredObject(top[field], "schemaFamily", "schemaVersion", "adapterVersion", "digest", "signature")
			if err != nil {
				return err
			}
			if _, err := requiredObject(binding["signature"], "trustDomain", "algorithm", "keyId", "value"); err != nil {
				return err
			}
		}
		if _, err := requiredObject(top["proofBindings"], "admissionLedgerSchemaVersion", "admissionLedgerDigest", "rawClassifierVersion", "rawClassifierDigest", "preparationVersion", "preparationDigest", "resourceProfileId", "inspectionProofFormat", "inspectionProofDigest"); err != nil {
			return err
		}
	}
	if raw, ok := top["buildContextManifest"]; ok {
		if _, err := requiredObject(raw, "path", "digest"); err != nil {
			return err
		}
	}
	if raw, ok := top["artifactManifest"]; ok {
		manifest, err := requiredObject(raw, "digest", "entries")
		if err != nil {
			return err
		}
		var entries []json.RawMessage
		if !isJSONArray(manifest["entries"]) || json.Unmarshal(manifest["entries"], &entries) != nil {
			return fmt.Errorf("artifact entries must be an array")
		}
		for _, entry := range entries {
			if _, err := requiredObject(entry, "path", "type", "size", "digest"); err != nil {
				return err
			}
		}
	}
	var mode string
	if err := json.Unmarshal(top["mode"], &mode); err != nil {
		return err
	}
	if mode == "pr" {
		if err := requireFields(source, "baseCommit", "mergeBase"); err != nil {
			return err
		}
	}
	if mode == "release" {
		if err := requireFields(top, "buildContextManifest", "artifactManifest"); err != nil {
			return err
		}
		if err := requireFields(source, "firstRelease"); err != nil || !isJSONBoolean(source["firstRelease"]) {
			return fmt.Errorf("firstRelease must be present as a boolean")
		}
		var first bool
		_ = json.Unmarshal(source["firstRelease"], &first)
		if !first {
			if err := requireFields(source, "baseCommit"); err != nil {
				return err
			}
		}
	}
	return nil
}

func requiredObject(data []byte, fields ...string) (map[string]json.RawMessage, error) {
	var object map[string]json.RawMessage
	if len(bytes.TrimSpace(data)) == 0 || bytes.Equal(bytes.TrimSpace(data), []byte("null")) || json.Unmarshal(data, &object) != nil || object == nil {
		return nil, fmt.Errorf("required object is invalid")
	}
	if err := requireFields(object, fields...); err != nil {
		return nil, err
	}
	return object, nil
}

func requireFields(object map[string]json.RawMessage, fields ...string) error {
	for _, field := range fields {
		raw, ok := object[field]
		if !ok || len(bytes.TrimSpace(raw)) == 0 || bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
			return fmt.Errorf("required field %s is absent", field)
		}
	}
	return nil
}

func isJSONBoolean(raw []byte) bool {
	trimmed := bytes.TrimSpace(raw)
	return bytes.Equal(trimmed, []byte("true")) || bytes.Equal(trimmed, []byte("false"))
}

func isJSONArray(raw []byte) bool {
	trimmed := bytes.TrimSpace(raw)
	return len(trimmed) > 0 && trimmed[0] == '['
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
