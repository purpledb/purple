package strato

import (
	"fmt"
	"strings"

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

func SelectKV(key string) (KV, error) {
	k := strings.ToLower(key)

	switch k {
	case "memory", "mem":
		return NewMemoryBackend(), nil
	case "disk":
		return NewDisk("tmp/strato")
	default:
		return nil, ErrBackendNotRecognized
	}
}

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
		Bucket: l.Bucket,
		Key:    l.Key,
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
