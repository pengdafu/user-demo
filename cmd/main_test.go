package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"user/api"
	v1 "user/api/v1"
	"user/controller"
	pkgRedis "user/pkg/redis"
)

var redisPool *redis.Pool
var route *gin.Engine

func TestMain(m *testing.M) {
	fmt.Println("before testing")
	redisAddr, ok := os.LookupEnv("GO_TEST_REDIS_URL")
	if !ok {
		panic("请添加环境变量：export GO_TEST_REDIS_URL=xxxxx:6379")
	}
	redisPool = pkgRedis.NewRedisPool(context.Background(), redisAddr)
	route = api.NewGinEngine()
	v1.SetV1Route(route, redisPool)
	code := m.Run()
	_ = redisPool.Close()
	fmt.Println("after testing")
	os.Exit(code)
}

func TestAddUser(t *testing.T) {
	w := httptest.NewRecorder()
	data := url.Values{}
	fmt.Println(data.Encode())
	req := httptest.NewRequest("POST", "/v1/user/", strings.NewReader(`{"nickName":"test user","role":1,"userId":"testAddUserId"}`))
	req.Header.Set("Content-Type", "application/json")

	route.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	res := controller.Result{}
	_ = json.Unmarshal([]byte(w.Body.String()), &res)
	assert.Equal(t, 0, res.Code)
	assert.Equal(t, "success", res.Msg)
	m := res.Res.(map[string]interface{})
	userId := m["userId"].(string)
	assert.Equal(t, "testAddUserId", userId)

	cases := []struct {
		caseName string
		f        func(t *testing.T, userId string, expect controller.Result)
		params   []interface{}
	}{
		{caseName: "AfterAddUser-TestGetCurUser", f: testGetUser, params: []interface{}{t, userId, res}},
		{caseName: "AfterAddUser-DestroyUser", f: testDestroyUser, params: []interface{}{t, userId, controller.Result{Code: 0, Msg: "success"}}},
		{caseName: "AfterAddUser-DestroyUser-GetCurUser", f: testGetUser, params: []interface{}{t, userId, controller.Result{Code: 404, Msg: "用户不存在"}}},
	}
	for _, c := range cases {
		t.Run(c.caseName, func(t *testing.T) {
			c.f(c.params[0].(*testing.T), c.params[1].(string), c.params[2].(controller.Result))
		})
	}
}

func TestGetUser(t *testing.T) {
	testGetUser(t, "0", controller.Result{
		Code: 404,
		Msg:  "用户不存在",
	})
}

func TestDestroyUser(t *testing.T) {
	testDestroyUser(t, "0", controller.Result{
		Code: 0,
		Msg:  "success",
	})
}

func TestGetUserList(t *testing.T) {
	users := []string{
		`{"nickName":"list1","role":1,"userId":"userId1"}`,
		`{"nickName":"list2","role":2,"userId":"userId2"}`,
		`{"nickName":"list3","role":1,"userId":"userId3"}`,
		`{"nickName":"list4","role":2,"userId":"userId4"}`,
	}
	for _, user := range users {
		user := user
		route.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "/v1/user/", strings.NewReader(user)))
	}
	t.Run("GetUserList", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/v1/user/?skip=0&limit=3", nil)
		route.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)

		res := controller.Result{}
		_ = json.Unmarshal([]byte(w.Body.String()), &res)

		assert.Equal(t, 0, res.Code)
		assert.Equal(t, "success", res.Msg)
		users := res.Res.([]interface{})
		assert.Equal(t, 3, len(users))
		for _, m := range users {
			m := m.(map[string]interface{})
			for _, field := range []string{ "userId", "loginTime", "role", "nickName" } {
				_, ok := m[field]
				assert.Equal(t, true, ok)
			}
		}
	})
}

func testGetUser(t *testing.T, userId string, expect controller.Result) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/v1/user/"+userId, nil)
	route.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	res := controller.Result{}
	_ = json.Unmarshal([]byte(w.Body.String()), &res)

	assert.Equal(t, expect, res)
}

func testDestroyUser(t *testing.T, userId string, expect controller.Result) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/v1/user/"+userId, nil)

	route.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	res := controller.Result{}
	_ = json.Unmarshal([]byte(w.Body.String()), &res)
	assert.Equal(t, expect, res)
}
