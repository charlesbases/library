//go:generate protoc --proto_path=${GOPATH}/src:. --gogo_out=paths=source_relative:. error.proto

package xerror

import (
	"fmt"
)

// Error .
func (e *Error) Error() string {
	return e.GetMessage()
}

// New .
func New(code Code, message interface{}) *Error {
	err := &Error{Code: int32(code)}
	switch message.(type) {
	case error:
		err.Message = message.(error).Error()
	case string:
		err.Message = message.(string)
	default:
		err.Message = fmt.Sprintf("%v", message)
	}
	return err
}
