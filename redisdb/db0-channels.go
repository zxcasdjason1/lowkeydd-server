package redisdb

import (
	"encoding/json"
	"log"
	. "lowkeydd-server/share"
)

func (this *Driver) GetChannel(cid string) (ChannelInfo, bool) {

	this.SelectDB("")

	ch := ChannelInfo{}
	if jsonStr := this.Get(cid); jsonStr != "" {
		json.Unmarshal([]byte(jsonStr), &ch)
		return ch, true
	} else {
		return ch, false
	}
}

func (this *Driver) GetAllChannels() ([]ChannelInfo, bool) {

	this.SelectDB("")

	if cidlist := this.Keys("*"); cidlist != nil {

		channels := make([]ChannelInfo, 0, len(cidlist))

		for _, cid := range cidlist {
			if ch, exist := this.GetChannel(cid); exist {
				channels = append(channels, ch)
			}
		}

		return channels, true
	} else {
		return []ChannelInfo{}, false
	}
}

func (this *Driver) GetChannelsByCondition(condition func(c ChannelInfo) bool) ([]ChannelInfo, bool) {

	this.SelectDB("")

	if cidlist := this.Keys("*"); cidlist != nil {

		channels := make([]ChannelInfo, 0, len(cidlist))

		for _, cid := range cidlist {
			if info, exist := this.GetChannel(cid); exist {
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

func (this *Driver) GetChannelsByConditionV2(condition func(c ChannelInfo) bool) ([]ChannelInfo, []ChannelInfo, bool) {

	this.SelectDB("")

	if cidlist := this.Keys("*"); cidlist != nil {

		all := make([]ChannelInfo, 0, len(cidlist))
		channels := make([]ChannelInfo, 0, len(cidlist))

		for _, cid := range cidlist {
			if info, exist := this.GetChannel(cid); exist {
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

func (this *Driver) SetChannel(ch ChannelInfo) {

	this.SelectDB("")

	bytes, err := json.Marshal(ch)
	if err != nil {
		log.Fatal("json.Marshal失敗")
	} else {
		this.Set(ch.Cid, bytes, 0)
	}
}
