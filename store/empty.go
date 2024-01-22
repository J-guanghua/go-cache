package store

import (
	"context"
	"time"
)

// Implement an empty storage interface
// 使用 empty 通过配置项在不改变任何业务代码情况下关闭服务缓存
type emptyStore string

func NewEmpty() Store {
	var empty emptyStore
	return &empty
}

func (empty *emptyStore) Get(_ context.Context, _ string) ([]byte, error) {
	return nil, ErrNotFound
}

func (empty *emptyStore) Set(_ context.Context, _ string, _ []byte, _ time.Duration) error {
	return nil
}

func (empty *emptyStore) Del(_ context.Context, _ string) error {
	return nil
}

func (empty *emptyStore) Flush(_ context.Context, _ string) error {
	return nil
}
