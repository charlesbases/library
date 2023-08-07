package library

import "time"

const RFC3339Z = "2006-01-02 15:04:05.000"

type TimeString string

// Time .
func (ts TimeString) Time() time.Time {
	t, _ := time.Parse(RFC3339Z, string(ts))
	return t
}

// Timestamp .
func (ts TimeString) Timestamp() int64 {
	return ts.Time().UnixMilli()
}

type TimeTimestamp int64

// Time .
func (tt TimeTimestamp) Time() time.Time {
	return time.UnixMilli(int64(tt))
}

// TimeString .
func (tt TimeTimestamp) TimeString() string {
	return tt.Time().Format(RFC3339Z)
}

// Now .
func Now() time.Time {
	return time.Now()
}

// NowString .
func NowString() string {
	return time.Now().Format(RFC3339Z)
}

// NowTimestamp .
func NowTimestamp() int64 {
	return time.Now().UnixMilli()
}
