package controller

import (
	"go-fly-muti/models"
	"time"
)

var (
	Port    string
	Address string
)

type Response struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	result interface{} `json:"result"`
}
type ChatMessage struct {
	MsgId      uint   `json:"msg_id"`
	Time       string `json:"time"`
	Content    string `json:"content"`
	MesType    string `json:"mes_type"`
	Name       string `json:"name"`
	Avator     string `json:"avator"`
	ReadStatus string `json:"read_status"`
}
type VisitorOnline struct {
	Id          uint      `json:"id"`
	VisitorId   string    `json:"visitor_id"`
	Username    string    `json:"username"`
	Avator      string    `json:"avator"`
	Ip          string    `json:"ip"`
	LastMessage string    `json:"last_message"`
	City        string    `json:"city"`
	UpdatedAt   time.Time `json:"updated_at"`
	UnreadNum   uint32    `json:"unread_num"`
	Status      uint      `json:"status"`
}
type GetuiResponse struct {
	Code float64                `json:"code"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}
type VisitorExtra struct {
	VisitorName   string `json:"visitorName"`
	VisitorAvatar string `json:"visitorAvatar"`
	VisitorId     string `json:"visitorId"`
}
type VisitorExtend struct {
	ID        uint   `gorm:"primary_key" json:"id"`
	VisitorId string `json:"visitor_id"`
	Url       string `json:"url"`
	Ua        string `json:"ua"`
	Title     string `json:"title"`
	ClientIp  string `json:"client_ip"`
	CreatedAt string `json:"created_at"`
	Browser   string `json:"browser"`
	OsVersion string `json:"os_version"`
}
type Visitor struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Avator    string `json:"avator"`
	ToId      string `json:"to_id"`
	VisitorId string `json:"visitor_id"`
	City      string `json:"city"`
	ClientIp  string `json:"client_ip"`
	EntId     string `json:"ent_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
type VisitorAttrParams struct {
	VisitorId   string              `json:"visitor_id"`
	VisitorAttr models.Visitor_attr `json:"visitor_attr"`
}
