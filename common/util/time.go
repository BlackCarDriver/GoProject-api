package util

import "time"

// ParseAnyTimeFormat 常见格式的时间转字符串
func ParseAnyTimeFormat(str string) (t int64) {
	layouts := []string{
		"2006-01-02",
		"2006/01/02",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"20060102 15:04",
		"20060102150405",
		"20060102",
	}
	for _, layout := range layouts {
		if timestamp, err := time.Parse(layout, str); err == nil {
			return timestamp.Unix()
		}
	}
	return 0
}
