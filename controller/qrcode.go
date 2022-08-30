package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"go-fly-muti/types"
)

func GetQrcode(c *gin.Context) {
	str := c.Query("str")
	var png []byte
	png, err := qrcode.Encode(str, qrcode.Medium, 256)
	if err != nil {
		c.JSON(200, gin.H{
			"code":   types.ApiCode.FAILED,
			"msg":    types.ApiCode.GetMessage(types.ApiCode.INVALID),
			"result": err.Error(),
		})
		return
	}
	c.Writer.Header().Set("Content-Type", "image/png")
	c.Writer.Write(png)
}
