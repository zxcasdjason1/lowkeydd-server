package services

import (
	"fmt"
	"log"
	"lowkeydd-server/redisdb"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LoginRequest struct {
	UserID string `json:"username"`
	Passwd string `json:"password"`
}

type LoginResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

// type SessionValue struct {
// 	UserID  string `json:"userid"`
// 	Timeout string `json:"timeout"`
// }

func LoginEndpoint(c *gin.Context) {
	userid := c.DefaultPostForm("username", "")
	passwd := c.DefaultPostForm("password", "")
	log.Printf("username %v", userid)
	log.Printf("password %v", passwd)

	resp := GetLoginMessage(userid, passwd)
	log.Printf("登入回應訊息:> %v\n", resp.Msg)

	// 寫入ssid
	if resp.Code == "success" {

		const timeout = 3600 //逾時時間

		if ssid, _ := c.Cookie("ssid"); ssid == "" {
			// cookie 還沒有設定的話
			if uid, err := uuid.NewUUID(); err == nil {
				// cookie也設
				c.SetCookie("ssid", uid.String(), timeout, "/", "localhost", false, true)
				// redis也設
				redisdb.GetInstance().SetSession(uid.String(), userid, timeout)
			} else {
				log.Fatal(err)
			}
		} else {
			log.Printf("ssid: %s ,已經設定了 \n", ssid)
			// 刷新存活時間
			// cookie也設
			c.SetCookie("ssid", ssid, timeout, "/", "localhost", false, true)
			// redis也設
			redisdb.GetInstance().SetSession(ssid, userid, timeout)
		}
	}

	LoginTransPort(c, resp)
}

func LoginTransPort(c *gin.Context, resp LoginResponse) {
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

func GetLoginMessage(userid string, passwd string) LoginResponse {

	// 空白檢查
	if userid == "" {
		log.Printf("用戶名不可空白\n")
		return LoginResponse{
			Code: "failure",
			Msg:  "用戶名不可空白",
		}
	}
	if passwd == "" {
		log.Printf("密碼不可空白\n")
		return LoginResponse{
			Code: "failure",
			Msg:  "密碼不可空白",
		}
	}

	// 如果該userid已經存在就返回失敗訊息
	if pgxpool, pass := checkAuthPass(userid, passwd); pass {
		defer pgxpool.Close()
		log.Printf("登入成功\n")
		return LoginResponse{
			Code: "success",
			Msg:  fmt.Sprintf("用戶: %s 登入成功", userid),
		}
	}

	log.Printf("登入失敗\n")
	return LoginResponse{
		Code: "failure",
		Msg:  "用戶名或密碼輸入錯誤",
	}
}
