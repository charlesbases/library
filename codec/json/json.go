package json

import (
	"encoding/json"

	"github.com/charlesbases/library/codec"
)

// Marshaler default codec.Marshaler
var Marshaler = NewMarshaler(func(o *codec.MarshalOptions) { o.Indent = true })

type c struct {
	*codec.MarshalOptions
}

// NewMarshaler .
func NewMarshaler(opts ...func(o *codec.MarshalOptions)) codec.Marshaler {
	var options = new(codec.MarshalOptions)
	for _, opt := range opts {
		opt(options)
	}

	return &c{MarshalOptions: options}
}

// Marshal .
func (c *c) Marshal(v interface{}) ([]byte, error) {
	if c.Indent {
		return json.MarshalIndent(v, "", "  ")
	}
	return json.Marshal(v)
}

// Unmarshal .
func (c *c) Unmarshal(d []byte, v interface{}) error {
	return json.Unmarshal(d, v)
}
