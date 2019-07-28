package strato

import (
	"testing"

	"github.com/lucperkins/strato/proto"

	"github.com/stretchr/testify/assert"
)

func TestKVTypes(t *testing.T) {
	is := assert.New(t)

	t.Run("Location", func(t *testing.T) {
		bucket, key := "some-bucket", "some-key"

		loc := &Location{}
		is.Equal(loc.validate(), ErrNoBucket)

		loc.Bucket = bucket
		is.Equal(loc.validate(), ErrNoKey)

		loc.Key = key
		is.NoError(loc.validate())
	})

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
