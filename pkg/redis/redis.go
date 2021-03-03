package redis

import (
	"context"
	"github.com/gomodule/redigo/redis"
)

func NewRedisPool(ctx context.Context, address string) *redis.Pool {
	pool := &redis.Pool{
		Wait: true,
		DialContext: func(ctx context.Context) (redis.Conn, error) {
			return redis.DialContext(ctx, "tcp", address)
		},
		MaxActive: 400,
		MaxIdle: 50,
	}
	return pool
}
