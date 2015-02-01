package mytime

import (
	"time"
)

// 通过时间戳获得通用时间格式的字符串
func GetDateTime(timestamp int64) (datetime string) {
	t := time.Unix(timestamp, 0)
	datetime = t.Format("2006-01-02 15:04:05")
	return
}
