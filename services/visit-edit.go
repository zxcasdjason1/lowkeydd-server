package services

import (
	"log"
	. "lowkeydd-server/share"

	"github.com/gin-gonic/gin"
)

type VisitEditRequest struct {
	UserID string `json:"username"`
}

type VisitEditResponse struct {
	Code  string    `json:"code"`
	Visit VisitList `json:"visit"`
}

func VisitEditEndpoint(c *gin.Context) {

	userid := c.DefaultPostForm("username", "")
	log.Printf("username %v", userid)

	code, visit := GetVisitList(userid)
	log.Printf("瀏覽追隨者清單:> %v\n", visit)
	VisitEditTransPort(c, code, visit)
}

func VisitEditTransPort(c *gin.Context, code string, visit VisitList) {
	switch code {
	case "success":
		c.JSON(200, gin.H{"code": code, "visit": visit})
		return
	case "failure":
		c.JSON(200, gin.H{"code": code, "visit": visit})
		return
	case "error":
		c.JSON(400, gin.H{"code": code, "visit": visit})
		return
	}
}
