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
	}
	return pool
}
