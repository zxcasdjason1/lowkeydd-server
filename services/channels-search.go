package services

import (
	"lowkeydd-server/crawlers"
	"lowkeydd-server/crawlers/twitch"
	. "lowkeydd-server/share"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

type SearchChannelResponse struct {
	Code     string
	Channels []ChannelInfo
}

func GetSearchChannelResponse(c *gin.Context) {

	url := c.DefaultPostForm("url", "")

	if url == "" {
		c.JSON(200, gin.H{"code": "error", "channels": []ChannelInfo{}})
	}
	// log.Printf("url %v\n", url)

	// https://www.youtube.com/watch?v=5cyW7zqhAO0
	// https://www.youtube.com/channel/UC1uv2Oq6kNxgATlCiez59hw
	// https://www.twitch.tv/uzra

	cid, method := "", ""
	if cid = GetYoutubeCid(url); cid != "" {
		method = "youtube"
	} else if cid = GetTwitchCid(url); cid != "" {
		method = "twitch"
	} else {
		c.JSON(200, gin.H{"code": "error", "channels": []ChannelInfo{}})
		return
	}

	// log.Printf("cid %v\n", cid)
	// log.Printf("method %v\n", method)

	// 此處用cname來判斷，因為yt搜尋結果中，即使url輸入錯誤，也會以該錯誤的cid返回。
	// 所以改用其他屬性，作為成功搜尋與否的依據。
	if ch := crawlers.GetInstance().GetSearchChannel(cid, method); ch.Thumbnail != "" { // 沒有背景圖表示獲取失敗
		// log.Printf("GetSearchChannel: ch %v\n", ch)
		c.JSON(200, gin.H{"code": "success", "channels": []ChannelInfo{ch}})
	} else {
		c.JSON(200, gin.H{"code": "failure", "channels": []ChannelInfo{ch}})
	}
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
