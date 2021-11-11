package encoder

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

// encHMAC .
type encHMAC struct {
	opts *options
}

// NewHMAC .
func NewHMAC(opts ...option) encoder {
	var options = new(options)
	for _, opt := range opts {
		opt(options)
	}
	return &encHMAC{opts: options}
}

func (enc *encHMAC) Encode(text string) string {
	hash := hmac.New(sha256.New, enc.opts.secretKey)
	hash.Write([]byte(text))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func (enc *encHMAC) Decode(text string) string {
	return ""
}

func (enc *encHMAC) Description() string {
	return "hmacsha256"
}
