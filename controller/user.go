package controller

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"net/http"
	"time"
	"user/model"
	"user/pkg/log"
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
	User model.UserDataI
}

func NewUserService(pool *redis.Pool) userService {
	return userService{
		User: model.NewUserData(pool),
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
		log.WithContext(ctx).Error("获取用户列表失败", err)
		ctx.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "获取用户列表失败",
		})
		return
	}
	users := make([]model.User, len(res))
	for i, v := range res {
		u := model.User{}
		v := string(v.([]uint8))
		_ = json.Unmarshal([]byte(v), &u)
		users[i] = u
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "success",
		Res:  users,
	})
}

func (svc userService) AddUser(ctx *gin.Context) {
	user := model.User{}
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 400,
			Msg:  err.Error(),
		})
		return
	}

	user.LoginTime = time.Now().Unix()
	if err := svc.User.AddUser(user); err != nil {
		log.WithContext(ctx).Error("添加用户失败", err)
		ctx.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "请求异常，添加用户失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "success",
		Res:  user,
	})
}

func (svc userService) GetUser(ctx *gin.Context) {
	userId := ctx.Param("userId")
	user, err := svc.User.GetUser(userId)

	if err != nil && errors.Is(err, model.UserNotExist) {
		ctx.JSON(http.StatusOK, Result{
			Code: 404,
			Msg:  "用户不存在",
		})
		return
	} else if err != nil {
		log.WithContext(ctx).Error("获取用户信息异常", err)
		ctx.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "获取用户信息异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "success",
		Res:  user,
	})
}

func (svc userService) DestroyUser(ctx *gin.Context) {
	userId := ctx.Param("userId")
	if err := svc.User.DestroyUser(userId); err != nil {
		log.WithContext(ctx).Error("删除用户失败", err)
		ctx.JSON(http.StatusOK, Result{
			Code: 500,
			Msg:  "删除异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "success",
	})
}
