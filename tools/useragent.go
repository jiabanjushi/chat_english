package tools

import "strings"

type UserAgentParser struct {
	UserAgent string
	OsVersion string
	Browser   string
}

func NewUaParser(userAgent string) *UserAgentParser {
	obj := &UserAgentParser{
		UserAgent: userAgent,
	}
	obj.OsVersion = obj.GetOsVersion()
	obj.Browser = obj.GetBrowser()
	return obj
}
func (this *UserAgentParser) GetOsVersion() string {
	osVersion := "unknow"
	if strings.Contains(this.UserAgent, "NT 10.0") {
		osVersion = "Windows 10"
	}
	if strings.Contains(this.UserAgent, "NT 6.2") {
		osVersion = "Windows 8"
	}
	if strings.Contains(this.UserAgent, "NT 6.1") {
		osVersion = "Windows 7"
	}
	if strings.Contains(this.UserAgent, "NT 6.0") {
		osVersion = "Windows Vista"
	}
	if strings.Contains(this.UserAgent, "NT 5.2") {
		osVersion = "Windows Server 2003"
	}
	if strings.Contains(this.UserAgent, "NT 5.1") {
		osVersion = "Windows XP"
	}
	if strings.Contains(this.UserAgent, "Mac") {
		osVersion = "Mac"
		if strings.Contains(this.UserAgent, "iPhone") {
			osVersion = "iPhone"
		}
	}
	if strings.Contains(this.UserAgent, "Linux") {
		osVersion = "Linux"
	}
	if strings.Contains(this.UserAgent, "Android") {
		osVersion = "Android"
	}
	return osVersion
}
func (this *UserAgentParser) GetBrowser() string {
	browserVersion := ""
	if strings.Contains(this.UserAgent, "Chrome") {
		browserVersion += "Chrome "
	}
	if strings.Contains(this.UserAgent, "UCBrowser") {
		browserVersion += "UCBrowser "
	}
	if strings.Contains(this.UserAgent, "Safari") {
		browserVersion += "Safari "
	}
	if strings.Contains(this.UserAgent, "QQ") {
		browserVersion += "QQBrowser "
	}
	if strings.Contains(this.UserAgent, "MicroMessenger") {
		browserVersion += "MicroMessenger "
	}
	if strings.Contains(this.UserAgent, "UCBrowser") {
		browserVersion += "UCBrowser "
	}
	if strings.Contains(this.UserAgent, "baiduboxapp") {
		browserVersion += "baiduboxapp "
	}
	if strings.Contains(this.UserAgent, "baidubrowser") {
		browserVersion += "baidubrowser "
	}
	if strings.Contains(this.UserAgent, "Weibo") {
		browserVersion += "Weibo "
	}
	if strings.Contains(this.UserAgent, "DingTalk") {
		browserVersion += "DingTalk "
	}
	if strings.Contains(this.UserAgent, "YaBrowser") {
		browserVersion += "YaBrowser "
	}
	if strings.Contains(this.UserAgent, "Baiduspider") {
		browserVersion += "Baiduspider "
	}
	if strings.Contains(this.UserAgent, "Firefox") {
		browserVersion += "Firefox "
	}
	if strings.Contains(this.UserAgent, "AhrefsBot") {
		browserVersion += "AhrefsBot Spider "
	}
	return browserVersion
}
