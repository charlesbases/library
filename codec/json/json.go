package json

import (
	"encoding/json"

	"library/codec"
)

type Marshaler struct {
	indent bool
}

// NewMarshaler .
func NewMarshaler() codec.Marshaler {
	return &Marshaler{indent: false}
}

// NewMarshalerIndent .
func NewMarshalerIndent() codec.Marshaler {
	return &Marshaler{indent: true}
}

// Marshal .
func (m *Marshaler) Marshal(v interface{}) ([]byte, error) {
	if m.indent {
		return json.MarshalIndent(v, "", "  ")
	}
	return json.Marshal(v)
}

// Unmarshal .
func (m *Marshaler) Unmarshal(d []byte, v interface{}) error {
	return json.Unmarshal(d, v)
}

// ContentType .
func (m *Marshaler) ContentType() codec.ContentType {
	return codec.ContentTypeJson
}
