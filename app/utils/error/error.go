package error

import (
	"encoding/json"
	"fmt"
	"gin-web/app/config"
	"gin-web/app/router/middleware/exception"
	"gin-web/app/utils/mail"
	"os"
	"runtime/debug"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}
func New(text string) error {
	alarm("INFO", text)
	return &errorString{text}
}
func Sms(text string) error {
	alarm("SMS", text)
	return &errorString{text}
}
func Email(text string) error {
	alarm("EMAIL", text)
	return &errorString{text}
}
func WeChat(text string) error {
	alarm("WX", text)
	return &errorString{text}
}
func alarm(level, text string) {
	if level == "INFO" { // 记录日志
		if f, err := os.OpenFile(config.AppErrorLogName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err != nil {
			log.Warn("open error log error: ", err)
		} else {
			errLogMap := make(map[string]interface{})
			errLogMap["time"] = time.Now().Format("2006-01-02 15:04:05")
			errLogMap["info"] = text
			errLogJson, err := json.Marshal(errLogMap)
			if err != nil {
				log.Warn(err)
			}
			f.WriteString(string(errLogJson) + "\n")
		}

	} else if level == "SMS" { // 发短信

	} else if level == "EMAIL" { // 发邮件
		DebugStack := ""
		for _, v := range strings.Split(string(debug.Stack()), "\n") {
			DebugStack += v + "<br>"
		}
		subject := fmt.Sprintf("【系统告警】%s 项目出错了!", config.AppName)
		body := strings.ReplaceAll(exception.MailTemplate, "{ErrorMsg}", fmt.Sprintf("%s", text))
		body = strings.ReplaceAll(body, "{RequestTime}", time.Now().Format("2006-01-02 15:04:05"))
		body = strings.ReplaceAll(body, "{RequestURL}", "---")
		body = strings.ReplaceAll(body, "{RequestUA}", "---")
		body = strings.ReplaceAll(body, "{RequestIP}", "---")
		body = strings.ReplaceAll(body, "{DebugStack}", DebugStack)
		options := &mail.Options{
			MailPort: config.SystemEmailPort,
			MailHost: config.SystemEmailHost,
			MailUser: config.SystemEmailUser,
			MailPass: config.SystemEmailPass,
			MailTo:   config.ErrorNotifyUser,
			Body:     body,
			Subject:  subject,
		}
		if err := mail.Send(options); err != nil {
			log.Warn("Error Send Mail Error", err)
		}
	} else if level == "WX" { // 发微信

	} else {
		log.Warnf("invalid level type %v,", level, "must INFO/SMS/EMAIL/WX")
	}
}
