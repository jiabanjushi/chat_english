package wechat

import "testing"

func TestGetAccessToken(t *testing.T) {
	GetAccessToken("wx2987986f14daa9fc", "0d5a220a35cb7edcf1943459258fea2c")
}
func TestCreateQrTicket(t *testing.T) {
	token := GetAccessToken("wx2987986f14daa9fc", "0d5a220a35cb7edcf1943459258fea2c")
	CreateQrImgUrl(token.AccessToken, "agent_20291289&2378278")
}
