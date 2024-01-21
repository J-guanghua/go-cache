package store

import (
	"context"
	"time"
)

// Implement an empty storage interface
// 使用 empty 通过配置项在不改变任何业务代码情况下关闭服务缓存
type emptyStore struct {}

func NewEmpty() Store {
	return &emptyStore{}
}

func (es *emptyStore) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, ErrNotFound
}

func (es *emptyStore) Set(ctx context.Context, key string, value []byte, expiration time.Duration)error {
	return nil
}

func (rs *emptyStore) Del(ctx context.Context, key string) error {
	return nil
}
func (rs *emptyStore) Flush(ctx context.Context, prefix string) error {
	return nil
}
