package logger

const (
	// DefaultSkip .
	DefaultSkip = 1
	// DefaultDateFormat date format
	DefaultDateFormat = "2006-01-02 15:04:05.000"
	// DefaultFilename default file name
	DefaultFilename = "./log.log"
)

// Options .
type Options struct {
	Service    string // 服务名
	Filename   string // 日志文件
	DateFormat string // 日期格式

	Skip int // 跳过的调用者数量. default: DefaultSkip
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
