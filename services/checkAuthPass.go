package services

import (
	"log"
	"lowkeydd-server/pgxdb"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
)

func checkAuthPass(userid string, passwd string) (*pgxpool.Pool, bool) {

	pgxpool, row := pgxdb.QueryRow(`SELECT * FROM auth WHERE userid = $1 AND passwd = $2 ;`, userid, passwd)
	if err := row.Scan(nil); err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			log.Printf("找不匹配的對象\n")
			return pgxpool, false
		} else if strings.Contains(err.Error(), "number of field descriptions must equal number of destinations,") {
			log.Printf("找到多個匹配對象\n 異常登入請求:> 用戶名: %s, 密碼: %s", userid, passwd)
			return pgxpool, true
		} else {
			log.Printf("意外狀況導致登入失敗\n, 異常登入請求:> 用戶名: %s, 密碼: %s", userid, passwd)
			log.Printf("意外狀況導致登入失敗\n, %s", err)
			return pgxpool, false
		}
	}
	log.Printf("用戶: %s 登入成功", userid)
	return pgxpool, true
}
