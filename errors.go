package purple

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNoKey                = errors.New("no resource key provided")
	ErrNoValue              = errors.New("no value provided")
	ErrNoAddress            = errors.New("no server address provided")
	ErrNoPort               = errors.New("no server port provided")
	ErrPortOutOfRange       = errors.New("port must be between 1024 and 49151")
	ErrBackendNotRecognized = errors.New("backend key not recognized")
	ErrNoBackend            = errors.New("no backend specified")
)

type NotFoundError struct {
	string
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
