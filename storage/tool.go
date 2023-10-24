package storage

import (
	"strings"
	"unicode/utf8"

	"github.com/pkg/errors"

	"github.com/charlesbases/library/regexp"
)

// validString .
func validString(v string) error {
	if len(v) == 0 {
		return errors.New("name cannot be empty")
	}
	if !utf8.ValidString(v) {
		return errors.New("name non UTF-8 strings are not supported")
	}
	return nil
}

// Validator .
type Validator interface {
	Error() error
}

// BucketName bucket name
type BucketName string

// Error .
func (n BucketName) Error() error {
	if err := validString(string(n)); err != nil {
		return errors.Wrap(err, "bucket")
	}
	if regexp.IP.MatchString(string(n)) {
		return errors.New("bucket: name cannot be an ip address")
	}
	return nil
}

// KeyName key name
type KeyName string

// Error .
func (n KeyName) Error() error {
	if err := validString(string(n)); err != nil {
		return errors.Wrap(err, "key")
	}
	if strings.HasSuffix(string(n), "/") {
		return errors.New("key: name cannot end with '/'")
	}
	return nil
}

// KeyPrefixName prefix of key
type KeyPrefixName string

// Error .
func (n KeyPrefixName) Error() error {
	return errors.Wrap(validString(string(n)), "prefix")
}

// ValidatorFunc error func
type ValidatorFunc func() error

// Error .
func (f ValidatorFunc) Error() error {
	return f()
}

// ErrorValidator .
func ErrorValidator(vs ...Validator) error {
	for _, v := range vs {
		if err := v.Error(); err != nil {
			return err
		}
	}
	return nil
}
