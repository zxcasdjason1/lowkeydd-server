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
	SSID   string `json:"ssid"`
}

type LetsddResponse struct {
	Code     string        `json:"code"`
	Channels []ChannelInfo `json:"channels"`
}

func LetsddEndpoint(c *gin.Context) {

	userid := c.DefaultPostForm("username", "")
	ssid := c.DefaultPostForm("ssid", "")
	log.Printf("username :> %s\n", userid)
	log.Printf("ssid :> %s\n", ssid)

	if userid == "" {
		c.JSON(200, gin.H{"msg": "userid is required"})
		return
	}

	// 從cookie裡面去取得 ssid作為驗證
	if ssid == "" {
		c.JSON(200, gin.H{"msg": "get cookie fail"})
	}

	if s, success := redisdb.GetInstance().GetSession(userid); success && s.SSID == ssid {

		log.Printf("ssid:> %s 驗證成功", ssid)

		// 驗證成功，獲取該使用者visit
		if code, visit := GetVisitList(userid); code == "success" {
			crawlers.GetInstance().Checked_VisitByList(visit.List)
		}
	} else {
		log.Printf("ssid:> %s 驗證失敗", ssid)
	}

	GetAllChannelsResponse(c)

	// 驗證成功
	// 自定義的channels

	// 要是沒有驗證到
	// redis目前的channels
}
