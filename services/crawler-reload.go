package services

import (
	"log"
	"lowkeydd-crawler/crawlers"

	"github.com/gin-gonic/gin"
)

func CrawlerReload(c *gin.Context) {

	log.Println("CrawlerReload")

	// 更新
	crawlers.NewCrawlers()

	c.JSON(200, gin.H{"code": "reload"})
}
