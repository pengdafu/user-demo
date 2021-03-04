package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"user/controller"
)

func SetV1Route(r *gin.Engine, pool *redis.Pool)  {
	userSvc := controller.NewUserService(pool)
	userGroup := r.Group("/v1/user")
	userGroup.GET("/", userSvc.GetUserList)
	userGroup.POST("/", userSvc.AddUser)
	userGroup.GET("/:userId", userSvc.GetUser)
	userGroup.DELETE("/:userId", userSvc.DestroyUser)
}
