package log

import "fmt"

type Color uint8

// foreground colors
const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

var (
	_color = map[Level]Color{
		LEVEL_TRACE: White,
		LEVEL_DEBUG: Cyan,
		LEVEL_INFO:  Green,
		LEVEL_WARN:  White,
		LEVEL_ERROR: Red,
		LEVEL_FATAL: Red,
	}

	_colorString = make(map[Level]string, len(_color))
)

// Add adds the coloring to the given string.
func (c Color) Add(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), s)
}

func init() {
	for level, color := range _color {
		_colorString[level] = color.Add("[" + level.String() + "]")
	}
}
