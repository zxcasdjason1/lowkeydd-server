package redisdb

import (
	"encoding/json"
	"fmt"
	"log"
	"lowkeydd-server/share"
	"time"
)

func (this *Driver) SetSession(userid string, ssid string, timeout time.Duration) {

	session := fmt.Sprintf(`{"ssid":"%s","timeout":"%d"}`, ssid, timeout)
	log.Printf("session :> %s\n", session)

	this.SelectDB("ssid")
	this.Set(userid, []byte(session), timeout*time.Second)
}

func (this *Driver) GetSession(userid string) (share.SessionValue, bool) {

	this.SelectDB("ssid")

	data := this.Get(userid)
	log.Printf("session :> %v\n", data)

	var sv share.SessionValue
	if data == "" {
		return sv, false
	}
	if err := json.Unmarshal([]byte(data), &sv); err != nil {
		return sv, false
	} else {
		return sv, true
	}

}
