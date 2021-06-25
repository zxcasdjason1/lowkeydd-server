package services

import (
	"encoding/json"
	"log"
	"lowkeydd-crawler/crawlers"
	"lowkeydd-crawler/redisdb"
	. "lowkeydd-crawler/share"

	"github.com/gin-gonic/gin"
)

type CrawlerVisitRequest struct {
	Cid    string `json:"cid"`
	Method string `json:"method"`
}
type CrawlerVisitResponse struct {
	Code    string      `json:"code"`
	Channel ChannelInfo `json:"channel"`
}

func CrawlerVisit(c *gin.Context) {

	// var req CrawlerVisitRequest
	// var resp CrawlerVisitResponse

	cid := c.DefaultPostForm("cid", "")
	method := c.DefaultPostForm("method", "")
	log.Printf("cid %v\n", cid)
	log.Printf("method %v\n", method)

	// 做爬蟲，資料會寫入到redis中
	crawlers.GetInstance().Visit(cid, method)
	// 再從redis取出資料作為回傳
	defer CrawlerVisitTransport(c, cid)
}

func CrawlerVisitTransport(c *gin.Context, cid string) {

	jsonStr := redisdb.GetInstance().Get(cid)
	log.Printf("ChannelInfo:> %v\n", jsonStr)

	var channel ChannelInfo
	if jsonStr != "" {
		json.Unmarshal([]byte(jsonStr), &channel)
		c.JSON(200, gin.H{"code": "success", "channel": channel})
	} else {
		c.JSON(400, gin.H{"code": "error", "channel": channel})
	}
}
