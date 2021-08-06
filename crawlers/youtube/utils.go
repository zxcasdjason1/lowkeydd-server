package youtube

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

func getWaitOrOffViewCountStr(ctx string) string {
	var suffix = ""
	res := gjson.Get(ctx, "viewCountText.simpleText").Raw
	if res != "" {
		// status == "off"
		suffix = " 已觀看"

	} else {
		// status == "wait"
		vc1 := gjson.Get(ctx, "viewCountText.runs.0.text")
		vc2 := gjson.Get(ctx, "viewCountText.runs.1.text")
		res = vc1.Raw + " " + vc2.Raw

		suffix = " 等待中"
	}
	// 修整取出純數字的觀眾數量。
	num := regexp.MustCompile("[0-9,]+").FindAllString(res, -1)
	if len(num) > 0 {
		res = num[0] + suffix
	} else {
		res = "0" + suffix
	}
	return res
}

func getStartTimeStr(ctx string) string {
	res := gjson.Get(ctx, "upcomingEventData.startTime").Raw // 檢查有沒有預定發布時間
	if res != "" {
		return timestampToDate(removeQuotes(res))
	}
	return ""
}

func removeSlash(s string) string {
	if len(s) > 0 && s[0] == '\\' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '\\' {
		s = s[:len(s)-1]
	}
	return s
}

func removeQuotes(s string) string {
	return strings.ReplaceAll(s, "\"", "")
}

func removeQuoteSlash(s string) string {
	return strings.ReplaceAll(s, "\\\"", "")
}

// 將時間戳記轉換成現實時間
func timestampToDate(s string) string {
	// 得到時戳
	timestamp, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Fatal("時間轉碼失敗")
		panic(err)
	}
	//
	// return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04")
}

func getTimeByStartTimeStr(str string) int {

	// log.Printf("[getTimeByStartTimeStr] startTime: %s\n", str)

	var priority = 0
	if ok := strings.Contains(str, "秒"); ok {
		priority = 1
	} else if ok := strings.Contains(str, "分"); ok {
		priority = 60
	} else if ok := strings.Contains(str, "小時"); ok {
		priority = 3600
	} else if ok := strings.Contains(str, "天"); ok {
		priority = 86400
	} else if ok := strings.Contains(str, "週"); ok {
		priority = 604800
	} else if ok := strings.Contains(str, "月"); ok {
		priority = 2592000
	} else if ok := strings.Contains(str, "年"); ok {
		priority = 31536000
	}

	re := regexp.MustCompile("[0-9]+")
	num := re.FindAllString(str, -1)
	if len(num) > 0 {
		val, _ := strconv.Atoi(num[0])
		log.Printf("[getTimeByStartTimeStr] val: %d\n", val*priority)
		return val * priority
	} else {
		return 2147483647
	}

}
