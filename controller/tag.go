package controller

import (
	"github.com/gin-gonic/gin"
	models "go-fly-muti/models/v2"
	"go-fly-muti/types"
	"log"
	"strconv"
)

type VisitorTag struct {
	VisitorId string `binding:"required" form:"visitor_id" json:"visitor_id" uri:"visitor_id" xml:"visitor_id"`
	TagName   string `binding:"required" form:"tag_name" json:"tag_name" uri:"tag_name" xml:"tag_name"`
}

func PostVisitorTag(c *gin.Context) {
	var form VisitorTag
	err := c.Bind(&form)
	if err != nil {
		c.JSON(200, gin.H{
			"code":   types.ApiCode.FAILED,
			"msg":    types.ApiCode.GetMessage(types.ApiCode.INVALID),
			"result": err.Error(),
		})
		return
	}

	kefuId, _ := c.Get("kefu_name")
	entIdStr, _ := c.Get("ent_id")
	entId, _ := strconv.Atoi(entIdStr.(string))
	tagModel := models.GetTag("name = ? and ent_id = ?", form.TagName, entId)
	if tagModel.ID == 0 {
		tagModel.Name = form.TagName
		tagModel.Kefu = kefuId.(string)
		tagModel.EntId = uint(entId)
		tagModel.InsertTag()
	}
	tagIds := models.GetVisitorTags("visitor_id = ? and ent_id = ?", form.VisitorId, entId)
	for _, tagId := range tagIds {
		if tagId.TagId == tagModel.ID {
			models.DelVisitorTags("visitor_id = ? and tag_id = ? and ent_id = ?", form.VisitorId,
				tagId.TagId, entId)
		}
	}

	visitorTagModel := models.VisitorTag{
		VisitorId: form.VisitorId,
		TagId:     tagModel.ID,
		Kefu:      kefuId.(string),
		EntId:     uint(entId),
	}
	visitorTagModel.InsertVisitorTag()

	c.JSON(200, gin.H{
		"code": types.ApiCode.SUCCESS,
		"msg":  types.ApiCode.GetMessage(types.ApiCode.SUCCESS),
	})
}
func DelVisitorTag(c *gin.Context) {
	var form VisitorTag
	err := c.Bind(&form)
	if err != nil {
		c.JSON(200, gin.H{
			"code":   types.ApiCode.FAILED,
			"msg":    types.ApiCode.GetMessage(types.ApiCode.INVALID),
			"result": err.Error(),
		})
		return
	}

	entIdStr, _ := c.Get("ent_id")
	entId, _ := strconv.Atoi(entIdStr.(string))
	tagModel := models.GetTag("name = ? and ent_id = ?", form.TagName, entId)

	models.DelVisitorTags("visitor_id = ? and tag_id = ? and ent_id = ?", form.VisitorId,
		tagModel.ID, entId)

	c.JSON(200, gin.H{
		"code": types.ApiCode.SUCCESS,
		"msg":  types.ApiCode.GetMessage(types.ApiCode.SUCCESS),
	})
}
func GetVisitorAllTags(c *gin.Context) {
	visitorId := c.Query("visitor_id")
	entIdStr, _ := c.Get("ent_id")
	entId, _ := strconv.Atoi(entIdStr.(string))

	alltags := models.GetTags("ent_id = ?", entId)
	tagIds := models.GetVisitorTags("visitor_id = ? and ent_id = ?", visitorId, entId)
	tagMap := make(map[uint]uint, 0)
	for _, tagId := range tagIds {
		tagMap[tagId.TagId] = 1
	}
	log.Println(tagIds, tagMap)

	for _, tag := range alltags {
		if _, ok := tagMap[tag.ID]; ok {
			tag.IsTaged = 1
		} else {
			tag.IsTaged = 0
		}
	}

	c.JSON(200, gin.H{
		"code":   types.ApiCode.SUCCESS,
		"msg":    types.ApiCode.GetMessage(types.ApiCode.SUCCESS),
		"result": alltags,
	})
}



func GetVisitorTags(c *gin.Context) {
	visitorId := c.Query("visitor_id")
	entIdStr, _ := c.Get("ent_id")
	entId, _ := strconv.Atoi(entIdStr.(string))

	tags := models.GetVisitorTagsByVisitorId(visitorId, uint(entId))

	c.JSON(200, gin.H{
		"code":   types.ApiCode.SUCCESS,
		"msg":    types.ApiCode.GetMessage(types.ApiCode.SUCCESS),
		"result": tags,
	})
}
func GetTags(c *gin.Context) {

	entIdStr, _ := c.Get("ent_id")
	entId, _ := strconv.Atoi(entIdStr.(string))
	tagName := c.Query("tag_name")
	var tags []*models.Tag
	if tagName != "" {
		tags = models.GetTags("ent_id = ? and name like ?", entId, tagName+"%")
	} else {
		tags = models.GetTags("ent_id = ?", entId)
	}

	c.JSON(200, gin.H{
		"code":   types.ApiCode.SUCCESS,
		"msg":    types.ApiCode.GetMessage(types.ApiCode.SUCCESS),
		"result": tags,
	})
}
