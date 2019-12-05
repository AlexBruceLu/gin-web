package util

import "time"

// 获取当前时间的字符时间
func GetCurrentDate() string {
	return time.Now().Format("2006/01/02 15:04:05")
}

// 获取当前的unix时间戳
func GetCurrentUnixTime() int64 {
	return time.Now().Unix()
}

// 获取当前时间的毫秒级时间戳
func GetCurrentMilliUnixTime() int64 {
	return time.Now().UnixNano() / 1000000
}

// 获取当前时间的纳秒级时间戳
func GetCurrentNanoUnixTime() int64 {
	return time.Now().UnixNano()
}
