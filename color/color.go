package color

import (
	"fmt"
	"strconv"
	"strings"
)

type Color []Attribute

// Attribute defines a single SGR Code
type Attribute int

const escape = "\x1b"

// Base attributes
const (
	Reset Attribute = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut
)

// Foreground text colors
const (
	FgBlack Attribute = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

// Foreground Hi-Intensity text colors
const (
	FgHiBlack Attribute = iota + 90
	FgHiRed
	FgHiGreen
	FgHiYellow
	FgHiBlue
	FgHiMagenta
	FgHiCyan
	FgHiWhite
)

// Background text colors
const (
	BgBlack Attribute = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

// Background Hi-Intensity text colors
const (
	BgHiBlack Attribute = iota + 100
	BgHiRed
	BgHiGreen
	BgHiYellow
	BgHiBlue
	BgHiMagenta
	BgHiCyan
	BgHiWhite
)

// New .
func New(v ...Attribute) *Color {
	var color Color = make([]Attribute, 0, len(v))
	color.Add(v...)
	return &color
}

// Add .
func (c *Color) Add(v ...Attribute) {
	*c = append(*c, v...)
}

// Sprint .
func (c *Color) Sprint(a ...interface{}) string {
	return c.wrap(fmt.Sprint(a...))
}

// Sprintf .
func (c *Color) Sprintf(format string, a ...interface{}) string {
	return c.wrap(fmt.Sprintf(format, a...))
}

// len .
func (c *Color) len() int {
	return len(*c)
}

// sequence .
func (c *Color) sequence() string {
	format := make([]string, c.len())
	for i, v := range *c {
		format[i] = strconv.Itoa(int(v))
	}
	return strings.Join(format, ";")
}

// format .
func (c *Color) format() string {
	return fmt.Sprintf("%s[%sm", escape, c.sequence())
}

// unformat .
func (c *Color) unformat() string {
	return fmt.Sprintf("%s[%dm", escape, Reset)
}

func (c *Color) wrap(s string) string {
	return c.format() + s + c.unformat()
}
