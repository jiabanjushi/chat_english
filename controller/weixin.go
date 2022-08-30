package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2/officialaccount"
	"github.com/silenceper/wechat/v2/officialaccount/basic"
	"go-fly-muti/lib"
	"go-fly-muti/models"
	"go-fly-muti/ws"
	"log"
	"net/http"
	"strings"
	"time"

	wechat "github.com/silenceper/wechat/v2"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

func GetCheckWeixinSign(c *gin.Context) {
	token := models.FindConfig("WeixinToken")
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")
	wechat := &lib.Wechat{
		AppId:     "",
		AppSecret: "",
		Token:     token,
	}
	res, err := wechat.CheckWechatSign(signature, timestamp, nonce, echostr)
	if err == nil {
		c.Writer.Write([]byte(res))
	} else {
		log.Println("微信API验证失败")
	}
}

//处理微信消息
func PostWechatServer(c *gin.Context) {
	kefuName := c.Param("kefuName")
	entId := c.Param("entId")
	serveWechat(c.Writer, c.Request, entId, kefuName)
}
func serveWechat(rw http.ResponseWriter, req *http.Request, entId, kefuName string) {
	//这里本地内存保存access_token，也可选择redis，memcache或者自定cache
	wechatConfig, _ := lib.NewWechatLib(entId)
	if wechatConfig == nil {
		return
	}

	cfg := &offConfig.Config{
		AppID:     wechatConfig.AppId,
		AppSecret: wechatConfig.AppSecret,
		Token:     wechatConfig.Token,
		//EncodingAESKey: "xxxx",
		Cache: memory,
	}
	wc := wechat.NewWechat()
	officialAccount := wc.GetOfficialAccount(cfg)

	// 传入request和responseWriter
	server := officialAccount.GetServer(req, rw)
	//设置接收消息的处理方法
	server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {
		//TODO

		visitorId := fmt.Sprintf("wx|%s|%s", entId, string(msg.CommonToken.FromUserName))
		avator := "/static/images/we-chat-wx.png"
		visitorName := "WechatUser"
		vistorInfo := models.FindVisitorByVistorId(visitorId)
		kefuInfo := models.FindUser(kefuName)
		isExist := true
		if vistorInfo.ID == 0 {
			//u := officialAccount.GetUser()
			//userInfo, _ := u.GetUserInfo(string(msg.CommonToken.FromUserName))
			//log.Printf(" %+v\n", userInfo)
			//if userInfo.Nickname != "" {
			//	visitorName = userInfo.Nickname
			//}
			//if userInfo.Headimgurl != "" {
			//	avator = userInfo.Headimgurl
			//}
			isExist = false
			vistorInfo = *models.CreateVisitor(visitorName, avator, "", kefuName, visitorId, "微信对话框", "", "", entId, "")
		} else {
			avator = vistorInfo.Avator
			visitorName = vistorInfo.Name
		}

		log.Printf("%+v \n", msg)
		if msg.MsgType == "text" {
			if isExist {
				go models.UpdateVisitorStatus(visitorId, 3)
			}

			//回复消息：演示回复用户发送的消息
			reply := models.FindArticleRow("ent_id = ? and find_in_set( ? , title)", vistorInfo.EntId, msg.Content)
			if reply.Content != "" {
				go ws.KefuMessage(vistorInfo.VisitorId, reply.Content, kefuInfo)
				go models.ReadMessageByVisitorId(vistorInfo.VisitorId, "kefu")
				go models.CreateMessage(kefuInfo.Name, vistorInfo.VisitorId, reply.Content, "kefu", vistorInfo.EntId, "read")
				text := message.NewText(reply.Content)
				return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
			}
		} else if msg.MsgType == "event" {
			msg.Content = string(msg.Event)
			if string(msg.Event) == "subscribe" {
				go sendWelcome(kefuName, string(msg.CommonToken.FromUserName), officialAccount)
				//客服带场景关注
				if msg.EventKey != "" {
					go models.CreateOauth(strings.Replace(msg.EventKey, "qrscene_", "", -1), string(msg.CommonToken.FromUserName))
					//访客关注不带场景
				} else {
					go models.CreateOauth(visitorId, string(msg.CommonToken.FromUserName))
				}
				go models.UpdateVisitorStatus(visitorId, 3)
			}
			//取关事件
			if string(msg.Event) == "unsubscribe" {
				go models.UpdateVisitorStatus(visitorId, 2)
				go models.DelOauth(string(msg.CommonToken.FromUserName))
			}
			if string(msg.Event) == "SCAN" {
				go models.UpdateVisitorStatus(visitorId, 3)
				go models.CreateOauth(msg.EventKey, string(msg.CommonToken.FromUserName))
			}
			//模板事件
			if msg.Content == "TEMPLATESENDJOBFINISH" || msg.Content == "VIEW" {
				return &message.Reply{}
			}

		} else if msg.MsgType == "image" {
			if isExist {
				go models.UpdateVisitorVisitorId(visitorId, visitorId)
			}
			msg.Content = fmt.Sprintf("<a href='%s' target='_blank'><img class='chatImagePic' src='%s'/></a>", msg.PicURL, msg.PicURL)
		} else {
			text := message.NewText("仅支持图片和文本")
			return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
		}

		go ws.VisitorOnline(kefuName, vistorInfo)
		go sendWeixinToKefu(kefuName, visitorName, avator, entId, visitorId, msg.Content)
		return &message.Reply{}
	})

	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		log.Println(err)
		return
	}
	//发送回复的消息
	server.Send()
}
func sendWelcome(kefuName, openId string, officialAccount *officialaccount.OfficialAccount) {
	welcomes := models.FindWelcomesByKeyword(kefuName, "wechat")
	if len(welcomes) > 0 {
		messager := officialAccount.GetCustomerMessageManager()
		for _, welcome := range welcomes {
			time.Sleep(time.Second * time.Duration(int64(welcome.DelaySecond)))
			messager.Send(message.NewCustomerTextMessage(openId, welcome.Content))
		}
	}
}
func sendWeixinToKefu(kefuName, visitorName, avator, entId, visitorId, content string) {

	msgId := models.CreateMessage(kefuName, visitorId, content, "visitor", entId, "read")
	msg := ws.TypeMessage{
		Type: "message",
		Data: ws.ClientMessage{
			MsgId:     msgId,
			Avator:    avator,
			Id:        visitorId,
			VisitorId: visitorId,
			Name:      visitorName,
			ToId:      kefuName,
			Content:   content,
			Time:      time.Now().Format("2006-01-02 15:04:05"),
			IsKefu:    "no",
		},
	}
	str, _ := json.Marshal(msg)
	ws.OneKefuMessage(kefuName, str)
	go SendAppGetuiPush(kefuName, "[信息]"+visitorName, content)
}

//查询绑定的oauth
func GetWechatOauth(c *gin.Context) {
	visitorId := c.Query("visitor_id")
	visitorIdArr := strings.Split(visitorId, "|")
	var oauth models.Oauth
	if len(visitorIdArr) >= 3 && visitorIdArr[0] == "wx" {
		args := []interface{}{
			visitorIdArr[2],
		}
		oauth = models.FindOauthsQuery("oauth_id = ? ", args)
	} else {
		oauth = models.FindOauthById(visitorId)
	}
	if oauth.OauthId == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "未绑定公众号",
		})
		return
	}
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": oauth.OauthId,
	})
}
func bindVisitorWechat(snowId, openId string) bool {
	visitor := models.FindVisitorByVistorId(snowId)
	if visitor.EntId == "" {
		return false
	}
	newId := fmt.Sprintf("%s|%s", visitor.EntId, openId)
	models.UpdateVisitorVisitorId(snowId, newId)
	models.UpdateMessageVisitorId(snowId, newId)
	msg := ws.TypeMessage{
		Type: "change_id",
		Data: ws.SimpleMessage{
			From:    snowId,
			To:      newId,
			Content: "绑定微信",
		},
	}
	ws.VisitorCustomMessage(snowId, msg)
	return true
}

//展示带参二维码
func GetShowQrCode(c *gin.Context) {

	entId := c.Query("entId")
	sceneName := c.Query("sceneName")
	configs := models.FindEntConfigs(entId)
	AppID := ""
	AppSecret := ""
	Token := ""
	for _, config := range configs {
		if config.ConfKey == "WechatAppId" {
			AppID = config.ConfValue
		}
		if config.ConfKey == "WechatAppSecret" {
			AppSecret = config.ConfValue
		}
		if config.ConfKey == "WechatAppToken" {
			Token = config.ConfValue
		}
	}
	if AppID == "" || AppSecret == "" || Token == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "AppID,AppSecret,Token不能为空",
		})
		return
	}
	cfg := &offConfig.Config{
		AppID:     AppID,
		AppSecret: AppSecret,
		Token:     Token,
		//EncodingAESKey: "xxxx",
		Cache: memory,
	}
	wc := wechat.NewWechat()
	officialAccount := wc.GetOfficialAccount(cfg)
	basicObj := officialAccount.GetBasic()
	tq := basic.NewTmpQrRequest(time.Duration(24*3600), sceneName)
	ticket, err := basicObj.GetQRTicket(tq)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	url := basic.ShowQRCode(ticket)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"result": gin.H{
			"ticket": ticket,
			"url":    url,
		},
	})
}
