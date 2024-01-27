package cache

import (
	"context"
	"fmt"
	"testing"
)

func benchmarkCache() Cache {
	return NewCache(
		Name("app"),
	)
}

func BenchmarkSet(b *testing.B) {
	c := benchmarkCache()
	ctx := context.Background()
	b.ReportAllocs()
	for i := 0; i <= b.N; i++ {
		key := Key("set").Join(i)
		if err := c.Set(ctx, key, i); err != nil {
			b.Errorf("Set(%v,%v,%v): %v", ctx, key, i, err)
		}
	}
}

func BenchmarkTake(b *testing.B) {
	c := benchmarkCache()
	ctx := context.Background()
	b.ReportAllocs()
	for i := 0; i <= b.N; i++ {
		key := Key("take").Join(i)
		var s string
		if err := c.Take(ctx, key, func(ctx context.Context) (interface{}, error) {
			return fmt.Sprintf("hello %d", i), nil
		}, &s); err != nil {
			b.Errorf("Set(%v,%v,%v): %v", ctx, key, i, err)
		}
	}
}

func BenchmarkSetAndGetAndDel(b *testing.B) {
	c := benchmarkCache()
	ctx := context.Background()
	b.ReportAllocs()
	for i := 0; i <= b.N; i++ {
		key := Key("set").Join(i)
		if err := c.Set(ctx, key, i); err != nil {
			b.Errorf("Set(%v,%v,%v): %v", ctx, key, i, err)
		}
		var v int
		if err := c.Get(ctx, key, &v); err != nil {
			b.Errorf("Get(%v,%v,%v): %v", ctx, key, i, err)
		}
		if err := c.Del(ctx, key); err != nil {
			b.Errorf("Del(%v,%v): %v", ctx, key, err)
		}
	}
}
