package services

import (
	"lowkeydd-server/crawlers"

	"github.com/gin-gonic/gin"
)

func CrawlerUpdate(c *gin.Context) {

	// 更新
	crawlers.GetInstance().UnChecked_Update()

	// 顯示結果
	GetAllVisitChannelsResponse(c)
}
