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
	value []byte
	expiration time.Time
}

func (v *value) valid() bool {
	return v.expiration.Sub(time.Now()).Microseconds() > 0
}

type memoryStore struct {
	mutex sync.RWMutex
	data map[string]*value
}

func NewMemory() Store {
	return &memoryStore{
		data: make(map[string]*value,10000),
	}
}

func (memory *memoryStore) Get(ctx context.Context, key string) ([]byte, error) {
	memory.mutex.RLock()
	defer memory.mutex.RUnlock()
	 if v,ok:= memory.data[key];ok && v.valid(){
	 	return v.value,nil
	 }
	 return nil,ErrNotFound
}

func (memory *memoryStore) Set(ctx context.Context, key string, v []byte, expiration time.Duration) error {
	memory.mutex.Lock()
	defer memory.mutex.Unlock()
	memory.data[key] = &value{
		value: v,
		expiration: time.Now().Add(expiration),
	}
	return nil
}

func (memory *memoryStore) Del(ctx context.Context, key string) error {
	memory.mutex.Lock()
	defer memory.mutex.Unlock()
	delete(memory.data,key)
	return nil
}

func (memory *memoryStore) Flush(ctx context.Context, pattern string) error {
	memory.mutex.Lock()
	defer memory.mutex.Unlock()
	pattern = strings.Trim(pattern,"*") + "*"
	for key := range memory.data {
		if strings.HasPrefix(key,pattern) {
			delete(memory.data,key)
		}
	}
	return nil
}
