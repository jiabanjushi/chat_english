package controller

import (
	"github.com/gin-gonic/gin"
	"go-fly-muti/models"
	"go-fly-muti/types"
	"strconv"
)

type VisitorBlackForm struct {
	Id        uint   `form:"id" json:"id" uri:"id" xml:"id"`
	VisitorId string `form:"visitor_id" json:"visitor_id" uri:"visitor_id" xml:"visitor_id" binding:"required"`
	Name      string `form:"name" json:"name" uri:"name" xml:"name" binding:"required"`
}

//列表
func GeVisitorBlacks(c *gin.Context) {
	entId, _ := c.Get("ent_id")
	page, _ := strconv.Atoi(c.Query("page"))
	if page <= 0 {
		page = 1
	}
	pagesize, _ := strconv.Atoi(c.Query("pagesize"))
	if pagesize <= 0 || pagesize > 50 {
		pagesize = 10
	}
	count := models.CountVisitorBlack("ent_id = ? ", entId)
	list := models.FindVisitorBlacks(page, pagesize, "ent_id = ? ", entId)
	c.JSON(200, gin.H{
		"code": types.ApiCode.SUCCESS,
		"msg":  types.ApiCode.GetMessage(types.ApiCode.SUCCESS),
		"result": gin.H{
			"list":     list,
			"count":    count,
			"pagesize": pagesize,
			"page":     page,
		},
	})
}

//添加
func PostVisitorBlack(c *gin.Context) {
	kefuName, _ := c.Get("kefu_name")
	entId, _ := c.Get("ent_id")
	var form VisitorBlackForm
	err := c.Bind(&form)

	if err != nil {
		c.JSON(200, gin.H{
			"code":   types.ApiCode.FAILED,
			"msg":    types.ApiCode.GetMessage(types.ApiCode.INVALID),
			"result": err.Error(),
		})
		return
	}
	model := &models.VisitorBlack{
		VisitorId: form.VisitorId,
		Name:      form.Name,
		EntId:     entId.(string),
		KefuName:  kefuName.(string),
	}

	model.AddVisitorBlack()

	c.JSON(200, gin.H{
		"code": types.ApiCode.SUCCESS,
		"msg":  types.ApiCode.GetMessage(types.ApiCode.SUCCESS),
	})
}

//删除
func DelVisitorBlack(c *gin.Context) {
	entId, _ := c.Get("ent_id")
	id := c.Query("id")
	err := models.DelVisitorBlack("id = ? and ent_id = ?", id, entId)
	if err != nil {
		c.JSON(200, gin.H{
			"code": types.ApiCode.FAILED,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code": types.ApiCode.SUCCESS,
		"msg":  types.ApiCode.GetMessage(types.ApiCode.SUCCESS),
	})
}
