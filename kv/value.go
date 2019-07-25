package kv

import (
	"fmt"
	"strato/proto"
)

type Value struct {
	Content []byte `json:"content"`
}

func (v *Value) String() string {
	return fmt.Sprintf(`Value<content: "%s">`, v.Content)
}

func (v *Value) Proto() *proto.Value {
	return &proto.Value{
		Content: v.Content,
	}
}