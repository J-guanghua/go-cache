package calls

import (
	"context"
	"time"
)

type (
	CallOption interface {
		Before(context.Context, Action) error
		After(context.Context, Action)
	}
	Action interface {
		Err() error
		Name() string
		Key() string
		Method() string
		Value() []byte
		Extpiex() time.Duration
		SetExtpiex(extpiex time.Duration)
	}
)
