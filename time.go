package library

import "time"

const RFC3339Z = "2006-01-02 15:04:05.000"

type TimeFormat string

// Time .
func (tf TimeFormat) Time() time.Time {
	t, _ := time.Parse(RFC3339Z, string(tf))
	return t
}

// Covert .
func (t TimeFormat) Covert() time.Duration {
	return time.Duration(t.Time().UnixMilli())
}

// TimeDuration 毫秒时间戳
type TimeDuration time.Duration

// Time .
func (td TimeDuration) Time() time.Time {
	return time.UnixMilli(int64(td))
}

// Covert .
func (td TimeDuration) Covert() TimeFormat {
	return TimeFormat(td.Time().Format(RFC3339Z))
}

// NowFormat .
func NowFormat() TimeFormat {
	return TimeFormat(time.Now().Format(RFC3339Z))
}

// NowDuration .
func NowDuration() TimeDuration {
	return TimeDuration(time.Now().UnixMilli())
}
