package calls

import (
	"context"
	"log"
	"time"
)

type (
	Envet interface {
		Name() string
		Err() error
		Key() string
		Method() string
		Value() []byte
		Extpiex ()time.Duration
	}
)

type Logs struct {
	log log.Logger
}

func (l *Logs) before(ctx context.Context,envet Envet) error {
	l.log.Printf("log.before: name = %v,key = %v ,method = %v,v = %v,err = %v",
		envet.Name(),envet.Name(),envet.Key(),envet.Method(),envet.Err())
	return nil
}

func (l *Logs) after(ctx context.Context,envet Envet) {
	l.log.Printf("log.before: name = %v,key = %v ,method = %v,v = %v,err = %v",
		envet.Name(),envet.Name(),envet.Key(),envet.Method(),envet.Err())
}
