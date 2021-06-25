package redisdb

import (
	"log"
	"sync"

	"github.com/go-redis/redis"
)

type Driver struct {
	client  *redis.Client
	Setting *Setting
}

type Setting struct {
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Password string `json:"password"`
	DBIndex  int    `json:"dbindex"`
}

var lock = &sync.Mutex{}
var driver *Driver

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

func (this *Driver) GetClient() *redis.Client {
	return driver.client
}

func (this *Driver) Connect(setting *Setting) {

	log.Println("創建資料庫驅動器")
	this.Setting = setting
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
		panic("資料庫連線狀態異常")
	}

	log.Println("資料庫連線成功 : ", pong)
}

func (this *Driver) Set(key string, val []byte) {

	client := *this.client

	err := client.Set(key, val, 0).Err() // => SET key value 0 數字代表過期秒數，在這裡0為永不過期
	if err != nil {
		panic(err)
	}

	log.Printf("保存數據: %s..............................寫入成功 ok\n", key)
}

func (this *Driver) Get(key string) string {

	client := *this.client

	val, err := client.Get(key).Result() // => GET key
	if err != nil {
		log.Printf("查無此數據: %s..............................讀取失敗 fail\n", key)
		return ""
	}

	log.Printf("取得數據: %s..............................讀取成功 ok\n", key)
	return val
}
