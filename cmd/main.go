package main

import (
	"context"
	"flag"
	"user/api"
	v1 "user/api/v1"
	"user/pkg/log"
	"user/pkg/redis"
)

var redisAddr string
var httpPort string

func init() {
	flag.StringVar(&redisAddr, "redisAddr", "", "获取redis地址")
	flag.StringVar(&httpPort, "httpPort", "8080", "http服务地址")
}

func main() {
	flag.Parse()
	redisPoll := redis.NewRedisPool(context.Background(), redisAddr)
	conn, err := redisPoll.GetContext(context.Background())
	if err != nil {
		panic("Redis 连接失败，请检查是否设置 redisAddr 参数")
	}
	_ = conn.Close()

	e := api.NewGinEngine()
	v1.SetV1Route(e, redisPoll)

	exit := make(chan error)
	go func() {
		exit <- e.Run(":" + httpPort)
	}()
	log.Logger().Infof("server start at :%v", httpPort)

	err = <-exit
	_ = redisPoll.Close()
	log.Logger().Fatalf("server shutdown: %v", err)
}
