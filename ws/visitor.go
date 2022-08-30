package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go-fly-muti/common"
	"go-fly-muti/models"
	"log"
	"time"
)

func NewVisitorServer(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)



	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	//获取GET参数,创建WS
	toId := c.Query("to_id")
	visitorId := c.Query("visitor_id")
	if toId == "" || visitorId == "" {
		log.Println("访客ws参数为空")
		conn.Close()
		return
	}
	vistorInfo := models.FindVisitorByVistorId(visitorId)
	if vistorInfo.VisitorId == "" {
		log.Println("访客visitorId不存在:", visitorId)
		conn.Close()
		return
	}
	user := &User{
		Conn:       conn,
		Name:       fmt.Sprintf("#%d%s", vistorInfo.ID, vistorInfo.Name),
		Avator:     vistorInfo.Avator,
		Id:         vistorInfo.VisitorId,
		To_id:      toId,
		Ent_id:     vistorInfo.EntId,
		UpdateTime: time.Now(),
	}
	AddVisitorToList(user)
	VisitorOnline(toId, vistorInfo)
	//判断是否有room_id,代表聊天室
	roomId := c.Query("room_id")
	if roomId != "" {
		AddVisitorToRoom(roomId, vistorInfo.VisitorId)
	}
	for {
		//接受消息
		var receive []byte
		messageType, receive, err := conn.ReadMessage()
		if err != nil {
			log.Println("ws/visitor.go conn.ReadMessage:", err, receive, messageType)
			conn.Close()
			//if visitor, ok := ClientList[vistorInfo.VisitorId]; ok {
			//	VisitorOffline(visitor.To_id, visitor.Id, visitor.Name)
			//	//time.Sleep(time.Duration(common.WsBreakTimeout) * time.Second)
			//	delete(ClientList, visitor.Id)
			//}
			for _, visitor := range ClientList {
				if visitor.Conn == conn {
					//conn.Close()
					delete(ClientList, visitor.Id)
					VisitorOffline(visitor.To_id, visitor.Id, visitor.Name)
					return
				}
			}
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
func AddVisitorToList(user *User) {
	//用户id对应的连接
	oldUser, ok := ClientList[user.Id]
	if oldUser != nil || ok {
		VisitorOffline(oldUser.To_id, oldUser.Id, oldUser.Name)
		closemsg := TypeMessage{
			Type: "close",
			Data: user.Id,
		}
		closeStr, _ := json.Marshal(closemsg)
		if err := oldUser.Conn.WriteMessage(websocket.TextMessage, closeStr); err != nil {
			oldUser.Conn.Close()
			user.UpdateTime = oldUser.UpdateTime
			delete(ClientList, user.Id)
		}
	}
	ClientList[user.Id] = user
}
func AddVisitorToRoom(roomId, visitorId string) {
	//用户id对应的连接
	user, ok := ClientList[visitorId]
	if user != nil && ok {
		members, _ := Room.GetMembers(roomId)
		members = append(members, user)
		Room.SetMembers(roomId, members)
	}
}
func VisitorOnline(kefuId string, visitor models.Visitor) {
	go models.UpdateUserRecNum(kefuId, 1)
	lastMessage := models.FindLastMessageByVisitorId(visitor.VisitorId)
	unreadMap := models.FindUnreadMessageNumByVisitorIds([]string{visitor.VisitorId}, "visitor")
	var unreadNum uint32
	if num, ok := unreadMap[visitor.VisitorId]; ok {
		unreadNum = num
	}

	userInfo := make(map[string]string)
	userInfo["uid"] = visitor.VisitorId
	userInfo["visitor_id"] = visitor.VisitorId
	userInfo["username"] = fmt.Sprintf("#%d%s", visitor.ID, visitor.Name)
	userInfo["avator"] = visitor.Avator
	userInfo["last_message"] = lastMessage.Content
	userInfo["unread_num"] = fmt.Sprintf("%d", unreadNum)
	if userInfo["last_message"] == "" {
		userInfo["last_message"] = "新访客"
	}
	msg := TypeMessage{
		Type: "userOnline",
		Data: userInfo,
	}
	str, _ := json.Marshal(msg)
	OneKefuMessage(kefuId, str)
}
func VisitorOffline(kefuId string, visitorId string, visitorName string) {
	go models.UpdateUserRecNum(kefuId, -1)
	userInfo := make(map[string]string)
	userInfo["uid"] = visitorId
	userInfo["visitor_id"] = visitorId
	userInfo["name"] = visitorName
	msg := TypeMessage{
		Type: "userOffline",
		Data: userInfo,
	}
	str, _ := json.Marshal(msg)
	//新版
	OneKefuMessage(kefuId, str)
}
func VisitorNotice(visitorId string, notice string) {
	msg := TypeMessage{
		Type: "notice",
		Data: notice,
	}
	str, _ := json.Marshal(msg)
	visitor, ok := ClientList[visitorId]
	if !ok || visitor == nil || visitor.Conn == nil {
		return
	}
	visitor.Conn.WriteMessage(websocket.TextMessage, str)
}
func VisitorCustomMessage(visitorId string, notice TypeMessage) {
	str, _ := json.Marshal(notice)
	visitor, ok := ClientList[visitorId]
	if !ok || visitor == nil || visitor.Conn == nil {
		return
	}
	visitor.Conn.WriteMessage(websocket.TextMessage, str)
}
func VisitorTransfer(visitorId string, kefuId string) {
	msg := TypeMessage{
		Type: "transfer",
		Data: kefuId,
	}
	str, _ := json.Marshal(msg)
	visitor, ok := ClientList[visitorId]
	if !ok || visitor == nil || visitor.Conn == nil {
		return
	}
	visitor.Conn.WriteMessage(websocket.TextMessage, str)
}
func VisitorMessage(visitorId, content string, kefuInfo models.User) {
	msg := TypeMessage{
		Type: "message",
		Data: ClientMessage{
			Name:    kefuInfo.Nickname,
			Avator:  kefuInfo.Avator,
			Id:      kefuInfo.Name,
			Time:    time.Now().Format("2006-01-02 15:04:05"),
			ToId:    visitorId,
			Content: content,
			IsKefu:  "no",
		},
	}
	str, _ := json.Marshal(msg)
	visitor, ok := ClientList[visitorId]
	if !ok || visitor == nil || visitor.Conn == nil {
		return
	}
	visitor.Conn.WriteMessage(websocket.TextMessage, str)
}
func VisitorAutoReply(vistorInfo models.Visitor, kefuInfo models.User, content string) {
	//var entInfo models.User
	//if fmt.Sprintf("%v", kefuInfo.Pid) == fmt.Sprintf("%d", 1) {
	//	entInfo = kefuInfo
	//} else {
	//	entInfo = models.FindUserByUid(kefuInfo.Pid)
	//}
	reply := models.FindArticleRow("ent_id = ? and find_in_set( ? , title)", vistorInfo.EntId, content)
	//reply = models.FindReplyItemByUserIdTitle(entInfo.Name, content)

	if reply.Content != "" {
		time.Sleep(1 * time.Second)
		VisitorMessage(vistorInfo.VisitorId, reply.Content, kefuInfo)
		KefuMessage(vistorInfo.VisitorId, reply.Content, kefuInfo)
		models.CreateMessage(kefuInfo.Name, vistorInfo.VisitorId, reply.Content, "kefu", vistorInfo.EntId, "read")
	}
}
func CleanVisitorExpire() {
	go func() {
		log.Println("cleanVisitorExpire start...")
		for {
			for _, user := range ClientList {
				diff := time.Now().Sub(user.UpdateTime).Seconds()
				if diff >= common.VisitorExpire {
					entConfig := models.FindEntConfig(user.Ent_id, "CloseVisitorMessage")
					if entConfig.ConfValue != "" {
						kefu := models.FindUserByUid(user.Ent_id)
						VisitorMessage(user.Id, entConfig.ConfValue, kefu)
					}
					msg := TypeMessage{
						Type: "auto_close",
						Data: user.Id,
					}
					str, _ := json.Marshal(msg)
					if err := user.Conn.WriteMessage(websocket.TextMessage, str); err != nil {
						user.Conn.Close()
						delete(ClientList, user.Id)
						VisitorOffline(user.To_id, user.Id, user.Name)
					}
					log.Println(user.Name + ":cleanVisitorExpire finshed")
				}
			}
			t := time.NewTimer(time.Second * 5)
			<-t.C
		}
	}()
}
