package log

const (
	// DefaultDateFormat date format
	DefaultDateFormat = "2006-01-02 15:04:05.000"
	// DefaultFilename default file name
	DefaultFilename = "./log/log.log"
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

// string .
func (l Level) String() string {
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

// Options .
type Options struct {
	Service    string // 服务名
	Filename   string // 日志文件
	DateFormat string // 日期格式

	LowestLevel  Level // 日志最低等级. default: LEVEL_TRACE
	HighestLevel Level // 日志最高等级. default: LEVEL_FATAL
}

type Option func(o *Options)

// WithService .
func WithService(service string) Option {
	return func(o *Options) {
		o.Service = "[" + service + "]"
	}
}

// WithFilename .
func WithFilename(filename string) Option {
	return func(o *Options) {
		o.Filename = filename
	}
}

// WithDateFormat .
func WithDateFormat(layout string) Option {
	return func(o *Options) {
		o.DateFormat = layout
	}
}

// WithLowestLevel .
func WithLowestLevel(level Level) Option {
	return func(o *Options) {
		o.LowestLevel = level
	}
}

// WithHighestLevel .
func WithHighestLevel(level Level) Option {
	return func(o *Options) {
		o.HighestLevel = level
	}
}
