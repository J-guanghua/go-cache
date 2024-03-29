package cache

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/J-guanghua/go-cache/calls"
	"github.com/J-guanghua/go-cache/store"
)

func testCache() Cache {
	return NewCache(
		Store(store.NewMemory()),
		Calls(calls.NewLog()),
	)
}

func TestCache(t *testing.T) {
	caches := []Cache{
		// NewCache(Name("redis")),
		NewCache(Name("memory")),
		NewCache(Name("file"), Store(store.NewFile())),
		NewCache(Name("empty"), Store(store.NewEmpty())),
	}
	ctx := context.Background()
	for i, c := range caches {
		err := c.Set(ctx, Key("cache").Join(i), i)
		if err != nil {
			t.Error(err)
		}
		err = c.Del(ctx, Key("cache").Join(i))
		if err != nil {
			t.Error(err)
		}
	}
}

func TestSet(t *testing.T) {
	c := testCache()
	ctx := context.Background()
	for i := 0; i < 10; i++ {
		err := c.Set(ctx, Key("set").Join(i), i)
		if err != nil {
			t.Error(err)
		}
		err = c.Del(ctx, Key("set").Join(i))
		if err != nil {
			t.Error(err)
		}
	}
}

func TestSetDuration(t *testing.T) {
	c := testCache()
	value := time.Now()
	ctx := context.Background()
	key := Key("time").Join(value)
	duration := SetDuration(time.Second)
	err := c.Set(ctx, key, value, duration)
	if err != nil {
		t.Error(err)
	}

	var result time.Time
	err = c.Get(ctx, key, &result)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(value.Format("2006-01-02 15:04:05"), result.Format("2006-01-02 15:04:05")) {
		t.Errorf("key = %v, result =%s,value = %s",
			key, result.Format("2006-01-02 15:04:05"), value.Format("2006-01-02 15:04:05"))
	}

	time.Sleep(time.Second)
	err = c.Get(ctx, key, &result)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v", err)
	}
}

func TestTake(t *testing.T) {
	c := testCache()
	ctx := context.Background()
	for i := 0; i < 10; i++ {
		var v int
		err := c.Take(ctx, Key("take").Join(i), func(ctx context.Context) (interface{}, error) {
			return i, nil
		}, &v)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(i, v) {
			t.Errorf("i =%v,v = %v", i, v)
		}
	}
}

func TestDifferentDataTypes(t *testing.T) {
	c := testCache()
	ctx := context.Background()
	tests := []struct {
		name     string
		key      string
		value    interface{}
		wantFunc func(t *testing.T, key string) interface{}
	}{
		{
			name: "struct 验证",
			key:  "struct",
			value: struct {
				ID   int64             `json:"id"`
				Name string            `json:"Name"`
				Data map[string]string `json:"data"`
			}{
				ID:   2,
				Name: "张三",
			},
			wantFunc: func(t *testing.T, key string) interface{} {
				want := struct {
					ID   int64             `json:"id"`
					Name string            `json:"Name"`
					Data map[string]string `json:"data"`
				}{}
				err := c.Get(ctx, key, &want)
				if err != nil {
					t.Error(err)
				}
				return want
			},
		},
		{
			name:  "map 验证",
			key:   "map",
			value: map[string]string{"Name": "张三"},
			wantFunc: func(t *testing.T, key string) interface{} {
				want := map[string]string{}
				err := c.Get(ctx, key, &want)
				if err != nil {
					t.Error(err)
				}
				return want
			},
		},
		{
			name:  "float32 验证",
			key:   "float32",
			value: float32(300.33),
			wantFunc: func(t *testing.T, key string) interface{} {
				var want float32
				err := c.Get(ctx, key, &want)
				if err != nil {
					t.Error(err)
				}
				return want
			},
		},
		{
			name:  "Name 验证",
			key:   "Name",
			value: "hello word",
			wantFunc: func(t *testing.T, key string) interface{} {
				var got string
				err := c.Get(ctx, key, &got)
				if err != nil {
					t.Error(err)
				}
				return got
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := c.Set(ctx, test.key, test.value)
			if err != nil {
				t.Error(err)
			}
			want := test.wantFunc(t, test.key)
			if !reflect.DeepEqual(test.value, want) {
				t.Errorf("Name = %v, value =%v,want = %v", test.name, test.value, want)
			}
			err = c.Del(ctx, test.key)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
