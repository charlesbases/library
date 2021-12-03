package encoder

import (
	"crypto/md5"
	"encoding/base64"
	"hash"
	"sync"
)

// encMD5 .
type encMD5 struct {
	opts *options
	hash hash.Hash

	lk sync.Mutex
}

// NewMD5 .
func NewMD5(opts ...option) encoder {
	var options = new(options)
	for _, opt := range opts {
		opt(options)
	}
	return &encMD5{opts: options, hash: md5.New()}
}

func (enc *encMD5) Encode(text string) string {
	enc.lk.Lock()
	enc.hash.Write([]byte(text))
	res := base64.StdEncoding.EncodeToString(enc.hash.Sum(nil))
	enc.hash.Reset()
	enc.lk.Unlock()
	return res
}

func (enc *encMD5) Decode(_ string) string {
	return ""
}

func (enc *encMD5) Description() string {
	return "md5"
}
