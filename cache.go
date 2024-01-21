package cache

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/J-guanghua/cache/store"
)

const(
	extpiexKey = iota
	labelKey = iota
)

type cacheCtx struct{ context.Context }
func Context(ctx context.Context) *cacheCtx {
	return &cacheCtx{ctx}
}

//自定义失效缓存时间
func (cc *cacheCtx) Extpiex(extpiex time.Duration) *cacheCtx {
	cc.Context = context.WithValue(cc.Context, extpiexKey, extpiex)
	return cc
}

type Cache interface {
	//获取key,并将key对于的值映射到v对象,v可以是任何数据类型
	Get(ctx context.Context, key string, v interface{}) error
	//设置key=v,并将存储到 store 持久化
	Set(ctx context.Context, key string, v interface{}) error
	// 删除key 对应得 v
	Del(ctx context.Context, key string) error
	//匹配清除带 pattern 前缀的可以集合
	Flush(ctx context.Context, pattern Key) error
	//获取key获取将key的fn结果映射到v对象
	Take(ctx context.Context, key string, fn VFunc, v interface{}) error
}

// Cache 对象
type cache struct {
	name    string                  // name 用于区分不同服务|模块key可能存在的冲突
	keyFunc func(key string) string // Key 自定义string算法
	store   store.Store                 // 只需要适配 store接口, 可支持yaml自定义配置
	extpiex time.Duration  // 默认缓存失效时间
	codec   Codec //自定义序列化对象
	calls   []CallOption
}

func NewCache(store store.Store,options ...Option) Cache {
	var codec = codec("json")
	c := &cache{
		name: "global",
		codec: &codec,
		store: store,
		extpiex: 5 * time.Second,
	}
	for _,opt:= range options {
		opt(c)
	}
	return c
}

func (c *cache) buildKey(ctx context.Context, key string) string {
	if c.keyFunc != nil {
		key = c.keyFunc(key)
	}
	return strings.Join([]string{c.name, key},"=")
}

func (c *cache) getStore() store.Store {
	return c.store
}

func (c *cache) before(ctx context.Context,in *envet) error {
	for _, o := range c.calls {
		if err := o.before(ctx,in); err != nil {
			return  err
		}
	}
	return nil
}

func (c *cache) after(ctx context.Context,in *envet) {
	for _, o := range c.calls {
		o.after(ctx,in)
	}
}

func (c *cache) Get(ctx context.Context, key string, v interface{}) error {
	key = c.buildKey(ctx, key)
	return func(ctx context.Context, in *envet) error {
		defer c.after(ctx,in)
		if err := c.before(ctx,in); err != nil {
			return err
		}
		in.value, in.err = c.store.Get(ctx, key)
		if in.err != nil {
			return in.err
		}
		return c.codec.Unmarshal(in.value, v)
	}(ctx,&envet{
		name: c.name,
		key: key,
		method: "GET",
	})
}

func (c *cache) Set(ctx context.Context, key string, v interface{}) error {
	key = c.buildKey(ctx, key)
	return func(ctx context.Context, in *envet) error {
		in.value, in.err = c.codec.Marshal(v)
		if in.err != nil{
			return in.err
		}
		defer c.after(ctx,in)
		if err := c.before(ctx,in); err != nil {
			return err
		}
		in.err = c.store.Set(ctx, key, in.value, in.extpiex)
		if in.err != nil {
			return in.err
		}
		return nil
	}(ctx,&envet{
		key: key,
		name: c.name,
		method: "SET",
		extpiex: c.extpiex,
	})
}

func (c *cache) Del(ctx context.Context, key string) error {
	key = c.buildKey(ctx, key)
	return func(ctx context.Context, in *envet) error {
		defer c.after(ctx,in)
		if err := c.before(ctx,in); err != nil {
			return err
		}
		in.err = c.store.Del(ctx, key)
		if in.err != nil {
			return in.err
		}
		return nil
	}(ctx,&envet{
		name: c.name,
		key: key,
		method: "DELITE",
	})
}

func (c *cache) Flush(ctx context.Context, pattern Key) error {
	key := c.buildKey(ctx, pattern.string())
	return c.store.Flush(ctx, key)
}

func (c *cache) Take(ctx context.Context, key string, fn VFunc, v interface{}) error {
	if err := c.Get(ctx, key, v); err != nil {
		fv, err := fn(ctx)
		if err != nil {
			return err
		}
		if fv != nil {
			_= c.Set(ctx,key,fv)
			return c.allocation(ctx,fv,v)
		}
	}
	return nil
}

func (c *cache) allocation(ctx context.Context,fv,v interface{}) (err error) {
	defer func() {
		if r := recover();r != nil {
			buf := make([]byte, 64<<10) //nolint:gomnd
			n := runtime.Stack(buf, false)
			buf = buf[:n]
			err = fmt.Errorf("%v: \n%s\n", r, buf)
		}
	}()
	if reflect.TypeOf(fv).Kind() != reflect.Ptr {
		reflect.ValueOf(v).Elem().Set(reflect.ValueOf(fv))
	}else {
		reflect.ValueOf(v).Elem().Set(reflect.ValueOf(fv).Elem())
	}
	return err
}