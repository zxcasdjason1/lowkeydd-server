package services

import (
	"encoding/json"
	"log"
	"lowkeydd-server/pgxdb"
	"lowkeydd-server/redisdb"
	. "lowkeydd-server/share"

	"github.com/gin-gonic/gin"
)

type VisitUpdateRequest struct {
	UserID string `json:"username"`
	Visit  []byte `json:"visit"`
}

type VisitUpdateResponse struct {
	Code  string    `json:"code"`
	Visit VisitList `json:"visit"`
}

func VisitUpdateEndpoint(c *gin.Context) {
	userid := c.DefaultPostForm("username", "")
	ssid := c.DefaultPostForm("ssid", "")
	data := c.DefaultPostForm("visit", "")
	log.Printf("username :> %s\n", userid)
	log.Printf("ssid :> %s\n", ssid)
	log.Printf("data :> %s\n", data)

	if userid == "" {
		log.Printf("userid is required")
		VisitUpdateTransPort(c, "failure", nil, "userid is required")
		return
	}

	if ssid == "" {
		log.Printf("userid is required")
		VisitUpdateTransPort(c, "failure", nil, "ssid is required")
		return
	}

	if s, success := redisdb.GetInstance().GetSession(userid); success && s.SSID == ssid {
		log.Printf("ssid:> %s 驗證成功", ssid)

		req := VisitUpdateRequest{
			UserID: userid,
			Visit:  []byte(data),
		}
		if code, visit := updateVisit(req); code == "success" {
			log.Printf("更新Visit:> %v\n", visit)
			VisitUpdateTransPort(c, "success", &visit, "auth is success")
		} else {
			log.Printf("ssid:> %s 驗證失敗", ssid)
			VisitUpdateTransPort(c, "error", nil, "database is locked, or visit update method is broken")
		}
	} else {
		log.Printf("ssid:> %s 驗證失敗", ssid)
		VisitUpdateTransPort(c, "failure", nil, "auth is fail")
	}

}

func study() {

	// 收到前端傳來的資料格式。
	s := []byte(`{
		"list": [
			{
				"cid":"UCqm3BQLlJfvkTsX_hvm0UmA",
				"cname":"Watame Ch. 角巻わため",
				"owner":"Watame Ch. 角巻わため",
				"avatar":"https://yt3.ggpht.com/ytc/AKedOLRWpyqOZzCmuSfmKGNo8TD2L_IRUYSw1wyhHXw-=s88-c-k-c0x00ffffff-no-rj",
				"method":"youtube",
				"group":"Favorite"
			}
		]
	}`)
	log.Printf("[]byte s:\n %s\n", s)

	// 資料的剖析，透過型態轉換成結構體
	var v VisitList
	err := json.Unmarshal(s, &v)
	if err != nil {
		panic(err)
	}
	log.Printf("結構體:\n %v\n", v)

	// to raw data
	buf, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("buf:\n %s\n", buf)
}

func updateVisit(req VisitUpdateRequest) (string, VisitList) {

	var v VisitList

	// 檢驗是否符合正確的JSON結構
	if err := json.Unmarshal(req.Visit, &v); err != nil || v.List == nil {
		log.Printf("輸入的JSON格式錯誤 \n %s", err)
		return "error", v
	}

	// 如果該userid不存在就返回失敗訊息
	if pgxpool, isExist := IsUseridExist(req.UserID); !isExist {
		log.Printf("找不到對應使用者名稱: %s\n", req.UserID)
		defer pgxpool.Close()
		return "error", v
	}

	// 檢查使用者的Visit是否存在
	if pgxpool, isExist, _ := IsVisitExist(req.UserID); !isExist {
		log.Printf("使用者%s的Visit 沒有紀錄 \n", req.UserID)
		defer pgxpool.Close()
		return "error", v
	}

	// 轉換成壓縮過的Raw資料
	buf, err := json.Marshal(v)
	if err != nil {
		log.Printf("輸入的JSON格式錯誤 \n %s", err)
		return "error", v
	}
	log.Printf("buf:\n %s\n", buf)

	pgxpool, row := pgxdb.QueryRow(`UPDATE visit SET data = $2 WHERE userid = $1 returning (userid);`, req.UserID, buf)

	var userid string
	if err := row.Scan(&userid); err != nil {
		log.Printf("sql錯誤,寫入失敗....\n %s", err)
		defer pgxpool.Close()
		return "error", v
	} else {
		log.Printf("sql正確,寫入成功.... %s", userid)
		defer pgxpool.Close()
		return "success", v
	}
}

func VisitUpdateTransPort(c *gin.Context, code string, visit *VisitList, msg string) {
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
