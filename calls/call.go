package calls

import (
	"context"
	"time"
)

type (
	CallOption interface {
		// Before call
		Before(context.Context, Action) error
		// post-call
		After(context.Context, Action)
	}
	Action interface {
		Err() error
		Name() string
		Key() string
		Method() string
		Value() []byte
		Duration() time.Duration
		SetDuration(extpiex time.Duration)
	}
)
