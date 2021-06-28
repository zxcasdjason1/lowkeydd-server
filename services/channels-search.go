package services

import (
	"log"
	"lowkeydd-server/crawlers"
	"lowkeydd-server/crawlers/twitch"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func GetSearchChannelResponse(c *gin.Context) {

	url := c.DefaultPostForm("url", "")

	if url == "" {
		GetSingleChannelResponse(c, "")
	}
	log.Printf("url %v\n", url)

	// https://www.youtube.com/watch?v=5cyW7zqhAO0
	// https://www.youtube.com/channel/UC1uv2Oq6kNxgATlCiez59hw
	// https://www.twitch.tv/uzra

	cid, method := "", ""
	if cid = GetYoutubeCid(url); cid != "" {
		method = "youtube"
	} else if cid = GetTwitchCid(url); cid != "" {
		method = "twitch"
	} else {
		GetSingleChannelResponse(c, "")
		return
	}

	log.Printf("cid %v\n", cid)
	log.Printf("method %v\n", method)
	// 做爬蟲，資料會寫入到redis中
	crawlers.GetInstance().Checked_Visit(cid, method)
	// 再從redis取出資料作為回傳
	GetSingleChannelResponse(c, cid)
}

func GetYoutubeCid(url string) string {
	r, _ := regexp.Compile("https?://.*.youtube.com/channel/(.*)")
	submatch := r.FindSubmatch([]byte(url))
	if len(submatch) == 0 {
		return ""
	}

	for strings.Contains(string(submatch[1]), "/") {
		r, _ := regexp.Compile("(.*)/")
		submatch = r.FindSubmatch([]byte(string(submatch[1])))
	}

	return string(submatch[1])
}

func GetTwitchCid(url string) string {
	r, _ := regexp.Compile("https?://www.twitch.tv/(.*)")
	submatch := r.FindSubmatch([]byte(url))
	if len(submatch) == 0 {
		return ""
	}

	for strings.Contains(string(submatch[1]), "/") {
		r, _ := regexp.Compile("(.*)/")
		submatch = r.FindSubmatch([]byte(string(submatch[1])))
	}

	loginName := string(submatch[1])
	users := crawlers.GetTwitchCrawler().GetUserInfo(loginName)
	return twitch.RemoveQuotes(gjson.Get(users, "users.0._id").Raw)
}
