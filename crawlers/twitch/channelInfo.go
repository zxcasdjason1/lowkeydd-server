package twitch

import (
	"strings"

	. "lowkeydd-crawler/share"

	"github.com/tidwall/gjson"
)

func getThumbnail(str string, width string, height string) string {
	return strings.ReplaceAll(strings.ReplaceAll(str, `{width}`, width), `{height}`, height)
}

func GetLiveChannelInfo(stream string) *ChannelInfo {
	//取直播中
	channel := gjson.Get(stream, "channel").Raw
	cid := gjson.Get(channel, "_id").Raw
	owner := gjson.Get(channel, "display_name").Raw
	avatar := gjson.Get(channel, "logo").Raw
	streamurl := gjson.Get(channel, "url").Raw
	thumbnail := getThumbnail(gjson.Get(stream, "preview.template").Raw, "500", "280")
	title := gjson.Get(channel, "status").Raw
	viewcount := gjson.Get(stream, "viewers").Raw
	// starttime := gjson.Get(stream, "created_at").Raw //正在直播

	return &ChannelInfo{
		Cid:        cid,
		Owner:      removeQuotes(owner),
		Avatar:     removeQuotes(avatar),
		Status:     "live",
		RenderType: "twitchapi-kraken-streams-[USERS_ID]",
		StreamURL:  removeQuotes(streamurl),
		Thumbnail:  removeQuotes(thumbnail),
		Title:      removeQuotes(title),
		ViewCount:  AddComma(removeQuotes(viewcount)),
		StartTime:  "",
	}
}

func GetOffChannelInfo(vedios string) *ChannelInfo {
	// 取第一筆
	vedio := gjson.Get(vedios, "videos.0").Raw
	cid := gjson.Get(vedio, "channel._id").Raw
	owner := gjson.Get(vedio, "channel.display_name").Raw
	avatar := gjson.Get(vedio, "channel.logo").Raw
	streamurl := gjson.Get(vedio, "url").Raw
	thumbnail := getThumbnail(gjson.Get(vedio, "preview.template").Raw, "500", "280")
	title := gjson.Get(vedio, "title").Raw
	viewcount := gjson.Get(vedio, "views").Raw
	starttime := gjson.Get(vedio, "published_at").Raw

	return &ChannelInfo{
		Cid:        cid,
		Owner:      removeQuotes(owner),
		Avatar:     removeQuotes(avatar),
		Status:     "off",
		RenderType: "twitchapi-kraken-channels-[USERS_ID]-videos",
		StreamURL:  removeQuotes(streamurl),
		Thumbnail:  removeQuotes(thumbnail),
		Title:      removeQuotes(title),
		ViewCount:  "觀看次數：" + AddComma(removeQuotes(viewcount)) + "次",
		StartTime:  removeQuotes(starttime),
	}

}
