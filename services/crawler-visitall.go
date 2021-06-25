package services

import (
	"log"
	"lowkeydd-crawler/crawlers"

	"github.com/gin-gonic/gin"
)

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
