package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go-fly-muti/common"
	"go-fly-muti/models"
	"go-fly-muti/tools"
	"go-fly-muti/types"
	"go-fly-muti/ws"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type VisitorMessageForm struct {
	ToId    string `form:"to_id" json:"to_id" uri:"to_id" xml:"to_id"  binding:"required"`
	FromId  string `form:"from_id" json:"from_id" uri:"from_id" xml:"from_id" binding:"required"`
	Content string `form:"content" json:"content" uri:"content" xml:"content" binding:"required"`
}

func SendMessageV2(c *gin.Context) {
	var form VisitorMessageForm
	err := c.Bind(&form)
	if err != nil {
		c.JSON(200, gin.H{
			"code":   types.ApiCode.FAILED,
			"msg":    types.ApiCode.GetMessage(types.ApiCode.INVALID),
			"result": err.Error(),
		})
		return
	}
	fromId := form.FromId
	toId := form.ToId
	content := form.Content

	//限流
	if !tools.LimitFreqSingle("sendmessage:"+c.ClientIP(), 1, 1) {
		c.JSON(200, gin.H{
			"code": types.ApiCode.FREQ_LIMIT,
			"msg":  c.ClientIP() + types.ApiCode.GetMessage(types.ApiCode.FREQ_LIMIT),
		})
		return
	}
	//验证黑名单
	//验证访客黑名单
	if !CheckVisitorBlack(fromId) {
		c.JSON(200, gin.H{
			"code": types.ApiCode.VISITOR_BAN,
			"msg":  types.ApiCode.GetMessage(types.ApiCode.VISITOR_BAN),
		})
		return
	}

	var kefuInfo models.User
	var vistorInfo models.Visitor

	vistorInfo = models.FindVisitorByVistorId(fromId)
	kefuInfo = models.FindUser(toId)

	if kefuInfo.ID == 0 || vistorInfo.ID == 0 {
		c.JSON(200, gin.H{
			"code": types.ApiCode.ACCOUNT_NO_EXIST,
			"msg":  types.ApiCode.GetMessage(types.ApiCode.ACCOUNT_NO_EXIST),
		})
		return
	}

	msgId := models.CreateMessage(kefuInfo.Name, vistorInfo.VisitorId, content, "visitor", vistorInfo.EntId, "unread")

	guest, ok := ws.ClientList[vistorInfo.VisitorId]
	if ok && guest != nil {
		guest.UpdateTime = time.Now()
	}

	msg := ws.TypeMessage{
		Type: "message",
		Data: ws.ClientMessage{
			MsgId:     msgId,
			Avator:    vistorInfo.Avator,
			Id:        vistorInfo.VisitorId,
			VisitorId: vistorInfo.VisitorId,
			Name:      vistorInfo.Name,
			ToId:      kefuInfo.Name,
			Content:   content,
			Time:      time.Now().Format("2006-01-02 15:04:05"),
			IsKefu:    "no",
		},
	}
	str, _ := json.Marshal(msg)
	ws.OneKefuMessage(kefuInfo.Name, str)
	go SendAppGetuiPush(kefuInfo.Name, "[信息]"+vistorInfo.Name, content)
	go SendWechatVisitorMessageTemplate(kefuInfo.Name, vistorInfo.Name, content, vistorInfo.EntId)
	//go SendWechatKefuNotice(kefuInfo.Name, "[访客]"+vistorInfo.Name+",说："+content, vistorInfo.EntId)
	kefus, ok := ws.KefuList[kefuInfo.Name]
	if !ok || len(kefus) == 0 {
		go SendNoticeEmail(vistorInfo.Name, "[留言]"+vistorInfo.Name, vistorInfo.EntId, content)
	}
	go ws.VisitorAutoReply(vistorInfo, kefuInfo, content)
	go models.ReadMessageByVisitorId(vistorInfo.VisitorId, "kefu")
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
	})

}
func SendKefuMessage(c *gin.Context) {
	entId, _ := c.Get("ent_id")
	fromId, _ := c.Get("kefu_name")
	toId := c.PostForm("to_id")
	content := c.PostForm("content")
	cType := c.PostForm("type")
	if content == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "内容不能为空",
		})
		return
	}

	var kefuInfo models.User
	var vistorInfo models.Visitor
	kefuInfo = models.FindUser(fromId.(string))
	vistorInfo = models.FindVisitorByVistorId(toId)

	if kefuInfo.ID == 0 || vistorInfo.ID == 0 {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "用户不存在",
		})
		return
	}

	var msg ws.TypeMessage
	guest, ok := ws.ClientList[vistorInfo.VisitorId]
	isRead := "unread"
	msgId := models.CreateMessage(kefuInfo.Name, vistorInfo.VisitorId, content, cType, entId.(string), isRead)
	if guest != nil && ok {
		conn := guest.Conn
		msg = ws.TypeMessage{
			Type: "message",
			Data: ws.ClientMessage{
				MsgId:   msgId,
				Name:    kefuInfo.Nickname,
				Avator:  kefuInfo.Avator,
				Id:      kefuInfo.Name,
				Time:    time.Now().Format("2006-01-02 15:04:05"),
				ToId:    vistorInfo.VisitorId,
				Content: content,
				IsKefu:  "no",
			},
		}
		str, _ := json.Marshal(msg)
		conn.WriteMessage(websocket.TextMessage, str)
	} else {
		go SendWechatKefuTemplate(vistorInfo.VisitorId, kefuInfo.Name, kefuInfo.Nickname, content, fmt.Sprintf("%v", entId))
	}

	go SendWechatVisitorMessage(vistorInfo.VisitorId, content, fmt.Sprintf("%v", entId))
	msg = ws.TypeMessage{
		Type: "message",
		Data: ws.ClientMessage{
			MsgId:   msgId,
			Name:    kefuInfo.Nickname,
			Avator:  kefuInfo.Avator,
			Id:      vistorInfo.VisitorId,
			Time:    time.Now().Format("2006-01-02 15:04:05"),
			ToId:    vistorInfo.VisitorId,
			Content: content,
			IsKefu:  "yes",
		},
	}
	str2, _ := json.Marshal(msg)
	ws.OneKefuMessage(kefuInfo.Name, str2)
	go models.ReadMessageByVisitorId(vistorInfo.VisitorId, "visitor")
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}
func SendVisitorNotice(c *gin.Context) {
	notice := c.Query("msg")
	if notice == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "msg不能为空",
		})
		return
	}
	msg := ws.TypeMessage{
		Type: "notice",
		Data: notice,
	}
	str, _ := json.Marshal(msg)
	for _, visitor := range ws.ClientList {
		visitor.Conn.WriteMessage(websocket.TextMessage, str)
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}

func SendCloseMessageV2(c *gin.Context) {
	ent_id, _ := c.Get("ent_id")
	visitorId := c.Query("visitor_id")
	if visitorId == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "visitor_id不能为空",
		})
		return
	}
	config := models.FindEntConfig(ent_id, "CloseVisitorMessage")
	oldUser, ok := ws.ClientList[visitorId]
	if oldUser != nil || ok {
		ws.VisitorOffline(oldUser.To_id, oldUser.Id, oldUser.Name)
		if config.ConfValue != "" {
			kefu := models.FindUserByUid(ent_id)
			ws.VisitorMessage(visitorId, config.ConfValue, kefu)
		}
		msg := ws.TypeMessage{
			Type: "force_close",
			Data: visitorId,
		}
		str, _ := json.Marshal(msg)
		err := oldUser.Conn.WriteMessage(websocket.TextMessage, str)
		oldUser.Conn.Close()
		delete(ws.ClientList, visitorId)
		log.Println("close_message", oldUser, err)
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}
func UploadImg(c *gin.Context) {
	f, err := c.FormFile("imgfile")
	if err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "上传失败!",
		})
		return
	} else {

		fileExt := strings.ToLower(path.Ext(f.Filename))
		if fileExt != ".png" && fileExt != ".jpg" && fileExt != ".gif" && fileExt != ".jpeg" {
			c.JSON(200, gin.H{
				"code": 400,
				"msg":  "上传失败!只允许png,jpg,gif,jpeg文件",
			})
			return
		}

		fileName := tools.Md5(fmt.Sprintf("%s%s", f.Filename, time.Now().String()))
		fildDir := fmt.Sprintf("%s%d%s/", common.Upload, time.Now().Year(), time.Now().Month().String())
		isExist, _ := tools.IsFileExist(fildDir)
		if !isExist {
			os.Mkdir(fildDir, os.ModePerm)
		}
		filepath := fmt.Sprintf("%s%s%s", fildDir, fileName, fileExt)
		c.SaveUploadedFile(f, filepath)
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "上传成功!",
			"result": gin.H{
				"path": filepath,
			},
		})
	}
}
func UploadFile(c *gin.Context) {
	SendAttachment, err := strconv.ParseBool(models.FindConfig("SendAttachment"))
	if !SendAttachment || err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "禁止上传附件!",
		})
		return
	}
	f, err := c.FormFile("realfile")
	if err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "上传失败!",
		})
		return
	} else {

		fileExt := strings.ToLower(path.Ext(f.Filename))
		if f.Size >= 90*1024*1024 {
			c.JSON(200, gin.H{
				"code": 400,
				"msg":  "上传失败!不允许超过90M",
			})
			return
		}

		fileName := tools.Md5(fmt.Sprintf("%s%s", f.Filename, time.Now().String()))
		fildDir := fmt.Sprintf("%s%d%s/", common.Upload, time.Now().Year(), time.Now().Month().String())
		isExist, _ := tools.IsFileExist(fildDir)
		if !isExist {
			os.Mkdir(fildDir, os.ModePerm)
		}
		filepath := fmt.Sprintf("%s%s%s", fildDir, fileName, fileExt)
		c.SaveUploadedFile(f, filepath)
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "上传成功!",
			"result": gin.H{
				"path": filepath,
			},
		})
	}
}
func UploadAudio(c *gin.Context) {
	SendAttachment, err := strconv.ParseBool(models.FindConfig("SendAttachment"))
	if !SendAttachment || err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "禁止上传附件!",
		})
		return
	}
	f, err := c.FormFile("realfile")
	if err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "上传失败!",
		})
		return
	} else {

		fileExt := ".wav"
		if f.Size >= 20*1024*1024 {
			c.JSON(200, gin.H{
				"code": 400,
				"msg":  "上传失败!不允许超过20M",
			})
			return
		}

		fileName := tools.Md5(fmt.Sprintf("%s%s", f.Filename, time.Now().String()))
		fildDir := fmt.Sprintf("%s%d%s/", common.Upload, time.Now().Year(), time.Now().Month().String())
		isExist, _ := tools.IsFileExist(fildDir)
		if !isExist {
			os.Mkdir(fildDir, os.ModePerm)
		}
		filepath := fmt.Sprintf("%s%s%s", fildDir, fileName, fileExt)
		c.SaveUploadedFile(f, filepath)
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "上传成功!",
			"result": gin.H{
				"path": filepath,
			},
		})
	}
}
func GetMessagesV2(c *gin.Context) {
	visitorId := c.Query("visitor_id")
	messages := models.FindMessageByVisitorId(visitorId)
	//result := make([]map[string]interface{}, 0)
	chatMessages := make([]ChatMessage, 0)

	for _, message := range messages {
		//item := make(map[string]interface{})
		var chatMessage ChatMessage
		chatMessage.Time = message.CreatedAt.Format("2006-01-02 15:04:05")
		chatMessage.Content = message.Content
		chatMessage.MesType = message.MesType
		if message.MesType == "kefu" {
			chatMessage.Name = message.KefuName
			chatMessage.Avator = message.KefuAvator
		} else {
			chatMessage.Name = message.VisitorName
			chatMessage.Avator = message.VisitorAvator
		}
		chatMessages = append(chatMessages, chatMessage)
	}
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": chatMessages,
	})
}
func PostMessagesVisitorRead(c *gin.Context) {
	visitorId := c.PostForm("visitor_id")
	toId := c.PostForm("kefu")
	models.ReadMessageByVisitorId(visitorId, "kefu")
	msg := ws.TypeMessage{
		Type: "read",
		Data: ws.ClientMessage{
			VisitorId: visitorId,
			ToId:      toId,
		},
	}
	str2, _ := json.Marshal(msg)
	ws.OneKefuMessage(toId, str2)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}
func PostMessagesKefuRead(c *gin.Context) {
	visitorId := c.PostForm("visitor_id")
	entId, _ := c.Get("ent_id")
	models.ReadMessageByEntIdVisitorId(visitorId, entId.(string), "visitor")
	msg := ws.TypeMessage{
		Type: "read",
		Data: ws.ClientMessage{
			VisitorId: visitorId,
		},
	}
	ws.VisitorCustomMessage(visitorId, msg)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}
func PostMessagesAsk(c *gin.Context) {
	entId := c.PostForm("ent_id")
	content := c.PostForm("content")
	if content == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "内容不能为空",
		})
		return
	}
	entInfo := models.FindUserByUid(entId)
	reply := models.FindReplyItemByUserIdTitle(entInfo.Name, content)
	if reply.Content != "" {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "ok",
			"result": gin.H{
				"content": reply.Content,
				"name":    entInfo.Nickname,
				"avator":  entInfo.Avator,
				"time":    time.Now().Format("2006-01-02 15:05:05"),
			},
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 400,
		"msg":  "no result!",
	})
}
func GetMessagesVisitorUnread(c *gin.Context) {
	visitorId := c.Query("visitor_id")
	messages := models.FindMessageByVisitorIdUnread(visitorId, "kefu")
	chatMessages := make([]ChatMessage, 0)

	for _, message := range messages {
		//item := make(map[string]interface{})
		var chatMessage ChatMessage
		chatMessage.Time = message.CreatedAt.Format("2006-01-02 15:04:05")
		chatMessage.Content = message.Content
		chatMessage.MesType = message.MesType
		if message.MesType == "kefu" {
			chatMessage.Name = message.KefuName
			chatMessage.Avator = message.KefuAvator
		} else {
			chatMessage.Name = message.VisitorName
			chatMessage.Avator = message.VisitorAvator
		}
		chatMessages = append(chatMessages, chatMessage)
	}
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": chatMessages,
	})
}
func GetAllVisitorMessagesByKefu(c *gin.Context) {
	visitorId := c.Query("visitor_id")
	messages := models.FindMessageByVisitorId(visitorId)
	//result := make([]map[string]interface{}, 0)
	chatMessages := make([]ChatMessage, 0)

	for _, message := range messages {
		//item := make(map[string]interface{})
		var chatMessage ChatMessage
		chatMessage.Time = message.CreatedAt.Format("2006-01-02 15:04:05")
		chatMessage.Content = message.Content
		chatMessage.MesType = message.MesType
		if message.MesType == "kefu" {
			chatMessage.Name = message.KefuName
			chatMessage.Avator = message.KefuAvator
		} else {
			chatMessage.Name = message.VisitorName
			chatMessage.Avator = message.VisitorAvator
		}
		chatMessages = append(chatMessages, chatMessage)
	}
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": chatMessages,
	})
}
func GetVisitorListMessagesPage(c *gin.Context) {
	visitorId := c.Query("visitor_id")
	entId := c.Query("ent_id")
	page, _ := strconv.Atoi(c.Query("page"))
	pagesize, _ := strconv.Atoi(c.Query("pagesize"))
	if pagesize == 0 {
		pagesize = int(common.PageSize)
	}
	count := models.CountMessage("message.ent_id= ? and message.visitor_id=?", entId, visitorId)
	chatMessages := CommonMessagesPage(uint(page), uint(pagesize), visitorId, entId)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"result": gin.H{
			"count":    count,
			"page":     page,
			"list":     chatMessages,
			"pagesize": common.PageSize,
		},
	})
}
func GetVisitorListMessagesPageBykefu(c *gin.Context) {
	visitorId := c.Query("visitor_id")
	entId, _ := c.Get("ent_id")
	page, _ := strconv.Atoi(c.Query("page"))
	pagesize, _ := strconv.Atoi(c.Query("pagesize"))
	if pagesize == 0 {
		pagesize = int(common.PageSize)
	}
	count := models.CountMessage("message.ent_id= ? and message.visitor_id=?", entId.(string), visitorId)
	chatMessages := CommonMessagesPage(uint(page), uint(pagesize), visitorId, entId.(string))
	go func() {
		models.ReadMessageByVisitorId(visitorId, "visitor")
		msg := ws.TypeMessage{
			Type: "read",
			Data: ws.ClientMessage{
				VisitorId: visitorId,
			},
		}
		ws.VisitorCustomMessage(visitorId, msg)
	}()
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"result": gin.H{
			"count":    count,
			"page":     page,
			"list":     chatMessages,
			"pagesize": common.PageSize,
		},
	})
}
func DeleteMessage(c *gin.Context) {
	entId, _ := c.Get("ent_id")
	visitorId := c.PostForm("visitor_id")
	msgId := c.PostForm("msg_id")
	models.DelMessage("id=? and visitor_id=? and ent_id=?", msgId, visitorId, entId)

	msgIdInt, _ := strconv.Atoi(msgId)
	msg := ws.TypeMessage{
		Type: "delete",
		Data: ws.ClientMessage{
			MsgId:     uint(msgIdInt),
			VisitorId: visitorId,
		},
	}
	ws.VisitorCustomMessage(visitorId, msg)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}

//删除访客聊天记录
func DeleteVisitorMessage(c *gin.Context) {
	entId, _ := c.Get("ent_id")
	visitorId := c.Query("visitor_id")
	models.DelMessage("visitor_id=? and ent_id=?", visitorId, entId)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}
func CommonMessagesPage(page, pagesize uint, visitorId, entId string) []ChatMessage {

	list := models.FindMessageByPage(page, pagesize, "message.ent_id= ? and message.visitor_id=?", entId, visitorId)
	chatMessages := make([]ChatMessage, 0)

	for _, message := range list {
		//item := make(map[string]interface{})
		var chatMessage ChatMessage
		chatMessage.MsgId = message.ID
		chatMessage.Time = message.CreatedAt.Format("2006-01-02 15:04:05")
		chatMessage.Content = message.Content
		chatMessage.MesType = message.MesType
		chatMessage.ReadStatus = message.Status
		if message.MesType == "kefu" {
			chatMessage.Name = message.KefuName
			chatMessage.Avator = message.KefuAvator
		} else {
			chatMessage.Name = message.VisitorName
			chatMessage.Avator = message.VisitorAvator
		}
		chatMessages = append(chatMessages, chatMessage)
	}
	return chatMessages
}
