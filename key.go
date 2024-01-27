package cache

import (
	"fmt"
	"strings"
)

type (
	Key   string
	Label string
)

func (k Key) Join(keys ...interface{}) string {
	return strings.Join([]string{k.Name(), fmt.Sprint(keys...)}, "")
}

func (k Key) Name() string {
	return fmt.Sprintf("%s#", k)
}
