package calls

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

func (l Logs) Before(_ context.Context, a Action) error {
	l.log.Printf("log.before: name(%v),key(%v) ,method(%v),v(%v),err(%v)",
		a.Name(), a.Key(), a.Method(), string(a.Value()), a.Err())
	return nil
}

func (l Logs) After(_ context.Context, _ Action) {
}
