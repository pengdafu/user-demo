package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"log"
	"net/http"
	"time"
	"user/internal/data"
)

type Result struct {
	Code int
	Msg  string
	Res  interface{}
}

type Params struct {
	Skip int `form:"skip" binding:"min=0"`
	Limit int `form:"limit" binding:"required,gte=1,lte=100"`
}

type userService struct {
	User data.UserDataI
}

func NewUserService(pool *redis.Pool) userService {
	return userService{
		User: data.NewUserData(pool),
	}
}

func (svc userService) GetUserList(ctx *gin.Context) {
	var params Params
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 400,
			Msg:  err.Error(),
		})
		return
	}
	res, err := svc.User.GetUserList(params.Skip, params.Limit)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  err.Error(),
		})
		return
	}
	users := make([]data.User, len(res))
	for i, v := range res {
		u := data.User{}
		v := string(v.([]uint8))
		_ = json.Unmarshal([]byte(v), &u)
		users[i] = u
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 200,
		Msg:  "success",
		Res:  users,
	})
}

func (svc userService) AddUser(ctx *gin.Context) {
	user := data.User{}
	if err := ctx.ShouldBindJSON(&user); err != nil {
		fmt.Println("----------------")
		ctx.JSON(http.StatusOK, Result{
			Code: 400,
			Msg:  err.Error(),
		})
		return
	}

	user.LoginTime = time.Now().Unix()
	if err := svc.User.AddUser(user); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "请求异常，添加用户失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 200,
		Msg:  "success",
		Res:  user,
	})
}

func (svc userService) GetUser(ctx *gin.Context) {
	userId := ctx.Param("userId")
	user, err := svc.User.GetUser(userId)

	if err != nil && errors.Is(err, data.UserNotExist) {
		ctx.JSON(http.StatusOK, Result{
			Code: 404,
			Msg:  "用户不存在",
		})
		return
	} else if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "获取用户信息异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 200,
		Msg:  "success",
		Res:  user,
	})
}

func (svc userService) DestroyUser(ctx *gin.Context) {
	userId := ctx.Param("userId")
	if err := svc.User.DestroyUser(userId); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "删除异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 200,
		Msg:  "success",
	})
}
