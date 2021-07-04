package services

import (
	"log"
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

func Letsddv2_grouped_Response(c *gin.Context, visit *VisitList, channels []ChannelInfo) {

	log.Printf("[Letsddv2_grouped_Response]\n")

	if visit == nil {
		log.Printf("[Letsddv2_grouped_Response] visit nill \n")
		c.JSON(200, gin.H{"code": "success", "channels": [][]ChannelInfo{channels}, "group": []string{}})
		return
	}

	groupMap := make(map[string]string)
	for _, item := range visit.List {
		groupMap[item.Cid] = string(item.Group)
	}

	indexMap := make(map[string]int)
	for i, groupName := range visit.Group {
		indexMap[string(groupName)] = i
	}

	len := len(visit.Group)
	groupedChannels := make([][]ChannelInfo, len+1)

	for _, ch := range channels {
		groupName := groupMap[ch.Cid]
		if ix, ok := indexMap[groupName]; ok {
			log.Printf("cid: %s , xi:%d , groupName: %s", ch.Cid, ix, groupName)
			groupedChannels[ix] = append(groupedChannels[ix], ch)
		} else {
			log.Printf("cid: %s , xi:%d , groupName: %s", ch.Cid, ix, groupName)
			groupedChannels[len] = append(groupedChannels[len], ch)
		}
	}

	// for i, chsgroup := range groupedChannels {
	// 	log.Printf("GetLetsddV2ChannelsResponse: %d , %v", i, chsgroup)
	// }

	c.JSON(200, gin.H{"code": "success", "channels": groupedChannels, "group": visit.Group})

}
