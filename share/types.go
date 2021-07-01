package share

// ChannelInfo 頻道資訊
type ChannelInfo struct {
	Cid        string `json:"cid"`
	Cname      string `json:"cname"`
	Owner      string `json:"owner"`
	Avatar     string `json:"avatar"`
	Status     string `json:"status"`
	RenderType string `json:"rendertype"`
	StreamURL  string `json:"streamurl"`
	Thumbnail  string `json:"thumbnail"`
	Title      string `json:"title"`
	ViewCount  string `json:"viewcount"`
	StartTime  string `json:"starttime"`  // 直播開始時間
	Method     string `json:"method"`     // 使用的crawler方式
	UpdateTime int64  `json:"updatetime"` // 最近一次的更新時間
}

type VisitItem struct {
	Cid    string `json:"cid"`    // 頻道索引
	Cname  string `json:"cname"`  // 頻道主真實ID
	Owner  string `json:"owner"`  // 自定義的頻道稱呼
	Group  string `json:"group"`  // 自定義的頻道分群
	Method string `json:"method"` // 使用的crawler方式
}

type VisitList struct {
	ClientID  string      `json:"clien_id"`
	AuthToken string      `json:"auth_token"`
	List      []VisitItem `json:"list"`
}

type Session struct {
	SSID       string `json:"ssid"`
	Expiration int64  `json:"expiration"`
}
