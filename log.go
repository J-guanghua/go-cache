package cache

import (
	"context"
	"log"
)


type Logs struct {
	log *log.Logger
}

func (l *Logs) before(ctx context.Context,envet Envet) error {
	return nil
}

func (l *Logs) after(ctx context.Context,envet Envet) {
	//l.log.Printf("log.after: name(%v),key(%v) ,method(%v),v(%v),err(%v)",
	//	envet.Name(),envet.Key(),envet.Method(),envet.Value(),envet.Err())
}
