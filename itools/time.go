package itools

import "time"

// FmtTime 格式化时间
func FmtTime(now time.Time) string {
	return now.Format("2006-01-02T15:04:05")
}

// FmtNow 格式化当前时间
func FmtNow() string {
	return time.Now().Format("2006-01-02T15:04:05")
}

// FmtTimeWithMill 带有毫秒
func FmtTimeWithMill(now time.Time) string {
	return now.Format("2006-01-02T15:04:05.000")
}

// FmtTimeWithMillOut 带有毫秒
func FmtTimeWithMillOut(now time.Time) string {
	return now.Format("20060102150405.000")
}
