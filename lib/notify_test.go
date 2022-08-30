package lib

import "testing"

func TestNotify_SendMail(t *testing.T) {
	notify := &Notify{
		Subject:     "测试主题",
		MainContent: "测试内容",
		EmailServer: NotifyEmail{
			Server:   "smtp.sina.cn",
			Port:     587,
			From:     "taoshihan1@sina.com",
			Password: "382e8a5e11cfae8c",
			To:       []string{"630892807@qq.com"},
			FromName: "GOFLY客服",
		},
	}
	ok, err := notify.SendMail()
	t.Log(ok, err)
}
