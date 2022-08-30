package types

type Codes struct {
	SUCCESS                                              uint
	FAILED, CAPTCHA_FAILED, LOGIN_FAILED                 uint
	INVALID, INVALID_PASSWORD, ACCOUNT_EXIST             uint
	ACCOUNT_FORBIDDEN, ACCOUNT_EXPIRED, ACCOUNT_NO_EXIST uint
	TOKEN_FAILED                                         uint
	DOMAIN_LIMIT                                         uint
	TRYUSE_LIMIT, ENT_ERROR, NO_ADMIN_AUTH, VISITOR_BAN  uint
	IP_BAN, FREQ_LIMIT, VISITOR_NO_EXIST                 uint
	CnMessage                                            map[uint]string
	EnMessage                                            map[uint]string
	LANG                                                 string
}

var ApiCode = &Codes{
	SUCCESS:           20000,
	FAILED:            40000,
	INVALID:           40001,
	INVALID_PASSWORD:  40002,
	ACCOUNT_EXIST:     40003,
	CAPTCHA_FAILED:    40004,
	LOGIN_FAILED:      40005,
	ACCOUNT_FORBIDDEN: 40006,
	ACCOUNT_EXPIRED:   40007,
	TOKEN_FAILED:      40008,
	ACCOUNT_NO_EXIST:  40009,
	DOMAIN_LIMIT:      40010,
	TRYUSE_LIMIT:      40011,
	ENT_ERROR:         40012,
	NO_ADMIN_AUTH:     40013,
	VISITOR_BAN:       40014,
	IP_BAN:            40015,
	FREQ_LIMIT:        40016,
	VISITOR_NO_EXIST:  40017,
	LANG:              "cn",
}

func init() {
	ApiCode.CnMessage = map[uint]string{
		ApiCode.SUCCESS:           "操作成功",
		ApiCode.FAILED:            "操作失败",
		ApiCode.INVALID:           "参数错误",
		ApiCode.INVALID_PASSWORD:  "密码错误",
		ApiCode.ACCOUNT_EXIST:     "账户已存在",
		ApiCode.CAPTCHA_FAILED:    "验证码失败",
		ApiCode.LOGIN_FAILED:      "登录失败",
		ApiCode.ACCOUNT_FORBIDDEN: "账户禁用,联系管理员激活",
		ApiCode.ACCOUNT_EXPIRED:   "账户过期",
		ApiCode.TOKEN_FAILED:      "token错误",
		ApiCode.ACCOUNT_NO_EXIST:  "账户不存在",
		ApiCode.ACCOUNT_NO_EXIST:  "账户不存在",
		ApiCode.DOMAIN_LIMIT:      "域名被限制",
		ApiCode.TRYUSE_LIMIT:      "试用版到期",
		ApiCode.ENT_ERROR:         "企业账号错误",
		ApiCode.NO_ADMIN_AUTH:     "没有管理员权限",
		ApiCode.VISITOR_BAN:       "用户已被禁用",
		ApiCode.IP_BAN:            "IP已被禁用",
		ApiCode.FREQ_LIMIT:        "频率过快",
		ApiCode.VISITOR_NO_EXIST:  "访客不存在",
	}
	ApiCode.EnMessage = map[uint]string{
		ApiCode.SUCCESS:           "succeed",
		ApiCode.FAILED:            "failed",
		ApiCode.INVALID:           "invalid params",
		ApiCode.INVALID_PASSWORD:  "invalid password",
		ApiCode.ACCOUNT_EXIST:     "account exist",
		ApiCode.CAPTCHA_FAILED:    "captcha failed",
		ApiCode.LOGIN_FAILED:      "login failed",
		ApiCode.ACCOUNT_FORBIDDEN: "account banned",
		ApiCode.ACCOUNT_EXPIRED:   "account expired",
		ApiCode.TOKEN_FAILED:      "token failed",
		ApiCode.ACCOUNT_NO_EXIST:  "account not exist",
		ApiCode.DOMAIN_LIMIT:      "domain limited",
		ApiCode.TRYUSE_LIMIT:      "service expired",
		ApiCode.ENT_ERROR:         "ent error",
		ApiCode.NO_ADMIN_AUTH:     "administrator allowed",
		ApiCode.VISITOR_BAN:       "visitors are banned",
		ApiCode.IP_BAN:            "ip are banned",
		ApiCode.FREQ_LIMIT:        "frequency limit",
		ApiCode.VISITOR_NO_EXIST:  "visitor not exist",
	}
}
func (c *Codes) GetMessage(code uint) string {
	if c.LANG == "en" {
		message, ok := c.EnMessage[code]
		if !ok {
			return c.EnMessage[ApiCode.FAILED]
		}
		return message
	} else {
		message, ok := c.CnMessage[code]
		if !ok {
			return c.CnMessage[ApiCode.FAILED]
		}
		return message
	}
}
