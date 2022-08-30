package controller

import (
	"github.com/gin-gonic/gin"
	"go-fly-muti/models"
	"go-fly-muti/types"
	"strconv"
)

type NewsForm struct {
	Id      uint   `form:"id" json:"id" uri:"id" xml:"id"`
	Tag     string `form:"tag" json:"tag" uri:"tag" xml:"tag" binding:"required"`
	Title   string `form:"title" json:"title" uri:"title" xml:"title" binding:"required"`
	Content string `form:"content" json:"content" uri:"content" xml:"content" binding:"required"`
}

//新闻列表
func GetNews(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	if page <= 0 {
		page = 1
	}
	pagesize, _ := strconv.Atoi(c.Query("pagesize"))
	if pagesize <= 0 || pagesize > 50 {
		pagesize = 10
	}
	count := models.CountNews("")
	list := models.FindNews(page, pagesize, "")
	c.JSON(200, gin.H{
		"code": types.ApiCode.SUCCESS,
		"msg":  types.ApiCode.GetMessage(types.ApiCode.SUCCESS),
		"result": gin.H{
			"list":     list,
			"count":    count,
			"pagesize": pagesize,
		},
	})
}

//添加新闻
func PostNews(c *gin.Context) {
	var form NewsForm
	err := c.Bind(&form)

	if err != nil {
		c.JSON(200, gin.H{
			"code":   types.ApiCode.FAILED,
			"msg":    types.ApiCode.GetMessage(types.ApiCode.INVALID),
			"result": err.Error(),
		})
		return
	}
	modelNews := &models.New{
		Title:   form.Title,
		Content: form.Content,
		Tag:     form.Tag,
	}
	//添加新闻
	if form.Id == 0 {
		err := modelNews.AddNews()
		if err != nil {
			c.JSON(200, gin.H{
				"code": types.ApiCode.FAILED,
				"msg":  err.Error(),
			})
			return
		}
	} else {
		//修改新闻
		err := modelNews.SaveNews("id = ?", form.Id)
		if err != nil {
			c.JSON(200, gin.H{
				"code": types.ApiCode.FAILED,
				"msg":  err.Error(),
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"code": types.ApiCode.SUCCESS,
		"msg":  types.ApiCode.GetMessage(types.ApiCode.SUCCESS),
	})
}

//删除新闻
func DelNews(c *gin.Context) {
	id := c.Query("id")
	err := models.DelNews("id = ?", id)
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
