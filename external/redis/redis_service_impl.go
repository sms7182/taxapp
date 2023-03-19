package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type ServiceImpl struct {
	Rdb *redis.Client
}

func (r ServiceImpl) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return r.Rdb.Set(ctx, key, value, expiration).Err()
}

func (r ServiceImpl) Get(ctx context.Context, key string) (string, error) {
	return r.Rdb.Get(ctx, key).Result()
}

func (r ServiceImpl) Delete(ctx context.Context, key string) error {
	return r.Rdb.Del(ctx, key).Err()
}
