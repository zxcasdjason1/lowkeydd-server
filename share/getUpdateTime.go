package share

import (
	"log"
	"time"
)

const (
	UpdateTimeInterval = 180 //秒
)

func GetNextUpdateTime() int64 {
	timeUnix := time.Now().Unix()
	timeUnix += UpdateTimeInterval

	// formatTimeStr := time.Unix(timeUnix, 0).Format("2006-01-02 15:04:05")
	// log.Println(formatTimeStr)

	return timeUnix
}

func ToDateStr(timeUnix int64) string {
	formatTimeStr := time.Unix(timeUnix, 0).Format("2006-01-02 15:04:05")
	log.Println(formatTimeStr)
	return formatTimeStr
}
