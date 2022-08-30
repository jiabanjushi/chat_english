package controller

import (
	"github.com/gin-gonic/gin"
	"go-fly-muti/common"
	"go-fly-muti/tools"
)

func JyKey(c *gin.Context) {
	key := c.PostForm("key")
	if len(key) != 32 {
		c.JSON(200, gin.H{
			"code":   -101,
			"msg":    "",
			"result": "",
		})
		return
	}
	bo, _ := tools.InArray(key, common.KFMYArray)

	if bo==false {
		c.JSON(200, gin.H{
			"code":   -101,
			"msg":    "",
			"result": "",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "",
		"result": "",
	})
	return

}
