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

func CreateCrawlers(v *VisitList) Crawlers {

	// 設置Crawlers獲取頻道資訊，當取得資訊時將ChannelInfo 轉換成 json後儲存到 redis中

	this := Crawlers{
		ytCrawler: youtube.NewCrawler(v),
		twCrawler: twitch.NewCrawler(v),
		wg:        &sync.WaitGroup{},
		List:      v.List,
	}

	return this
}

func (c *Crawlers) Visit(cid string, method string) {
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

func (c *Crawlers) VisitAll() {

	log.Printf("Start to VisitAll...")
	log.Printf("%v", c.List)

	// 根據vist list中的紀錄，以並行方式去各頻道蒐集資訊，並寫入到redis中。
	c.wg.Add(len(c.List))
	for _, item := range c.List {
		go c.Visit(item.Cid, item.Method)
	}
	c.wg.Wait()

	// time.Sleep(time.Second * 2)
	log.Printf("Time Complete...VisitAll is done ")

}
