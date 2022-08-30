package lib

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"go-fly-muti/models"
	"go-fly-muti/tools"
	"sort"
)

type Wechat struct {
	AppId, AppSecret, Token, WechatVisitorTemplateId,
	WechatMessageTemplateId, WechatKefuTemplateId,
	WechatKefu, WechatHost string
}

//AppId,AppSecret,Token
func NewWechatLib(entId string) (*Wechat, error) {
	configs := models.FindEntConfigs(entId)
	AppID := ""
	AppSecret := ""
	Token := ""
	messageTemplateId := ""
	visitorTemplateId := ""
	kefuTemplateId := ""
	host := ""
	wechatKefu := ""
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
		if config.ConfKey == "WechatVisitorTemplateId" {
			visitorTemplateId = config.ConfValue
		}
		if config.ConfKey == "WechatKefuTemplateId" {
			kefuTemplateId = config.ConfValue
		}
		if config.ConfKey == "WechatMessageTemplateId" {
			messageTemplateId = config.ConfValue
		}
		if config.ConfKey == "WechatHost" {
			host = config.ConfValue
		}
		if config.ConfKey == "WechatKefu" {
			wechatKefu = config.ConfValue
		}
	}
	if AppID == "" || AppSecret == "" || Token == "" {
		return nil, errors.New("AppID,AppSecret,Token获取失败!")
	}
	wechat := &Wechat{
		AppId:                   AppID,
		AppSecret:               AppSecret,
		Token:                   Token,
		WechatKefuTemplateId:    kefuTemplateId,
		WechatVisitorTemplateId: visitorTemplateId,
		WechatMessageTemplateId: messageTemplateId,
		WechatHost:              host,
		WechatKefu:              wechatKefu,
	}
	return wechat, nil
}
func (this *Wechat) GetAccess() (string, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", this.AppId, this.AppSecret)
	tools.Get(url)
	return "", errors.New("获取access_token失败")
}
func (this *Wechat) CheckWechatSign(signature, timestamp, nonce, echostr string) (string, error) {
	token := this.Token
	//将token、timestamp、nonce三个参数进行字典序排序
	var tempArray = []string{token, timestamp, nonce}
	sort.Strings(tempArray)
	//将三个参数字符串拼接成一个字符串进行sha1加密
	var sha1String string = ""
	for _, v := range tempArray {
		sha1String += v
	}
	h := sha1.New()
	h.Write([]byte(sha1String))
	sha1String = hex.EncodeToString(h.Sum([]byte("")))
	//获得加密后的字符串可与signature对比
	if sha1String == signature {
		return echostr, nil
	}
	return "", errors.New("微信API验证失")
}
