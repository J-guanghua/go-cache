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

    // default new memory cache
    cache := cache.NewCache()
    // new file cache
    cache := cache.NewCache(
    	cache.Name("app"),
    	cache.Store(store.NewFile())
    )
    // new redis cache
    cache := cache.NewCache(
        cache.Name("app"),
        cache.Store(store.NewRedis())
    )
    // new empty cache
	cache := cache.NewCache(
        cache.Name("app"),
        cache.Store(store.NewEmpty())
    )
```
### New
```go
    //本地测试过程f1le效率高于redis很多
    type user struct{
        ID int `json:"id"`
        Age int `json:"age"`
        Name string`json:"name"`
    }
    
    cache := cache.NewCache(cache.Name("app"))

    var user = user{}
    err := cache.Set(ctx,"user1",user)
    if err != nil {
        t.Error(err)
    }
    err = cache.Get(ctx,"user1",&user)
    if err != nil {
    	t.Error(err)
    }
    err = cache.Del(ctx,"user1")

```
## Example
```go
package main

import (
	"fmt"
	"log"
	"context"
	"github.com/J-guanghua/cache"
)

// Keys & values in cache2go can be of arbitrary types, e.g. a struct.
type user struct{
	ID int `json:"id"`
	Age int `json:"age"`
	Name string`json:"name"`
}
type Users struct{
	user cache.Key
	cache cache.Cache
	data map[int]user
}
// 获取suer对象
func(u *Users) getUser(ctx context.Context,id int)(user,error){
	var user user
	return user,u.cache.Take(ctx,u.user.Join(id),func(ctx context.Context)(interface{},error){
		if us,ok :=u.users[id];ok {
			defer u.cache.Set(ctx,"user=test24555",id*900)
			defer u.cache.Set(ctx,u.table.Join("test2",us),us.Age/3)
			log.Println(id,"执行数据user查询…………")
			return us,nil
		}
		return nil,fmt.Errorf("未找到用户%v",id)
	},&user)
}

func main() {
	users := &Users{
		user:cache.Key("users"),
		cache:cache.NewCache(cache.Calls(&cache.Logs{log: log.New(os.Stderr,"",1)},cache.NewStat())),
		data: map[int]user{
			1:  user{1, 11, "test1"},
			2:  user{2, 12, "test2"},
		},
    }
	//清除user缓存
	defer users.cache.Flush(ctx,users.user)
	// 获取suer对象
	user,err := users.getUser(ctx,users.user.Jion(2))
	if err != nil {
		panic(err)
    }
	log.Println(user)
}
```

