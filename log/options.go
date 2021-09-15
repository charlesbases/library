package log

const (
	// DefaultSkip .
	DefaultSkip = 2
	// DefaultDateFormat date format
	DefaultDateFormat = "2006-01-02 15:04:05.000"
	// DefaultFilename default file name
	DefaultFilename = "./log/log.log"
)

// Options .
type Options struct {
	Service    string // 服务名
	Filename   string // 日志文件
	DateFormat string // 日期格式

	LowestLevel  Level // 日志最低等级. default: LEVEL_TRACE
	HighestLevel Level // 日志最高等级. default: LEVEL_FATAL

	Skip int // 跳过的调用者数量. default: 2
}

type Option func(o *Options)

// WithSkip .
func WithSkip(skip int) Option {
	return func(o *Options) {
		o.Skip = skip
	}
}

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
