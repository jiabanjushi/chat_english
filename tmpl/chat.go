package tmpl

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/oauth"
	"go-fly-muti/common"
	"go-fly-muti/lib"
	"go-fly-muti/models"
	"net/http"
)

var memory = cache.NewMemory()

//聊天室界面
func PageChatRoom(c *gin.Context) {
	c.HTML(http.StatusOK, "chat_room.html", gin.H{})
}

//咨询界面
func PageChat(c *gin.Context) {
	kefuId := c.Query("kefu_id")
	lang, _ := c.Get("lang")
	refer := c.Query("refer")
	entId := c.Query("ent_id")
	if refer == "" {
		refer = c.Request.Referer()
	}
	if refer == "" {
		refer = ""
	}

	visitorId := c.Query("visitor_id")
	visitorName := c.Query("visitor_name")
	avator := c.Query("avator")
	errMsg := ""
	//获取微信用户消息
	weixinCode := c.Query("code")
	if weixinCode != "" && entId != "" {
		userInfo, err := GetWechatUserInfo(weixinCode, entId)
		if err != nil {
			errMsg = err.Error()
		} else {
			visitorId = "wx|" + entId + "|" + userInfo.OpenID
			avator = userInfo.HeadImgURL
			visitorName = userInfo.Nickname
		}
	}
	//落地域名跳转
	landHost := models.FindConfig("LandHost")
	if landHost != "" && landHost != c.Request.Host {
		c.Redirect(302, fmt.Sprintf("//%s/chatIndex?kefu_id=%s&ent_id=%s&lang=%s&visitor_id=%s&visitor_name=%s&avator=%s", landHost,
			kefuId, entId, lang, visitorId, visitorName, avator))
		return
	}
	SystemNotice := models.FindConfig("SystemNotice")
	SystemTitle := models.FindConfig("SystemTitle")
	SystemKeywords := models.FindConfig("SystemKeywords")
	SystemDesc := models.FindConfig("SystemDesc")
	CopyrightTxt := models.FindConfig("CopyrightTxt")
	CopyrightUrl := models.FindConfig("CopyrightUrl")
	ShowKefuName := models.FindConfig("ShowKefuName")


fmt.Println("????????????????????????????????????????????")
	fmt.Println(CopyrightTxt)

	entInfo := models.FindUserById(entId)
	title := "GOFLY客服-免费在线客服系统源码-网站开源在线客服系统-私有化部署网页在线客服软件代码下载"
	if entInfo.Nickname != "" {
		title = entInfo.Nickname
	}


	c.HTML(http.StatusOK, "chat_page.html", gin.H{
		"KEFU_ID":        kefuId,
		"Lang":           lang.(string),
		"Refer":          refer,
		"ENT_ID":         entId,
		"visitorId":      visitorId,
		"visitorName":    visitorName,
		"avator":         avator,
		"errMsg":         errMsg,
		"IS_TRY":         common.IsTry,
		"SystemTitle":    SystemTitle,
		"SystemKeywords": SystemKeywords,
		"SystemDesc":     SystemDesc,
		"CopyrightTxt":   CopyrightTxt,
		"CopyrightUrl":   CopyrightUrl,
		"SystemNotice":   SystemNotice,
		"Title":          title,
		"ShowKefuName":   ShowKefuName,
	})
}

// PageWechatChat 微信跳转界面
func PageWechatChat(c *gin.Context) {
	kefuId := c.Query("kefu_id")
	lang, _ := c.Get("lang")
	refer := c.Query("refer")
	entId := c.Query("ent_id")
	if refer == "" {
		refer = c.Request.Referer()
	}
	if refer == "" {
		refer = ""
	}
	appId := ""
	host := ""
	wechatConfig, _ := lib.NewWechatLib(entId)
	if wechatConfig != nil {
		appId = wechatConfig.AppId
		host = wechatConfig.WechatHost
	}
	entInfo := models.FindUserById(entId)
	c.HTML(http.StatusOK, "chat_wechat_page.html", gin.H{
		"Title":  entInfo.Nickname,
		"kefuId": kefuId,
		"Lang":   lang.(string),
		"Refer":  refer,
		"EntId":  entId,
		"AppId":  appId,
		"Host":   host,
	})
}

// GetWechatUserInfo 获取微信用户信息
func GetWechatUserInfo(weixinCode, entId string) (oauth.UserInfo, error) {
	var userinfo oauth.UserInfo
	wechatConfig, err := lib.NewWechatLib(entId)
	if wechatConfig == nil {
		return userinfo, err
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
	oauth := officialAccount.GetOauth()
	accessToken, err := oauth.GetUserAccessToken(weixinCode)
	if err != nil {
		return userinfo, err
	}
	userinfo, err = oauth.GetUserInfo(accessToken.AccessToken, accessToken.OpenID, "")
	if err != nil {
		return userinfo, err
	}
	return userinfo, nil
}
func PageKfChat(c *gin.Context) {
	kefuId := c.Query("kefu_id")
	visitorId := c.Query("visitor_id")
	token := c.Query("token")
	c.HTML(http.StatusOK, "chat_kf_page.html", gin.H{
		"KefuId":    kefuId,
		"VisitorId": visitorId,
		"Token":     token,
	})
}
