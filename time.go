package library

import "time"

const (
	// DefaultTimeFormatLayout 默认时间格式
	DefaultTimeFormatLayout = "2006-01-02 15:04:05"
)

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
