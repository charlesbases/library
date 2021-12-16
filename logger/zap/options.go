package zap

const (
	// DefaultSkip .
	DefaultSkip = 1
	// DefaultMaxRolls 日志保留时间
	DefaultMaxRolls = 7
	// DefaultDateFormat date format
	DefaultDateFormat = "2006-01-02 15:04:05.000"
	// DefaultFilePath default file path
	DefaultFilePath = "./log/log.log"
)

// Options .
type Options struct {
	Service  string // 服务名
	FilePath string // 日志文件
	MaxRolls int    // 日志保留天数

	Skip int // 跳过的调用者数量. default: DefaultSkip
}

// defaultOption .
func defaultOption() *Options {
	return &Options{
		FilePath: DefaultFilePath,
		MaxRolls: DefaultMaxRolls,
		Skip:     DefaultSkip,
	}
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

// WithFilePath .
func WithFilePath(filename string) Option {
	return func(o *Options) {
		o.FilePath = filename
	}
}

// WithMaxRolls .
func WithMaxRolls(rolls int) Option {
	return func(o *Options) {
		o.MaxRolls = rolls
	}
}
