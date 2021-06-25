package services

import (
	"log"
	"lowkeydd-crawler/crawlers"

	"github.com/gin-gonic/gin"
)

func CrawlerVisitAll(c *gin.Context) {
	log.Println("CrawlerVisitAll")
	// 做爬蟲，資料會寫入到redis中
	crawlers.GetInstance().VisitAll_Conditionally()
	// 再從redis取出資料作為回傳
	defer GetAllChannelsResponse(c)

}
