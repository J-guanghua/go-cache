package cache

import (
	"context"
	"log"
	"net/http"
)

type statTop struct {
	total    uint64 // 总次数
	errTotal uint64 // 失败次数
	keyTop   map[string]int64
}

func (s *statTop) before (envet Envet) {
	s.total++
	s.keyTop[envet.Key()]++
}

func (s *statTop) after (envet Envet) {
	if envet.Err() != nil {
		s.errTotal++
	}
}

type cacheStat struct {
	total    uint64 // 总次数
	errTotal uint64 // 失败次数
	errs     map[string]int64
	method   map[string]*statTop // key top 统计分析
}

func NewStat() *cacheStat {
	return &cacheStat{
		errs: map[string]int64{},
		method: map[string]*statTop{
			"GET":    {keyTop: make(map[string]int64, 100)},
			"SET":    {keyTop: make(map[string]int64, 100)},
			"DELITE": {keyTop: make(map[string]int64, 100)},
		},
	}
}

func (stat *cacheStat) before(ctx context.Context, envet Envet) error {
	stat.total++
	stat.method[envet.Method()].before(envet)
	return nil
}

func (stat *cacheStat) after(ctx context.Context, envet Envet) {
	if err := envet.Err(); err != nil {
		stat.errTotal++
		stat.errs[err.Error()]++
	}
}

func (stat *cacheStat) HandleFunc(writer http.ResponseWriter, request *http.Request) {
	log.SetOutput(writer)
	log.Printf("stat: 总计数(%d),错误计数(%d),errs(%v)", stat.total, stat.errTotal, stat.errs)
	log.Printf("stat: method(%v)", stat.method["GET"])
	log.Printf("stat: method(%v)", stat.method["SET"])
}
