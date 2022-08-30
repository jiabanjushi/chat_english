package lib

import (
	"encoding/base64"
	"fmt"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/go-playground/validator/v10"
	"strings"
)

/**
发送通知类
*/
type Notify struct {
	ReceiveName, Subject, MainContent string
	EmailServer                       NotifyEmail
}
type NotifyEmail struct {
	Server, From, Password string `validate:"required"`
	FromName               string
	To                     []string
	Port                   uint `validate:"required"`
}

/*
发送邮件
*/
func (this *Notify) SendMail() (bool, error) {
	validate := validator.New()

	emailServer := this.EmailServer
	err := validate.Struct(emailServer)
	if err != nil {
		return false, err
	}
	err = SendSmtp(fmt.Sprintf("%s:%d", this.EmailServer.Server,
		this.EmailServer.Port), this.EmailServer.From,
		this.EmailServer.Password,
		this.EmailServer.To,
		this.Subject, this.MainContent, this.EmailServer.FromName)
	if err != nil {
		return false, err
	}
	return true, nil
}

func SendSmtp(server string, from string, password string, to []string, subject string, body string, fromName string) error {
	auth := sasl.NewPlainClient("", from, password)
	subjectBase := base64.StdEncoding.EncodeToString([]byte(subject))
	fromHeader := fmt.Sprintf("From:=?UTF-8?B?%s?= <%s>\r\n", base64.StdEncoding.EncodeToString([]byte(fromName)), from)
	msg := strings.NewReader(
		fromHeader +
			"To: " + strings.Join(to, ",") + "\r\n" +
			"Subject: =?UTF-8?B?" + subjectBase + "?=\r\n" +
			"Content-Type: text/html; charset=UTF-8" +
			"\r\n\r\n" +
			body + "\r\n")
	err := smtp.SendMail(server, auth, from, to, msg)
	if err != nil {
		return err
	}
	return nil
}
