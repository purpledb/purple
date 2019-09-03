package kv

import (
	"fmt"

	"github.com/lucperkins/strato/proto"
)

type (
	KV interface {
		KVGet(key string) (*Value, error)
		KVPut(key string, value *Value) error
		KVDelete(string) error
	}

	Value struct {
		Content []byte `json:"content"`
	}
)

func (v *Value) String() string {
	return fmt.Sprintf(`Value<content: "%s">`, v.Content)
}

func (v *Value) Proto() *proto.Value {
	return &proto.Value{
		Content: v.Content,
	}
}
