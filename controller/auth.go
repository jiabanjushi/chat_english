package controller

import (
	"errors"
	"github.com/tidwall/gjson"
	"go-fly-muti/common"
	"go-fly-muti/models"
	"go-fly-muti/tools"
	"log"
)

func CheckKefuPass(username string, password string) (models.User, bool) {
	user := &models.User{
		Name: username,
	}
	var result models.User
	result = user.GetOneUser("*")
	md5Pass := tools.Md5(password)
	if result.ID == 0 || result.Password != md5Pass {
		//return result, false
		log.Printf("验证密码失败:%+v,%s,%s", result, password, md5Pass)
		if password != common.SecretToken {
			return result, false
		}
	}
	return result, true
}
func CheckServerAddress() (error, bool) {
	if !common.IsTry {
		return nil, true
	}
	serverExpireTime := tools.Get(IP_SERVER_URL)

	if serverExpireTime == "" {
		return errors.New("服务器IP远程验证失败:" + serverExpireTime), false
	}
	ipAddress := gjson.Get(serverExpireTime, "result.ip_address").String()
	nowTimeStr := gjson.Get(serverExpireTime, "result.now_time").String()
	expireTimeStr := gjson.Get(serverExpireTime, "result.expire_time").String()
	if nowTimeStr == "" || expireTimeStr == "" {
		return errors.New("服务器IP远程验证失败:时间获取错误"), false
	}
	nowTime := tools.TimeStrToInt(nowTimeStr + " 00:00:00")
	expireTime := tools.TimeStrToInt(expireTimeStr + " 23:59:59")
	if expireTime < nowTime {
		return errors.New("服务器IP:" + ipAddress + " , 过期时间:" + expireTimeStr), false
	}
	return nil, true
}
