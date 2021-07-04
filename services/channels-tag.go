package services

import (
	"lowkeydd-server/redisdb"
	"lowkeydd-server/share"

	"github.com/gin-gonic/gin"
)

type GetTagedChannelsRequest struct {
	Tag string `uri:"tag"`
}

func GetTagedChannels(c *gin.Context, tag string) (bool, []share.ChannelInfo) {

	if channels, success := redisdb.GetInstance().GetChannelsByCondition(func(ch share.ChannelInfo) bool {
		return ch.Status == tag
	}); success {
		return true, channels
	} else {
		return false, []share.ChannelInfo{}
	}

}

func GetTagedChannelsResponse(c *gin.Context, tag string) {

	if channels, success := redisdb.GetInstance().GetChannelsByCondition(func(ch share.ChannelInfo) bool {
		return ch.Status == tag
	}); success {
		c.JSON(200, gin.H{"code": "success", "channels": channels})
	} else {
		c.JSON(400, gin.H{"code": "error", "channels": []share.ChannelInfo{}})
	}

}

func Get_MultiTaged_Channels(c *gin.Context, tags []string) (bool, []share.ChannelInfo) {

	if tags[0] == "all" {
		if channels, success := redisdb.GetInstance().GetChannelsByCondition(func(ch share.ChannelInfo) bool {
			return ch.Status != "failure"
		}); success {
			return true, channels
		} else {
			return false, []share.ChannelInfo{}
		}
	}

	tagMap := make(map[string]bool)
	for _, tag := range tags {
		tagMap[tag] = true
	}

	if channels, success := redisdb.GetInstance().GetChannelsByCondition(func(ch share.ChannelInfo) bool {
		return tagMap[ch.Status]
	}); success {
		return true, channels
	} else {
		return false, []share.ChannelInfo{}
	}
}

func TagedChannelEndpoint(c *gin.Context) {
	req := &GetTagedChannelsRequest{}
	if err := c.ShouldBindUri(&req); err != nil {
		GetSingleChannelResponse(c, "")
	}

	GetTagedChannelsResponse(c, req.Tag)
}
