package v2

import (
	"github.com/gin-gonic/gin"
	"go-fly-muti/models"
	"go-fly-muti/types"
)

type ClientForm struct {
	ClientId string `form:"client_id" json:"client_id" uri:"client_id" xml:"client_id" binding:"required"`
}

/**
注册app client_id
*/
func PostAppKefuClient(c *gin.Context) {
	kefuName, _ := c.Get("kefu_name")
	var form ClientForm
	err := c.Bind(&form)
	if err != nil {
		c.JSON(200, gin.H{
			"code":   types.ApiCode.FAILED,
			"msg":    types.ApiCode.GetMessage(types.ApiCode.INVALID),
			"result": err.Error(),
		})
		return
	}
	clientModel := &models.User_client{
		Kefu:      kefuName.(string),
		Client_id: form.ClientId,
	}
	clientInfo := clientModel.FindClient()
	if clientInfo.ID != 0 {
		c.JSON(200, gin.H{
			"code": types.ApiCode.FAILED,
			"msg":  types.ApiCode.GetMessage(types.ApiCode.INVALID),
		})
		return
	}
	models.CreateUserClient(kefuName.(string), form.ClientId)
	c.JSON(200, gin.H{
		"code": types.ApiCode.SUCCESS,
		"msg":  types.ApiCode.GetMessage(types.ApiCode.SUCCESS),
	})
}
