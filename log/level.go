package log

import (
	"library/color"
)

type Level int8

const (
	LEVEL_TRACE Level = iota
	LEVEL_DEBUG
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
	LEVEL_FATAL
)

var colors = map[Level]*color.Color{
	LEVEL_TRACE: color.New(color.FgWhite),
	LEVEL_DEBUG: color.New(color.FgCyan),
	LEVEL_INFO:  color.New(color.FgGreen),
	LEVEL_WARN:  color.New(color.FgWhite),
	LEVEL_ERROR: color.New(color.FgRed),
	LEVEL_FATAL: color.New(color.FgRed),
}

// short .
func (l Level) short() string {
	switch l {
	case LEVEL_TRACE:
		return "TRC"
	case LEVEL_DEBUG:
		return "DBG"
	case LEVEL_INFO:
		return "INF"
	case LEVEL_WARN:
		return "WRN"
	case LEVEL_ERROR:
		return "ERR"
	case LEVEL_FATAL:
		return "FAT"
	default:
		return "UNK"
	}
}

// Sprint .
func (l Level) Sprint(data string) string {
	if color, found := colors[l]; found {
		return color.Sprintf(data)
	}
	return data
}
