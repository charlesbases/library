package logger

const (
	// DefaultDateFormat date format
	DefaultDateFormat = "2006-01-02 15:04:05.000"
	// DefaultFilename default file name
	DefaultFilename = "./log/log.log"
)

type level int8

const (
	LEVEL_TRACE level = iota
	LEVEL_DEBUG
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
	LEVEL_FATAL
)

// Options .
type Options struct {
	Service    string // 服务名
	Filename   string // 日志文件
	DateFormat string // 日期格式

	LowestLevel  level // 日志最低等级. default: LEVEL_TRACE
	HighestLevel level // 日志最高等级. default: LEVEL_FATAL
}

type Option func(o *Options)

// WithService .
func WithService(service string) Option {
	return func(o *Options) {
		o.Service = service
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
func WithLowestLevel(level level) Option {
	return func(o *Options) {
		o.LowestLevel = level
	}
}

// WithHighestLevel .
func WithHighestLevel(level level) Option {
	return func(o *Options) {
		o.HighestLevel = level
	}
}
