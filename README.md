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
### New
```go
    import (
        "github.com/J-guanghua/go-cache"
        "github.com/J-guanghua/go-cache/store"
    )

    // new a cache object
    caches := cache.NewCache()
    
    // new a file cache object
    caches = cache.NewCache(cache.Store(store.NewFile(store.Directory("./cache"))))
    
    // new a redis cache object
    caches = cache.NewCache(cache.Store(
        store.NewRedis(redis.NewClient(&redis.Options{
            Network: "tcp",
            Addr: "127.0.0.1:3306",
        }))),
    )
    // Turn off cache
    caches = cache.NewCache(cache.Store(store.NewEmpty()))
        
```
### New
```go

ctx := context.Background()
var value = map[string]interface{}{}
var user = map[string]interface{}{"name":"张三"}

caches.Set(ctx,cache.Key("user").Join(1),user)
err := caches.Get(ctx,cache.Key("user").Join(1),&value)
if err != nil {
    panic(err)
}
fmt.Println(value)

var text string
err = caches.Take(ctx,"message", func(ctx context.Context) (interface{}, error) {
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

// 获取suer对象
func(u *Users) getUser(ctx context.Context,id int)(user,error){
	var user user
	return user,u.cache.Take(ctx,u.user.Join(id),func(ctx context.Context)(interface{},error){
		if us,ok :=u.data[id];ok {
			defer u.cache.Set(ctx,"name",us.Name)
			defer u.cache.Set(ctx,u.user.Join("age",id),us.Age)
			log.Println(id,"执行数据user查询…………")
			return us,nil
		}
		return nil,fmt.Errorf("未找到用户%v",id)
	},&user)
}

func main() {
	users := &Users{
		user:cache.Key("users"),
		cache:cache.NewCache(cache.Store(store.NewFile()),cache.Calls(cache.NewLog())),
		data: map[int]user{
			1:  user{1, 11, "test1"},
			2:  user{2, 12, "test2"},
		},
	}
	ctx := context.Background()
	// 清除user缓存
	defer users.cache.Flush(ctx,users.user)
	// 获取suer对象
	user,err := users.getUser(ctx,2)
	if err != nil {
		panic(err)
	}
	log.Println(user)
}
```

