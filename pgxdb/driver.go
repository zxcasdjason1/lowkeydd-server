package pgxdb

import (
	"context"
	"fmt"
	"log"
	. "lowkeydd-server/share"
	"os"
	"sync"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// type Driver struct {
// }

type Setting struct {
	IP       string `json:"ip"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

var (
	lock = &sync.Mutex{}
	// driver  *Driver
	setting *Setting
)

// func GetInstance() *Driver {
// 	if driver == nil {
// 		// 只允許一個goroutine訪問
// 		lock.Lock()
// 		defer lock.Unlock()
// 		if driver == nil {
// 			driver = &Driver{}
// 		}
// 	}
// 	return driver
// }

func NewDriver() {

	JSONFileLoader("setting/postgres.json", &setting)

	// 照理說應該跟 SERVICE_IP 區隔開比較好，但是目前service與Postgres資料庫都在同一台機器裡。
	if serviceIP := os.Getenv("SERVICE_IP"); serviceIP != "" {
		setting.IP = serviceIP
		log.Printf("[POSTGRES] SERVICE_IP :> %s \n", serviceIP)
	}

	// connent to psql database
	pgxpool, err := pgxpool.Connect(context.Background(), GetEnv())
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pgxpool.Close()

	log.Println("[POSTGRES] connection success")
}

func GetEnv() string {

	env := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		setting.IP,
		setting.Port,
		setting.User,
		setting.Password,
		setting.DBName,
	)
	return env
}

func QueryRow(sqlStmt string, argv ...interface{}) (*pgxpool.Pool, pgx.Row) {

	// driver := GetInstance()
	ctx := context.Background()

	pgxpool, err := pgxpool.Connect(ctx, GetEnv())
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	row := pgxpool.QueryRow(ctx, sqlStmt, argv...)

	return pgxpool, row
}

func Exec(sqlStmt string, argv ...interface{}) (*pgxpool.Pool, error) {

	// driver := GetInstance()
	ctx := context.Background()

	pgxpool, err := pgxpool.Connect(ctx, GetEnv())
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	_, err = pgxpool.Exec(ctx, sqlStmt, argv...)
	return pgxpool, err
}
