package json

import (
	"encoding/json"

	"github.com/charlesbases/library/codec"
	"github.com/charlesbases/library/content"
)

// Marshaler default codec.Marshaler
var Marshaler = NewMarshaler()

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

func (c *c) Marshal(v interface{}) ([]byte, error) {
	if c.Indent {
		return json.MarshalIndent(v, "", "  ")
	}
	return json.Marshal(v)
}

func (c *c) Unmarshal(d []byte, v interface{}) error {
	return json.Unmarshal(d, v)
}

func (c *c) ShowMessage(data []byte) string {
	return string(data)
}

func (c *c) ContentType() content.Type {
	return content.Json
}
