package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"math"
	"os"
	"time"
)

var (
	Logger *logrus.Logger
)

func init(){
	Logger = InitLogger()
}

func InitLogger() *logrus.Logger {
	logName := "agent"

	Log := logrus.New()
	Log.Hooks.Add(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.DebugLevel: initWriter(logName,"debug"),
			logrus.InfoLevel:  initWriter(logName,"info"),
			logrus.WarnLevel:  initWriter(logName,"warn"),
			logrus.ErrorLevel: initWriter(logName,"error"),
			logrus.FatalLevel: initWriter(logName,"fatal"),
			logrus.PanicLevel: initWriter(logName,"panic"),
		},
		&logrus.JSONFormatter{
			PrettyPrint: true,
		},
	))
	return Log
}

func InitGinLogger() *logrus.Logger {
	//pathMap := lfshook.PathMap{
	//	logrus.InfoLevel:  "logs/gin-info.log",
	//	logrus.ErrorLevel: "logs/gin-error.log",
	//	logrus.WarnLevel:  "logs/gin-warn.log",
	//}
	//Log := logrus.New()
	//Log.Hooks.Add(lfshook.NewHook(
	//	pathMap,
	//	&logrus.JSONFormatter{},
	//))

	logName := "gin"

	Log := logrus.New()
	Log.Hooks.Add(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.DebugLevel: initWriter(logName,"debug"),
			logrus.InfoLevel:  initWriter(logName,"info"),
			logrus.WarnLevel:  initWriter(logName,"warn"),
			logrus.ErrorLevel: initWriter(logName,"error"),
			logrus.FatalLevel: initWriter(logName,"fatal"),
			logrus.PanicLevel: initWriter(logName,"panic"),
		},
		&logrus.JSONFormatter{
			PrettyPrint: true,
		},
	))
	return Log
}

func initWriter(logName string, logType string)*rotatelogs.RotateLogs{
	writer, _ := rotatelogs.New(
		"logs/" + logName + "-" + logType + ".%Y-%m-%d-%H.log",
		// WithLinkName为最新的日志建立软连接，以方便随着找到当前日志文件
		//rotatelogs.WithLinkName(logName),

		// WithRotationTime设置日志分割的时间，这里设置为一小时分割一次
		rotatelogs.WithRotationTime(time.Hour),

		// WithMaxAge和WithRotationCount二者只能设置一个，
		// WithMaxAge设置文件清理前的最长保存时间，
		// WithRotationCount设置文件清理前最多保存的个数。
		rotatelogs.WithMaxAge(time.Hour*24),
		//rotatelogs.WithRotationCount(10),
	)
	return writer
}

func GinLoggerHandler(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// other handler can change c.Path so:
		path := c.Request.URL.Path
		start := time.Now()
		c.Next()
		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		clientUserAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknow"
		}
		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}

		entry := logrus.NewEntry(log).WithFields(logrus.Fields{
			"hostname":   hostname,
			"statusCode": statusCode,
			"latency":    latency, // time to process
			"clientIP":   clientIP,
			"method":     c.Request.Method,
			"path":       path,
			"referer":    referer,
			"dataLength": dataLength,
			"userAgent":  clientUserAgent,
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			if path != "/agent/version" {
				msg := fmt.Sprintf("%s - %s [%s] \"%s %s\" %d %d \"%s\" \"%s\" (%dms)", clientIP, hostname, time.Now().Format("02/Jan/2006:15:04:05 -0700"), c.Request.Method, path, statusCode, dataLength, referer, clientUserAgent, latency)
				if statusCode > 499 {
					entry.Error(msg)
				} else if statusCode > 399 {
					entry.Warn(msg)
				} else {
					entry.Info(msg)
				}
			}
		}
	}
}
