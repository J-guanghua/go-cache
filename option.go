package cache

import (
	"context"
	"encoding/json"
	"time"
)

type (
	Options func(c *cache)
	VFunc func(ctx context.Context) (interface{}, error)
	CallOption interface {
		before(context.Context,Envet) error
		after(context.Context,Envet)
	}
)
func Name(name string) Options {
	return func(c *cache) {
		c.name = name
	}
}

func Calls(call ...CallOption) Options {
	return func(c *cache){
		c.calls = call
	}
}

func Extpiex(extpiex time.Duration) Options {
	return func(c *cache){
		c. extpiex = extpiex
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

type envet struct {
	err     error
	name    string
	key     string
	method  string
	value   []byte
	extpiex time.Duration
}
func (e *envet) Name() string {
	return e.name
}
func (e *envet) Key() string {
	return e.key
}
func (e *envet) Method() string {
	return e.method
}
func (e *envet) Err() error {
	return e.err
}
func (e *envet) Value() []byte {
	return e.value
}
func (e *envet) Extpiex() time.Duration {
	return e.extpiex
}
func (e *envet) SetExtpiex(extpiex time.Duration) {
	e.extpiex = extpiex
}