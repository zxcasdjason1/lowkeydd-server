package main

import (
	"fmt"
	"log"
	"lowkeydd-server/consul"
	"lowkeydd-server/crawlers"
	"lowkeydd-server/pgxdb"
	"lowkeydd-server/redisdb"
	"lowkeydd-server/services"
	. "lowkeydd-server/share"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {

	// 建立 Redis Diver，並透過設定檔取得連線
	redisdb.GetInstance().Connect()
	pgxdb.NewDriver()

	//配置爬蟲
	crawlers.NewCrawlers()

	// 設定GIN路由器
	router := gin.Default()

	// 解決Cors問題
	router.Use(CORSMiddleware())

	// health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Crawler後臺控制
	router.GET("/crawler/reload/", services.CrawlerReload)

	router.GET("/crawler/visitall/", services.CrawlerVisitAll)

	router.GET("/crawler/update/", services.CrawlerUpdate)

	router.GET("/crawler/:method/:cid", services.CrawlerVisit)

	// GetChannels前端API
	router.POST("/channels/search/", services.GetSearchChannelResponse)

	router.GET("/channels/all", services.GetAllChannelsResponse)

	router.GET("/channels/:tag", services.GetTagedChannelsResponse)

	// pgx
	router.POST("/login/", services.LoginEndpoint)

	router.GET("/cookie/:key/:value", services.SetCookie)
	router.GET("/cookie/:key", services.GetCookie)

	router.POST("/register/", services.RegisterEndpoint)

	router.POST("/visit/edit", services.VisitEditEndpoint)

	router.POST("/visit/update", services.VisitUpdateEndpoint)

	router.GET("/letsdd", services.LetsddEndpoint)

	cs := consul.GetInstance()
	cs.RegisterService()

	errChan := make(chan error)
	go (func() {
		err := router.Run(":8002")
		if err != nil {
			log.Println(err)
			errChan <- err
		}
	})()
	go (func() {
		sig_c := make(chan os.Signal, 1)

		// 信號名	value  說明
		// SIGINT    2    發送ctrl+c
		// SIGTERM   15   结束程序時
		signal.Notify(sig_c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-sig_c)
	})()
	// 主程序因為 errChan 已阻塞，直到收到 (如錯誤、程序停止、ctrl+c等)
	// 訊息才會繼續往下執行，藉此實現服務對 consult "優雅地退出"。
	getErr := <-errChan
	cs.KillService() //
	log.Println(getErr)

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
