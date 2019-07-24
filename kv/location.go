package kv

import "fmt"

type Location struct {
	Key string
}

func (l *Location) String() string {
	return fmt.Sprintf("location<key: %s>", l.Key)
}
