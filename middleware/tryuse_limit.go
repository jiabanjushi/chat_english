package middleware

import (
	"github.com/gin-gonic/gin"
	"go-fly-muti/common"
	"go-fly-muti/tools"
	"go-fly-muti/types"
	"io/ioutil"
)

/**
域名限制中间件
*/
func TryUseLimitMiddleware(c *gin.Context) {
	if !common.IsTry {
		return
	}
	installTimeByte, _ := ioutil.ReadFile("./install.lock")
	installTimeByte, _ = tools.AesDecrypt(installTimeByte, []byte(common.AesKey))
	installTime := tools.ByteToInt64(installTimeByte)
	if tools.Now()-installTime >= common.TryDeadline {
		c.JSON(200, gin.H{
			"code": types.ApiCode.TRYUSE_LIMIT,
			"msg":  types.ApiCode.GetMessage(types.ApiCode.TRYUSE_LIMIT),
		})
		c.Abort()
		return
	}
}
