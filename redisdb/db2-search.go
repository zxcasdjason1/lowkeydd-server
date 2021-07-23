package redisdb

import (
	"encoding/json"
	"log"
	. "lowkeydd-server/share"
	"time"
)

func (this *Driver) SetSearchChannel(ch ChannelInfo, expiration time.Duration) {

	this.SelectDB("search")

	bytes, err := json.Marshal(ch)
	if err != nil {
		log.Fatal("json.Marshal失敗")
	} else {
		this.Set(ch.Cid, bytes, expiration)
	}

	this.SelectDB("")
}

func (this *Driver) GetSearchChannel(cid string) (ChannelInfo, bool) {

	this.SelectDB("search")

	ch := ChannelInfo{}
	if jsonStr := this.Get(cid); jsonStr != "" {
		json.Unmarshal([]byte(jsonStr), &ch)
		log.Printf("GetSearchChannel %v\n", ch)
		this.SelectDB("")
		return ch, true
	} else {
		this.SelectDB("")
		return ch, false
	}

}
