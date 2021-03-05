package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type User struct {
	UserId    string `json:"userId" binding:"required,min=1"`
	NickName  string `json:"nickName" binding:"required,min=1"`
	Role      uint32 `json:"role" binding:"required,oneof=1 2"` //1 anchor, 2 audience
	LoginTime int64  `json:"loginTime"`
}

type UserListRes struct {
	Rows []interface{} `json:"rows"`
	Count int64 `json:"count"`
}

const (
	BaseIdIncKey = "u:inc:id:"
	BaseIdKey    = "u:id:"
	AllIdsZSet   = "u:ids"
)

// LuaScripts 提供一个 getUserList 的原子操作，由于需要
// 先获取ids再去查询值，所以不能用MULTI处理
const LuaScripts   = `
local rcall = redis.call
local key   = KEYS[1]
local start = ARGV[1]
local _end  = ARGV[2]

local userIds = rcall("ZREVRANGE", key, start, _end)
if #userIds == 0 then 
	return nil
end
local users = rcall("MGET", unpack(userIds))
return { users, rcall("ZCARD", key) }
`

var UserNotExist = errors.New("用户不存在")

type UserData struct {
	redisPool *redis.Pool
}

func NewUserData(pool *redis.Pool) UserDataI {
	return UserData{
		redisPool: pool,
	}
}

func (u UserData) GetUserList(skip, limit int) (*UserListRes, error) {
	r := u.redisPool.Get()
	defer r.Close()
	scripts := redis.NewScript(1, LuaScripts)
	reply, err := scripts.Do(r, AllIdsZSet, skip, skip + limit - 1)

	if err != nil {
		return nil, err
	}
	if reply == nil {
		return &UserListRes{
			Rows: []interface{}{},
			Count: 0,
		}, nil
	}
	_res := reply.([]interface{})
	res := &UserListRes{
		Rows: _res[0].([]interface{}),
		Count: _res[1].(int64),
	}
	return res, nil
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
	GetUserList(skip, limit int) (*UserListRes, error)
	AddUser(user User) error
	GetUser(userId string) (*User, error) // 不存在返回error
	DestroyUser(userId string) error
	RandomUserId() (string, error)
}
