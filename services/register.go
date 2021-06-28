package services

import (
	"fmt"
	"log"
	"lowkeydd-server/pgxdb"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	UserID string `json:"username"`
	Passwd string `json:"password"`
}

type RegisterResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

func RegisterEndpoint(c *gin.Context) {

	userid := c.DefaultPostForm("username", "")
	passwd := c.DefaultPostForm("password", "")
	log.Printf("username %v", userid)
	log.Printf("password %v", passwd)

	msg := GetRegisterMessage(userid, passwd)
	log.Printf("會員註冊資訊:> %s\n", msg.Msg)

	RegisterTransPort(c, msg)
}

func RegisterTransPort(c *gin.Context, resp RegisterResponse) {
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

func GetRegisterMessage(userid string, passwd string) RegisterResponse {

	// 特殊字元檢查
	if isTrue := IsContainsSpecial(userid); isTrue {
		log.Printf("用戶名只能為數字或英文\n")
		return RegisterResponse{
			Code: "failure",
			Msg:  "用戶名只能為數字或英文",
		}
	}

	// 如果該userid已經存在就返回失敗訊息
	if pgxpool, isExist := IsUseridExist(userid); isExist {
		defer pgxpool.Close()
		log.Printf("註冊失敗了\n")
		return RegisterResponse{
			Code: "failure",
			Msg:  "用戶名已經存在",
		}
	}

	// 如果該userid尚不存在，創建一個新的帳戶
	pgxpool, row := pgxdb.QueryRow("INSERT INTO auth (userid, passwd) values ($1,$2) returning (userid) ;", userid, passwd)

	var uid string
	if err := row.Scan(&uid); err == nil {
		defer pgxpool.Close()
		log.Printf("註冊成功了\n")
		return RegisterResponse{
			Code: "success",
			Msg:  fmt.Sprintf("用戶: %s 註冊成功", uid),
		}

	} else {

		if strings.Contains(err.Error(), "auth_passwd_check") {
			defer pgxpool.Close()
			log.Printf("無效密碼: %s\n", passwd)
			return RegisterResponse{
				Code: "failure",
				Msg:  "無效用戶密碼",
			}
		} else if strings.Contains(err.Error(), "auth_userid_check") {
			defer pgxpool.Close()
			log.Printf("無效會員名稱: %s\n", userid)
			return RegisterResponse{
				Code: "failure",
				Msg:  "無效會員名稱",
			}
		}
	}

	defer pgxpool.Close()
	return RegisterResponse{
		Code: "error",
		Msg:  "系統異常",
	}
}

func IsContainsSpecial(taget string) bool {
	// 特殊字元檢查
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	if replaced := reg.ReplaceAllString(taget, ""); replaced != taget {
		return true
	}
	return false
}
