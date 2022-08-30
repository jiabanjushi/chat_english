package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go-fly-muti/models"
	"go-fly-muti/ws"
	"time"
)

func PostCallKefu(c *gin.Context) {
	kefuId := c.PostForm("kefu_id")
	visitorId := c.PostForm("visitor_id")
	kefuInfo := models.FindUser(kefuId)
	vistorInfo := models.FindVisitorByVistorId(visitorId)
	if kefuInfo.ID == 0 || vistorInfo.ID == 0 {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "用户不存在",
		})
		return
	}
	msg := ws.TypeMessage{
		Type: "callpeer",
		Data: ws.ClientMessage{
			Avator:  vistorInfo.Avator,
			Id:      vistorInfo.VisitorId,
			Name:    vistorInfo.Name,
			ToId:    kefuInfo.Name,
			Content: "请求通话",
			Time:    time.Now().Format("2006-01-02 15:04:05"),
			IsKefu:  "no",
		},
	}
	str, _ := json.Marshal(msg)
	ws.OneKefuMessage(kefuInfo.Name, str)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}
func PostKefuPeerId(c *gin.Context) {
	peerId := c.PostForm("peer_id")
	visitorId := c.PostForm("visitor_id")
	kefuName, _ := c.Get("kefu_name")
	if peerId == "" {
		c.JSON(200, gin.H{
			"code":   400,
			"msg":    "peer_id不能为空",
			"result": "",
		})
		return
	}
	msg := ws.TypeMessage{
		Type: "peerid",
		Data: peerId,
	}
	str, _ := json.Marshal(msg)
	visitor, ok := ws.ClientList[visitorId]
	if !ok || visitor.Name == "" || kefuName != visitor.To_id {
		c.JSON(200, gin.H{
			"code":   400,
			"msg":    "客户不存在",
			"result": "",
		})
		return
	}
	visitor.Conn.WriteMessage(websocket.TextMessage, str)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}
