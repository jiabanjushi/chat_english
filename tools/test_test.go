package tools

import (
	"log"
	"net/url"
	"testing"
)

func TestUrlEncode(t *testing.T) {
	str := "aabb,&8?%s"
	ret := url.QueryEscape(str)
	log.Println(ret)

	encodedValue := "Hell%C3%B6+W%C3%B6rld%40Golang"
	decodedValue, err := url.QueryUnescape(encodedValue)
	log.Println(decodedValue, err)
}
