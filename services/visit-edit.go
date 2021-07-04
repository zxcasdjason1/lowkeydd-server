package services

import (
	"log"
	"lowkeydd-server/redisdb"
	. "lowkeydd-server/share"

	"github.com/gin-gonic/gin"
)

type VisitEditRequest struct {
	UserID string `json:"username"`
}

type VisitEditResponse struct {
	Code  string    `json:"code"`
	Visit VisitList `json:"visit"`
}

func VisitEditEndpoint(c *gin.Context) {

	userid := c.DefaultPostForm("username", "")
	ssid := c.DefaultPostForm("ssid", "")
	log.Printf("username:> %v\n", userid)
	log.Printf("ssid :> %s\n", ssid)

	if userid == "" {
		c.JSON(200, gin.H{"code": "failure", "msg": "userid is required"})
		return
	}

	if ssid == "" {
		c.JSON(200, gin.H{"code": "failure", "msg": "get cookie fail"})
		return
	}

	if s, success := redisdb.GetInstance().GetSession(userid); success && s.SSID == ssid {
		log.Printf("ssid:> %s 驗證成功", ssid)

		code, visit := GetVisitList(userid)
		log.Printf("瀏覽追隨者清單:> %v\n", visit)
		VisitEditTransPort(c, code, &visit)
	} else {
		log.Printf("ssid:> %s 驗證失敗", ssid)
		VisitEditTransPort(c, "failure", nil)
	}

}

func VisitEditTransPort(c *gin.Context, code string, visit *VisitList) {
	switch code {
	case "success":
		c.JSON(200, gin.H{"code": code, "visit": visit})
		return
	case "failure":
		c.JSON(200, gin.H{"code": code, "visit": visit})
		return
	case "error":
		c.JSON(400, gin.H{"code": code, "visit": visit})
		return
	}
}
