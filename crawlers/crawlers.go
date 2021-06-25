package crawlers

import (
	"log"
	"lowkeydd-crawler/crawlers/twitch"
	"lowkeydd-crawler/crawlers/youtube"
	. "lowkeydd-crawler/share"
	"sync"
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

// 根據vist list中的紀錄，以並行方式去各頻道蒐集資訊，並寫入到redis中。
func (c *Crawlers) request(cid string, method string) {
	defer c.wg.Done()
	switch method {
	case "youtube":
		c.ytCrawler.Visit(cid)
		break
	case "twitch":
		c.twCrawler.Visit(cid)
		break
	}
	log.Printf("cid :> %v", cid)
}

// 對單一個目標去蒐集頻道資訊，並寫入到redis中。
func (c *Crawlers) Visit(cid string, method string) {
	log.Printf("Start to VisitAll...")
	log.Printf("%v", c.List)

	c.wg.Add(1)
	go c.request(cid, method)
	c.wg.Wait()

	log.Printf("Time Complete...VisitAll is done ")
}

// 根據vist list中的紀錄，對多個目標以並行方式去各頻道蒐集資訊，並寫入到redis中。
func (c *Crawlers) VisitAll() {

	log.Printf("Start to VisitAll...")
	log.Printf("%v", c.List)

	c.wg.Add(len(c.List))
	for _, item := range c.List {
		go c.request(item.Cid, item.Method)
	}
	c.wg.Wait()

	log.Printf("Time Complete...VisitAll is done ")

}
