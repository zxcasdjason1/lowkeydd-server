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

func GetSingleChannelResponse(c *gin.Context, cid string) {

	log.Printf("單筆查詢:> %v ", cid)

	if cid == "" {
		c.JSON(400, gin.H{"code": "error", "channels": []ChannelInfo{}})
		return
	}

	if info, exist := redisdb.GetInstance().GetChannel(cid); exist {
		c.JSON(200, gin.H{"code": "success", "channels": []ChannelInfo{info}})
	} else {
		c.JSON(200, gin.H{"code": "failure", "channels": []ChannelInfo{}})
	}
}
