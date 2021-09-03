package log

type Color []byte

// foreground colors
var (
	Black  Color = []byte{27 /* 0x1B */, 91 /* [ */, 51 /* 3 (前景色) */, 48 /* 0 */, 109 /* m */} // 0x1B[30m
	Red    Color = []byte{27 /* 0x1B */, 91 /* [ */, 51 /* 3 (前景色) */, 49 /* 1 */, 109 /* m */} // 0x1B[31m
	Green  Color = []byte{27 /* 0x1B */, 91 /* [ */, 51 /* 3 (前景色) */, 50 /* 2 */, 109 /* m */} // 0x1B[32m
	Yellow Color = []byte{27 /* 0x1B */, 91 /* [ */, 51 /* 3 (前景色) */, 51 /* 3 */, 109 /* m */} // 0x1B[33m
	Blue   Color = []byte{27 /* 0x1B */, 91 /* [ */, 51 /* 3 (前景色) */, 52 /* 4 */, 109 /* m */} // 0x1B[34m
	Purple Color = []byte{27 /* 0x1B */, 91 /* [ */, 51 /* 3 (前景色) */, 53 /* 5 */, 109 /* m */} // 0x1B[35m
	Cyan   Color = []byte{27 /* 0x1B */, 91 /* [ */, 51 /* 3 (前景色) */, 54 /* 6 */, 109 /* m */} // 0x1B[36m
	White  Color = []byte{27 /* 0x1B */, 91 /* [ */, 51 /* 3 (前景色) */, 55 /* 7 */, 109 /* m */} // 0x1B[37m

	bound Color = []byte{27 /* 0x1B */, 91 /* [ */, 48 /* 0 */, 109 /* m */}
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

// String color output
func (c Color) String(s string) string {
	return string(c) + s + string(bound)
}

func init() {
	for level, color := range _color {
		_colorString[level] = color.String("[" + level.String() + "]")
	}
}
