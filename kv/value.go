package kv

import "fmt"

type Value struct {
	Content []byte `json:"content"`
}

func (v *Value) String() string {
	return fmt.Sprintf(`Value<content: "%s">`, v.Content)
}