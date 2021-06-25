package youtube

import (
	"encoding/json"
	"log"
	"lowkeydd-crawler/redisdb"
	. "lowkeydd-crawler/share"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

type Crawler struct {
	collector *colly.Collector
	List      []VisitItem
}

func NewCrawler(v *VisitList) *Crawler {

	this := &Crawler{
		List:      v.List,
		collector: colly.NewCollector(),
	}

	this.collector.AllowURLRevisit = true

	this.collector.OnHTML("script", func(h *colly.HTMLElement) {

		isValid, ytJSONstr := getYtInitialData(h.Text)

		if isValid {

			// 以youtube頻道的URL作為儲存資訊時的Key值。
			// 當 CreateChannelInfo 失敗時才直接從輸入的url取。
			info := CreateChannelInfo(ytJSONstr)
			if info.Cid == "" {
				if channelID := getCidByUrl(h.Request.URL.String()); channelID == "" {
					panic("獲取CID失敗，無法設定KEY值")
				} else {
					info.Cid = channelID
				}
			}
			info.Method = "youtube"
			info.UpdateTime = GetNextUpdateTime()

			// 打印出獲取到的頻道資訊
			// log.Println("========================================")
			// log.Printf("Cid:>>> 		%s", info.Cid)
			// log.Printf("Status:>>> 		%s", info.Status)
			// log.Printf("Owner:>>> 		%s", info.Owner)
			// log.Printf("Avatar:>>> 		%s", info.Avatar)
			// log.Printf("RenderType:>>> 	%s", info.RenderType)
			// log.Printf("StreamURL:>>>  	%s", info.StreamURL)
			// log.Printf("Thumbnail:>>>  	%s", info.Thumbnail)
			// log.Printf("Title:>>>      	%s", info.Title)
			// log.Printf("ViewCount:>>>  	%s", info.ViewCount)
			// log.Printf("StartTime:>>>  	%s", info.StartTime)
			// log.Println("========================================")

			// 寫入到 Redis中
			bytes, err := json.Marshal(info)
			if err != nil {
				log.Fatal("json.Marshal失敗")
				panic(err)
			} else {
				redisdb.GetInstance().Set(info.Cid, bytes)
			}
		}
	})

	this.collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.85 Safari/537.36")
	})

	return this
}

func (c *Crawler) Visit(cid string) {
	c.collector.Visit("https://www.youtube.com/channel/" + cid)
	log.Printf("[Youtube] cid :> %v", cid)
}

func getCidByUrl(urlStr string) string {
	m, _ := regexp.MatchString("/channel/", urlStr)
	if m {
		st := strings.Index(urlStr, "/channel/") + 9
		ed := len(urlStr)
		urlStr = urlStr[st:ed]
	}
	return urlStr
}

func getYtInitialData(ytStr string) (matched bool, result string) {
	// fetch ytInitialData str from response
	m, _ := regexp.MatchString("var ytInitialData", ytStr)
	if m {
		st := strings.Index(ytStr, "{")
		ed := len(ytStr) - 1
		for ytStr[ed] != ('}') {
			ed--
		}
		ytStr = ytStr[st : ed+1]
	}
	return m, ytStr
}
