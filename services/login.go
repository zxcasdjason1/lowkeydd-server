package services

import (
	"fmt"
	"log"
	"lowkeydd-server/redisdb"
	"lowkeydd-server/share"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LoginRequest struct {
	UserID string `json:"username"`
	Passwd string `json:"password"`
}

type LoginResponse struct {
	Code    string        `json:"code"`
	Session share.Session `json:"session"`
	Msg     string        `json:"msg"`
}

const (
	SSID_Expiration = 3600 //逾時時間
)

func getSSID() string {
	if uid, err := uuid.NewUUID(); err == nil {
		return uid.String()
	} else {
		log.Fatal("getUUID失敗, 錯誤如下: \n", err)
		return ""
	}
}

func LoginEndpoint(c *gin.Context) {
	userid := c.DefaultPostForm("username", "")
	passwd := c.DefaultPostForm("password", "")
	log.Printf("username %v", userid)
	log.Printf("password %v", passwd)

	resp := GetLoginMessage(userid, passwd)
	log.Printf("登入回應訊息:> %v\n", resp.Msg)

	LoginTransPort(c, resp)
}

func LoginTransPort(c *gin.Context, resp LoginResponse) {
	switch resp.Code {
	case "success":
		c.JSON(200, gin.H{"code": resp.Code, "session": resp.Session, "msg": resp.Msg})
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

		// 添加驗證session；必須重複登入時對session刷新時效。
		var ssid string
		if session, success := redisdb.GetInstance().GetSession(userid); success {
			// session已經存在時，表示為重複登入，則刷新session的時效
			ssid = session.SSID
			redisdb.GetInstance().SetSession(userid, ssid, SSID_Expiration)
		} else {
			// session不存在，表示首次登入，則生成session
			ssid = getSSID()
			redisdb.GetInstance().SetSession(userid, ssid, SSID_Expiration)
		}

		return LoginResponse{
			Code: "success",
			Msg:  fmt.Sprintf("用戶: %s 登入成功", userid),
			Session: share.Session{
				SSID:       ssid,
				Expiration: SSID_Expiration,
			},
		}
	}

	log.Printf("登入失敗\n")
	return LoginResponse{
		Code: "failure",
		Msg:  "用戶名或密碼輸入錯誤",
	}
}
