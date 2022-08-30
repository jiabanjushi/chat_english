package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-fly-muti/tools"
)

func SendEntEmail(c *gin.Context) {
	ent_id := c.PostForm("ent_id")
	email := c.PostForm("email")
	weixin := c.PostForm("weixin")
	msg := c.PostForm("msg")
	name := c.PostForm("name")
	if msg == "" || name == "" || ent_id == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "内容/姓名不能为空",
		})
		return
	}
	content := fmt.Sprintf("[%s] [%s] [%s] [%s]", name, weixin, email, msg)
	go SendEntSmtpEmail("[留言]"+name, content, ent_id)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}

//翻译
func GetTranslate(c *gin.Context) {
	words := c.PostForm("words")
	ret := tools.Get("http://translate.google.cn/translate_a/single?client=gtx&dt=t&dj=1&ie=UTF-8&sl=auto&tl=en&q=" + words)
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": ret,
	})
}
