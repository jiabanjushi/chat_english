package tools

import (
	"fmt"
	"log"
	"testing"
)

func TestSendSmtp(t *testing.T) {
	for i := 0; i < 1; i++ {
		body := fmt.Sprintf("<h1>hello,body %d</h1>", i)
		subject := fmt.Sprintf("hello subject %d", i)
		err := SendSmtp("smtp.qq.com:465", "78084388@qq.com", "ysrttzaiwcmobjjd", []string{"630892807@qq.com"}, subject, body)
		log.Println(err)
	}
}
