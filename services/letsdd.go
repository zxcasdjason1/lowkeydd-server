package services

import (
	"log"
	"lowkeydd-server/crawlers"
	"lowkeydd-server/redisdb"
	. "lowkeydd-server/share"

	"github.com/gin-gonic/gin"
)

type LetsddRequest struct {
	UserID string `json:"username"`
}

type LetsddResponse struct {
	Code     string        `json:"code"`
	Channels []ChannelInfo `json:"channels"`
}

func LetsddEndpoint(c *gin.Context) {
	// 從cookie裡面去取得 ssid作為驗證
	ssid, err := c.Cookie("ssid")
	if err != nil {
		c.JSON(200, gin.H{"msg": "get cookie fail"})
		return
	}
	log.Printf("ssid :> %s\n", ssid)
	if session, success := redisdb.GetInstance().GetSession(ssid); success {

		log.Printf("ssid:> %s 驗證成功", ssid)
		log.Printf("userid : %s , timeout : %s ", session.UserID, session.Timeout)

		// 驗證成功，獲取該使用者visit
		if code, visit := GetVisitList(session.UserID); code == "success" {
			crawlers.GetInstance().Checked_VisitByList(visit.List)
		}

		// 根據list去顯示response
		GetAllChannelsResponse(c)

	} else {
		log.Printf("ssid:> %s 驗證失敗", ssid)
		GetSingleChannelResponse(c, "")
	}

	// 驗證成功
	// 自定義的channels

	// 要是沒有驗證到
	// redis目前的channels
}
