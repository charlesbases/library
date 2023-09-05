package webserver

import (
	"fmt"
)

// WebError .
type WebError struct {
	ID      string `json:"id"`
	Code    Code   `json:"code"`
	Message string `json:"message"`
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
