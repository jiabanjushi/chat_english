package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"go-fly-muti/common"
	"go-fly-muti/lib"
	"go-fly-muti/models"
	"log"
	"strconv"
	"strings"
	"time"
)

var IP_SERVER_URL = ""
var memory = cache.NewMemory()

/**
处理分页页码
*/
func HandlePagePageSize(c *gin.Context) (uint, uint) {
	var page uint
	pagesize := common.VisitorPageSize
	myPage, _ := strconv.Atoi(c.Query("page"))
	if myPage != 0 {
		page = uint(myPage)
	}
	myPagesize, _ := strconv.Atoi(c.Query("pagesize"))
	if myPagesize != 0 {
		pagesize = uint(myPagesize)
	}
	return page, pagesize
}

//发送访客微信消息
func SendWechatVisitorMessage(visitorId, content, entId string) bool {
	visitorIdArr := strings.Split(visitorId, "|")
	if len(visitorIdArr) < 3 || visitorIdArr[0] != "wx" {
		return false
	}
	return SendWechatMesage(visitorIdArr[2], content, entId)
}

//发送客服微信消息
func SendWechatKefuNotice(kefuName, content, entId string) bool {
	oauth := models.FindOauthById(kefuName)
	if oauth.OauthId == "" {
		return false
	}
	return SendWechatMesage(oauth.OauthId, content, entId)
}

//发送新访客提醒模板消息
func SendWechatVisitorTemplate(kefuName, visitorName, content, entId string) bool {
	oauths := models.FindOauthsById(kefuName)
	if len(oauths) == 0 {
		return false
	}
	wechatConfig, _ := lib.NewWechatLib(entId)
	if wechatConfig.WechatVisitorTemplateId == "" {
		return false
	}
	msgData := make(map[string]*message.TemplateDataItem)
	msgData["keyword1"] = &message.TemplateDataItem{
		Value: visitorName,
		Color: "",
	}
	msgData["keyword2"] = &message.TemplateDataItem{
		Value: time.Now().Format("2006-01-02 15:04:05"),
		Color: "",
	}
	msgData["keyword3"] = &message.TemplateDataItem{
		Value: content,
		Color: "",
	}
	for _, oauth := range oauths {
		msg := &message.TemplateMessage{
			ToUser:     oauth.OauthId,
			Data:       msgData,
			TemplateID: wechatConfig.WechatVisitorTemplateId,
			URL:        "https://" + wechatConfig.WechatHost + "/static/h5/#/",
		}
		SendWechatTemplate(wechatConfig, msg)
	}
	return true
}

//发送访客新消息提醒模板消息
func SendWechatVisitorMessageTemplate(kefuName, visitorName, content, entId string) bool {
	oauths := models.FindOauthsById(kefuName)
	if len(oauths) == 0 {
		return false
	}
	wechatConfig, _ := lib.NewWechatLib(entId)
	if wechatConfig.WechatMessageTemplateId == "" {
		return false
	}
	msgData := make(map[string]*message.TemplateDataItem)
	msgData["keyword1"] = &message.TemplateDataItem{
		Value: visitorName,
		Color: "",
	}
	msgData["keyword2"] = &message.TemplateDataItem{
		Value: time.Now().Format("2006-01-02 15:04:05"),
		Color: "",
	}
	msgData["keyword3"] = &message.TemplateDataItem{
		Value: content,
		Color: "",
	}
	msgData["remark"] = &message.TemplateDataItem{
		Value: models.FindConfig("WechatTemplateRemark"),
		Color: "",
	}
	for _, oauth := range oauths {
		msg := &message.TemplateMessage{
			ToUser:     oauth.OauthId,
			Data:       msgData,
			TemplateID: wechatConfig.WechatMessageTemplateId,
			URL:        "https://" + wechatConfig.WechatHost + "/static/h5/#/",
		}
		SendWechatTemplate(wechatConfig, msg)
	}
	return true
}

//发送客服回复模板消息
func SendWechatKefuTemplate(visitorId, kefuName, kefuNickname, content, entId string) bool {
	visitorIdArr := strings.Split(visitorId, "|")
	if len(visitorIdArr) < 3 || visitorIdArr[0] != "wx" {
		return false
	}
	wechatConfig, _ := lib.NewWechatLib(entId)
	if wechatConfig.WechatKefuTemplateId == "" {
		return false
	}
	msgData := make(map[string]*message.TemplateDataItem)
	msgData["keyword1"] = &message.TemplateDataItem{
		Value: kefuNickname,
		Color: "",
	}
	msgData["keyword2"] = &message.TemplateDataItem{
		Value: time.Now().Format("2006-01-02 15:04:05"),
		Color: "",
	}
	msgData["keyword3"] = &message.TemplateDataItem{
		Value: content,
		Color: "",
	}
	msgData["remark"] = &message.TemplateDataItem{
		Value: models.FindConfig("WechatTemplateRemark"),
		Color: "",
	}
	msg := &message.TemplateMessage{
		ToUser:     visitorIdArr[2],
		Data:       msgData,
		TemplateID: wechatConfig.WechatKefuTemplateId,
		URL: "https://" + wechatConfig.WechatHost + "/chatIndex" +
			"?ent_id=" + entId +
			"&kefu_id=" + kefuName +
			"&visitor_id=" + visitorId,
	}
	return SendWechatTemplate(wechatConfig, msg)
}

//发送微信模板消息
func SendWechatTemplate(wechatConfig *lib.Wechat, msg *message.TemplateMessage) bool {

	if wechatConfig == nil {
		return false
	}

	wc := wechat.NewWechat()
	cfg := &offConfig.Config{
		AppID:     wechatConfig.AppId,
		AppSecret: wechatConfig.AppSecret,
		Token:     wechatConfig.Token,
		//EncodingAESKey: "xxxx",
		Cache: memory,
	}
	officialAccount := wc.GetOfficialAccount(cfg)
	template := officialAccount.GetTemplate()

	msgId, err := template.Send(msg)
	log.Println("发送微信模板消息：", msgId, err)
	return true
}

//发送微信客服消息
func SendWechatMesage(openId, content, entId string) bool {
	wechatConfig, _ := lib.NewWechatLib(entId)
	if wechatConfig == nil || wechatConfig.WechatKefu == "" || wechatConfig.WechatKefu == "off" {
		return false
	}
	wc := wechat.NewWechat()
	cfg := &offConfig.Config{
		AppID:     wechatConfig.AppId,
		AppSecret: wechatConfig.AppSecret,
		Token:     wechatConfig.Token,
		//EncodingAESKey: "xxxx",
		Cache: memory,
	}
	officialAccount := wc.GetOfficialAccount(cfg)
	messager := officialAccount.GetCustomerMessageManager()
	err := messager.Send(message.NewCustomerTextMessage(openId, content))
	log.Println("发送微信客服消息：", openId, content, err)
	return true
}

//验证访客黑名单
func CheckVisitorBlack(visitorId string) bool {
	black := models.FindVisitorBlack("visitor_id = ?", visitorId)
	if black.Id != 0 {
		return false
	}
	return true
}
