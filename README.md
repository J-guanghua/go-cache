# go-cache
微服务 缓存方案


# go-cache

[![Latest Release](https://img.shields.io/github/release/muesli/cache2go.svg)](https://github.com/muesli/cache2go/releases)
[![Build Status](https://github.com/muesli/cache2go/workflows/build/badge.svg)](https://github.com/muesli/cache2go/actions)
[![Coverage Status](https://coveralls.io/repos/github/muesli/cache2go/badge.svg?branch=master)](https://coveralls.io/github/muesli/cache2go?branch=master)
[![Go ReportCard](https://goreportcard.com/badge/muesli/cache2go)](https://goreportcard.com/report/muesli/cache2go)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://pkg.go.dev/github.com/muesli/cache2go)

Concurrency-safe golang caching library with expiration capabilities.

## Installation

Make sure you have a working Go environment (Go 1.2 or higher is required).
See the [install instructions](https://golang.org/doc/install.html).

To install go-cache, simply run:

    go get github.com/J-guanghua/go-cache

To compile it from source:

    cd $GOPATH/src/github.com/J-guanghua/go-cache
    go get -u -v
    go build && go test -v
### New cache interface
```go
    import (
        "github.com/J-guanghua/go-cache"
        "github.com/J-guanghua/go-cache/store"
    )

    // new a cache object
    c := cache.NewCache()
    
    // new a file cache object
    c = cache.NewCache(cache.Store(store.NewFile(store.Directory("./cache"))))
    
    // new a redis cache object
    c = cache.NewCache(cache.Store(
        store.NewRedis(redis.NewClient(&redis.Options{
            Network: "tcp",
            Addr: "127.0.0.1:3306",
        }))),
    )
    // Turn off cache
    c = cache.NewCache(cache.Store(store.NewEmpty()))
        
```
### reference
```go
c := cache.NewCache(
	Name("app"), // 当前服务名称 隔离不同服务缓存key,默认app
	Calls(calls.NewLog(),calls.NewStat(), // 打印日志信息,缓存统计
	cache.Duration(10 * time.Second), // 默认失效时间 10秒后
)
	
ctx := context.Background()
user := map[string]interface{}{"name":"张三"}

// 设置 app.default#user-1  1秒后失效
c.Set(ctx,"user-1",user,cache.SetDuration(time.Second))

value := map[string]interface{}{}
err := caches.Get(ctx,"user-1",&value)
if err != nil {
    panic(err)
}
fmt.Println(value)

// app.default#message  默认10秒 后失效
var text string
err = c.Take(ctx,"message", func(ctx context.Context) (interface{}, error) {
    return "Hello world",nil
},&text)
if err != nil {
    panic(err)
}
fmt.Println(text)
```
## Example
```go
package main

import (
	"context"
	"fmt"
	"github.com/J-guanghua/go-cache"
	"github.com/J-guanghua/go-cache/calls"
	"github.com/J-guanghua/go-cache/store"
	"log"
)

// user struct
type user struct{
	ID int `json:"id"`
	Age int `json:"age"`
	Name string`json:"name"`
}

// user model
type Users struct{
	user cache.Key
	cache cache.Cache
	data map[int]user
}

func model() *Users {
	return &Users{
		user:cache.Key("users"),
		cache:cache.NewCache(cache.Store(store.NewFile()),cache.Calls(calls.NewLog())),
		data: map[int]user{
			1:  user{1, 11, "test1"},
			2:  user{2, 12, "test2"},
		},
	}
}

// 获取suer对象
func(u *Users) getUser(ctx context.Context,id int)(user,error){
	var user user
	return user,u.cache.Take(ctx,u.user.Join(id),func(ctx context.Context)(interface{},error){
		if us,ok :=u.data[id];ok {
			defer u.cache.Set(ctx,u.user.Join("name"),us.Name)
			defer u.cache.Set(ctx,u.user.Join("age",id),us.Age)
			log.Println(id,"执行数据user查询…………")
			return us,nil
		}
		return nil,fmt.Errorf("未找到用户%v",id)
	},&user)
}

func main() {
	users := model()
	ctx := context.Background()
	// 获取 app.users#2 对象
	user,err := users.getUser(ctx,2)
	if err != nil {
		panic(err)
	}
	// 删除 app.user#2 缓存
	defer users.cache.Del(ctx,users.user.Join(2))
	// 清除前缀 app.users# 所有缓存（app.users#name,app.users#age）...
	defer users.cache.Flush(ctx,users.user)
	log.Println(user)
}
```

