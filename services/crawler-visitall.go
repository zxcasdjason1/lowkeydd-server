package services

import (
	"encoding/json"
	"log"
	"lowkeydd-crawler/crawlers"
	"lowkeydd-crawler/redisdb"

	. "lowkeydd-crawler/share"

	"github.com/gin-gonic/gin"
)

type ChannelInfoResponse = ChannelInfo

func GetChannelsByCondition(c *gin.Context, condition func(info *ChannelInfoResponse) bool) {
	if cidlist := redisdb.GetInstance().GetClient().Keys("*").Val(); cidlist != nil {

		resp := make([]ChannelInfoResponse, 0)
		log.Printf("多筆查詢:> %v ", cidlist)

		for _, cid := range cidlist {
			var info ChannelInfoResponse
			if jsonStr := redisdb.GetInstance().Get(cid); jsonStr != "" {
				json.Unmarshal([]byte(jsonStr), &info)
				if condition(&info) {
					resp = append(resp, info)
				}
			}
		}
		c.JSON(200, gin.H{"channels": resp})

	} else {
		c.JSON(400, gin.H{"msg": nil})
	}
}

func GetAllChannels(c *gin.Context) {
	GetChannelsByCondition(c, func(info *ChannelInfoResponse) bool {
		return info.Status != "failure"
	})
}

func CrawlerVisitAll(c *gin.Context) {
	log.Println("CrawlerVisitAll")
	crawlers.GetInstance().VisitAll()

	// 做爬蟲，資料會寫入到redis中
	crawlers.GetInstance().VisitAll()
	// 再從redis取出資料作為回傳
	defer CrawlerVisitAllTransport(c)
}

func CrawlerVisitAllTransport(c *gin.Context) {

	GetAllChannels(c)
}
