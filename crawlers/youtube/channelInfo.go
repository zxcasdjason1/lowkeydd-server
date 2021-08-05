package youtube

import (
	"log"
	. "lowkeydd-server/share"
	"regexp"
	"strconv"

	"github.com/tidwall/gjson"
)

func getSectionContents(ytStr string, index int) gjson.Result {

	path := "contents.twoColumnBrowseResultsRenderer.tabs.0" +
		".tabRenderer.content.sectionListRenderer.contents." + strconv.Itoa(index) +
		".itemSectionRenderer.contents.0"
	return gjson.Get(ytStr, path)
}

func CreateChannelInfo(ytStr string) ChannelInfo {

	header := gjson.Get(ytStr, "header.c4TabbedHeaderRenderer").Raw
	// 獲取Header
	channelId := gjson.Get(header, "channelId").Raw
	if channelId == "" {
		log.Println("取得頻道代碼(cid)失敗, " + channelId)
		return ChannelInfo{}
	}
	cname := gjson.Get(header, "title").Raw
	if cname == "" {
		log.Printf("取得頻道所有者(cname)失敗,\n %v", header)
	}
	avatar := gjson.Get(header, "avatar.thumbnails.1.url").Raw //88x88
	if avatar == "" {
		log.Printf("取得頭像圖片(avatar)失敗,\n %v", header)
	}

	var ch = ChannelInfo{}
	var lasttime = 2147483647

	for i := 0; i < 3; i++ {

		ctx := getSectionContents(ytStr, i)
		if ctx.Raw == "" {
			// log.Printf("index: %d \n\n\n\n\n  ctx:> %v", i, ctx.Raw)
			break
		}

		// 直播開啟，會含有 channelFeaturedContentRenderer
		if r := gjson.Get(ctx.Raw, "channelFeaturedContentRenderer.items.0.videoRenderer").Raw; r != "" {
			// log.Printf("index: %d\n\n channelFeaturedContentRenderer:> %s\n", i, r)
			channelFeaturedContentRenderer(r, &ch)
			break
		}

		// 已排定的直播訊息 (upcomingEventData)
		if r := gjson.Get(ctx.Raw, "shelfRenderer.content.expandedShelfContentsRenderer.items.0.videoRenderer").Raw; r != "" {
			// log.Printf("index: %d\n\n expandedShelf_ContentsRenderer:> %s\n", i, r)
			shelf_expandedshelf_VideoRenderer(r, &ch)
			if gjson.Get(r, "upcomingEventData").Raw != "" {
				break //有直播預訂訊息
			} else {
				continue //沒有就繼續找
			}
		}

		// 為一整排的影片欄位，獲取欄位中第一個作為最近一次的重播紀錄。
		if r := gjson.Get(ctx.Raw, "shelfRenderer.content.horizontalListRenderer.items.0.gridVideoRenderer").Raw; r != "" {
			// log.Printf("index: %d\n\n shelf_gridVideoRenderer:> %s\n", i, r)
			var curChannel = ChannelInfo{}
			shelf_gridVideoRenderer(r, &curChannel)
			cur := getTimeByStartTimeStr(curChannel.StartTime)
			if cur < lasttime {
				lasttime = cur
				ch = curChannel
			}
			continue
		}

		// // 之下的層級是我們不關注的類型，先當作候選選項。
		// // 為多個播放清單所形成的欄位 (這類人通常沒有自己的任何影片)
		// if r := gjson.Get(ctx.Raw, "shelfRenderer.content.horizontalListRenderer.items.0.gridPlaylistRenderer").Raw; r != "" {
		// 	// log.Printf("index: %d\n\n shelf_gridPlaylistRenderer:> %s\n", i, r)
		// 	shelf_gridPlaylistRenderer(r, &ch)
		// 	continue
		// }
	}

	// 設定Header
	if ch.RenderType == "" {
		log.Println("缺少對應此類頻道資訊的方式 cid:> " + channelId)
		ch.RenderType = "nilRenderer"
		ch.Status = "failure"
	}
	ch.Cid = removeQuotes(channelId)
	ch.Cname = removeQuotes(cname)
	ch.Owner = removeQuotes(cname)
	ch.Avatar = removeQuotes(avatar)
	return ch
}

// 直播中
func channelFeaturedContentRenderer(ctx string, ch *ChannelInfo) {

	// 這個項目照理說不會被列出，有影片必定有播放清單，播放清單的

	title := gjson.Get(ctx, "title.runs.0.text")              //影片標題
	videoID := gjson.Get(ctx, "videoId")                      // 影片連結
	thumbnail := gjson.Get(ctx, "thumbnail.thumbnails.3.url") //圖片
	viewCount := gjson.Get(ctx, "viewCountText.runs.0.text")  //影片當前觀看人數

	ch.Cid = ""
	ch.Status = "live"
	ch.RenderType = "channelFeaturedContentRenderer"
	ch.StreamURL = "https://www.youtube.com/watch?v=" + removeQuotes(videoID.Raw)
	ch.Thumbnail = removeQuotes(thumbnail.Raw)
	ch.Title = removeQuotes(title.Raw)
	ch.ViewCount = removeQuotes(viewCount.Raw)
	ch.StartTime = ""

}

// 等待中
func shelf_expandedshelf_VideoRenderer(ctx string, ch *ChannelInfo) {

	title := gjson.Get(ctx, "title.simpleText")               //影片標題
	videoID := gjson.Get(ctx, "videoId")                      //影片連結
	thumbnail := gjson.Get(ctx, "thumbnail.thumbnails.3.url") //圖片
	viewCountStr := getWaitOrOffViewCountStr(ctx)             //影片當前觀看人數
	startTimeStr := getStartTimeStr(ctx)                      //預定發布時間

	ch.Cid = ""
	ch.Status = "wait"
	ch.RenderType = "expandedShelfContentsRenderer"
	ch.StreamURL = "https://www.youtube.com/watch?v=" + removeQuotes(videoID.Raw)
	ch.Thumbnail = removeQuotes(thumbnail.Raw)
	ch.Title = removeQuotes(title.Raw)
	ch.ViewCount = viewCountStr
	ch.StartTime = startTimeStr

}

// 等待中
func shelf_gridVideoRenderer(ctx string, ch *ChannelInfo) {

	title := gjson.Get(ctx, "title.simpleText")
	videoId := gjson.Get(ctx, "videoId")
	thumbnail := gjson.Get(ctx, "thumbnail.thumbnails.3.url")
	viewCountStr := getWaitOrOffViewCountStr(ctx)
	startTimeStr := getStartTimeStr(ctx)
	status := "wait"
	if startTimeStr == "" {
		// 觀看時間: xx ...前
		if publishedTimeStr := gjson.Get(ctx, "publishedTimeText.simpleText").Raw; publishedTimeStr != "" {
			log.Printf("videoId: %s publishedTimeStr %s", videoId, publishedTimeStr)
			if startTime := regexp.MustCompile("[0-9]+.*").FindAllString(publishedTimeStr, -1); len(startTime) > 0 {
				startTimeStr = startTime[0]
			} else {
				startTimeStr = publishedTimeStr
			}
			status = "off"
		}
	}

	ch.Cid = ""
	ch.Status = status
	ch.RenderType = "shelfRenderer+gridVideoRenderer"
	ch.StreamURL = "https://www.youtube.com/watch?v=" + removeQuotes(videoId.Raw)
	ch.Thumbnail = removeQuotes(thumbnail.Raw)
	ch.Title = removeQuotes(title.Raw)
	ch.ViewCount = removeQuotes(viewCountStr)
	ch.StartTime = removeQuotes(startTimeStr)

}

// func shelf_gridPlaylistRenderer(ctx string, ch *ChannelInfo) {

// 	playlistId := gjson.Get(ctx, "playlistId")
// 	title := gjson.Get(ctx, "title.runs.0.text")
// 	thumbnail := gjson.Get(ctx, "thumbnail.thumbnails.0.url")

// 	ch.Cid = ""
// 	ch.Status = "off"
// 	ch.RenderType = "shelfRenderer+gridPlaylistRenderer"
// 	ch.StreamURL = "https://www.youtube.com/watch?v=" + removeQuotes(playlistId.Raw)
// 	ch.Thumbnail = removeQuotes(thumbnail.Raw)
// 	ch.Title = removeQuotes(title.Raw)
// 	ch.ViewCount = ""
// 	ch.StartTime = ""

// }
