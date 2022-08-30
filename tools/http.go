package tools

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func Get(url string) string {

	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return ""
	}
	var reader io.ReadCloser
	if res.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(res.Body)
		if err != nil {
			return ""
		}
	} else {
		reader = res.Body
	}
	//utf8Reader := transform.NewReader(reader,
	//	simplifiedchinese.GBK.NewDecoder())
	robots, err := ioutil.ReadAll(reader)
	res.Body.Close()
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(robots)
}

//Post("http://xxxx","application/json;charset=utf-8",[]byte("{'aaa':'bbb'}"))
func Post(url string, contentType string, body []byte) (string, error) {
	res, err := http.Post(url, contentType, strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
func PostHeader(url string, msg []byte, headers map[string]string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(msg)))
	if err != nil {
		return "", err
	}
	for key, header := range headers {
		req.Header.Set(key, header)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func IsMobile(userAgent string) bool {
	mobileRe, _ := regexp.Compile("(?i:Mobile|iPod|iPhone|Android|Opera Mini|BlackBerry|webOS|UCWEB|Blazer|PSP)")
	str := mobileRe.FindString(userAgent)
	if str != "" {
		return true
	}
	return false
}
