package log

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{})
}

func WithContext(ctx *gin.Context) *logrus.Entry {
	return logger.WithFields(logrus.Fields{
		"reqMethod": ctx.Request.Method,
		"reqUrl": ctx.Request.RequestURI,
	})
}

func Logger() *logrus.Logger {
	return logger
}