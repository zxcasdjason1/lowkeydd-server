package main

import (
	"lowkeydd-crawler/crawlers"
	"lowkeydd-crawler/redisdb"
	"lowkeydd-crawler/services"
	. "lowkeydd-crawler/share"

	"github.com/gin-gonic/gin"
)

var (
	visitList    *VisitList
	redisSetting *redisdb.Setting
)

const (
	Interval = 10000
	During   = 300000 //單位為ms
)

func main() {

	// 載入設定檔
	JSONFileLoader("setting/redis.json", &redisSetting)
	JSONFileLoader("setting/visit.json", &visitList)

	// 建立 Redis Diver，並透過設定檔取得連線
	redisdb.GetInstance().Connect(redisSetting)

	//配置爬蟲
	crawlers.GetInstance(visitList)

	// 設定GIN路由器
	router := gin.Default()

	// 解決Cors問題
	router.Use(CORSMiddleware())

	// 配置API

	router.GET("/crawler/update/", services.CrawlerUpdate)

	router.GET("/crawler/visitall/", services.CrawlerVisitAll)

	router.POST("/crawler/visit/", services.CrawlerVisit)

	router.Run(":8002")

	// // 流程: 先執行一次 VisitAll, 之後間隔再循環執行
	// crawlers := crawlers.GetInstance(visitList)
	// crawlers.VisitAll()

	// 啟動 Schedule 循環
	// for {

	// 	// 設置流程
	// 	remaining := During
	// 	scheUpdate := func() {
	// 		remaining -= Interval
	// 		if remaining <= 0 {
	// 			remaining = 0
	// 		}
	// 		log.Printf("Event remaining, %d", remaining)
	// 	}
	// 	scheEnd := func() {
	// 		crawlers.VisitAll()
	// 	}

	// 	sche := NewSchedule(Interval * time.Millisecond)
	// 	sche.Event.AddListener(SCHE_UPDATE, scheUpdate)
	// 	sche.Event.AddListener(SCHE_END, scheEnd)
	// 	time.Sleep(time.Millisecond * During)
	// 	sche.Stop()
	// 	sche.Event.RemoveListener(SCHE_UPDATE, scheUpdate)
	// 	sche.Event.RemoveListener(SCHE_END, scheEnd)
	// }

}

// func CORSMiddleware(c *gin.Context) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		h := &http.Header{}
// 		h.Add("Content-Type", "application/json")
// 		h.Add("Access-Control-Allow-Origin", "*")
// 		h.Add("Access-Control-Max-Age", "86400")
// 		h.Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
// 		h.Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
// 		h.Add("Access-Control-Allow-Credentials", "true")
// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(200)
// 		} else {
// 			c.Next()
// 		}
// 	}
// }
