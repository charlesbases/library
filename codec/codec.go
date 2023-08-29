package codec

import (
	"io"

	"github.com/charlesbases/library/content"
)

type Decoder interface {
	Decode(v interface{}) error
}

type Marshaler interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
	RawMessage(data []byte) string
	ContentType() content.Type
}

// DecodeOptions 从文件或 io.Reader 中加载
type DecodeOptions struct {
	FileName string
	Reader   io.Reader
}

// MarshalOptions .
type MarshalOptions struct {
	// Indent 格式化输出
	Indent bool
}
