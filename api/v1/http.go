package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"user/internal/service"
)

func NewRoute(pool *redis.Pool) *gin.Engine {
	r := gin.Default()
	{
		userSvc := service.NewUserService(pool)
		userGroup := r.Group("/v1/user")
		userGroup.GET("/", userSvc.GetUserList)
		userGroup.POST("/", userSvc.AddUser)
		userGroup.GET("/:userId", userSvc.GetUser)
		userGroup.DELETE("/:userId", userSvc.DestroyUser)
	}
	return r
}
