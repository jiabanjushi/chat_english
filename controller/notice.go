package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-fly-muti/models"
	"go-fly-muti/types"
	"go-fly-muti/ws"
	"strconv"
	"time"
)

type KefuNoticeForm struct {
	EntId     string `form:"ent_id" json:"ent_id" uri:"ent_id" xml:"ent_id"  binding:"required"`
	KefuName  string `form:"kefu_name" json:"kefu_name" uri:"kefu_name" xml:"kefu_name"`
	VisitorId string `form:"visitor_id" json:"visitor_id" uri:"visitor_id" xml:"visitor_id"`
}

func GetNotice(c *gin.Context) {
	var form KefuNoticeForm
	err := c.Bind(&form)
	if err != nil {
		c.JSON(200, gin.H{
			"code":   types.ApiCode.FAILED,
			"msg":    types.ApiCode.GetMessage(types.ApiCode.INVALID),
			"result": err.Error(),
		})
		return
	}

	if form.VisitorId != "" {
		vistorInfo := models.FindVisitorByVistorId(form.VisitorId)
		if vistorInfo.ID == 0 {
			c.JSON(200, gin.H{
				"code": types.ApiCode.ACCOUNT_NO_EXIST,
				"msg":  types.ApiCode.GetMessage(types.ApiCode.ACCOUNT_NO_EXIST),
			})
			return
		}
	}

	entId := form.EntId
	kefuId := form.KefuName
	var welcomes []models.Welcome
	var kefu models.User
	var ent models.User
	var onlineUser models.User
	var allOffline = true
	users := models.FindUsersWhere("pid = ? or id=?", entId, entId)
	//users := models.FindUsersByEntId(entId)
	agents := make([]gin.H, 0)
	for _, user := range users {
		h := gin.H{
			"name":   user.Nickname,
			"avator": user.Avator,
		}
		agents = append(agents, h)
		if fmt.Sprintf("%d", user.ID) == entId {
			ent = user
		}
		if user.Name == kefuId {
			kefu = user
		}
		if _, ok := ws.KefuList[user.Name]; ok {
			allOffline = false
			onlineUser = user
		}
	}

	if kefuId == "" && allOffline {
		kefu = ent
		onlineUser = ent
		welcomes = models.FindWelcomesByKeyword(ent.Name, "welcome")
	}
	if kefuId == "" && !allOffline {
		kefu = onlineUser
		welcomes = models.FindWelcomesByKeyword(onlineUser.Name, "welcome")
	}
	if kefuId != "" {
		welcomes = models.FindWelcomesByKeyword(kefu.Name, "welcome")
	}
	if _, ok := ws.KefuList[ent.Name]; ok {
		allOffline = false
		onlineUser = ent
	}
	result := make([]gin.H, 0)
	for _, welcome := range welcomes {
		h := gin.H{
			"name":         kefu.Nickname,
			"avator":       kefu.Avator,
			"is_kefu":      false,
			"content":      welcome.Content,
			"delay_second": welcome.DelaySecond,
			"time":         time.Now().Format("2006-01-02 15:04:05"),
		}
		result = append(result, h)
		if form.KefuName != "" && form.VisitorId != "" {
			models.CreateMessage(kefu.Name, form.VisitorId, welcome.Content, "kefu", entId, "unread")
		}
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"result": gin.H{
			"welcome":     result,
			"username":    onlineUser.Nickname,
			"avatar":      onlineUser.Avator,
			"agents":      agents,
			"all_offline": allOffline,
		},
	})
}
func GetNotices(c *gin.Context) {
	kefuId, _ := c.Get("kefu_name")
	welcomes := models.FindWelcomesByUserId(kefuId)
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": welcomes,
	})
}

//获取客服的自动欢迎信息
func GetKefuNotice(c *gin.Context) {
	kefuId, _ := c.Get("kefu_name")
	welcomes := models.FindWelcomesByUserId(kefuId)
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": welcomes,
	})
}
func PostNotice(c *gin.Context) {
	kefuId, _ := c.Get("kefu_name")
	content := c.PostForm("content")
	keyword := c.PostForm("keyword")
	delaySecond, _ := strconv.Atoi(c.PostForm("delay_second"))
	id := c.PostForm("id")
	if id != "" {
		id := c.PostForm("id")
		models.UpdateWelcome(fmt.Sprintf("%s", kefuId), id, content, uint(delaySecond))
	} else {
		models.CreateWelcome(fmt.Sprintf("%s", kefuId), content, keyword)
	}
	c.JSON(200, gin.H{
		"code": types.ApiCode.SUCCESS,
		"msg":  types.ApiCode.GetMessage(types.ApiCode.SUCCESS),
	})
}
func PostNoticeSave(c *gin.Context) {
	kefuId, _ := c.Get("kefu_name")
	content := c.PostForm("content")
	delaySecond, _ := strconv.Atoi(c.PostForm("delay_second"))
	id := c.PostForm("id")
	models.UpdateWelcome(fmt.Sprintf("%s", kefuId), id, content, uint(delaySecond))
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": "",
	})
}
func DelNotice(c *gin.Context) {
	kefuId, _ := c.Get("kefu_name")
	id := c.Query("id")
	models.DeleteWelcome(kefuId, id)
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": "",
	})
}
