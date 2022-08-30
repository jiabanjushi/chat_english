package lib

import (
	"testing"
)

func TestGetui(t *testing.T) {
	getui := &Getui{
		AppId:           "fndEwCoJbZ9u9iC2zbmCDA",
		AppKey:          "WnoWdaKPhC9gs9PzI7e0L3",
		AppSecret:       "UtqBkuRJBv5s5dzDEiqQM6",
		AppMasterSecret: "4KtIXIPN507GdMPwtiXkN4",
	}
	token, err := getui.GetGetuiToken()
	t.Logf("%+v,%+v", token, err)
}
func TestGetui_PushSingle(t *testing.T) {
	getui := &Getui{
		AppId:           "fndEwCoJbZ9u9iC2zbmCDA",
		AppKey:          "WnoWdaKPhC9gs9PzI7e0L3",
		AppSecret:       "UtqBkuRJBv5s5dzDEiqQM6",
		AppMasterSecret: "4KtIXIPN507GdMPwtiXkN4",
	}
	token := "614e57028322896f090b32a5d0e43db59d96c41b300a40675cadaba9269fa8ca,"
	clientId := "27e0b1d016432a65504ac234a5cccad0"
	title := "你好"
	content := "你好啊"
	res, err := getui.PushSingle(token, clientId, title, content)
	t.Logf("%+v,%+v", res, err)

	if res == 10001 {
		token, _ := getui.GetGetuiToken()
		res, err := getui.PushSingle(token, clientId, title, content)
		t.Logf("%+v,%+v", res, err)
	}
}
