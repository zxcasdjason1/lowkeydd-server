package services

import (
	"lowkeydd-server/redisdb"
	. "lowkeydd-server/share"

	"github.com/gin-gonic/gin"
)

// type LetsddResponse struct {
// 	Code     string        `json:"code"`
// 	Channels []ChannelInfo `json:"channels"`
// 	Visit    VisitList     `json:"visit"`
// }

func GetLetsddChannelsResponse(c *gin.Context, visit *VisitList) {

	if channels, success := redisdb.GetInstance().GetChannelsByCondition(func(info ChannelInfo) bool {
		return info.Status != "failure"
	}); success {
		c.JSON(200, gin.H{"code": "success", "channels": channels, "visit": visit})
	} else {
		c.JSON(400, gin.H{"code": "error", "channels": []ChannelInfo{}, "visit": visit})
	}
}
