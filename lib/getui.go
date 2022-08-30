package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-fly-muti/tools"
	"log"
	"strconv"
	"time"
)

type Getui struct {
	AppId, AppKey, AppSecret, AppMasterSecret string
}
type GetuiResponse struct {
	Code float64                `json:"code"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}
type GetuiReq struct {
	Sign      string `json:"sign"`
	Timestamp string `json:"timestamp"`
	Appkey    string `json:"appkey"`
}

func (this *Getui) GetGetuiToken() (string, error) {
	appid := this.AppId
	appkey := this.AppKey
	appmastersecret := this.AppMasterSecret
	if appid == "" || appkey == "" {
		return "", errors.New("appid appkey failed")
	}
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	reqJson := GetuiReq{
		Sign:      tools.Sha256(appkey + timestamp + appmastersecret),
		Timestamp: timestamp,
		Appkey:    appkey,
	}
	reqStr, _ := json.Marshal(reqJson)
	url := "https://restapi.getui.com/v2/" + appid + "/auth"
	res, err := tools.Post(url, "application/json;charset=utf-8", reqStr)
	log.Println(url, string(reqStr), err, res)
	if err != nil || res == "" {
		return "", err
	}
	var pushRes GetuiResponse
	json.Unmarshal([]byte(res), &pushRes)
	if pushRes.Code != 0 {
		return "", errors.New(pushRes.Msg)
	}
	token, ok := pushRes.Data["token"]
	if !ok {
		return "", errors.New("token not exist")
	}
	return token.(string), nil
}
func (this *Getui) PushSingle(token, clientId, title, content string) (int, error) {
	appid := this.AppId
	if appid == "" {
		return 400, errors.New("appid failed")
	}
	url := "https://restapi.getui.com/v2/" + appid + "/push/single/cid"
	format := `
{
    "request_id":"%s",
    "settings":{
        "ttl":3600000
    },
    "audience":{
        "cid":[
            "%s"
        ]
    },
    "push_message":{
        "notification":{
            "title":"%s",
            "body":"%s",
            "click_type":"startapp"
        }
    }
}
`
	req := fmt.Sprintf(format, tools.Md5(tools.Uuid()), clientId, title, content)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json;charset=utf-8"
	headers["token"] = token
	res, err := tools.PostHeader(url, []byte(req), headers)
	log.Println(url, string(req), err, res)
	if err != nil && res == "" {
		return 400, err
	}
	var pushRes GetuiResponse
	json.Unmarshal([]byte(res), &pushRes)

	if pushRes.Code != 0 {
		return int(pushRes.Code), errors.New(pushRes.Msg)
	}
	return int(pushRes.Code), nil
}
