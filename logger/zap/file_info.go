package zap

import (
	"os"
	"time"
)

const defaultDateLayou = "2006-01-02"

// fileDate .
type fileDate struct {
	createAt time.Time
	expireAt time.Time
}

// parseZero .
func parseZero(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// newFileDate .
func newFileDate(ttl int) *fileDate {
	return newFileDateWithTime(time.Now(), ttl)
}

// newFileDateWithTime .
func newFileDateWithTime(t time.Time, ttl int) *fileDate {
	var zero = parseZero(t)
	fileDate := &fileDate{
		createAt: zero,
		expireAt: zero.AddDate(0, 0, ttl),
	}
	return fileDate
}

// newFileDateWithStr .
func newFileDateWithStr(s string, ttl int) *fileDate {
	t, err := time.Parse(defaultDateLayou, s)
	if err == nil && !t.IsZero() {
		return newFileDateWithTime(t, ttl)
	}
	return nil
}

// string .
func (fileDate *fileDate) string() string {
	return fileDate.createAt.Format(defaultDateLayou)
}

// fileInfo .
type fileInfo struct {
	fileDate *fileDate
	fileName string
	filePath string
}

// fileStack 历史日志文件
type fileStack struct {
	files []*fileInfo

	headIndex int
	tailIndex int

	len int
	cap int
}

// push .
func (stack *fileStack) push(fileInfo *fileInfo) {
	switch {
	case stack.len == 0:
		stack.len++
		stack.headIndex = 0
		stack.tailIndex = 0
		stack.files[0] = fileInfo
	case stack.len < stack.cap:
		stack.len++
		stack.headIndex++
		stack.files[stack.headIndex] = fileInfo
	case stack.len == stack.cap:
		// 删除最旧的日志文件
		os.Remove(stack.files[stack.tailIndex].filePath)
		// 插入最新的文件
		stack.files[stack.tailIndex] = fileInfo
		// 更新头尾下标
		if stack.headIndex != stack.len-1 {
			stack.headIndex++
		} else {
			stack.headIndex = 0
		}
		if stack.tailIndex != stack.len-1 {
			stack.tailIndex++
		} else {
			stack.tailIndex = 0
		}
	}
}
