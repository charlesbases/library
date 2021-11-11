package encoder

import (
	"errors"
)

var (
	ErrInvalidPrivateKey = errors.New("invalid private key")
	ErrInvalidPublicKey  = errors.New("invalid public key")
)

type encoder interface {
	Encode(text string) string
	Decode(text string) string
	Description() string
}

type Bytes []byte

// options .
type options struct {
	secretKey  Bytes // 密钥
	publicKey  Bytes // 公钥
	privateKey Bytes // 私钥
}

type option func(o *options)

// WithSecretKey .
func WithSecretKey(key Bytes) option {
	return func(o *options) {
		o.secretKey = key
	}
}

// WithPrivateKey .
func WithPrivateKey(publicKey, privateKey Bytes) option {
	return func(o *options) {
		o.publicKey = publicKey
		o.privateKey = privateKey
	}
}
