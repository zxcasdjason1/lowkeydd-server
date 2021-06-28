package services

import (
	"encoding/json"
	"log"
	"lowkeydd-server/pgxdb"
	. "lowkeydd-server/share"
)

// 當你按下新增建時，會先檢查資料庫有沒有Vist，沒有就創建一份空白新的
func GetVisitList(userid string) (string, VisitList) {

	// 如果該userid不存在就返回失敗訊息
	if pgxpool, isExist := IsUseridExist(userid); !isExist {
		log.Printf("找不到對應使用者名稱: %s\n", userid)
		defer pgxpool.Close()
		return "error", VisitList{}
	}

	if pgxpool, isExist, data := IsVisitExist(userid); isExist {
		log.Printf("使用者%s的Visit 已有紀錄 \n", userid)
		var visit VisitList
		json.Unmarshal(data, &visit)
		defer pgxpool.Close()
		return "success", visit
	}

	pgxpool, err := pgxdb.Exec(`INSERT INTO visit (userid) values ($1) returning (userid) ;`, userid)
	if err != nil {
		log.Printf("使用者%s的Visit 創建失敗 \n", userid)
		defer pgxpool.Close()
		return "error", VisitList{}
	}
	log.Printf("使用者%s的Visit 創建成功 \n", userid)
	defer pgxpool.Close()
	return "success", VisitList{}
}
