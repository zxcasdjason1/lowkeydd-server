package services

import (
	"log"
	. "lowkeydd-server/share"

	"github.com/gin-gonic/gin"
)

func GetLetsddv2Response(c *gin.Context, visit *VisitList, channels []ChannelInfo, allChannels []ChannelInfo, authMsg string) {

	log.Printf("[Letsddv2_grouped_Response]\n")

	if visit == nil || len(visit.Group) == 0 {
		// log.Printf("[Letsddv2_grouped_Response] visit nill \n")
		c.JSON(200, gin.H{"code": "success", "channels": [][]ChannelInfo{channels}, "group": []string{}, "authMsg": authMsg})
		return
	}

	grouplen := len(visit.Group)
	indexMap := make(map[string]int)
	// groupIndexMap := make(map[string]int)
	listMap := make(map[string]VisitItem)
	groupNumCount := make([]int, grouplen)
	result := make([][]ChannelInfo, grouplen+1)
	last := grouplen

	// groupIndexMap 只要透過item.cid，就能知道對應的groupindex
	for i, groupName := range visit.Group {
		indexMap[string(groupName)] = i
	}
	for _, item := range visit.List {
		ix := indexMap[item.Group]
		// groupIndexMap[item.Cid] = ix
		listMap[item.Cid] = item
		groupNumCount[ix] += 1
		// log.Printf("ix: %d , count: %d, cname: %s, group: %s", ix, groupNumCount[ix], item.Cname, item.Group)
	}
	// log.Printf("[channelslen]: %d", len(channels))
	// log.Printf("[groupName]: %v", visit.Group)
	// log.Printf("[groupNumCount]: %v", groupNumCount)

	// 按群組分配
	for _, ch := range channels {
		if item, ok := listMap[ch.Cid]; ok {
			ix := indexMap[item.Group]
			if result[ix] == nil {
				result[ix] = make([]ChannelInfo, 0, groupNumCount[ix])
				// log.Printf("xi: %d , groupNumCount: %d", ix, groupNumCount[ix])
			}
			ch.Owner = item.Owner //改成自定義的頻道名稱
			result[ix] = append(result[ix], ch)
			// log.Printf("last: %d , cname: %s, status: %s", ix, ch.Cname, ch.Status)
		} else {
			if result[last] == nil {
				result[last] = make([]ChannelInfo, 0, len(channels))
				// log.Printf("xi: %d , groupNumCount: %d", last, len(channels))
			}
			result[last] = append(result[last], ch)
			// log.Printf("last: %d , cname: %s, status: %s", last, ch.Cname, ch.Status)
		}
	}

	// 透過 allChannels 來更新 List中的資料
	channelsMap := make(map[string]ChannelInfo)
	for _, ch := range allChannels {
		channelsMap[ch.Cid] = ch
	}
	newList := []VisitItem{}
	for _, item := range visit.List {
		ch := channelsMap[item.Cid]
		item.Avatar = ch.Avatar
		item.Cname = ch.Cname
		newList = append(newList, item)
	}

	c.JSON(200, gin.H{"code": "success", "channels": result, "group": visit.Group, "list": newList, "authMsg": authMsg})
}
