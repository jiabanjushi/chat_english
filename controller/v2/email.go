package v2

import (
	"github.com/gin-gonic/gin"
	"go-fly-muti/lib"
	"go-fly-muti/types"
)

func PostEmailCode(c *gin.Context) {
	email := c.PostForm("email")
	notify := &lib.Notify{
		Subject:     "测试主题",
		MainContent: "测试内容",
		EmailServer: lib.NotifyEmail{
			Server:   "smtp.sina.cn",
			Port:     587,
			From:     "taoshihan1@sina.com",
			Password: "382e8a5e11cfae8c",
			To:       []string{email},
			FromName: "GOFLY客服",
		},
	}
	_, err := notify.SendMail()
	if err != nil {
		c.JSON(200, gin.H{
			"code": types.ApiCode.FAILED,
			"msg":  types.ApiCode.GetMessage(types.ApiCode.FAILED),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": types.ApiCode.SUCCESS,
		"msg":  types.ApiCode.GetMessage(types.ApiCode.SUCCESS),
	})
	return
}
