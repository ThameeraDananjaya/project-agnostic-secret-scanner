package outcome

import (
	"bytes"
	"encoding/json"
	"io"
)

func Serialize(w io.Writer, value Outcome) error {
	if err := value.Validate(); err != nil {
		return err
	}
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(true)
	if err := encoder.Encode(value); err != nil {
		return err
	}
	_, err := w.Write(buffer.Bytes())
	return err
}
