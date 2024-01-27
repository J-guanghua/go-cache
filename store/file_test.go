package store

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	file := NewFile()
	ctx := context.Background()
	err := file.Set(ctx, "key", []byte("Hello world"), time.Second)
	if err != nil {
		t.Error(err)
	}
	_, err = file.Get(ctx, "key")
	if err != nil {
		t.Error(err)
	}
}

func TestDel(t *testing.T) {
	file := NewFile()
	ctx := context.Background()
	err := file.Del(ctx, "key")
	if err != nil {
		t.Error(err)
	}
}

func TestGet(t *testing.T) {
	file := NewFile()
	ctx := context.Background()
	_, err := file.Get(ctx, "key")
	if !errors.Is(err, ErrNotFound) {
		t.Error(err)
	}
}

func TestFlush(t *testing.T) {
	file := NewFile()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_ = file.Set(ctx, "key-1", []byte("Hello 1"), time.Second)
	_ = file.Set(ctx, "key-2", []byte("Hello 2"), time.Second)
	_ = file.Set(ctx, "key-3", []byte("Hello 3"), time.Second)
	_ = file.Flush(ctx, "key")
	_, err := file.Get(ctx, "key-3")
	if !errors.Is(err, ErrNotFound) {
		t.Error(err)
	}
}
