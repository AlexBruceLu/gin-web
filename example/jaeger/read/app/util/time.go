package util

import "time"

func GetCurrentDate() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
func GetCurrentUnixTime() int64 {
	return time.Now().Unix()
}
func GetCurrentMilliUnixTime() int64 {
	return time.Now().UnixNano() / 1000000
}
func GetCurrentNanoUnixTime() int64 {
	return time.Now().UnixNano()
}
