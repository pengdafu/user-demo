package main

import (
	"context"
	"flag"
	v1 "user/api/v1"
	"user/pkg/redis"
)

var redisAddr string

func init() {
	flag.StringVar(&redisAddr, "redisAddr", "", "获取redis地址")
}

func main() {
	flag.Parse()
	redisPoll := redis.NewRedisPool(context.Background(), redisAddr)
	r := v1.NewRoute(redisPoll)
	r.Run(":8080")
}
