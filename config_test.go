package purple

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigInstantiation(t *testing.T) {
	is := assert.New(t)

	t.Run("Client", func(t *testing.T) {
		testCases := []struct {
			config *ClientConfig
			err    error
		}{
			{&ClientConfig{}, ErrNoAddress},
			{&ClientConfig{Address: "example.com:2222"}, nil},
		}

		for _, tc := range testCases {
			is.Equal(tc.config.Validate(), tc.err)
		}
	})

	t.Run("Server", func(t *testing.T) {
		testCases := []struct {
			config *ServerConfig
			err    error
		}{
			{&ServerConfig{}, ErrNoPort},
			{&ServerConfig{Port: 1234}, ErrNoBackend},
			{&ServerConfig{Port: 10}, ErrPortOutOfRange},
			{&ServerConfig{Port: 1234, Backend: "disk"}, nil},
		}

		for _, tc := range testCases {
			is.Equal(tc.config.Validate(), tc.err)
		}
	})
}
