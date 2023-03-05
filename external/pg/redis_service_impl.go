package pg

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisServiceImpl struct {
	Rdb *redis.Client
}

func (r RedisServiceImpl) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return r.Rdb.Set(ctx, key, value, expiration).Err()
}

func (r RedisServiceImpl) Get(ctx context.Context, key string) (string, error) {
	return r.Rdb.Get(ctx, key).Result()
}

func (r RedisServiceImpl) Delete(ctx context.Context, key string) error {
	return r.Rdb.Del(ctx, key).Err()
}
