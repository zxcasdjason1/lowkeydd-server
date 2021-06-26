package services

import (
	"log"
	"lowkeydd-crawler/crawlers"
	. "lowkeydd-crawler/share"

	"github.com/gin-gonic/gin"
)

type CrawlerVisitRequest struct {
	Cid    string `uri:"cid"`
	Method string `uri:"method"`
}

type CrawlerVisitResponse struct {
	Code    string      `json:"code"`
	Channel ChannelInfo `json:"channel"`
}

func CrawlerVisit(c *gin.Context) {

	req := &CrawlerVisitRequest{}
	if err := c.ShouldBindUri(&req); err != nil {
		GetSingleChannelResponse(c, "")
	}

	if req.Cid == "" || req.Method == "" {
		GetSingleChannelResponse(c, "")
	}

	log.Printf("cid %v\n", req.Cid)
	log.Printf("method %v\n", req.Method)
	// 做爬蟲，資料會寫入到redis中
	crawlers.GetInstance().Visit_Conditionally(req.Cid, req.Method)
	// 再從redis取出資料作為回傳
	GetSingleChannelResponse(c, req.Cid)
}
