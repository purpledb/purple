package server

import "fmt"

type ConfigError struct {
	string
}

func (e ConfigError) Error() string {
	return fmt.Sprintf("server config error: %s", e.string)
}

func ErrNoPort() ConfigError {
	return ConfigError{"no server port supplied"}
}

func IsConfigError(err error) bool {
	_, ok := err.(ConfigError)
	return ok
}
