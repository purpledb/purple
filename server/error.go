package server

import "fmt"

var (
	ErrNoPort         = ConfigError{"no server port supplied"}
	ErrPortOutOfRange = ConfigError{"port must be between 1024 and 49151"}
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
