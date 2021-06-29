package youtube

import (
	"log"
	. "lowkeydd-server/share"
	"strconv"

	"github.com/tidwall/gjson"
)

func getSectionContents(ytStr string, index int) gjson.Result {

	path := "contents.twoColumnBrowseResultsRenderer.tabs.0" +
		".tabRenderer.content.sectionListRenderer.contents." + strconv.Itoa(index) +
		".itemSectionRenderer.contents.0"
	return gjson.Get(ytStr, path)
}

func CreateChannelInfo(ytStr string) *ChannelInfo {

	header := gjson.Get(ytStr, "header.c4TabbedHeaderRenderer").Raw
	// 獲取Header
	channelId := gjson.Get(header, "channelId").Raw
	if channelId == "" {
		log.Println("取得頻道代碼(cid)失敗, " + channelId)
		return &ChannelInfo{}
	}
	cname := gjson.Get(header, "title").Raw
	if cname == "" {
		panic("取得頻道所有者(cname)失敗, " + cname)
	}
	avatar := gjson.Get(header, "avatar.thumbnails.1.url").Raw //88x88
	if avatar == "" {
		panic("取得頭像圖片(avatar)失敗, " + avatar)
	}

	var info *ChannelInfo = nil

	for i := 0; i < 3; i++ {

		ctx := getSectionContents(ytStr, i)
		if ctx.Raw == "" {
			// log.Printf("index: %d \n\n\n\n\n  ctx:> %v", i, ctx.Raw)
			break
		}

		// 直播開啟，會含有 channelFeaturedContentRenderer
		if r := gjson.Get(ctx.Raw, "channelFeaturedContentRenderer.items.0.videoRenderer").Raw; r != "" {
			// log.Printf("index: %d\n\n channelFeaturedContentRenderer:> %s\n", i, r)
			info = channelFeaturedContentRenderer(r)
			break
		}

		// 已排定的直播訊息
		if r := gjson.Get(ctx.Raw, "shelfRenderer.content.expandedShelfContentsRenderer.items.0.videoRenderer").Raw; r != "" {
			// log.Printf("index: %d\n\n expandedShelf_ContentsRenderer:> %s\n", i, r)
			info = shelf_expandedshelf_VideoRenderer(r)
			if gjson.Get(r, "upcomingEventData").Raw != "" {
				break //有直播預訂訊息
			} else {
				continue //沒有就繼續找
			}
		}

		// 為一整排的影片欄位，獲取欄位中第一個作為，最近一次的重播紀錄。
		if r := gjson.Get(ctx.Raw, "shelfRenderer.content.horizontalListRenderer.items.0.gridVideoRenderer").Raw; r != "" {
			// log.Printf("index: %d\n\n shelf_gridVideoRenderer:> %s\n", i, r)
			info = shelf_gridVideoRenderer(r)
			break
		}

		// 之下的層級是我們不關注的類型，先當作候選選項。
		// 為多個播放清單所形成的欄位 (這類人通常沒有自己的任何影片)
		if r := gjson.Get(ctx.Raw, "shelfRenderer.content.horizontalListRenderer.items.0.gridPlaylistRenderer").Raw; r != "" {
			// log.Printf("index: %d\n\n shelf_gridPlaylistRenderer:> %s\n", i, r)
			info = shelf_gridPlaylistRenderer(r)
			continue
		}
	}

	// 設定Header
	if info == nil {
		log.Println("缺少對應此類頻道資訊的方式 cid:> " + channelId)
		info = &ChannelInfo{
			RenderType: "nilRenderer",
			Status:     "failure",
		}
	}
	info.Cid = removeQuotes(channelId)
	info.Cname = removeQuotes(cname)
	info.Owner = removeQuotes(cname)
	info.Avatar = removeQuotes(avatar)
	return info
}

func channelFeaturedContentRenderer(ctx string) *ChannelInfo {

	// 這個項目照理說不會被列出，有影片必定有播放清單，播放清單的

	title := gjson.Get(ctx, "title.runs.0.text")              //影片標題
	videoID := gjson.Get(ctx, "videoId")                      // 影片連結
	thumbnail := gjson.Get(ctx, "thumbnail.thumbnails.3.url") //圖片
	viewCount := gjson.Get(ctx, "viewCountText.runs.0.text")  //影片當前觀看人數

	return &ChannelInfo{
		Cid:        "",
		Status:     "live",
		RenderType: "channelFeaturedContentRenderer",
		StreamURL:  "https://www.youtube.com/watch?v=" + removeQuotes(videoID.Raw),
		Thumbnail:  removeQuotes(thumbnail.Raw),
		Title:      removeQuotes(title.Raw),
		ViewCount:  removeQuotes(viewCount.Raw),
		StartTime:  "",
	}

}
func getViewCountStr(ctx string) string {
	res := gjson.Get(ctx, "viewCountText.simpleText").Raw //影片當前觀看人數
	if res == "" {
		vc1 := gjson.Get(ctx, "viewCountText.runs.0.text")
		vc2 := gjson.Get(ctx, "viewCountText.runs.1.text")
		res = vc1.Raw + " " + vc2.Raw
	}
	return removeQuotes(res)
}
func getStartTimeStr(ctx string) string {
	res := gjson.Get(ctx, "upcomingEventData.startTime").Raw // 檢查有沒有預定發布時間
	if res != "" {
		return timestampToDate(removeQuotes(res))
	}
	return ""
}
func shelf_expandedshelf_VideoRenderer(ctx string) *ChannelInfo {

	title := gjson.Get(ctx, "title.simpleText")               //影片標題
	videoID := gjson.Get(ctx, "videoId")                      //影片連結
	thumbnail := gjson.Get(ctx, "thumbnail.thumbnails.3.url") //圖片
	viewCountStr := getViewCountStr(ctx)                      //影片當前觀看人數
	startTimeStr := getStartTimeStr(ctx)                      //預定發布時間

	return &ChannelInfo{
		Cid:        "",
		Status:     "wait",
		RenderType: "expandedShelfContentsRenderer",
		StreamURL:  "https://www.youtube.com/watch?v=" + removeQuotes(videoID.Raw),
		Thumbnail:  removeQuotes(thumbnail.Raw),
		Title:      removeQuotes(title.Raw),
		ViewCount:  viewCountStr,
		StartTime:  startTimeStr,
	}
}

func shelf_gridVideoRenderer(ctx string) *ChannelInfo {

	title := gjson.Get(ctx, "title.simpleText")
	videoId := gjson.Get(ctx, "videoId")
	thumbnail := gjson.Get(ctx, "thumbnail.thumbnails.3.url")
	viewCountStr := getViewCountStr(ctx)
	startTimeStr := getStartTimeStr(ctx)
	status := "wait"
	if startTimeStr == "" {
		startTimeStr = gjson.Get(ctx, "publishedTimeText.simpleText").Raw //改成發布時間
		status = "off"
	}

	return &ChannelInfo{
		Cid:        "",
		Status:     status,
		RenderType: "shelfRenderer+gridVideoRenderer",
		StreamURL:  "https://www.youtube.com/watch?v=" + removeQuotes(videoId.Raw),
		Thumbnail:  removeQuotes(thumbnail.Raw),
		Title:      removeQuotes(title.Raw),
		ViewCount:  removeQuotes(viewCountStr),
		StartTime:  removeQuotes(startTimeStr),
	}

}

func shelf_gridPlaylistRenderer(ctx string) *ChannelInfo {

	playlistId := gjson.Get(ctx, "playlistId")
	title := gjson.Get(ctx, "title.runs.0.text")
	thumbnail := gjson.Get(ctx, "thumbnail.thumbnails.0.url")

	return &ChannelInfo{
		Cid:        "",
		Status:     "off",
		RenderType: "shelfRenderer+gridPlaylistRenderer",
		StreamURL:  "https://www.youtube.com/watch?v=" + removeQuotes(playlistId.Raw),
		Thumbnail:  removeQuotes(thumbnail.Raw),
		Title:      removeQuotes(title.Raw),
		ViewCount:  "",
		StartTime:  "",
	}
}
