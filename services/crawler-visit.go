package services

import (
	"log"
	"lowkeydd-crawler/crawlers"
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
	crawlers.GetInstance().Visit_Conditionally(cid, method)
	// 再從redis取出資料作為回傳
	defer GetSingleChannelResponse(c, cid)
}
