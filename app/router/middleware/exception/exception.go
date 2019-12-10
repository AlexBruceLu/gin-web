package exception

import (
	"fmt"
	"gin-web/app/config"
	"gin-web/app/utils/mail"
	"gin-web/app/utils/response"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func SetUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				DebugStack := ""
				for _, v := range strings.Split(string(debug.Stack()), "\n") {
					DebugStack += v + "<br>"
				}
				subject := fmt.Sprintf("【重要错误】%s 项目出错了!", config.AppName)

				body := strings.ReplaceAll(MailTemplate, "{ErrorMsg}", fmt.Sprintf("%s", err))
				body = strings.ReplaceAll(body, "{RequestTime}", time.Now().Format("2006-01-02 15:04:05"))
				body = strings.ReplaceAll(body, "{RequestURL}", c.Request.Method+" "+c.Request.Host+c.Request.RequestURI)
				body = strings.ReplaceAll(body, "{RequestUA}", c.Request.UserAgent())
				body = strings.ReplaceAll(body, "{RequestIP}", c.ClientIP())
				body = strings.ReplaceAll(body, "{DebugStack}", DebugStack)

				options := &mail.Options{
					MailHost: config.SystemEmailHost,
					MailPort: config.SystemEmailPort,
					MailUser: config.SystemEmailUser,
					MailPass: config.SystemEmailPass,
					MailTo:   config.ErrorNotifyUser,
					Subject:  subject,
					Body:     body,
				}
				if err := mail.Send(options); err != nil {
					log.Warn("Send mail error: ", err)
				}
				utilGin := response.Gin{Ctx: c}
				utilGin.Response(500, "系统异常,请联系管理员!", nil)
			}
		}()
		c.Next()
	}
}
