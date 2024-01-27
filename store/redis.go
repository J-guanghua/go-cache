package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
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

func (red *redisStore) Get(ctx context.Context, key string) ([]byte, error) {
	b, err := red.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return b, nil
}

func (red *redisStore) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	return red.client.Set(ctx, key, value, expiration).Err()
}

func (red *redisStore) Del(ctx context.Context, key string) error {
	return red.client.Del(ctx, key).Err()
}

func (red *redisStore) Flush(ctx context.Context, prefix string) error {
	prefix = strings.Trim(prefix, "*") + "*"
	keys, err := red.client.Keys(ctx, prefix).Result()
	if err != nil {
		return err
	}
	return red.client.Del(ctx, keys...).Err()
}
