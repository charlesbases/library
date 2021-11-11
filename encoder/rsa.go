package encoder

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
)

// encRSA .
type encRSA struct {
	opts       *options
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

// NewRSA .
func NewRSA(opts ...option) (encoder, error) {
	var options = new(options)
	for _, opt := range opts {
		opt(options)
	}
	enc := &encRSA{opts: options}
	{
		block, _ := pem.Decode(enc.opts.publicKey)
		if block == nil {
			return nil, ErrInvalidPublicKey
		}
		key, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		enc.publicKey = key.(*rsa.PublicKey)
	}
	{
		block, _ := pem.Decode(enc.opts.privateKey)
		if block == nil {
			return nil, ErrInvalidPrivateKey
		}
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		enc.privateKey = key
	}
	return enc, nil
}

func (env *encRSA) Encode(text string) string {
	data, _ := rsa.EncryptPKCS1v15(rand.Reader, env.publicKey, []byte(text))
	return base64.StdEncoding.EncodeToString(data)
}

func (env *encRSA) Decode(text string) string {
	data, _ := base64.StdEncoding.DecodeString(text)
	source, _ := rsa.DecryptPKCS1v15(rand.Reader, env.privateKey, []byte(data))
	return string(source)
}

func (env *encRSA) Description() string {
	return "rsa"
}
