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
	return new(Marshaler)
}

// Marshal .
func (m *Marshaler) Marshal(v interface{}, options ...codec.MarshalOption) ([]byte, error) {
	var opts = new(codec.MarshalOptions)
	for _, o := range options {
		o(opts)
	}

	if opts.Indent {
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
