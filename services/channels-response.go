package services

import (
	"encoding/json"
	"log"
	"lowkeydd-crawler/redisdb"

	. "lowkeydd-crawler/share"

	"github.com/gin-gonic/gin"
)

type ChannelInfoResponse struct {
	Code     string
	Channels []ChannelInfo
}

func GetSingleChannelResponse(c *gin.Context, cid string) {

	jsonStr := redisdb.GetInstance().Get(cid)
	log.Printf("單筆查詢:> %v ", jsonStr)

	if jsonStr != "" {
		info := ChannelInfo{}
		json.Unmarshal([]byte(jsonStr), &info)
		c.JSON(200, gin.H{"code": "success", "channels": []ChannelInfo{info}})
	} else {
		c.JSON(400, gin.H{"code": "error", "channels": []ChannelInfo{}})
	}
}

func GetChannelsResponseByCondition(c *gin.Context, condition func(c ChannelInfo) bool) {
	if cidlist := redisdb.GetInstance().GetClient().Keys("*").Val(); cidlist != nil {

		channels := make([]ChannelInfo, 0, len(cidlist))
		log.Printf("多筆查詢:> %v ", cidlist)

		for _, cid := range cidlist {
			var info ChannelInfo
			if jsonStr := redisdb.GetInstance().Get(cid); jsonStr != "" {
				json.Unmarshal([]byte(jsonStr), &info)
				if condition(info) {
					channels = append(channels, info)
				}
			}
		}
		c.JSON(200, gin.H{"code": "success", "channels": channels})
	} else {
		c.JSON(400, gin.H{"code": "error", "channels": []ChannelInfo{}})
	}
}

func GetAllChannelsResponse(c *gin.Context) {
	GetChannelsResponseByCondition(c, func(info ChannelInfo) bool {
		return info.Status != "failure"
	})
}

func GetFailureChannelsResponse(c *gin.Context) {
	GetChannelsResponseByCondition(c, func(info ChannelInfo) bool {
		return info.Status == "failure"
	})
}

func GetLiveChannelsResponse(c *gin.Context) {
	GetChannelsResponseByCondition(c, func(info ChannelInfo) bool {
		return info.Status == "live"
	})
}

func GetWaitingChannelsResponse(c *gin.Context) {
	GetChannelsResponseByCondition(c, func(info ChannelInfo) bool {
		return info.Status == "wait"
	})
}

func GetOfflineChannelsResponse(c *gin.Context) {
	GetChannelsResponseByCondition(c, func(info ChannelInfo) bool {
		return info.Status == "off"
	})
}
