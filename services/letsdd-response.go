package services

import (
	"lowkeydd-server/redisdb"
	. "lowkeydd-server/share"

	"github.com/gin-gonic/gin"
)

func GetLetsddChannelsResponse(c *gin.Context, visit *VisitList) {

	if channels, success := redisdb.GetInstance().GetVisitChannelsByCondition(func(info ChannelInfo) bool {
		return info.Status != "failure"
	}); success {
		c.JSON(200, gin.H{"code": "success", "channels": channels, "visit": visit})
	} else {
		c.JSON(400, gin.H{"code": "error", "channels": []ChannelInfo{}, "visit": visit})
	}
}
