package services

import (
	"log"
	"lowkeydd-server/pgxdb"
	"os"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
)

func IsUseridExist(userid string) (*pgxpool.Pool, bool) {
	var msg interface{}
	pgxpool, row := pgxdb.QueryRow(`SELECT * FROM auth WHERE userid = $1 ;`, userid)
	if err := row.Scan(&msg); err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			log.Printf("會員名稱: %s 還沒有人使用過\n", userid)
			return pgxpool, false
		} else if strings.Contains(err.Error(), "number of field descriptions must equal number of destinations,") {
			log.Printf("會員名稱: %s 已被註冊過\n", userid)
			return pgxpool, true
		} else {
			log.Printf("[IsUseridExist] %v", err)
			os.Exit(1)
		}
	}
	panic("IsUseridExist 發生例外狀況")
}

func IsVisitExist(userid string) (*pgxpool.Pool, bool, []byte) {

	var data []byte

	pgxpool, row := pgxdb.QueryRow("SELECT data FROM visit WHERE userid = $1 ;", userid)

	if err := row.Scan(&data); err != nil {

		if strings.Contains(err.Error(), "no rows in result set") {

			log.Printf("會員 %s 的Visit還沒設定\n", userid)
			return pgxpool, false, []byte{}

		} else if strings.Contains(err.Error(), "number of field descriptions must equal number of destinations,") {

			log.Printf("會員 %s 的Visit已經設定\n", userid)
			return pgxpool, true, data

		} else {

			log.Printf("[IsVisitExist] %v", err)
			os.Exit(1)
		}
	}

	return pgxpool, true, data
}

// func QueryByUserId(userid string) {
// 	var id uint16
// 	var user string
// 	var pw string
// 	var is_del bool
// 	var date time.Time
// 	var visit interface{}

// 	pgxpool, row := QueryRow(`SELECT * FROM users WHERE userid = $1 ;`, userid)
// 	err := row.Scan(&id, &user, &pw, &is_del, &date, &visit)
// 	if err != nil {
// 		log.Printf("QueryRow failed: %v\n", err)
// 		os.Exit(1)
// 	}
// 	log.Printf("%d %s %s %t %s %v", id, user, pw, is_del, date, visit)
// 	defer pgxpool.Close()
// }
