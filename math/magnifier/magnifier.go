package magnifier

import "math"

var m *magnifier

// magnifier 浮点数放大和复原
type magnifier struct {
	opts *Options
}

func init() {
	useDefaultMagnifier()
}

// New .
func New(opts ...Option) {
	var options = new(Options)
	for _, o := range opts {
		o(options)
	}

	magnifier := new(magnifier)
	magnifier.opts = options

	if magnifier.opts.Multiple == 0 {
		magnifier.opts.Multiple = DefaultMultiple
	}

	m = magnifier
}

// useDefaultMagnifier .
func useDefaultMagnifier() {
	m = new(magnifier)

	m.opts = new(Options)
	m.opts.Multiple = DefaultMultiple
}

// Magnify 原始浮点数扩大 Multiple 倍
func Magnify(source float64) int64 {
	number := source * m.opts.Multiple

	var transformed int64
	switch m.opts.Scheme {
	case Trunc:
		transformed = toint64(math.Trunc(number))
	case Round:
		transformed = toint64(math.Round(number))
	case Ceil:
		transformed = toint64(math.Ceil(number))
	case Floor:
		transformed = toint64(math.Floor(number))
	}

	return transformed
}

// Restore 浮点数复原
func Restore(transformed int64) float64 {
	return tofloat64(transformed) * m.opts.Multiple
}

// tofloat64 int64 to float64
func tofloat64(number int64) float64 {
	return float64(number)
}

// toint64 float64 to int64
func toint64(number float64) int64 {
	return int64(number)
}
