package redisdb

import (
	"encoding/json"
	"log"
	. "lowkeydd-server/share"
	"time"
)

func (this *Driver) GetVisitChannel(cid string) (ChannelInfo, bool) {

	this.SelectDB("")

	ch := ChannelInfo{}
	if jsonStr := this.Get(cid); jsonStr != "" {
		json.Unmarshal([]byte(jsonStr), &ch)
		return ch, true
	} else {
		return ch, false
	}
}

func (this *Driver) GetAllVisitChannels() ([]ChannelInfo, bool) {

	this.SelectDB("")

	if cidlist := this.Keys("*"); cidlist != nil {

		channels := make([]ChannelInfo, 0, len(cidlist))

		for _, cid := range cidlist {
			if ch, exist := this.GetVisitChannel(cid); exist {
				channels = append(channels, ch)
			}
		}

		return channels, true
	} else {
		return []ChannelInfo{}, false
	}
}

func (this *Driver) GetVisitChannelsByCondition(condition func(c ChannelInfo) bool) ([]ChannelInfo, bool) {

	this.SelectDB("")

	if cidlist := this.Keys("*"); cidlist != nil {

		channels := make([]ChannelInfo, 0, len(cidlist))

		for _, cid := range cidlist {
			if info, exist := this.GetVisitChannel(cid); exist {
				if condition(info) {
					channels = append(channels, info)
				}
			}
		}

		return channels, true
	} else {
		return []ChannelInfo{}, false
	}
}

func (this *Driver) GetVisitChannelsByConditionV2(condition func(c ChannelInfo) bool) ([]ChannelInfo, []ChannelInfo, bool) {

	this.SelectDB("")

	if cidlist := this.Keys("*"); cidlist != nil {

		all := make([]ChannelInfo, 0, len(cidlist))
		channels := make([]ChannelInfo, 0, len(cidlist))

		for _, cid := range cidlist {
			if info, exist := this.GetVisitChannel(cid); exist {
				all = append(all, info)
				if condition(info) {
					channels = append(channels, info)
				}
			}
		}

		return channels, all, true
	} else {
		return []ChannelInfo{}, []ChannelInfo{}, false
	}
}

func (this *Driver) SetVisitChannel(ch ChannelInfo, expiration time.Duration) {

	this.SelectDB("")

	bytes, err := json.Marshal(ch)
	if err != nil {
		log.Fatal("json.Marshal失敗")
	} else {
		this.Set(ch.Cid, bytes, expiration)
	}
}
