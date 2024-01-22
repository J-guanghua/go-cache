package cache

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/J-guanghua/go-cache/store"
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

var(
	users = Key("users")
	stat = NewStat()
	//本地测试过程f1le效率高于redis很多
	caches =NewCache(
	//store.NewRedis(redis.NewClient(&redis.Option{
	//	Network: "tcp",
	//	Addr:    "192.168.43.151:6379",
	//})),
	store.NewMemory(),
	Name("integration"),
	Calls(&Logs{log: log.New(os.Stderr,"",1)},stat),
	//Store(store.NewFile()), //文件存储
	//Store(NewRedis(redis.NewClient(&redis.Option{//Network:"tcp",//Addr:"21.98.162.62:6379",
	////Addr:"service-redis.acp-dev.svc.cluster.local",//Password:“Acp@1234”,////Password:sm2.MustDecrypt("Acp@1234"),//ReadTimeout:2 *time.Second,//WriteTimeout:2 *time.Second,//}))),
	//Store(store.NewFile()), //关闭缓存
	)
	repo = UserRepo{
		table: users,
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

type UserRepo struct{
	table Key
	c Cache
	users map[int]user
}
type user struct{
	ID int `json:"id"`
	Age int `json:"age"`
	Name string`json:"name"`
}

// 获取suer对象
func(u *UserRepo) getUser(ctx context.Context,id int)(user,error){
	var user user
	return user,u.c.Take(ctx,u.table.Join(id),func(ctx context.Context)(interface{},error){
		if us,ok :=u.users[id];ok {
			defer u.c.Set(ctx,"user=test24555",id*900)
			defer u.c.Set(ctx,u.table.Join("test2",us),us.Age/3)
			log.Println(id,"执行数据user查询…………")
			return us,nil
		}
		return nil,fmt.Errorf("未找到用户%v",id)
	},&user)
}

//获取users集合
func TestCetUsers(t *testing.T){
	ctx :=Context(context.Background()).Extpiex(20 *time.Second)
	//清除table缓存
	defer repo.c.Flush(ctx,repo.table)
	for i :=0;i<20;i++{
		user,err := repo.getUser(ctx,i)
		if !reflect.DeepEqual(repo.users[i],user){
			t.Errorf("TestCetUsers()got =%v,want = %v,err= %v",repo.users[i],user,err)
		}
	}
}
//user login
func TestUserlogin(t *testing.T){
	var uid=10
	ctx :=Context(context.Background()).Extpiex(20 *time.Second)
	user,err :=repo.getUser(ctx,uid)
	if err !=nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(repo.users[uid],user){
		t.Errorf("TestUserLogin()got =%v,want %v",repo.users[uid],user)
	}

}

//user login
func TestGetOrSet(t *testing.T){
	var orders = Key("orders")
	ctx :=Context(context.Background()).Extpiex(20 *time.Second)
	//清除table缓存
	//defer repo.c.Flush(ctx,orders)
	for i :=0;i<20000;i++{
		err := caches.Set(ctx,orders.Join(i),i)
		if err != nil {
			t.Error(err)
		}
		var v int
		err = caches.Get(ctx,orders.Join(i),&v)
		if err != nil {
			t.Error(err)
		}
		err = caches.Del(ctx,orders.Join(i))
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(v,i){
			t.Errorf("TestGetOrSet()got =%v,want = %v",v,i)
		}
	}

}

func TestHttp(t *testing.T)  {
	// Detailed reference https://github.com/go-kratos/examples/tree/main/metrics
	_metricSeconds := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "server",
		Subsystem: "requests",
		Name:      "duration_sec",
		Help:      "server requests duratio(sec).",
		Buckets:   []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.250, 0.5, 1},
	}, []string{"kind", "operation"})
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "app_cache_tetal",
		Help: "test help",
	})
	_metricRequests := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "client",
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "The total number of processed requests",
	}, []string{"kind", "operation", "code", "reason"})

	prometheus.MustRegister(_metricSeconds, _metricRequests,counter)
	http.HandleFunc("/metrics", stat.HandleFunc)
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		TestGetOrSet(t)
		return
	})
	http.ListenAndServe(":8000",nil)
}

func BenchmarkSast(b *testing.B)  {
	var benchmark = Key("benchmark")
	ctx:=context.Background()
	for i:=0;i<=b.N;i++ {
		b.ReportAllocs()
		err := caches.Set(ctx, benchmark.Join(i),i)
		if err != nil {
			b.Error(err)
		}
		//var v int
		//err = caches.Get(ctx, benchmark.Join(i),&v)
		//if err != nil {
		//	//b.Error(err)
		//}
	}
}