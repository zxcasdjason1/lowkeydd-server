package services

import (
	"log"
	"lowkeydd-server/redisdb"

	. "lowkeydd-server/share"

	"github.com/gin-gonic/gin"
)

type ChannelInfoResponse struct {
	Code     string
	Channels []ChannelInfo
}

func GetSingleVisitChannelResponse(c *gin.Context, cid string) {

	log.Printf("單筆查詢:> %v ", cid)

	if cid == "" {
		c.JSON(400, gin.H{"code": "error", "channels": []ChannelInfo{}})
		return
	}

	if info, exist := redisdb.GetInstance().GetVisitChannel(cid); exist {
		c.JSON(200, gin.H{"code": "success", "channels": []ChannelInfo{info}})
	} else {
		c.JSON(200, gin.H{"code": "failure", "channels": []ChannelInfo{}})
	}
}

func GetAllVisitChannelsResponse(c *gin.Context) {

	if channels, success := redisdb.GetInstance().GetVisitChannelsByCondition(func(info ChannelInfo) bool {
		return info.Status != "failure"
	}); success {
		c.JSON(200, gin.H{"code": "success", "channels": channels})
	} else {
		c.JSON(400, gin.H{"code": "error", "channels": []ChannelInfo{}})
	}
}
