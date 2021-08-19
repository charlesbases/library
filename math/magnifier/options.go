package magnifier

// 原始浮点数扩大 Multiple 倍后，仍存在的小数位处理方案
type schme int8

const (
	Trunc schme = iota // 返回整数部分
	Round              // 四舍五入
	Ceil               // 向上取整
	Floor              // 向下取整
)

// DefaultMultiple 默认扩大十倍
var DefaultMultiple = 1e10

// Options .
type Options struct {
	Multiple float64 // 倍数. default: 1e10
	Scheme   schme   // default: Trunc
}

type Option func(o *Options)

// WithMultiple 放大倍数
func WithMultiple(multiple int) Option {
	return func(o *Options) {
		o.Multiple = float64(multiple)
	}
}

// WithScheme 小数处理方案
func WithScheme(s schme) Option {
	return func(o *Options) {
		o.Scheme = s
	}
}
