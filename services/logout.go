package services

import (
	"fmt"
	"log"
	"lowkeydd-server/redisdb"

	"github.com/gin-gonic/gin"
)

type LogoutRequest struct {
	UserID string `json:"username"`
	SSID   string `json:"ssid"`
}

type LogoutResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

func LogoutEndpoint(c *gin.Context) {
	userid := c.DefaultPostForm("username", "")
	ssid := c.DefaultPostForm("ssid", "")
	log.Printf("username %v", userid)
	log.Printf("ssid %v", ssid)

	resp := GetLogoutMessage(userid, ssid)
	log.Printf("登出回應訊息:> %v\n", resp.Msg)

	LogoutTransPort(c, resp)
}

func GetLogoutMessage(userid string, ssid string) LogoutResponse {

	// 空白檢查
	if userid == "" {
		log.Printf("用戶名不可空白\n")
		return LogoutResponse{
			Code: "failure",
			Msg:  "用戶名不可空白",
		}
	}
	if ssid == "" {
		log.Printf("ssid不可空白\n")
		return LogoutResponse{
			Code: "failure",
			Msg:  "ssid不可空白",
		}
	}

	if s, success := redisdb.GetInstance().GetSession(userid); success && s.SSID == ssid {
		log.Printf("ssid:> %s 驗證成功", ssid)

		// remove 該session
		numOfDeleted := redisdb.GetInstance().Del(userid)
		log.Printf("numOfDeleted %d", numOfDeleted)
		if numOfDeleted > 0 {
			return LogoutResponse{
				Code: "success",
				Msg:  fmt.Sprintf("用戶: %s 的Session 刪除...成功", userid),
			}
		} else {
			return LogoutResponse{
				Code: "error",
				Msg:  fmt.Sprintf("用戶: %s 的Session 刪除...失敗", userid),
			}
		}

	} else {
		return LogoutResponse{
			Code: "failure",
			Msg:  "Session 驗證失敗",
		}
	}
}

func LogoutTransPort(c *gin.Context, resp LogoutResponse) {
	switch resp.Code {
	case "success":
		c.JSON(200, gin.H{"code": resp.Code, "msg": resp.Msg})
		return
	case "failure":
		c.JSON(200, gin.H{"code": resp.Code, "msg": resp.Msg})
		return
	case "error":
		c.JSON(400, gin.H{"code": resp.Code, "msg": resp.Msg})
		return
	}
}
