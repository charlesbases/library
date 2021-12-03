package encoder

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
	"sync"
)

// encHMAC .
type encHMAC struct {
	opts *options
	hash hash.Hash

	lk sync.Mutex
}

// NewHMAC .
func NewHMAC(opts ...option) encoder {
	var options = new(options)
	for _, opt := range opts {
		opt(options)
	}

	return &encHMAC{opts: options, hash: hmac.New(sha256.New, options.secretKey)}
}

func (enc *encHMAC) Encode(text string) string {
	enc.lk.Lock()
	enc.hash.Write([]byte(text))
	res := base64.StdEncoding.EncodeToString(enc.hash.Sum(nil))
	enc.hash.Reset()
	enc.lk.Unlock()
	return res
}

func (enc *encHMAC) Decode(_ string) string {
	return ""
}

func (enc *encHMAC) Description() string {
	return "hmacsha256"
}
