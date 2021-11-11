package encoder

import (
	"crypto/md5"
	"encoding/base64"
)

// encMD5 .
type encMD5 struct {
	opts *options
}

// NewMD5 .
func NewMD5(opts ...option) encoder {
	var options = new(options)
	for _, opt := range opts {
		opt(options)
	}
	return &encMD5{opts: options}
}

func (enc *encMD5) Encode(text string) string {
	hash := md5.New()
	hash.Write([]byte(text))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func (enc *encMD5) Decode(text string) string {
	return ""
}

func (enc *encMD5) Description() string {
	return "md5"
}
