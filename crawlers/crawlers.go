package crawlers

import (
	"log"
	"lowkeydd-server/crawlers/twitch"
	"lowkeydd-server/crawlers/youtube"
	"lowkeydd-server/redisdb"
	. "lowkeydd-server/share"
	"sync"
	"time"
)

type Crawlers struct {
	ytCrawler *youtube.Crawler
	twCrawler *twitch.Crawler
	wg        *sync.WaitGroup
}

type Authenticators struct {
	TwitchAuth  twitch.Authenticator  `json:"twitch"`
	YoutubeAuth youtube.Authenticator `json:"youtube"`
}

var (
	lock     = &sync.Mutex{}
	crawlers *Crawlers
	visit    *VisitList
	auths    *Authenticators
)

func GetInstance() *Crawlers {
	if crawlers == nil {
		// 只允許一個goroutine訪問
		lock.Lock()
		defer lock.Unlock()
		if crawlers == nil {
			NewCrawlers()
		}
	}
	return crawlers
}

func NewCrawlers() {

	// crawlers 針對 twith、youtube兩直播平台，能以api或是爬蟲方式獲取頻道資訊，並存放到redis資料庫中
	JSONFileLoader("setting/visit.json", &visit)
	JSONFileLoader("setting/authenticator.json", &auths)

	crawlers = &Crawlers{
		ytCrawler: youtube.NewCrawler(visit, &auths.YoutubeAuth),
		twCrawler: twitch.NewCrawler(visit, &auths.TwitchAuth),
		wg:        &sync.WaitGroup{},
	}

}

// 透過cid, method, 直接對對應的平台進行訪問，並將解析後的資訊寫入到redis中。
// 一般的頻道訪問保存期限一天。
func (c *Crawlers) request(cid string, method string) ChannelInfo {
	defer c.wg.Done()
	switch method {
	case "youtube":
		return c.ytCrawler.Visit(cid, 86400*time.Second)
	case "twitch":
		return c.twCrawler.Visit(cid, 86400*time.Second)
	default:
		return ChannelInfo{}
	}
}

// 對當前資料庫的資料，根據更新時間重新獲取。
func (c *Crawlers) updated_Request(curr int64, item ChannelInfo) {
	if curr > item.UpdateTime {
		c.request(item.Cid, item.Method)
	} else {
		c.wg.Done()
	}
}

// Checked 先驗證更新時間與是否存在才進行訪問。
// 減少實際對平台訪問次數，降低被平台當作機器人的機會。
func (c *Crawlers) checked_Request(curr int64, cid string, method string) {
	if info, success := redisdb.GetInstance().GetVisitChannel(cid); success {
		c.updated_Request(curr, info)
	} else {
		c.request(cid, method)
	}
}

// 訪問目標前，先檢查其更新狀態，再蒐集資訊後寫入到redis中。
func (c *Crawlers) Checked_Visit(cid string, method string) {
	log.Printf("[crawlers] Start to Visit: %v", cid)

	curr := time.Now().Unix()
	c.wg.Add(1)
	go c.checked_Request(curr, cid, method)
	c.wg.Wait()

	log.Printf("[crawlers] Time Complete...Visit is done ")
}

// 根據訪問清單，訪問訪問多個目標後，再將蒐集資訊後寫入到redis中。
func (c *Crawlers) UnChecked_VisitByDefaultList() {

	log.Printf("[crawlers] Start to VisitAll: \n%v", visit.List)

	c.wg.Add(len(visit.List))
	for _, item := range visit.List {
		go c.request(item.Cid, item.Method)
	}
	c.wg.Wait()

	log.Printf("[crawlers] Time Complete...VisitAll is done ")

}

func (c *Crawlers) Checked_VisitByDefaultList() {

	log.Printf("[crawlers] Start to VisitAll: \n%v", visit.List)

	curr := time.Now().Unix()
	c.wg.Add(len(visit.List))
	for _, item := range visit.List {
		go c.checked_Request(curr, item.Cid, item.Method)
	}
	c.wg.Wait()

	log.Printf("[crawlers] Time Complete...VisitAll is done ")

}

func (c *Crawlers) Checked_VisitByList(list []VisitItem) {

	log.Printf("[crawlers] Start to VisitAll: \n%v", list)

	curr := time.Now().Unix()
	c.wg.Add(len(list))
	for _, item := range list {
		go c.checked_Request(curr, item.Cid, item.Method)
	}
	c.wg.Wait()

	log.Printf("[crawlers] Time Complete...VisitAll is done ")

}

func (c *Crawlers) UnChecked_Update() {

	// 為當前Redis中所有的頻道資訊建立副本
	channels, _ := redisdb.GetInstance().GetAllVisitChannels()

	log.Println("[crawlers] 所有頻道資訊更新作業開始....")

	curr := time.Now().Unix()
	c.wg.Add(len(channels))
	for _, item := range channels {
		go c.checked_Request(curr, item.Cid, item.Method)
	}
	c.wg.Wait()

	log.Println("[crawlers] 所有頻道資訊更新作業結束....")
}

// 搜尋頻道資訊，搜尋結果一樣被保存在redis中
func (c *Crawlers) search(cid string, method string) ChannelInfo {
	switch method {
	case "youtube":
		return c.ytCrawler.Search(cid, 300*time.Second)
	case "twitch":
		return c.twCrawler.Search(cid, 300*time.Second)
	default:
		return ChannelInfo{}
	}
}

func (c *Crawlers) GetSearchChannel(cid string, method string) ChannelInfo {

	log.Printf("[crawlers] Start to get search channel: %v", cid)

	result := ChannelInfo{}
	curr := time.Now().Unix()
	if ch, success := redisdb.GetInstance().GetSearchChannel(cid); success {
		if curr > ch.UpdateTime {
			result = c.search(ch.Cid, ch.Method)
		} else {
			result = ch
		}
	} else {
		result = c.search(cid, method)
	}
	// log.Printf("[crawlers] Time Complete...get search channel is done, \n%v ", result)
	return result
}

func GetTwitchCrawler() *twitch.Crawler {
	return GetInstance().twCrawler
}

func GetYouttubeCrawler() *youtube.Crawler {
	return GetInstance().ytCrawler
}
