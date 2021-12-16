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
func (*Marshaler) Marshal(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

// Unmarshal .
func (*Marshaler) Unmarshal(data []byte, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}

// String .
func (*Marshaler) String() string {
	return codec.MarshalerType_Proto.String()
}
