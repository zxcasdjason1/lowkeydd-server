package main

import (
	"fmt"
	"log"
	"lowkeydd-server/consul"
	"lowkeydd-server/crawlers"
	"lowkeydd-server/pgxdb"
	"lowkeydd-server/redisdb"
	"lowkeydd-server/services"
	"lowkeydd-server/share"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {

	// 建立 Redis Diver，並透過設定檔取得連線
	redisdb.GetInstance().Connect()
	pgxdb.NewDriver()

	//配置爬蟲
	serviceType := strings.ToUpper(os.Getenv("SERVICE_TYPE"))
	if serviceType == "SERVER" {
		crawlers.NewCrawlers()
		crawlers.GetInstance().UnChecked_Update()
	}

	// 設定GIN路由器
	router := gin.Default()

	// 解決Cors問題
	router.Use(share.CORSMiddleware())

	// health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Crawler僅以SERVER擔任工作，由後臺控制。
	if serviceType == "SERVER" {

		router.GET("/crawler/reload", services.CrawlerReload)

		router.GET("/crawler/visitall", services.CrawlerVisitAll)

		router.GET("/crawler/update", services.CrawlerUpdate)

		router.GET("/crawler/:method/:cid", services.CrawlerVisit)
	}

	router.GET("/cookie/:key/:value", services.SetCookie)

	router.GET("/cookie/:key", services.GetCookie)

	// GetChannels前端API
	router.POST("/channels/search", services.GetSearchChannelResponse)

	router.GET("/channels/all", services.GetAllVisitChannelsResponse)

	router.GET("/channels/:tag", services.GetTagedVisitChannelEndpoint)

	// pgx-GetVisits
	router.POST("/visit/edit", services.VisitEditEndpoint)

	router.POST("/visit/update", services.VisitUpdateEndpoint)

	router.POST("/auth/login", services.LoginEndpoint)

	router.POST("/auth/register", services.RegisterEndpoint)

	router.POST("/auth/logout", services.LogoutEndpoint)

	// router.POST("/letsdd/v1", services.LetsddEndpoint)
	router.POST("/letsdd/v2", services.Letsddv2Endpoint)

	// router.Run(":8000")

	cs := consul.GetInstance()
	cs.RegisterService()

	errChan := make(chan error)
	go (func() {
		err := router.Run(":8000")
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
}
