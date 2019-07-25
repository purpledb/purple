package client

import "fmt"

var (
	ErrNoAddress = ConfigError{"no Strato server address provided"}
)

type ConfigError struct {
	string
}

func (e ConfigError) Error() string {
	return fmt.Sprintf("server config error: %s", e.string)
}

func IsConfigError(err error) bool {
	_, ok := err.(ConfigError)
	return ok
}
