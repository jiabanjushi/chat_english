package controller

import (
	"encoding/json"
	"fmt"
	"go-fly-muti/lib"
	"go-fly-muti/models"
	"go-fly-muti/tools"
	"go-fly-muti/ws"
	"log"
	"strconv"
)

func SendServerJiang(title string, content string, domain string) string {
	noticeServerJiang, err := strconv.ParseBool(models.FindConfig("NoticeServerJiang"))
	serverJiangAPI := models.FindConfig("ServerJiangAPI")
	if err != nil || !noticeServerJiang || serverJiangAPI == "" {
		log.Println("do not notice serverjiang:", serverJiangAPI, noticeServerJiang)
		return ""
	}
	sendStr := fmt.Sprintf("%s%s", title, content)
	desp := title + ":" + content + "[登录](http://" + domain + "/main)"
	url := serverJiangAPI + "?text=" + sendStr + "&desp=" + desp
	//log.Println(url)
	res := tools.Get(url)
	return res
}
func SendVisitorLoginNotice(kefuName, visitorName, avator, content, visitorId string) {
	if !tools.LimitFreqSingle("sendnotice:"+visitorId, 1, 5) {
		log.Println("SendVisitorLoginNotice limit")
		return
	}
	userInfo := make(map[string]string)
	userInfo["username"] = visitorName
	userInfo["avator"] = avator
	userInfo["content"] = content
	msg := ws.TypeMessage{
		Type: "notice",
		Data: userInfo,
	}
	str, _ := json.Marshal(msg)
	ws.OneKefuMessage(kefuName, str)
}
func SendNoticeEmail(username, subject, entId, content string) {
	if !tools.LimitFreqSingle("send_notice_email:"+username, 1, 10) {
		log.Println("send_notice_email limit")
		return
	}
	configs := models.FindEntConfigByEntid(entId)
	smtp := ""
	email := ""
	password := ""
	for _, config := range configs {
		if config.ConfKey == "NoticeEmailAddress" {
			email = config.ConfValue
		}
		if config.ConfKey == "NoticeEmailPassword" {
			password = config.ConfValue
		}
		if config.ConfKey == "NoticeEmailSmtp" {
			smtp = config.ConfValue
		}
	}
	if smtp == "" || email == "" || password == "" {
		return
	}
	log.Println("发送访客通知邮件:" + smtp + "," + email + "," + password)
	err := tools.SendSmtp(smtp, email, password, []string{email}, subject, content)
	if err != nil {
		log.Println("发送访客通知邮件失败:")
		log.Println(err)
	}
}
func SendEntSmtpEmail(subject, msg, entId string) {
	if !tools.LimitFreqSingle("send_ent_email:"+entId, 1, 2) {
		log.Println("send_ent_email limit")
		return
	}
	configs := models.FindEntConfigByEntid(entId)
	smtp := ""
	email := ""
	password := ""
	for _, config := range configs {
		if config.ConfKey == "NoticeEmailAddress" {
			email = config.ConfValue
		}
		if config.ConfKey == "NoticeEmailPassword" {
			password = config.ConfValue
		}
		if config.ConfKey == "NoticeEmailSmtp" {
			smtp = config.ConfValue
		}
	}
	if smtp == "" || email == "" || password == "" {
		return
	}
	err := tools.SendSmtp(smtp, email, password, []string{email}, subject, msg)
	if err != nil {
		log.Println(err)
	}
}
func SendAppGetuiPush(kefu string, title, content string) {
	clientModel := &models.User_client{
		Kefu: kefu,
	}
	clientInfos := clientModel.FindClients()
	if len(clientInfos) == 0 {
		return
	}
	appid := models.FindConfig("GetuiAppID")
	appkey := models.FindConfig("GetuiAppKey")
	appsecret := models.FindConfig("GetuiAppSecret")
	appmastersecret := models.FindConfig("GetuiMasterSecret")
	token := models.FindConfig("GetuiToken")
	getui := &lib.Getui{
		AppId:           appid,
		AppKey:          appkey,
		AppSecret:       appsecret,
		AppMasterSecret: appmastersecret,
	}
	for _, client := range clientInfos {
		res, err := getui.PushSingle(token, client.Client_id, title, content)
		//不正确的账号
		if res == 20001 && err.Error() == "target user is invalid" {
			clientModel2 := &models.User_client{
				Kefu:      kefu,
				Client_id: client.Client_id,
			}
			clientModel2.DeleteClient()
		}
		if res == 10001 {
			token, _ = getui.GetGetuiToken()
			models.UpdateConfig("GetuiToken", token)
			getui.PushSingle(token, client.Client_id, title, content)
		}
	}
}
