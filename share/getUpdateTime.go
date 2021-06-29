package share

import (
	"log"
	"math/rand"
	"time"
)

func GetNextUpdateTime() int64 {
	timeUnix := time.Now().Unix()
	timeUnix += 120 + rand.Int63n(180) //給予 120~300 秒的更新間隔時間。

	// ToDateStr(timeUnix)

	return timeUnix
}

func ToDateStr(timeUnix int64) string {
	formatTimeStr := time.Unix(timeUnix, 0).Format("2006-01-02 15:04:05")
	log.Println(formatTimeStr)
	return formatTimeStr
}
