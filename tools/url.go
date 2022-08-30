package tools

import "net/url"

func UrlEncode(str string) string {
	return url.QueryEscape(str)
}
func UrlDecode(str string) string {
	res, err := url.QueryUnescape(str)
	if err != nil {
		return ""
	}
	return res
}
