package cache

import (
	"context"
	"log"
	"os"
)

type Logs struct {
	log *log.Logger
}

func NewLog() *Logs {
	return &Logs{log: log.New(os.Stderr, "", 1)}
}
func (l *Logs) before(ctx context.Context, envet Action) error {
	return nil
}

func (l *Logs) after(ctx context.Context, envet Action) {
	l.log.Printf("log.after: name(%v),key(%v) ,method(%v),v(%v),err(%v)",
		envet.Name(), envet.Key(), envet.Method(), string(envet.Value()), envet.Err())
}
