package tools

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/go-gomail/gomail"
	"strconv"
	"strings"
)

func SendSmtp2(server string, from string, password string, to []string, subject string, body string) error {
	auth := sasl.NewPlainClient("", from, password)
	subjectBase := base64.StdEncoding.EncodeToString([]byte(subject))
	msg := strings.NewReader(
		"From: " + from + "\r\n" +
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
func SendSmtp(server string, from string, password string, to []string, subject string, body string) error {
	servers := strings.Split(server, ":")
	if len(servers) < 2 {
		return errors.New("server addr error")
	}
	host := servers[0]
	port, _ := strconv.ParseInt(servers[1], 10, 32)
	m := gomail.NewMessage()
	m.SetAddressHeader("From", from, "GOFLY在线客服通知")
	m.SetAddressHeader("To", to[0], to[0])
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(host, int(port), from, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
