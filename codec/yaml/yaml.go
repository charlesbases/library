package yaml

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/charlesbases/library/codec"
	"github.com/charlesbases/library/content"
)

// defaultConfigurationFilePath 默认配置文件路径
const defaultConfigurationFilePath = "config.yaml"

// Marshaler default codec.Marshaler
var Marshaler = NewMarshaler()

type c struct {
	*codec.DecodeOptions
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

// NewDecoder .
func NewDecoder(opts ...func(o *codec.DecodeOptions)) codec.Decoder {
	var options = &codec.DecodeOptions{FileName: defaultConfigurationFilePath}
	for _, opt := range opts {
		opt(options)
	}
	return &c{DecodeOptions: options}
}

func (c *c) Decode(v interface{}) error {
	if c.Reader != nil {
		return yaml.NewDecoder(c.Reader).Decode(v)
	} else {
		file, err := os.Open(c.FileName)
		if err != nil {
			return err
		}

		err = yaml.NewDecoder(file).Decode(v)
		file.Close()
		return err
	}
}

func (c *c) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func (c *c) Unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

func (c *c) RawMessage(data []byte) string {
	return string(data)
}

func (c *c) ContentType() content.Type {
	return content.Yaml
}
