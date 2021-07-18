package services

import (
	"log"
	"lowkeydd-server/crawlers"
	"lowkeydd-server/redisdb"
	. "lowkeydd-server/share"

	"github.com/gin-gonic/gin"
)

type Letsddv2Request struct {
	UserID string   `json:"username"`
	SSID   string   `json:"ssid"`
	Tags   []string `json:"tags"`
}

type Letsddv2Response struct {
	Code     string        `json:"code"`
	Channels []ChannelInfo `json:"channels"`
	Visit    VisitList     `json:"visit"`
}

func Letsddv2_all_Response(c *gin.Context, tags []string, visit *VisitList) {

	log.Printf("[Letsddv2_all_Response] tags: [all]")

	if channels, all, success := redisdb.GetInstance().GetChannelsByConditionV2(func(info ChannelInfo) bool {
		return info.Status != "failure"
	}); success {

		Letsddv2_grouped_Response(c, visit, channels, all)
	} else {

		c.JSON(400, gin.H{"code": "error", "channels": [][]ChannelInfo{}, "group": visit.Group})
	}
}

func Letsddv2_taged_Response(c *gin.Context, tags []string, visit *VisitList) {

	log.Printf("[Letsddv2_taged_Response] tags: %v", tags)

	tagMap := make(map[string]bool)
	for _, tag := range tags {
		tagMap[string(tag)] = true
	}
	if channels, all, success := redisdb.GetInstance().GetChannelsByConditionV2(func(ch ChannelInfo) bool {
		return tagMap[ch.Status]
	}); success {

		Letsddv2_grouped_Response(c, visit, channels, all)
	} else {

		c.JSON(400, gin.H{"code": "error", "channels": [][]ChannelInfo{}, "group": visit.Group})
	}
}

func Letsddv2_auth_failure_Response(c *gin.Context, tags []string) {

	log.Printf("[Letsddv2_auth_failure_Response] tags: %v", tags)

	if tags[0] == "all" {
		Letsddv2_all_Response(c, tags, nil)
		return
	}
	Letsddv2_taged_Response(c, tags, nil)

}

func Letsddv2_auth_success_Response(c *gin.Context, tags []string, visit *VisitList) {

	log.Printf("[Letsddv2_auth_success_Response] tags: %v", tags)

	if tags[0] == "all" {
		Letsddv2_all_Response(c, tags, visit)
		return
	}
	Letsddv2_taged_Response(c, tags, visit)

}

func Letsddv2Endpoint(c *gin.Context) {

	var req Letsddv2Request
	if c.BindJSON(&req) != nil {
		c.JSON(200, gin.H{"code": "failure", "msg": "錯誤的JSON格式"})
		return
	}
	log.Printf("username :> %s\n", req.UserID)
	log.Printf("ssid :> %s\n", req.SSID)
	userid := req.UserID
	ssid := req.SSID
	tags := req.Tags

	if userid == "" {
		log.Printf("userid 沒有\n")
		Letsddv2_auth_failure_Response(c, tags)
		return
	}

	if ssid == "" {
		log.Printf("ssid 沒有\n")
		Letsddv2_auth_failure_Response(c, tags)
		return
	}

	if s, success := redisdb.GetInstance().GetSession(userid); success && s.SSID == ssid {

		log.Printf("ssid:> %s 驗證成功", ssid)

		// 驗證成功，獲取該使用者visit
		if code, visit := GetVisitList(userid); code == "success" {
			// 將讀取的visit傳入
			crawlers.GetInstance().Checked_VisitByList(visit.List)
			Letsddv2_auth_success_Response(c, tags, &visit)
			return
		} else {
			log.Printf("無法取得Visitlist\n userid: %s", userid)
			Letsddv2_auth_failure_Response(c, tags)
			return
		}
	} else {
		log.Printf("ssid:> %s 驗證失敗", ssid)
		Letsddv2_auth_failure_Response(c, tags)
		return
	}
}