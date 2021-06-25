package twitch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"lowkeydd-crawler/redisdb"
	. "lowkeydd-crawler/share"
	"net/http"

	"github.com/tidwall/gjson"
)

type Crawler struct {
	client    *http.Client
	List      []VisitItem
	Event     *Emitter
	ClientID  string
	AuthToken string
}

type Header struct {
	Accept        string
	ClientID      string
	Authorization string
}

const (
	ON_RESPOSE         = "crawler@@on_response"
	ON_STREAM_RESPOSE  = "crawler@@on_stream_response"
	ON_CHANNEL_RESPOSE = "crawler@@on_channel_response"
	ON_VEDIOS_RESPOSE  = "crawler@@on_vedios_response"
)

func NewCrawler(v *VisitList) *Crawler {

	this := &Crawler{
		client: &http.Client{},

		List:      v.List,
		ClientID:  v.ClientID,
		AuthToken: v.AuthToken,
		Event:     NewEmitter(),
	}

	return this
}

func (c *Crawler) Fetch(url string, h Header) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}
	if h.Accept != "" {
		req.Header.Set("Accept", h.Accept)
	}
	if h.Authorization != "" {
		req.Header.Set("Authorization", h.Authorization)
	}
	if h.ClientID != "" {
		req.Header.Set("Client-ID", h.ClientID)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	return string(body)
}

func (c *Crawler) SearchChannels(userid_list []string) string {

	//可以為多值檢索

	len_list := len(userid_list)
	if len_list == 0 {
		return ""
	}

	var str = ""
	for i := 0; i < len_list; i++ {
		str += fmt.Sprintf("user_id=%s&", userid_list[i])
	}
	str = str[0 : len(str)-1]

	return c.Fetch("https://api.twitch.tv/helix/streams?"+str, Header{
		Accept:        "application/vnd.twitchtv.v5+json",
		Authorization: "Bearer " + c.AuthToken,
		ClientID:      c.ClientID,
	})
}

func (c *Crawler) TodoSearch(userid string) {

	var info *ChannelInfo

	if stream := c.GetStream(userid); stream != "null" {
		// 發起獲取正在直播的串流資訊
		info = GetLiveChannelInfo(stream)
	} else {
		// 假如取得串流訊息的結果為 {"stream":"null"}，表示當前沒有直播
		// 或直播已經結束，所以再重新請求來獲取最近一次的影片記錄檔。
		info = GetOffChannelInfo(c.GetVedios(userid))
	}

	// log.Println("================正在重播==================")
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

func (c *Crawler) GetStream(userid string) string {
	stream := c.Fetch("https://api.twitch.tv/kraken/streams/"+userid, Header{
		Accept:        "application/vnd.twitchtv.v5+json",
		Authorization: "",
		ClientID:      c.ClientID,
	})
	return gjson.Get(stream, "stream").Raw
}

func (c *Crawler) GetChannel(userid string) string {

	return c.Fetch("https://api.twitch.tv/kraken/channels/"+userid, Header{
		Accept:        "application/vnd.twitchtv.v5+json",
		Authorization: "",
		ClientID:      c.ClientID,
	})
}

func (c *Crawler) GetVedios(userid string) string {

	return c.Fetch("https://api.twitch.tv/kraken/channels/"+userid+"/videos", Header{
		Accept:        "application/vnd.twitchtv.v5+json",
		Authorization: "",
		ClientID:      c.ClientID,
	})
}

func (c *Crawler) Visit(cid string) {
	c.TodoSearch(cid)
}
