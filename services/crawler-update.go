package services

import (
	"log"
	"lowkeydd-crawler/crawlers"

	"time"

	"github.com/gin-gonic/gin"
)

func CrawlerUpdate(c *gin.Context) {

	channels := GetAllChannelInfo()

	// 更新
	log.Println("[Services] 頻道資訊更新作業開始....")

	crawler := crawlers.GetInstance()
	wg := crawler.GetWg()
	curr := time.Now().Unix()

	wg.Add(len(channels))
	for _, item := range channels {
		go crawler.Update(item, curr)
	}
	wg.Wait()

	log.Println("[Services] 頻道資訊更新作業結束....")

	// 顯示結果
	GetAllChannelsResponse(c)
}
