package tools

import (
	"log"
	"testing"
)

func TestPostHeader(t *testing.T) {
	url := "https://jd.sopans.com/check"
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	res, err := PostHeader(url, []byte("username=kefu2&password=1234526"), headers)
	log.Println(res, err)
}
func TestIsMobile(t *testing.T) {
	IsMobile("aaaa")
}
