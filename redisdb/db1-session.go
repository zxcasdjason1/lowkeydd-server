package redisdb

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type SessionValue struct {
	SSID       string `json:"ssid"`
	Expiration string `json:"expiration"`
}

func (this *Driver) SetSession(userid string, ssid string, timeout time.Duration) {

	session := fmt.Sprintf(`{"ssid":"%s","expiration":"%d"}`, ssid, timeout)
	log.Printf("session :> %s\n", session)

	this.SelectDB("ssid")
	this.Set(userid, []byte(session), timeout*time.Second)
}

func (this *Driver) GetSession(userid string) (SessionValue, bool) {

	this.SelectDB("ssid")

	data := this.Get(userid)
	log.Printf("session :> %v\n", data)

	var sv SessionValue
	if data == "" {
		return sv, false
	}
	if err := json.Unmarshal([]byte(data), &sv); err != nil {
		return sv, false
	} else {
		return sv, true
	}

}
