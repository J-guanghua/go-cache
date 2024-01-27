package cache

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/J-guanghua/go-cache/store"
)

var ErrNotFound = store.ErrNotFound

type Cache interface {
	// 获取key,并将key对于的值映射到v对象,v可以是任何数据类型
	Get(ctx context.Context, key string, v interface{}) error
	// 设置key=v,并将存储到 store 持久化
	Set(ctx context.Context, key string, v interface{}, as ...func(*action)) error
	// 删除key 对应得 v
	Del(ctx context.Context, key string) error
	// 匹配清除带 pattern 前缀的可以集合
	Flush(ctx context.Context, pattern Key) error
	// 获取key获取将key的fn结果映射到v对象
	Take(ctx context.Context, key string, fn VFunc, v interface{}) error
}

// Cache 对象
type cache struct {
	name    string                  // name 标识名称
	keyFunc func(key string) string // Key 自定义string算法
	store   store.Store             // 只需要适配 store接口, 可支持yaml自定义配置
	extpiex time.Duration           // 默认缓存失效时间
	codec   Codec                   // 自定义序列化对象
	calls   []CallOption
}

func NewCache(options ...Option) Cache {
	codec := codec("json")
	c := &cache{
		name:    "app",
		codec:   &codec,
		store:   store.NewMemory(),
		extpiex: 5 * time.Second,
	}
	for _, opt := range options {
		opt(c)
	}
	return c
}

func (c *cache) buildKey(ctx context.Context, key string) string {
	if keys := strings.Split(key, "#"); len(keys) <= 1 {
		key = "default#" + key
	}
	if c.keyFunc != nil {
		key = c.keyFunc(key)
	}
	return strings.Join([]string{c.name, key}, ".")
}

func (c *cache) before(ctx context.Context, in *action) error {
	for _, o := range c.calls {
		if err := o.before(ctx, in); err != nil {
			return err
		}
	}
	return nil
}

func (c *cache) after(ctx context.Context, in *action) {
	for _, o := range c.calls {
		o.after(ctx, in)
	}
}

func (c *cache) Get(ctx context.Context, key string, v interface{}) error {
	key = c.buildKey(ctx, key)
	return func(ctx context.Context, in *action) error {
		if err := c.before(ctx, in); err != nil {
			return err
		}
		defer c.after(ctx, in)
		in.value, in.err = c.store.Get(ctx, key)
		if in.err != nil {
			return in.err
		}
		return c.codec.Unmarshal(in.value, v)
	}(ctx, &action{
		name:   c.name,
		key:    key,
		method: "GET",
	})
}

func (c *cache) Set(ctx context.Context, key string, v interface{}, opts ...func(*action)) error {
	key = c.buildKey(ctx, key)
	in := &action{key: key, name: c.name, method: "SET", extpiex: c.extpiex}
	for _, opt := range opts {
		opt(in)
	}
	if err := c.before(ctx, in); err != nil {
		return err
	}
	defer c.after(ctx, in)
	in.value, in.err = c.codec.Marshal(v)
	if in.err != nil {
		return in.err
	}
	in.err = c.store.Set(ctx, key, in.value, in.extpiex)
	if in.err != nil {
		return in.err
	}
	return nil
}

func (c *cache) Del(ctx context.Context, key string) error {
	key = c.buildKey(ctx, key)
	return func(ctx context.Context, in *action) error {
		if err := c.before(ctx, in); err != nil {
			return err
		}
		defer c.after(ctx, in)
		in.err = c.store.Del(ctx, key)
		if in.err != nil {
			return in.err
		}
		return nil
	}(ctx, &action{
		name:   c.name,
		key:    key,
		method: "DELITE",
	})
}

func (c *cache) Flush(ctx context.Context, pattern Key) error {
	in := &action{
		name:   c.name,
		key:    pattern.Name(),
		method: "FLUSH",
	}
	key := c.buildKey(ctx, pattern.Name())
	if err := c.before(ctx, in); err != nil {
		return err
	}
	defer c.after(ctx, in)
	in.err = c.store.Flush(ctx, key)
	return in.err
}

func (c *cache) Take(ctx context.Context, key string, fn VFunc, v interface{}) error {
	if err := c.Get(ctx, key, v); err != nil {
		fv, err := fn(ctx)
		if err != nil {
			return err
		}
		if fv != nil {
			_ = c.Set(ctx, key, fv)
			return c.allocation(ctx, fv, v)
		}
	}
	return nil
}

func (c *cache) allocation(ctx context.Context, fv, v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 64<<10) //nolint:gomnd
			n := runtime.Stack(buf, false)
			buf = buf[:n]
			err = fmt.Errorf("%v: %s\n", r, buf)
		}
	}()
	if reflect.TypeOf(fv).Kind() != reflect.Ptr {
		reflect.ValueOf(v).Elem().Set(reflect.ValueOf(fv))
	} else {
		reflect.ValueOf(v).Elem().Set(reflect.ValueOf(fv).Elem())
	}
	return err
}
