package artifact

import (
	"bytes"
	"encoding/json"
	"io"
)

func rejectDuplicateJSON(raw []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	if duplicateJSONValue(decoder) || decoder.Decode(&struct{}{}) != io.EOF {
		return Rejection{Code: RejectMalformed}
	}
	return nil
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
