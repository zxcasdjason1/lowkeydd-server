package services

import (
	"lowkeydd-server/redisdb"
	"lowkeydd-server/share"

	"github.com/gin-gonic/gin"
)

type GetTagedChannelsRequest struct {
	Tag string `uri:"tag"`
}

func GetTagedChannelsResponse(c *gin.Context) {

	req := &GetTagedChannelsRequest{}
	if err := c.ShouldBindUri(&req); err != nil {
		GetSingleChannelResponse(c, "")
	}

	if channels, success := redisdb.GetInstance().GetChannelsByCondition(func(ch share.ChannelInfo) bool {
		return ch.Status == req.Tag
	}); success {
		c.JSON(200, gin.H{"code": "success", "channels": channels})
	} else {
		c.JSON(400, gin.H{"code": "error", "channels": []share.ChannelInfo{}})
	}

}
