package store

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("record not found")

// This is the cache storage interface
type Store interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, data []byte, expiration time.Duration) error
	Del(ctx context.Context, key string) error
	Flush(ctx context.Context, prefix string) error
}
