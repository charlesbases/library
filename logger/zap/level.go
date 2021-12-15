package zap

import (
	"fmt"
	"strconv"
)

type level int8

const (
	level_trace level = iota
	level_debug
	level_info
	level_warn
	level_error
	level_fatal
)

type attribute int

const (
	reset  = 0
	escape = "\x1b"
)

const (
	black attribute = iota + 30
	red
	green
	yellow
	blue
	magenta
	cyan
	white
)

var colors = map[level]attribute{
	level_trace: yellow,
	level_debug: magenta,
	level_info:  green,
	level_warn:  blue,
	level_error: red,
	level_fatal: red,
}

var shorts = map[level]string{
	level_trace: "[TRC]",
	level_debug: "[DBG]",
	level_info:  "[INF]",
	level_warn:  "[WRN]",
	level_error: "[ERR]",
	level_fatal: "[FAT]",
}

// format .
func (c attribute) format() string {
	return fmt.Sprintf("%s[%dm", escape, c)
}

// unformat .
func (c attribute) unformat() string {
	return fmt.Sprintf("%s[%dm", escape, reset)
}

// string .
func (c attribute) string() string {
	return strconv.Itoa(int(c))
}

// wrap
func (c attribute) wrap(s string) string {
	return c.format() + s + c.unformat()
}

// short .
func (l level) short() string {
	if name, find := shorts[l]; find {
		return name
	} else {
		return "UNK"
	}
}

// color .
func (l level) color() string {
	if color, find := colors[l]; find {
		return color.string()
	}
	return white.string()
}

// sprint .
func (l level) sprint(data string) string {
	if color, found := colors[l]; found {
		return color.wrap(data)
	}
	return data
}
