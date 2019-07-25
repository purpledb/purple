package strato

import (
	"fmt"
	"strato/proto"
)

type (
	KV interface {
		Get(location *Location) (*Value, error)
		Put(location *Location, value *Value) error
		Delete(location *Location) error
	}

	Location struct {
		Key string
	}

	Value struct {
		Content []byte `json:"content"`
	}
)

func (l *Location) String() string {
	return fmt.Sprintf("Location<key: %s>", l.Key)
}

func (l *Location) Proto() *proto.Location {
	return &proto.Location{
		Key: l.Key,
	}
}

func (v *Value) String() string {
	return fmt.Sprintf(`Value<content: "%s">`, v.Content)
}

func (v *Value) Proto() *proto.Value {
	return &proto.Value{
		Content: v.Content,
	}
}
