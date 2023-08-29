package proto

import (
	"errors"

	"github.com/golang/protobuf/proto"

	"github.com/charlesbases/library/codec"
	"github.com/charlesbases/library/content"
)

var ErrInvalidType = errors.New("proto: not implemented")

const mess = "[ProtoMessage]"

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

func (*c) Marshal(v interface{}) ([]byte, error) {
	if pv, ok := v.(proto.Message); ok {
		return proto.Marshal(pv)
	} else {
		return nil, ErrInvalidType
	}
}

func (*c) Unmarshal(data []byte, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}

func (c *c) RawMessage(data []byte) string {
	return mess
}

func (c *c) ContentType() content.Type {
	return content.Proto
}
