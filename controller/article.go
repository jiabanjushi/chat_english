package controller

import (
	"github.com/gin-gonic/gin"
	"go-fly-muti/models"
	"go-fly-muti/types"
	"strconv"
)

func PostArticleCate(c *gin.Context) {
	kefuName, _ := c.Get("kefu_name")
	entIdStr, _ := c.Get("ent_id")
	entId, _ := entIdStr.(string)
	catName := c.PostForm("name")
	if catName == "" {
		c.JSON(200, gin.H{
			"code": types.ApiCode.FAILED,
			"msg":  types.ApiCode.GetMessage(types.ApiCode.INVALID),
		})
		return
	}
	cateModel := &models.ArticleCate{
		CatName: catName,
		UserId:  kefuName.(string),
		IsTop:   0,
		EntId:   entId,
	}
	cateModel.AddArticleCate()
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}
func PostArticle(c *gin.Context) {
	kefuName, _ := c.Get("kefu_name")
	entIdStr, _ := c.Get("ent_id")
	entId, _ := entIdStr.(string)
	title := c.PostForm("title")
	content := c.PostForm("content")
	catId := c.PostForm("cat_id")
	id := c.PostForm("id")
	if title == "" || content == "" || catId == "" {
		c.JSON(200, gin.H{
			"code": types.ApiCode.FAILED,
			"msg":  types.ApiCode.GetMessage(types.ApiCode.INVALID),
		})
		return
	}
	catIdInt, _ := strconv.Atoi(catId)
	//编辑文章
	if id != "" {
		articleModel := &models.Article{
			Title:   title,
			Content: content,
			CatId:   uint(catIdInt),
			UserId:  kefuName.(string),
			EntId:   entId,
		}
		articleModel.SaveArticle("ent_id = ? and id = ?", entId, id)
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "ok",
		})
		return
	}
	articleModel := &models.Article{
		Title:   title,
		Content: content,
		CatId:   uint(catIdInt),
		UserId:  kefuName.(string),
		EntId:   entId,
	}
	articleModel.AddArticle()
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}
func DelArticle(c *gin.Context) {
	entIdStr, _ := c.Get("ent_id")
	articleId := c.Query("id")
	if articleId == "" {
		c.JSON(200, gin.H{
			"code": types.ApiCode.FAILED,
			"msg":  types.ApiCode.GetMessage(types.ApiCode.INVALID),
		})
		return
	}
	models.DelArticles("ent_id = ? and id = ? ", entIdStr, articleId)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}
func DelArticleCate(c *gin.Context) {
	entIdStr, _ := c.Get("ent_id")
	catId := c.Query("id")
	if catId == "" {
		c.JSON(200, gin.H{
			"code": types.ApiCode.FAILED,
			"msg":  types.ApiCode.GetMessage(types.ApiCode.INVALID),
		})
		return
	}
	models.DelArticleCate("ent_id = ? and id = ? ", entIdStr, catId)
	models.DelArticles("ent_id = ? and cat_id = ? ", entIdStr, catId)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}
func GetArticleCates(c *gin.Context) {
	entId, _ := c.Get("ent_id")
	list := models.FindArticleCatesByEnt(entId)
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": list,
	})
}
func GetArticleList(c *gin.Context) {
	catId := c.Query("cat_id")
	entId, _ := c.Get("ent_id")
	page, _ := strconv.Atoi(c.Query("page"))
	if page <= 0 {
		page = 1
	}
	pagesize, _ := strconv.Atoi(c.Query("pagesize"))
	if pagesize <= 0 || pagesize > 50 {
		pagesize = 10
	}
	search := "ent_id = ? "
	args := []interface{}{
		entId,
	}

	if catId != "" {
		search += "and cat_id = ?"
		args = append(args, catId)
	}
	count := models.CountArticleList(search, args...)
	list := models.FindArticleList(uint(page), uint(pagesize), search, args...)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"result": gin.H{
			"list":     list,
			"count":    count,
			"pagesize": uint(pagesize),
			"page":     page,
		},
	})
}
