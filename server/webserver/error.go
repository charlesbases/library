package webserver

import (
	"fmt"

	"github.com/charlesbases/library/sonyflake"
)

// WebError .
type WebError struct {
	ID      sonyflake.ID `json:"id"`
	Code    Code         `json:"code"`
	Message string       `json:"message"`
}

// Error .
func (e *WebError) Error() string {
	return e.Message
}

// NewWebError .
func NewWebError(code Code, message interface{}) *WebError {
	err := &WebError{Code: code}
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
