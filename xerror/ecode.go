package xerror

import (
	"fmt"
	"strconv"

	"github.com/charlesbases/logger"
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
type Code int

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
		return ServerErr
	} else {
		return Code(e)
	}
}
