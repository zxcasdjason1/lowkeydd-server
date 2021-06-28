package services

import "github.com/gin-gonic/gin"

type SetCookieRequest struct {
	Key   string `uri:"key"`
	Value string `uri:"value"`
}

func SetCookie(c *gin.Context) {
	var req SetCookieRequest
	c.ShouldBindUri(&req)
	c.SetCookie(req.Key, req.Value, 3600, "/", "localhost", false, true)
}

type GetCookieRequest struct {
	Key string `uri:"key"`
}

func GetCookie(c *gin.Context) {
	var req GetCookieRequest
	c.ShouldBindUri(&req)
	value, err := c.Cookie(req.Key)
	if err != nil {
		c.JSON(200, gin.H{"msg": "get cookie fail"})
	} else {
		c.JSON(200, gin.H{req.Key: value})
	}
}
