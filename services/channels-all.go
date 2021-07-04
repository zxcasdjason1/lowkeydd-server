package services

import (
	"lowkeydd-server/redisdb"
	. "lowkeydd-server/share"

	"github.com/gin-gonic/gin"
)

func GetAllChannelsResponse(c *gin.Context) {

	if channels, success := redisdb.GetInstance().GetChannelsByCondition(func(info ChannelInfo) bool {
		return info.Status != "failure"
	}); success {
		c.JSON(200, gin.H{"code": "success", "channels": channels})
	} else {
		c.JSON(400, gin.H{"code": "error", "channels": []ChannelInfo{}})
	}
}

func GetAllChannels(c *gin.Context) (bool, []ChannelInfo) {

	if channels, success := redisdb.GetInstance().GetChannelsByCondition(func(info ChannelInfo) bool {
		return info.Status != "failure"
	}); success {
		return true, channels
	} else {
		return false, channels
	}
}
