package cache

import (
	"fmt"
	"strings"
)

type Key string
type Label string

func (k Key) Join(keys ...interface{}) string {
	return strings.Join([]string{k.string(), fmt.Sprint(keys...)}, "")
}

func (k Key) string() string {
	return fmt.Sprintf("%s.", k)
}
