package webserver

import (
	"fmt"
	"strconv"

	"github.com/charlesbases/library/logger"
)

// ecodes register codes
var ecodes = map[int]string{}

// add .
func add(e int, m string) Code {
	if _, found := ecodes[e]; found {
		logger.Fatalf("ecode: %d already existed.", e)
	}
	ecodes[e] = m
	return Int(e)
}

// Code is an int error code spec
type Code int32

// Int .
func (e Code) Int() int {
	return int(e)
}

// Error .
func (e Code) Error() string {
	if m, found := ecodes[e.Int()]; found {
		return m
	}
	return fmt.Sprintf("unknown ecode: %d", e)
}

// WebError .
func (e Code) WebError(v interface{}) *WebError {
	return NewWebError(e, v)
}

// Int parse code int to error
func Int(e int) Code {
	return Code(e)
}

// String parse code string to error
func String(e string) Code {
	if len(e) == 0 {
		return StatusOK
	}
	if e, err := strconv.Atoi(e); err != nil {
		return StatusServerError
	} else {
		return Code(e)
	}
}
