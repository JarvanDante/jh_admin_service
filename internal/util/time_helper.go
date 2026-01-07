package util

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// FormatTime 统一格式化时间，处理 gtime.Time 的零值问题
// 使用标准库的 time.Time 来格式化，避免 gtime.Time.Format() 的问题
func FormatTime(t *gtime.Time) string {
	if t == nil {
		return ""
	}

	// 检查是否为零值时间
	if t.IsZero() {
		return ""
	}

	// 使用标准库的 time.Time 来格式化，而不是 gtime.Time.Format()
	stdTime := t.Time
	if stdTime.IsZero() || stdTime.Unix() <= 0 {
		return ""
	}

	return stdTime.Format("2006-01-02 15:04:05")
}
