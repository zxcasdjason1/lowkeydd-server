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
		log.Printf("userid is required")
		VisitEditTransPort(c, "failure", nil, "userid is required")
		return
	}

	if ssid == "" {
		log.Printf("userid is required")
		VisitEditTransPort(c, "failure", nil, "ssid is required")
		return
	}

	if s, success := redisdb.GetInstance().GetSession(userid); success && s.SSID == ssid {
		log.Printf("ssid:> %s 驗證成功", ssid)

		if code, visit := GetVisitList(userid); code == "success" {
			log.Printf("瀏覽追隨者清單:> %v\n", visit)
			VisitEditTransPort(c, "success", &visit, "auth is success")
		} else {
			// 取得Visitlist失敗
			log.Printf("取得Visitlist失敗\n")
			VisitEditTransPort(c, "error", nil, "visit edit method is broken")
		}
	} else {
		log.Printf("ssid:> %s 驗證失敗", ssid)
		VisitEditTransPort(c, "failure", nil, "auth is fail")
	}

}

func VisitEditTransPort(c *gin.Context, code string, visit *VisitList, msg string) {
	switch code {
	case "success":
		c.JSON(200, gin.H{"code": code, "visit": visit, "msg": msg})
		return
	case "failure":
		c.JSON(200, gin.H{"code": code, "visit": visit, "msg": msg})
		return
	case "error":
		c.JSON(400, gin.H{"code": code, "visit": visit, "msg": msg})
		return
	}
}
