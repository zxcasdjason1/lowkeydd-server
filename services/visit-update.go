package services

import (
	"encoding/json"
	"log"
	"lowkeydd-server/pgxdb"
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
	data := c.DefaultPostForm("visit", "")
	log.Printf("username %v", userid)

	req := VisitUpdateRequest{
		UserID: userid,
		Visit:  []byte(data),
	}

	resp := updateVisit(req)
	log.Printf("更新Visit:> %v\n", resp.Visit)

	VisitUpdateTransPort(c, resp)
}

func study() {

	// 收到前端傳來的資料格式。
	s := []byte(`{
		"list": [
			{
				"cid": " UCJFZiqLMntJufDCHc6bQixg",
				"owner": "hololive ホロライブ - VTuber Group"
			},
			{
				"cid": "UCp6993wxpyDPHUpavwDFqgg",
				"owner": "SoraCh. ときのそらチャンネル"
			},
			{
				"cid": "UCp6993wxpyDPHUpavwDFqgg",
				"owner": "SoraCh. ときのそらチャンネル"
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

func updateVisit(req VisitUpdateRequest) VisitUpdateResponse {

	var v VisitList

	// 檢驗是否符合正確的JSON結構
	if err := json.Unmarshal(req.Visit, &v); err != nil || v.List == nil {
		log.Printf("輸入的JSON格式錯誤 \n %s", err)
		return VisitUpdateResponse{
			Code:  "error",
			Visit: v,
		}
	}

	// 如果該userid不存在就返回失敗訊息
	if pgxpool, isExist := IsUseridExist(req.UserID); !isExist {
		log.Printf("找不到對應使用者名稱: %s\n", req.UserID)
		defer pgxpool.Close()
		return VisitUpdateResponse{
			Code:  "error",
			Visit: v,
		}
	}

	// 檢查使用者的Visit是否存在
	if pgxpool, isExist, _ := IsVisitExist(req.UserID); !isExist {
		log.Printf("使用者%s的Visit 沒有紀錄 \n", req.UserID)
		defer pgxpool.Close()
		return VisitUpdateResponse{
			Code:  "error",
			Visit: v,
		}
	}

	// 轉換成壓縮過的Raw資料
	buf, err := json.Marshal(v)
	if err != nil {
		log.Printf("輸入的JSON格式錯誤 \n %s", err)
		return VisitUpdateResponse{
			Code:  "error",
			Visit: v,
		}
	}
	log.Printf("buf:\n %s\n", buf)

	pgxpool, row := pgxdb.QueryRow(`UPDATE visit SET data = $2 WHERE userid = $1 returning (userid);`, req.UserID, buf)

	var userid string
	if err := row.Scan(&userid); err != nil {
		log.Printf("sql錯誤,寫入失敗....\n %s", err)
		defer pgxpool.Close()
		return VisitUpdateResponse{
			Code:  "error",
			Visit: v,
		}
	} else {
		log.Printf("sql正確,寫入成功.... %s", userid)
		defer pgxpool.Close()
		return VisitUpdateResponse{
			Code:  "success",
			Visit: v,
		}
	}
}

func VisitUpdateTransPort(c *gin.Context, resp VisitUpdateResponse) {
	switch resp.Code {
	case "success":
		c.JSON(200, gin.H{"code": resp.Code, "visit": resp.Visit})
		return
	case "failure":
		c.JSON(200, gin.H{"code": resp.Code, "visit": resp.Visit})
		return
	case "error":
		c.JSON(400, gin.H{"code": resp.Code, "visit": resp.Visit})
		return
	}
}
