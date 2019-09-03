package kv

import (
	"testing"

	"github.com/lucperkins/strato/proto"

	"github.com/stretchr/testify/assert"
)

func TestKVTypes(t *testing.T) {
	is := assert.New(t)

	t.Run("Value", func(t *testing.T) {
		content := []byte("some test content")

		val := &Value{
			Content: content,
		}

		is.Equal(val.String(), `Value<content: "some test content">`)
		is.Equal(val.Proto(), &proto.Value{Content: content})
		is.Equal(val.Content, content)
	})
}
