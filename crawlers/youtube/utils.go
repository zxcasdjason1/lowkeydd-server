package youtube

import (
	"log"
	"strconv"
	"strings"
	"time"
)

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
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Fatal("時間轉碼失敗")
		panic(err)
	}
	tm := time.Unix(i, 0)
	return tm.String()
}
