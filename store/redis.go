package store

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strings"
	"time"
)

// redis implements the storage interface
// redis is suitable for distributed service caching
type redisStore struct {
	client redis.Cmdable
}

func NewRedis(client redis.Cmdable) Store {
	return &redisStore{
		client: client,
	}
}

func (redis *redisStore) Get(ctx context.Context, key string) ([]byte, error) {
	return redis.client.Get(ctx,key).Bytes()
}

func (redis *redisStore) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	return redis.client.Set(ctx,key,value,expiration).Err()
}

func (redis *redisStore) Del(ctx context.Context, key string) error {
	return redis.client.Del(ctx,key).Err()
}

func (redis *redisStore) Flush(ctx context.Context, pattern string) error {
	pattern = strings.Trim(pattern,"*") + "*"
	keys,err := redis.client.Keys(ctx,pattern).Result()
	if err != nil {
		return err
	}
	return redis.client.Del(ctx,keys...).Err()
}
