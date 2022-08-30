package middleware

import (
	"github.com/gin-gonic/gin"
	"go-fly-muti/models"
	"go-fly-muti/types"
)

// Ipblack TODO  ip白名单
func Ipblack(c *gin.Context) {
	ip := c.ClientIP()
	ipblack := models.FindIp(ip)
	if ipblack.IP != "" {
		c.JSON(200, gin.H{
			"code": types.ApiCode.IP_BAN,
			"msg":  types.ApiCode.GetMessage(types.ApiCode.IP_BAN),
		})
		c.Abort()
		return
	}
}
