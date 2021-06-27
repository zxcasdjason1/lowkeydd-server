package crawlers

import (
	"log"
	"lowkeydd-crawler/crawlers/twitch"
	"lowkeydd-crawler/crawlers/youtube"
	"lowkeydd-crawler/redisdb"
	. "lowkeydd-crawler/share"
	"sync"
	"time"
)

type Crawlers struct {
	ytCrawler *youtube.Crawler
	twCrawler *twitch.Crawler
	wg        *sync.WaitGroup
	List      []VisitItem
}

var lock = &sync.Mutex{}
var crawlers *Crawlers

func GetInstance(v ...*VisitList) *Crawlers {
	if crawlers == nil {
		// 只允許一個goroutine訪問
		lock.Lock()
		defer lock.Unlock()
		if crawlers == nil {
			NewCrawlers(v[0])
		}
	}
	return crawlers
}

func NewCrawlers(v *VisitList) {

	// crawlers 針對 twith、youtube兩直播平台，能以api或是爬蟲方式獲取頻道資訊，並存放到redis資料庫中

	crawlers = &Crawlers{
		ytCrawler: youtube.NewCrawler(v),
		twCrawler: twitch.NewCrawler(v),
		wg:        &sync.WaitGroup{},
		List:      v.List,
	}
}

func (c *Crawlers) GetWg() *sync.WaitGroup {
	return c.wg
}

// 根據vist list中的紀錄，以並行方式去各頻道蒐集資訊，並寫入到redis中。
func (c *Crawlers) Request(cid string, method string) {
	defer c.wg.Done()
	switch method {
	case "youtube":
		c.ytCrawler.Visit(cid)
	case "twitch":
		c.twCrawler.Visit(cid)
	}
}

func (c *Crawlers) Update(item ChannelInfo, curr int64) {
	defer c.wg.Done()
	if curr > item.UpdateTime {
		switch item.Method {
		case "youtube":
			c.ytCrawler.Visit(item.Cid)
		case "twitch":
			c.twCrawler.Visit(item.Cid)
		}
	}
}

func (c *Crawlers) Checked_Request(curr int64, cid string, method string) {
	if info, exist := redisdb.GetInstance().GetChannelInfo(cid); exist {
		if curr > info.UpdateTime {
			c.Request(cid, method)
		} else {
			c.wg.Done()
		}
	} else {
		c.Request(cid, method)
	}
}

// 訪問目標後蒐集資訊後寫入到redis中。
func (c *Crawlers) Visit(cid string, method string) {
	log.Printf("[crawlers] Start to Visit: %v", cid)

	c.wg.Add(1)
	go c.Request(cid, method)
	c.wg.Wait()

	log.Printf("[crawlers] Time Complete...Visit is done ")
}

// 訪問目標前，先檢查其更新狀態，再蒐集資訊後寫入到redis中。
func (c *Crawlers) Checked_Visit(cid string, method string) {
	log.Printf("[crawlers] Start to Visit: %v", cid)

	curr := time.Now().Unix()
	c.wg.Add(1)
	go c.Checked_Request(curr, cid, method)
	c.wg.Wait()

	log.Printf("[crawlers] Time Complete...Visit is done ")
}

// 根據訪問清單，訪問訪問多個目標後，再將蒐集資訊後寫入到redis中。
func (c *Crawlers) VisitAll() {
	log.Printf("[crawlers] Start to VisitAll: \n%v", c.List)

	c.wg.Add(len(c.List))
	for _, item := range c.List {
		go c.Request(item.Cid, item.Method)
	}
	c.wg.Wait()

	log.Printf("[crawlers] Time Complete...VisitAll is done ")

}

func (c *Crawlers) Checked_VisitAll() {

	log.Printf("[crawlers] Start to VisitAll: \n%v", c.List)

	curr := time.Now().Unix()
	c.wg.Add(len(c.List))
	for _, item := range c.List {
		go c.Checked_Request(curr, item.Cid, item.Method)
	}
	c.wg.Wait()

	log.Printf("[crawlers] Time Complete...VisitAll is done ")

}

func GetTwitchCrawler() *twitch.Crawler {
	return crawlers.twCrawler
}

func GetYouttubeCrawler() *youtube.Crawler {
	return crawlers.ytCrawler
}
