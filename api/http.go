package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
	"user/pkg/log"
)

func NewGinEngine() *gin.Engine {
	e := gin.New()
	gin.Recovery()
	e.Use(loggerM(), gin.Recovery())
	return e
}

func loggerM() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger := log.WithContext(ctx)
		startTime := time.Now()
		ctx.Next()
		endTime := time.Now()
		logger.WithFields(logrus.Fields{
			"latencyTime": endTime.Sub(startTime),
			"statusCode": ctx.Writer.Status(),
			"clientIP": ctx.ClientIP(),
		}).Info()
	}
}