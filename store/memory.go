package store

import (
	"context"
	"strings"
	"sync"
	"time"
)

// 适合单机内存缓存，高效
// Implement memory storage interface
type value struct {
	data []byte
	time time.Time
}

func (v *value) valid() bool {
	return time.Until(v.time).Microseconds() > 0
}

type memoryStore struct {
	mutex sync.RWMutex
	data  map[string]*value
}

func NewMemory() Store {
	return &memoryStore{
		data: make(map[string]*value, 10000),
	}
}

func (memory *memoryStore) Get(_ context.Context, key string) ([]byte, error) {
	memory.mutex.RLock()
	defer memory.mutex.RUnlock()
	if v, ok := memory.data[key]; ok && v.valid() {
		return v.data, nil
	}
	return nil, ErrNotFound
}

func (memory *memoryStore) Set(_ context.Context, key string, v []byte, expiration time.Duration) error {
	memory.mutex.Lock()
	defer memory.mutex.Unlock()
	memory.data[key] = &value{
		data: v,
		time: time.Now().Add(expiration),
	}
	return nil
}

func (memory *memoryStore) Del(_ context.Context, key string) error {
	memory.mutex.Lock()
	defer memory.mutex.Unlock()
	delete(memory.data, key)
	return nil
}

func (memory *memoryStore) Flush(_ context.Context, prefix string) error {
	memory.mutex.Lock()
	defer memory.mutex.Unlock()
	for key := range memory.data {
		if strings.HasPrefix(key, prefix) {
			delete(memory.data, key)
		}
	}
	return nil
}
