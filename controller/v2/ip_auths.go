package v2

import (
	"github.com/gin-gonic/gin"
	"go-fly-muti/models"
	"go-fly-muti/tools"
)

func GetIpAuth(c *gin.Context) {
	ip := c.ClientIP()
	ipModel := models.FindServerIpAddress(ip)
	expireTime := tools.Now() + 7*24*3600
	if ipModel.ID == 0 {
		ipModel = models.CreateIpAuth(ip, tools.IntToTimeStr(expireTime, "2006-01-02"))
	}
	ipModel.NowTime = tools.IntToTimeStr(tools.Now(), "2006-01-02")
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": ipModel,
	})
}
