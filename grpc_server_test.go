package strato

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGrpcServer(t *testing.T) {
	is := assert.New(t)

	srv, err := NewGrpcServer(goodServerCfg)

	t.Run("Instantiation", func(_ *testing.T) {
		is.NoError(err)
		is.NotNil(srv)
	})

	is.NoError(srv.mem.Close())
}
