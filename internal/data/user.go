package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type User struct {
	UserId    string `json:"userId"`
	NickName  string `json:"nickName" binding:"required,min=1"`
	Role      uint32 `json:"role" binding:"required,oneof=1 2"`//1 anchor, 2 audience
	LoginTime int64 `json:"loginTime"`
}

const (
	BaseIdIncKey = "u:inc:id:"
	BaseIdKey    = "u:id:"
	AllIdsZSet   = "u:ids"
)

var UserNotExist = errors.New("用户不存在")

type UserData struct {
	redisPool *redis.Pool
}

func NewUserData(pool *redis.Pool) UserDataI {
	return UserData{
		redisPool: pool,
	}
}

func (u UserData) GetUserList(skip, limit int) ([]interface{}, error) {
	r := u.redisPool.Get()
	defer r.Close()
	reply, err := r.Do("ZREVRANGE", AllIdsZSet, skip, skip + limit-1)
	if err != nil {
		return nil, err
	}
	ids := reply.([]interface{})
	if len(ids) == 0 {
		return []interface{}{}, nil
	}
	reply, err = r.Do("MGET", ids...)
	if err != nil {
		return nil, err
	}

	return reply.([]interface{}), nil
}

func (u UserData) AddUser(user User) error {
	r := u.redisPool.Get()
	defer r.Close()
	bytes, err := json.Marshal(user)
	if err != nil {
		return err
	}
	r.Send("MULTI")
	r.Send("SET", BaseIdKey+user.UserId, string(bytes))
	r.Send("ZADD", AllIdsZSet, user.LoginTime, BaseIdKey+user.UserId)
	_, err = r.Do("EXEC")
	if err != nil {
		return err
	}
	return nil
}

func (u UserData) GetUser(userId string) (*User, error) {
	r := u.redisPool.Get()
	defer r.Close()
	reply, err := r.Do("GET", BaseIdKey+userId)
	if err != nil {
		return nil, err
	}
	if reply == nil {
		return nil, UserNotExist
	}
	userStr := string(reply.([]uint8))
	user := &User{}
	err = json.Unmarshal([]byte(userStr), user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u UserData) DestroyUser(userId string) error {
	r := u.redisPool.Get()
	defer r.Close()

	r.Send("MULTI")
	r.Send("DEL", BaseIdKey+userId)
	r.Send("ZREM", AllIdsZSet, BaseIdKey+userId)
	_, err := r.Do("EXEC")
	return err
}

func (u UserData) RandomUserId() (string, error) {
	r := u.redisPool.Get()
	defer r.Close()
	reply, err := r.Do("INCR", BaseIdIncKey)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", reply), nil
}

type UserDataI interface {
	GetUserList(skip, limit int) ([]interface{}, error)
	AddUser(user User) error
	GetUser(userId string) (*User, error) // 不存在返回error
	DestroyUser(userId string) error
	RandomUserId() (string, error)
}
