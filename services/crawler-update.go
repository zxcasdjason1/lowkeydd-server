package services

import (
	"log"
	"lowkeydd-crawler/crawlers"
	"lowkeydd-crawler/redisdb"
	. "lowkeydd-crawler/share"

	"time"

	"github.com/gin-gonic/gin"
)

func CrawlerUpdate(c *gin.Context) {

	channels := getAllChannelInfo()

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

func getAllChannelInfo() []ChannelInfo {
	if cidlist := redisdb.GetInstance().Keys("*"); cidlist != nil {

		log.Printf("多筆查詢:> %v ", cidlist)
		channels := make([]ChannelInfo, 0, len(cidlist))

		for _, cid := range cidlist {
			if info, exist := redisdb.GetInstance().GetChannelInfo(cid); exist {
				channels = append(channels, info)
			}
		}

		return channels
	} else {
		return []ChannelInfo{}
	}
}
