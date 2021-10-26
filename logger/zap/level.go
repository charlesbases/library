package zap

import (
	"fmt"
	"strconv"
)

type level int8

const (
	LEVEL_TRC level = iota
	LEVEL_DBG
	LEVEL_INF
	LEVEL_WRN
	LEVEL_ERR
	LEVEL_FAT
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
	LEVEL_TRC: yellow,
	LEVEL_DBG: magenta,
	LEVEL_INF: green,
	LEVEL_WRN: blue,
	LEVEL_ERR: red,
	LEVEL_FAT: red,
}

var shortNames = map[level]string{
	LEVEL_TRC: "[TRC]",
	LEVEL_DBG: "[DBG]",
	LEVEL_INF: "[INF]",
	LEVEL_WRN: "[WRN]",
	LEVEL_ERR: "[ERR]",
	LEVEL_FAT: "[FAT]",
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
	if name, find := shortNames[l]; find {
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

// Sprint .
func (l level) Sprint(data string) string {
	if color, found := colors[l]; found {
		return color.wrap(data)
	}
	return data
}
