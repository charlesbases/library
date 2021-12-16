package zap

import (
	"fmt"
)

type level int8

const (
	levelTRC level = iota
	levelDBG
	levelINF
	levelWRN
	levelERR
	levelFAT
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

var colors = [6]attribute{}

var shorts = [6]string{}

func init() {
	colors[levelTRC] = white
	colors[levelDBG] = magenta
	colors[levelINF] = green
	colors[levelWRN] = blue
	colors[levelERR] = red
	colors[levelFAT] = red

	shorts[levelTRC] = levelTRC.wrap("[TRC]")
	shorts[levelDBG] = levelDBG.wrap("[DBG]")
	shorts[levelINF] = levelINF.wrap("[INF]")
	shorts[levelWRN] = levelWRN.wrap("[WRN]")
	shorts[levelERR] = levelERR.wrap("[ERR]")
	shorts[levelFAT] = levelFAT.wrap("[FAT]")
}

// wrap .
func (l level) wrap(v string) string {
	return fmt.Sprintf("%s[%dm%s%s[%dm", escape, l.color(), v, escape, reset)
}

// color .
func (l level) color() attribute {
	return colors[l]
}

// short .
func (l level) short() string {
	return shorts[l]
}
