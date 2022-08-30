package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-fly-muti/models"
	"go-fly-muti/tools"
	"go-fly-muti/ws"
	"math/rand"
	"time"
)

func PostRoomLogin(c *gin.Context) {
	ipcity := tools.ParseIp(c.ClientIP())
	avator := fmt.Sprintf("/static/images/%d.jpg", rand.Intn(14))
	toId := c.PostForm("to_id")
	entId := c.PostForm("ent_id")
	id := c.PostForm("visitor_id")
	if id == "" {
		id = tools.Uuid()
	}
	var (
		city string
		name string
	)
	if ipcity != nil {
		city = ipcity.CountryName + ipcity.RegionName + ipcity.CityName
		name = ipcity.CountryName + ipcity.RegionName + ipcity.CityName + "网友"
	} else {
		city = "未识别地区"
		name = "匿名网友"
	}
	client_ip := c.ClientIP()
	extra := c.PostForm("extra")
	extraJson := tools.Base64Decode(extra)
	if extraJson != "" {
		var extraObj VisitorExtra
		err := json.Unmarshal([]byte(extraJson), &extraObj)
		if err == nil {
			if extraObj.VisitorName != "" {
				name = extraObj.VisitorName
			}
			if extraObj.VisitorAvatar != "" {
				avator = extraObj.VisitorAvatar
			}
		}
	}
	if name == "" || avator == "" || toId == "" || id == "" || entId == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "error",
		})
		return
	}

	entUsers := models.FindUsersByEntId(entId)
	var flag = false
	for _, user := range entUsers {
		if user.Name == toId {
			flag = true
			break
		}
	}
	if !flag {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "room_id不存在",
		})
		return
	}

	visitor := models.FindVisitorByVistorId(id)
	if visitor.Name == "" {
		visitor = *models.CreateVisitor(name, avator, c.ClientIP(), toId, id, "", city, client_ip, entId, extra)
	}
	visitor.Name = fmt.Sprintf("#%d%s", visitor.ID, visitor.Name)
	result := Visitor{
		ID:        visitor.ID,
		City:      visitor.City,
		Name:      visitor.Name,
		Avator:    visitor.Avator,
		ToId:      visitor.ToId,
		ClientIp:  visitor.ClientIp,
		VisitorId: visitor.VisitorId,
		EntId:     visitor.EntId,
		CreatedAt: visitor.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: visitor.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	userInfo := make(map[string]string)
	userInfo["username"] = visitor.Name
	userInfo["avator"] = visitor.Avator
	msg := ws.TypeMessage{
		Type: "userOnline",
		Data: userInfo,
	}
	str, _ := json.Marshal(msg)
	go func() {
		time.Sleep(3 * time.Second)
		ws.Room.SendMessageToRoom(toId, str)
	}()
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": result,
	})
}
func PostRoomMessage(c *gin.Context) {
	fromId := c.PostForm("from_id")
	toId := c.PostForm("to_id")
	content := c.PostForm("content")
	if content == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "内容不能为空",
		})
		return
	}
	//限流
	if !tools.LimitFreqSingle("sendmessage:"+c.ClientIP(), 1, 2) {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  c.ClientIP() + "发送频率过快",
		})
		return
	}

	vistorInfo := models.FindVisitorByVistorId(fromId)
	kefuInfo := models.FindUser(toId)
	if kefuInfo.ID == 0 || vistorInfo.ID == 0 {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "room不存在",
		})
		return
	}
	models.CreateMessage(kefuInfo.Name, vistorInfo.VisitorId, content, "visitor", vistorInfo.EntId, "unread")
	msg := ws.TypeMessage{
		Type: "message",
		Data: ws.ClientMessage{
			Avator:  vistorInfo.Avator,
			Id:      vistorInfo.VisitorId,
			Name:    fmt.Sprintf("#%d%s", vistorInfo.ID, vistorInfo.Name),
			ToId:    kefuInfo.Name,
			Content: content,
			Time:    time.Now().Format("2006-01-02 15:04:05"),
			IsKefu:  "no",
		},
	}
	guest, ok := ws.ClientList[vistorInfo.VisitorId]
	if ok && guest != nil {
		guest.UpdateTime = time.Now()
	}
	str, _ := json.Marshal(msg)
	go ws.OneKefuMessage(kefuInfo.Name, str)
	go ws.Room.SendMessageToRoom(toId, str)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}
