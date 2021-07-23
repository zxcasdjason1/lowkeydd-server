package youtube

import (
	"log"
	"lowkeydd-server/redisdb"
	. "lowkeydd-server/share"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

var (
	result = make(map[string]ChannelInfo)
	lock   = &sync.Mutex{}
)

type Crawler struct {
	List      []VisitItem
	collector *colly.Collector // for redis
}

type Authenticator struct {
}

func NewCrawler(v *VisitList, auth *Authenticator) *Crawler {

	this := &Crawler{
		List:      v.List,
		collector: colly.NewCollector(),
	}

	this.collector.AllowURLRevisit = true

	this.collector.OnHTML("script", func(h *colly.HTMLElement) {

		if isValid, ytJSONstr := getYtInitialData(h.Text); isValid {

			// 以youtube頻道的URL作為儲存資訊時的Key值。
			// 當 CreateChannelInfo 失敗時才直接從輸入的url取。
			ch := CreateChannelInfo(ytJSONstr)
			if ch.Cid == "" {
				if channelID := getCidByUrl(h.Request.URL.String()); channelID == "" {
					panic("獲取CID失敗，無法設定KEY值")
				} else {
					ch.Cid = channelID
				}
			}
			ch.Method = "youtube"
			ch.UpdateTime = GetNextUpdateTime()

			// 打印出獲取到的頻道資訊
			// log.Println("========================================")
			// log.Printf("Cid:>>> 		%s", ch.Cid)
			// log.Printf("Cname:>>> 		%s", ch.Cname)
			// log.Printf("Owner:>>> 		%s", ch.Owner)
			// log.Printf("Status:>>> 		%s", ch.Status)
			// log.Printf("Avatar:>>> 		%s", ch.Avatar)
			// log.Printf("RenderType:>>> 	%s", ch.RenderType)
			// log.Printf("StreamURL:>>>  	%s", ch.StreamURL)
			// log.Printf("Thumbnail:>>>  	%s", ch.Thumbnail)
			// log.Printf("Title:>>>      	%s", ch.Title)
			// log.Printf("ViewCount:>>>  	%s", ch.ViewCount)
			// log.Printf("StartTime:>>>  	%s", ch.StartTime)
			// log.Println("========================================")

			// 先暫存到 result 中
			lock.Lock()
			defer lock.Unlock()
			result[ch.Cid] = ch
		}
	})

	this.collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.85 Safari/537.36")
	})

	return this
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

// 訪問並獲取頻道資訊後寫入Redis
func (c *Crawler) Visit(cid string, expiration time.Duration) ChannelInfo {

	// 獲取頻道資訊
	log.Printf("[Youtube] cid :> %v", cid)
	c.collector.Visit("https://www.youtube.com/channel/" + cid)
	// 先取被暫存在result中的result數據, 再寫入到 Redis
	// 不可略過這步驟，因為result不會主動被清除，必須每次訪問後清除。
	lock.Lock()
	defer lock.Unlock()
	ch := result[cid]
	redisdb.GetInstance().SetVisitChannel(ch, expiration)
	delete(result, cid)

	return ch
}

// 搜尋頻道資訊後寫入Redis
func (c *Crawler) Search(cid string, expiration time.Duration) ChannelInfo {

	// 獲取頻道資訊
	log.Printf("[Youtube] cid :> %v", cid)
	c.collector.Visit("https://www.youtube.com/channel/" + cid)
	// 先取被暫存在result中的result數據, 再寫入到 Redis
	// 不可略過這步驟，因為result不會主動被清除，必須每次訪問後清除。
	lock.Lock()
	defer lock.Unlock()
	ch := result[cid]
	redisdb.GetInstance().SetSearchChannel(ch, expiration)
	delete(result, cid)

	return ch
}
