package services

import (
	"encoding/json"
	"log"
	"lowkeydd-crawler/redisdb"

	. "lowkeydd-crawler/share"
)

func GetChannelInfo(cid string) []ChannelInfo {

	jsonStr := redisdb.GetInstance().Get(cid)
	log.Printf("單筆查詢:> %v ", jsonStr)

	if jsonStr != "" {
		info := ChannelInfo{}
		json.Unmarshal([]byte(jsonStr), &info)
		return []ChannelInfo{info}
	} else {
		return []ChannelInfo{}
	}
}

func GetAllChannelInfo() []ChannelInfo {
	if cidlist := redisdb.GetInstance().GetClient().Keys("*").Val(); cidlist != nil {

		log.Printf("多筆查詢:> %v ", cidlist)
		channels := make([]ChannelInfo, 0, len(cidlist))

		for _, cid := range cidlist {
			var info ChannelInfo
			if jsonStr := redisdb.GetInstance().Get(cid); jsonStr != "" {
				json.Unmarshal([]byte(jsonStr), &info)
				channels = append(channels, info)
			}
		}
		return channels
	} else {
		return []ChannelInfo{}
	}
}
