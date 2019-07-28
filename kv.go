package strato

import (
	"fmt"

	"github.com/lucperkins/strato/proto"
)

type (
	KV interface {
		KVGet(location *Location) (*Value, error)
		KVPut(location *Location, value *Value) error
		KVDelete(location *Location) error
	}

	Location struct {
		Bucket string
		Key    string
	}

	Value struct {
		Content []byte `json:"content"`
	}
)

func (l *Location) validate() error {
	if l.Bucket == "" {
		return ErrNoBucket
	}

	if l.Key == "" {
		return ErrNoKey
	}

	return nil
}

func (l *Location) String() string {
	return fmt.Sprintf("Location<bucket: %s, key: %s>", l.Bucket, l.Key)
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
