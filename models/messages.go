package models

import "fmt"

type Message struct {
	Model
	KefuId    string `json:"kefu_id"`
	VisitorId string `json:"visitor_id"`
	Content   string `json:"content"`
	MesType   string `json:"mes_type"`
	Status    string `json:"status"`
	EntId     string `json:"ent_id"`
}
type MessageKefu struct {
	Model
	KefuId        string `json:"kefu_id"`
	VisitorId     string `json:"visitor_id"`
	Content       string `json:"content"`
	MesType       string `json:"mes_type"`
	Status        string `json:"status"`
	VisitorName   string `json:"visitor_name"`
	VisitorAvator string `json:"visitor_avator"`
	KefuName      string `json:"kefu_name"`
	KefuAvator    string `json:"kefu_avator"`
	CreateTime    string `json:"create_time"`
}
type VisitorUnread struct {
	VisitorId string `json:"visitor_id"`
	Num       uint32 `json:"num"`
}

func CreateMessage(kefu_id string, visitor_id string, content string, mes_type string, ent_id string, status string) uint {
	v := &Message{
		KefuId:    kefu_id,
		VisitorId: visitor_id,
		Content:   content,
		MesType:   mes_type,
		Status:    status,
		EntId:     ent_id,
	}
	DB.Create(v)
	return v.ID
}
func FindMessageByVisitorIdUnread(visitor_id, mes_type string) []MessageKefu {
	var messages []MessageKefu
	messages = FindMessageByWhere("message.visitor_id=? and message.status='unread' and message.mes_type=?", visitor_id, mes_type)
	return messages
}
func FindMessageByVisitorId(visitor_id string) []MessageKefu {
	var messages []MessageKefu
	messages = FindMessageByWhere("message.visitor_id=?", visitor_id)
	return messages
}
func FindMessageByWhere(query interface{}, args ...interface{}) []MessageKefu {
	var messages []MessageKefu
	DB.Table("message").Where(query, args...).Select("message.*,visitor.avator visitor_avator,visitor.name visitor_name,user.avator kefu_avator,user.nickname kefu_name").Joins("left join user on message.kefu_id=user.name").Joins("left join visitor on visitor.visitor_id=message.visitor_id").Order("message.id asc").Find(&messages)
	return messages
}

//修改消息状态
func ReadMessageByVisitorId(visitor_id, mesType string) {
	message := &Message{
		Status: "read",
	}
	DB.Model(&message).Where("visitor_id=? and mes_type=?", visitor_id, mesType).Update(message)
}

//修改消息状态
func ReadMessageByEntIdVisitorId(visitor_id, ent_id, mesType string) {
	message := &Message{
		Status: "read",
	}
	DB.Model(&message).Where("visitor_id=? and ent_id=? and mes_type=?", visitor_id, ent_id, mesType).Update(message)
}

//修改消息状态
func UpdateMessageVisitorId(visitorId, newId string) {
	message := &Message{
		VisitorId: newId,
	}
	DB.Model(&message).Where("visitor_id=? ", visitorId).Update(message)
}

//获取未读数
func FindUnreadMessageNumByVisitorIds(visitor_ids []string, messageFrom string) map[string]uint32 {
	var count []VisitorUnread
	DB.Table("message").Select("count(id) num,visitor_id").Where("visitor_id in(?) and status=? and mes_type=?", visitor_ids, "unread", messageFrom).Group("visitor_id").Find(&count)
	result := make(map[string]uint32)
	for _, v := range count {
		result[v.VisitorId] = v.Num
	}
	return result
}

//查询最后一条消息
func FindLastMessage(visitorIds []string) []Message {
	var messages []Message
	if len(visitorIds) <= 0 {
		return messages
	}
	var ids []Message
	DB.Select("MAX(id) id").Where(" visitor_id in (? )", visitorIds).Group("visitor_id").Find(&ids)
	if len(ids) <= 0 {
		return messages
	}
	var idStr = make([]string, 0, 0)
	for _, mes := range ids {
		idStr = append(idStr, fmt.Sprintf("%d", mes.ID))
	}
	DB.Select("visitor_id,id,content").Where(" id in (? )", idStr).Find(&messages)
	//subQuery := DB.
	//	Table("message").
	//	Where(" visitor_id in (? )", visitorIds).
	//	Order("id desc").
	//	Limit(1024).
	//	SubQuery()
	//DB.Raw("SELECT ANY_VALUE(visitor_id) visitor_id,ANY_VALUE(id) id,ANY_VALUE(content) content FROM ? message_alia GROUP BY visitor_id", subQuery).Scan(&messages)
	//DB.Select("ANY_VALUE(visitor_id) visitor_id,MAX(ANY_VALUE(id)) id,ANY_VALUE(content) content").Group("visitor_id").Find(&messages)
	return messages
}
func FindLastMessageMap(visitorIds []string) map[string]string {
	lastMessages := FindLastMessage(visitorIds)
	temp := make(map[string]string, 0)
	for _, mes := range lastMessages {
		temp[mes.VisitorId] = mes.Content
	}
	return temp
}

//查询最后一条消息
func FindLastMessageByVisitorId(visitorId string) Message {
	var m Message
	DB.Select("content").Where("visitor_id=?", visitorId).Order("id desc").First(&m)
	return m
}

//查询条数
func CountMessage(query interface{}, args ...interface{}) uint {
	var count uint
	DB.Model(&Message{}).Where(query, args...).Count(&count)
	return count
}
func DelMessage(query interface{}, args ...interface{}) {
	DB.Model(&Message{}).Where(query, args...).Delete(&Message{})
}
func FindMessageByPage(page uint, pagesize uint, query interface{}, args ...interface{}) []*MessageKefu {
	offset := (page - 1) * pagesize
	if offset < 0 {
		offset = 0
	}
	var messages []*MessageKefu
	DB.Table("message").Select("message.*,visitor.avator visitor_avator,visitor.name visitor_name,user.avator kefu_avator,user.nickname kefu_name").Offset(offset).Joins("left join user on message.kefu_id=user.name").Joins("left join visitor on visitor.visitor_id=message.visitor_id").Where(query, args...).Limit(pagesize).Order("message.id desc").Find(&messages)
	for _, mes := range messages {
		mes.CreateTime = mes.CreatedAt.Format("2006-01-02 15:04:05")
	}
	return messages
}
