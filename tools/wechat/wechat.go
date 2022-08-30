package wechat

import (
	"encoding/json"
	"fmt"
	"go-fly-muti/tools"
	"log"
	"net/url"
)

type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   uint   `json:"expires_in"`
}
type Ticket struct {
	Ticket        string `json:"ticket"`
	ExpireSeconds uint   `json:"expire_seconds"`
	Url           string `json:"url"`
}

/**
获取access_token
*/
func GetAccessToken(appId, appSecret string) *AccessToken {
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		appId,
		appSecret)
	accessToken := tools.Get(url)
	token := &AccessToken{}
	json.Unmarshal([]byte(accessToken), token)
	return token
}

/**
获取临时二维码ticket
*/
func CreateQrImgUrl(accessToken, accountInfo string) {
	api := "https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=" + accessToken
	data := fmt.Sprintf(
		"{\"expire_seconds\": 604800, \"action_name\": \"QR_STR_SCENE\", \"action_info\": {\"scene\": {\"scene_str\": \"%s\"}}}",
		accountInfo)
	result, _ := tools.Post(api, "application/json;charset=utf-8", []byte(data))
	ticket := &Ticket{}
	json.Unmarshal([]byte(result), ticket)
	imgUrl := "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=" + url.QueryEscape(ticket.Ticket)
	log.Println(ticket.Ticket, imgUrl)
}
