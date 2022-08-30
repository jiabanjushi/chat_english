package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go-fly-muti/models"
	"log"
	"time"
)

func NewKefuServer(c *gin.Context) {

	fmt.Println("======")
	kefuId, _ := c.Get("kefu_id")
	kefuInfo := models.FindUserById(kefuId)
	if kefuInfo.ID == 0 {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "用户不存在",
		})
		return
	}




	fmt.Print("----------------------------------------------")


	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	//获取GET参数,创建WS
	var kefu User
	kefu.Id = kefuInfo.Name
	kefu.Name = kefuInfo.Nickname
	kefu.Avator = kefuInfo.Avator
	kefu.Role_id = kefuInfo.RoleId
	kefu.Ent_id = fmt.Sprintf("%d", kefuInfo.ID)
	kefu.Conn = conn
	AddKefuToList(&kefu)
	go models.UpdateUserRecNumZero(kefuInfo.Name)
	for {
		//接受消息
		var receive []byte
		messageType, receive, err := conn.ReadMessage()
		if err != nil {
			log.Println("ws/user.go ", err)
			conn.Close()
			//go SendPingToKefuClient()
			//var newKefuConns = []*User{}
			//kfConns := KefuList[kefu.Id]
			//	for _, kefuConn := range kfConns {
			//		if kefuConn == nil {
			//			continue
			//		}
			//		if kefuConn.Conn != conn {
			//			newKefuConns = append(newKefuConns, kefuConn)
			//		}
			//	}
			//if len(newKefuConns) > 0 {
			//	KefuList[kefu.Id] = newKefuConns
			//} else {
			//	delete(KefuList, kefu.Id)
			//}
			return
		}

		message <- &Message{
			conn:        conn,
			content:     receive,
			context:     c,
			messageType: messageType,
		}
	}
}
func AddKefuToList(kefu *User) {
	var newKefuConns = []*User{kefu}
	kefuConns := KefuList[kefu.Id]
	if kefuConns != nil {
		for _, otherKefu := range kefuConns {
			msg := TypeMessage{
				Type: "many pong",
			}
			str, _ := json.Marshal(msg)
			err := otherKefu.Conn.WriteMessage(websocket.TextMessage, str)
			if err == nil {
				newKefuConns = append(newKefuConns, otherKefu)
			}
		}
	}

	KefuList[kefu.Id] = newKefuConns
}

//给超管发消息
func SuperAdminMessage(str []byte) {
	return
	//给超管发
	for _, kefuUsers := range KefuList {
		for _, kefuUser := range kefuUsers {
			if kefuUser.Role_id == "2" {
				kefuUser.Conn.WriteMessage(websocket.TextMessage, str)
			}
		}
	}
}

//给指定客服发消息
func OneKefuMessage(toId string, str []byte) {
	//新版
	mKefuConns, ok := KefuList[toId]
	if ok && len(mKefuConns) > 0 {
		for _, kefu := range mKefuConns {
			kefu.Mux.Lock()
			defer kefu.Mux.Unlock()
			error := kefu.Conn.WriteMessage(websocket.TextMessage, str)
			if error != nil {
				log.Println("send_kefu_message", error, string(str))
			}
		}
	}
	SuperAdminMessage(str)
}
func KefuMessage(visitorId, content string, kefuInfo models.User) {
	msg := TypeMessage{
		Type: "message",
		Data: ClientMessage{
			Name:    kefuInfo.Nickname,
			Avator:  kefuInfo.Avator,
			Id:      visitorId,
			Time:    time.Now().Format("2006-01-02 15:04:05"),
			ToId:    visitorId,
			Content: content,
			IsKefu:  "yes",
		},
	}
	str, _ := json.Marshal(msg)
	OneKefuMessage(kefuInfo.Name, str)
}

//给客服客户端发送消息判断客户端是否在线
func SendPingToKefuClient() {
	msg := TypeMessage{
		Type: "many pong",
	}
	str, _ := json.Marshal(msg)
	for kefuId, kfConns := range KefuList {
		var newKefuConns = []*User{}
		for _, kefuConn := range kfConns {
			if kefuConn == nil {
				continue
			}
			kefuConn.Mux.Lock()
			err := kefuConn.Conn.WriteMessage(websocket.TextMessage, str)
			kefuConn.Mux.Unlock()
			if err == nil {
				newKefuConns = append(newKefuConns, kefuConn)
			}
		}
		if len(newKefuConns) > 0 {
			KefuList[kefuId] = newKefuConns
		} else {
			delete(KefuList, kefuId)
		}
	}
}

//获取企业下在线的客服
func GetEntOnlineKefuId(entId string) string {
	for kefuId, kefuConn := range KefuList {
		if len(kefuConn) > 0 {
			if kefuConn[0].Ent_id == entId {
				return kefuId
			}
		}
	}
	return ""
}
