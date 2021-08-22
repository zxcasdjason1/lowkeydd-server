package redisdb

import (
	"log"
	. "lowkeydd-server/share"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

type Driver struct {
	client *redis.Client
}

type Setting struct {
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Password string `json:"password"`
	DBIndex  int    `json:"dbindex"`
}

var (
	lock    = &sync.Mutex{}
	lock2   = &sync.Mutex{}
	driver  *Driver
	setting *Setting
)

func GetInstance() *Driver {
	if driver == nil {
		// 只允許一個goroutine訪問
		lock.Lock()
		defer lock.Unlock()
		if driver == nil {
			driver = &Driver{}
		}
	}
	return driver
}

func (this *Driver) SelectDB(dbName string) {

	lock2.Lock()
	defer lock2.Unlock()

	var dbindex int
	switch dbName {
	case "ssid":
		dbindex = 1
	case "visit":
		dbindex = 0
	case "":
		dbindex = 0
	case "search":
		dbindex = 2
	default:
		dbindex = 0
	}
	if setting.DBIndex != dbindex {
		log.Printf("[RedisBD] 切換資料庫目錄: %d \n", dbindex)
		setting.DBIndex = dbindex
		driver.client.Do("SELECT", dbindex)
	}
}

func (this *Driver) Keys(p string) []string {
	return this.client.Keys(p).Val()
}

func (this *Driver) Connect() {

	JSONFileLoader("setting/redis.json", &setting)

	log.Println("[RedisBD] 創建資料庫驅動器")
	ip := setting.IP           // localhost 為預設本地連線，若為遠端請自行輸入
	port := setting.Port       // "6379" 為預設port
	passwd := setting.Password // "" 表示無密碼進入
	DB := setting.DBIndex      // 0  表示預設的資料庫

	this.client = redis.NewClient(&redis.Options{
		Addr:     ip + ":" + port,
		Password: passwd,
		DB:       DB,
	})

	pong, err := this.client.Ping().Result()
	if err != nil {
		panic("[RedisBD] 資料庫連線狀態異常")
	}

	log.Println("[RedisBD] 資料庫連線成功 : ", pong)
}

func (this *Driver) Set(key string, val []byte, expiration time.Duration) {

	err := this.client.Set(key, val, expiration).Err() // => SET key value 0 數字代表過期秒數，在這裡0為永不過期
	if err != nil {
		panic(err)
	}

	log.Printf("[RedisBD] 保存數據: %s..............................寫入成功 ok\n", key)
}

func (this *Driver) Get(key string) string {

	val, err := this.client.Get(key).Result() // => GET key
	if err != nil {
		log.Printf("[RedisBD] 查無此數據: %s............................讀取失敗 fail\n", key)
		return ""
	}

	// log.Printf("[RedisBD] 取得數據: %s..............................讀取成功 ok\n", key)
	return val
}

func (this *Driver) Del(key string) int64 {

	val, err := this.client.Del(key).Result() // => GET key
	if err != nil {
		// log.Printf("[RedisBD] 查無此數據: %s............................刪除失敗 fail\n", key)
		return -1
	}
	return val
}
