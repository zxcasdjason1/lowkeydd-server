package redisdb

import (
	"encoding/json"
	"fmt"
	"lowkeydd-server/share"
	"time"
)

// type SessionValue struct {
// 	UserID  string `json:"userid"`
// 	Timeout string `json:"timeout"`
// }

func (this *Driver) SetSession(ssid string, userid string, timeout time.Duration) {

	session := fmt.Sprintf(`{"userid":"%s","timeout":"%d"}`, userid, timeout)
	// log.Printf("session :> %s\n", session)

	this.SelectDB("ssid")
	this.Set(ssid, []byte(session), timeout*time.Second)
}

func (this *Driver) GetSession(ssid string) (share.SessionValue, bool) {

	this.SelectDB("ssid")

	data := this.Get(ssid)
	// log.Printf("session :> %v\n", data)

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
