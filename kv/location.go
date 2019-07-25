package kv

import (
	"fmt"
	"strato/proto"
)

type Location struct {
	Key string
}

func (l *Location) String() string {
	return fmt.Sprintf("Location<key: %s>", l.Key)
}

func (l *Location) Proto() *proto.Location {
	return &proto.Location{
		Key: l.Key,
	}
}