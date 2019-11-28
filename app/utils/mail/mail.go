package mail

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type Options struct {
	MailPort int
	MailHost string
	MailUser string //发件人
	MailPass string // 发件人密码
	MailTo   string // 收件人 多个用分隔符 ,
	Subject  string // 邮件主题
	Body     string // 邮件正文
}

func Send(o *Options) error {
	m := gomail.NewMessage()
	m.SetHeader("From", o.MailUser)
	mailAddrs := strings.Split(o.MailTo, ",")
	m.SetHeader("To", mailAddrs...)
	m.SetHeader("Subject", o.Subject)
	m.SetHeader("text/html", o.Body)

	d := gomail.NewDialer(o.MailHost, o.MailPort, o.MailUser, o.MailPass)
	if err := d.DialAndSend(m); err != nil {
		log.Warn("Send mail error: ", err)
		return err
	}
	return nil
}
