package kv

import (
	"github.com/lucperkins/strato/proto"
)

type (
	KV interface {
		KVGet(key string) (*Value, error)
		KVPut(key string, value *Value) error
		KVDelete(key string) error
	}

	Value struct {
		Metadata map[string]string `json:"metadata"`
		Content  []byte            `json:"content"`
	}
)

func (v *Value) Proto() *proto.Value {
	return &proto.Value{
		Content: v.Content,
	}
}
