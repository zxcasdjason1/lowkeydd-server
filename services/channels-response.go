package services

import (
	"encoding/json"
	"log"
	"lowkeydd-crawler/crawlers"
	"lowkeydd-crawler/crawlers/twitch"
	"lowkeydd-crawler/redisdb"
	"regexp"
	"strings"

	. "lowkeydd-crawler/share"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

type ChannelInfoResponse struct {
	Code     string
	Channels []ChannelInfo
}

func GetSingleChannelResponse(c *gin.Context, cid string) {

	if cid == "" {
		c.JSON(400, gin.H{"code": "error", "channels": []ChannelInfo{}})
		return
	}

	jsonStr := redisdb.GetInstance().Get(cid)
	log.Printf("單筆查詢:> %v ", jsonStr)

	if jsonStr != "" {
		info := ChannelInfo{}
		json.Unmarshal([]byte(jsonStr), &info)
		c.JSON(200, gin.H{"code": "success", "channels": []ChannelInfo{info}})
	} else {
		c.JSON(400, gin.H{"code": "error", "channels": []ChannelInfo{}})
	}
}

func GetChannelsResponseByCondition(c *gin.Context, condition func(c ChannelInfo) bool) {
	if cidlist := redisdb.GetInstance().Keys("*"); cidlist != nil {

		channels := make([]ChannelInfo, 0, len(cidlist))
		log.Printf("多筆查詢:> %v ", cidlist)

		for _, cid := range cidlist {
			var info ChannelInfo
			if jsonStr := redisdb.GetInstance().Get(cid); jsonStr != "" {
				json.Unmarshal([]byte(jsonStr), &info)
				if condition(info) {
					channels = append(channels, info)
				}
			}
		}
		c.JSON(200, gin.H{"code": "success", "channels": channels})
	} else {
		c.JSON(400, gin.H{"code": "error", "channels": []ChannelInfo{}})
	}
}

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
	crawlers.GetInstance().Visit_Conditionally(cid, method)
	// 再從redis取出資料作為回傳
	GetSingleChannelResponse(c, cid)
}

func GetYoutubeCid(url string) string {

	r, _ := regexp.Compile("https?://www.youtube.com/channel/(.*)")
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

func GetAllChannelsResponse(c *gin.Context) {
	GetChannelsResponseByCondition(c, func(info ChannelInfo) bool {
		return info.Status != "failure"
	})
}

// func GetFailureChannelsResponse(c *gin.Context) {
// 	GetChannelsResponseByCondition(c, func(info ChannelInfo) bool {
// 		return info.Status == "failure"
// 	})
// }

// func GetLiveChannelsResponse(c *gin.Context) {
// 	GetChannelsResponseByCondition(c, func(info ChannelInfo) bool {
// 		return info.Status == "live"
// 	})
// }

// func GetWaitingChannelsResponse(c *gin.Context) {
// 	GetChannelsResponseByCondition(c, func(info ChannelInfo) bool {
// 		return info.Status == "wait"
// 	})
// }

// func GetOfflineChannelsResponse(c *gin.Context) {
// 	GetChannelsResponseByCondition(c, func(info ChannelInfo) bool {
// 		return info.Status == "off"
// 	})
// }

type GetTagedChannelsRequest struct {
	Tag string `uri:"tag"`
}

func GetTagedChannelsResponse(c *gin.Context) {

	req := &GetTagedChannelsRequest{}
	if err := c.ShouldBindUri(&req); err != nil {
		GetSingleChannelResponse(c, "")
	}

	GetChannelsResponseByCondition(c, func(info ChannelInfo) bool {
		return info.Status == req.Tag
	})
}
