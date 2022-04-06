package proto

import (
	"library/codec"

	"github.com/golang/protobuf/proto"
)

type Marshaler struct{}

// NewMarshaler .
func NewMarshaler() codec.Marshaler {
	return new(Marshaler)
}

// Marshal .
func (m *Marshaler) Marshal(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

// Unmarshal .
func (m *Marshaler) Unmarshal(data []byte, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}

// ContentType .
func (m *Marshaler) ContentType() codec.ContentType {
	return codec.ContentTypeProto
}
