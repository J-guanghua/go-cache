package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/J-guanghua/go-cache/calls"
	"github.com/J-guanghua/go-cache/store"
)

type (
	Option func(c *cache)
	VFunc  func(ctx context.Context) (interface{}, error)
)

func Name(name string) Option {
	return func(c *cache) {
		c.name = name
	}
}

func Calls(call ...calls.CallOption) Option {
	return func(c *cache) {
		c.calls = call
	}
}

func Extpiex(extpiex time.Duration) Option {
	return func(c *cache) {
		c.extpiex = extpiex
	}
}

func Store(store store.Store) Option {
	return func(c *cache) {
		c.store = store
	}
}

type codec string

type Codec interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

func (*codec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (*codec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

type action struct {
	err     error
	name    string
	key     string
	method  string
	value   []byte
	extpiex time.Duration
}

func SetDuration(duration time.Duration) func(*action) {
	return func(a *action) {
		a.extpiex = duration
	}
}

func (a *action) Name() string {
	return a.name
}

func (a *action) Key() string {
	return a.key
}

func (a *action) Method() string {
	return a.method
}

func (a *action) Err() error {
	return a.err
}

func (a *action) Value() []byte {
	return a.value
}

func (a *action) Extpiex() time.Duration {
	return a.extpiex
}

func (a *action) SetExtpiex(extpiex time.Duration) {
	a.extpiex = extpiex
}
