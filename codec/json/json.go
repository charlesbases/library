package json

import (
	"encoding/json"

	"library/codec"
)

type Marshaler struct{}

// NewMarshaler .
func NewMarshaler() codec.Marshaler {
	return new(Marshaler)
}

// Marshal .
func (*Marshaler) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal .
func (*Marshaler) Unmarshal(d []byte, v interface{}) error {
	return json.Unmarshal(d, v)
}

// String .
func (*Marshaler) String() string {
	return "json"
}
