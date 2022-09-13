package library

import (
	"context"
	"time"
)

const (
	// DefaultTimeFormatLayout 默认时间格式
	DefaultTimeFormatLayout = "2006-01-02 15:04:05"
)

// Now 当前时间
func Now() string {
	return time.Now().Format(DefaultTimeFormatLayout)
}

// NowTimestamp 当前毫秒时间戳
func NowTimestamp() int64 {
	return time.Now().UnixMilli()
}

// ParseTime2String 时间格式化输出
func ParseTime2String(t time.Time) string {
	return t.Format(DefaultTimeFormatLayout)
}

// ParseTime2Timestamp 毫秒时间戳
func ParseTime2Timestamp(t time.Time) int64 {
	return t.UnixMilli()
}

// ParseTimeStr2Time 时间格式化字符串转时间
func ParseTimeStr2Time(s string) time.Time {
	t, _ := time.Parse(DefaultTimeFormatLayout, s)
	return t
}

// ParseTimeStr2Timestamp 时间格式化字符串转毫秒时间戳
func ParseTimeStr2Timestamp(s string) int64 {
	return ParseTime2Timestamp(ParseTimeStr2Time(s))
}

type Duration time.Duration

// Shrink will decrease the duration by comparing with context's timeout duration
// and return new timeout\context\CancelFunc.
func (d Duration) Shrink(c context.Context) (Duration, context.Context, context.CancelFunc) {
	if deadline, ok := c.Deadline(); ok {
		if ctimeout := time.Until(deadline); ctimeout < time.Duration(d) {
			// deliver small timeout
			return Duration(ctimeout), c, func() {}
		}
	}
	ctx, cancel := context.WithTimeout(c, time.Duration(d))
	return d, ctx, cancel
}
