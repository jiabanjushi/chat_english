package middleware

import (
	"github.com/gin-gonic/gin"
	"go-fly-muti/common"
	"go-fly-muti/tools"
	"go-fly-muti/types"
	"strings"
)

// DomainLimitMiddleware /**  TODO  域名白名单
func DomainLimitMiddleware(c *gin.Context) {
	if common.DomainWhiteList == "*" {
		return
	}
	host := c.Request.Host
	hostPort := strings.Split(host, ":")
	if len(hostPort) == 0 {
		c.JSON(200, gin.H{
			"code": types.ApiCode.DOMAIN_LIMIT,
			"msg":  types.ApiCode.GetMessage(types.ApiCode.DOMAIN_LIMIT),
		})
		c.Abort()
		return
	}

	allowHost := strings.Split(common.DomainWhiteList, ",")
	if exist, _ := tools.InArray(hostPort[0], allowHost); !exist {
		c.JSON(200, gin.H{
			"code": types.ApiCode.DOMAIN_LIMIT,
			"msg":  types.ApiCode.GetMessage(types.ApiCode.DOMAIN_LIMIT),
		})
		c.Abort()
		return
	}
}
