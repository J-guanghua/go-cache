package cache

import (
	"context"
	"fmt"
	"github.com/J-guanghua/go-cache/store"
	"log"
	"os"
	"reflect"
	"testing"
)

//config.yaml配置上示列
//value:
//cache:
//  store:redis //支持file,empty,redis
//   extpiex:2s
//   ......
//实例化
//wire.NewSet(NewCache)
//func NewCache(r *conf.Data)cache.Cache {
//return cache.NewCache(
//
//)

var (
	stat  = NewStat()
	caches = NewCache(
		Store(store.NewMemory()),
		Name("integration"),
		Calls(&Logs{log: log.New(os.Stderr, "", 1)}, stat),
	)
	repo = UserRepo{
		table: Key("users"),
		c:     caches,
		users: map[int]user{
			1:  user{1, 11, "test1"},
			2:  user{2, 12, "test2"},
			3:  user{3, 13, "test3"},
			4:  user{4, 14, "test4"},
			5:  user{5, 15, "test5"},
			6:  user{6, 16, "test6"},
			7:  user{7, 17, "test7"},
			8:  user{8, 18, "test8"},
			9:  user{9, 19, "test9"},
			10: user{10, 20, "test10"},
		},
	}
)

type UserRepo struct {
	table Key
	c     Cache
	users map[int]user
}

type user struct {
	ID   int    `json:"id"`
	Age  int    `json:"age"`
	Name string `json:"name"`
}

// 获取suer对象
func (u *UserRepo) getUser(ctx context.Context, id int) (user, error) {
	var user user
	return user, u.c.Take(ctx, u.table.Join(id), func(ctx context.Context) (interface{}, error) {
		if us, ok := u.users[id]; ok {
			defer u.c.Set(ctx, "user=test24555", id*900)
			defer u.c.Set(ctx, u.table.Join("test2", us), us.Age/3)
			log.Println(id, "执行数据user查询…………")
			return us, nil
		}
		return nil, fmt.Errorf("未找到用户%v", id)
	}, &user)
}

//获取users集合
func TestCetUsers(t *testing.T) {
	ctx := context.Background()
	//清除table缓存
	defer repo.c.Flush(ctx, repo.table)
	for i := 0; i < 20; i++ {
		user, err := repo.getUser(ctx, i)
		if !reflect.DeepEqual(repo.users[i], user) {
			t.Errorf("TestCetUsers()got =%v,want = %v,err= %v", repo.users[i], user, err)
		}
	}
}

// user login
func TestUserlogin(t *testing.T) {
	var uid = 10
	ctx := context.Background()
	user, err := repo.getUser(ctx, uid)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(repo.users[uid], user) {
		t.Errorf("TestUserLogin()got =%v,want %v", repo.users[uid], user)
	}

}

func TestGetOrSet(t *testing.T) {
	orders := Key("orders")
	ctx := context.Background()
	//清除table缓存
	defer repo.c.Flush(ctx,orders)
	for i := 0; i < 20000; i++ {
		err := caches.Set(ctx, orders.Join(i), i)
		if err != nil {
			t.Error(err)
		}
		var v int
		err = caches.Get(ctx, orders.Join(i), &v)
		if err != nil {
			t.Error(err)
		}
		err = caches.Del(ctx, orders.Join(i))
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(v, i) {
			t.Errorf("TestGetOrSet()got =%v,want = %v", v, i)
		}
	}

}
