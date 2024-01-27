package store

import (
	"context"
	"errors"
	"testing"
	"time"
)

func newRedis() Store {
	//return NewRedis(redis.NewClient(&redis.Options{
	//	Network: "tcp",
	//	Addr: "192.168.43.151:6379",
	//}))
	return NewMemory()
}

func TestRedisSet(t *testing.T) {
	redis := newRedis()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := redis.Set(ctx, "key", []byte("Hello world"), time.Second)
	if err != nil {
		t.Error(err)
	}
	_, err = redis.Get(ctx, "key")
	if err != nil {
		t.Error(err)
	}
}

func TestRedisDel(t *testing.T) {
	redis := newRedis()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := redis.Del(ctx, "key")
	if err != nil {
		t.Error(err)
	}
}

func TestRedisGet(t *testing.T) {
	redis := newRedis()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := redis.Get(ctx, "key")
	if !errors.Is(err, ErrNotFound) {
		t.Error(err)
	}
}

func TestRedisFlush(t *testing.T) {
	redis := newRedis()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_ = redis.Set(ctx, "key-1", []byte("Hello 1"), time.Second)
	_ = redis.Set(ctx, "key-2", []byte("Hello 2"), time.Second)
	_ = redis.Set(ctx, "key-3", []byte("Hello 3"), time.Second)
	_ = redis.Flush(ctx, "key")
	_, err := redis.Get(ctx, "key-3")
	if !errors.Is(err, ErrNotFound) {
		t.Error(err)
	}
}
