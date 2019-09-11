package strato

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNoKey                = errors.New("no resource key provided")
	ErrNoValue              = errors.New("no value provided")
	ErrHttpUnavailable      = errors.New("could not connect to Strato HTTP server")
	ErrNoAddress            = ConfigError{"no server address provided"}
	ErrNoPort               = ConfigError{"no server port provided"}
	ErrPortOutOfRange       = ConfigError{"port must be between 1024 and 49151"}
	ErrBackendNotRecognized = ConfigError{"backend key not recognized"}
	ErrNoBackend            = ConfigError{"no backend specified"}
)

type (
	ConfigError struct {
		string
	}

	NotFoundError struct {
		string
	}

	SetError struct {
		string
	}
)

func (e ConfigError) Error() string {
	return fmt.Sprintf("config error: %s", e.string)
}

func IsConfigError(err error) bool {
	_, ok := err.(ConfigError)
	return ok
}

func (e SetError) Error() string {
	return fmt.Sprintf("set error: %s", e.string)
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf(`no value found for %s`, e.string)
}

func (e NotFoundError) AsProtoStatus() error {
	return status.Error(codes.NotFound, e.Error())
}

func NotFound(key string) NotFoundError {
	return NotFoundError{key}
}

func IsNotFound(err error) bool {
	_, ok := err.(NotFoundError)
	return ok
}
