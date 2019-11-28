package logger

import (
	"bytes"
	"encoding/json"
	"gin-web/app/config"
	"gin-web/app/utils/response"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

var accessChannel = make(chan string, 100)

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func handleAccessChannel() {
	if f, err := os.OpenFile(config.AppAccessLogName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666); err != nil {
		log.Warn(err)
	} else {
		for accessLog := range accessChannel {
			f.WriteString(accessLog + "\n")
		}
	}
	return
}

func Setup() gin.HandlerFunc {
	// go handleAccessChannel()
	src, err := os.OpenFile(config.AppAccessLogName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Warn("open log file error: ", err.Error())
	}

	logger := log.New()
	logger.Out = src
	logger.SetLevel(log.DebugLevel)
	logger.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logWriter, err := rotatelogs.New(config.AppAccessLogName+"%Y%m%d.log",
		rotatelogs.WithLinkName(config.AppAccessLogName),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		log.Warn(err)
	}

	writeMap := lfshook.WriterMap{
		log.InfoLevel:  logWriter,
		log.WarnLevel:  logWriter,
		log.DebugLevel: logWriter,
		log.ErrorLevel: logWriter,
		log.FatalLevel: logWriter,
		log.PanicLevel: logWriter,
	}

	lfHook := lfshook.NewHook(writeMap, &log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logger.AddHook(lfHook)

	return func(c *gin.Context) {
		var (
			responseCode int
			responseMsg  string
			responseData interface{}
		)
		bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyLogWriter
		startTime := time.Now()
		c.Next()
		responseBody := bodyLogWriter.body.String()
		if responseBody != "" {
			res := response.Response{}
			if err := json.Unmarshal([]byte(responseBody), &res); err == nil {
				responseCode = res.Code
				responseData = res.Data
				responseMsg = res.Message
			}
		}
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		if c.Request.Method == "POST" {
			c.Request.ParseForm()
		}

		// 日志格式
		logger.WithFields(log.Fields{
			"request_time":      startTime.Unix(),
			"request_method":    c.Request.Method,
			"request_uri":       c.Request.RequestURI,
			"request_proto":     c.Request.Proto,
			"request_ua":        c.Request.UserAgent(),
			"request_referer":   c.Request.Referer,
			"request_post_data": c.Request.PostForm.Encode(),
			"request_client_ip": c.ClientIP(),
			"response_time":     endTime.Unix(),
			"response_code":     responseCode,
			"responseMsg":       responseMsg,
			"response_data":     responseData,
			"cost_time":         latencyTime.Milliseconds(),
		}).Info()

	}
}
