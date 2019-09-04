package oops

import (
	"errors"
	"fmt"
)

var (
	ErrNoKey   = errors.New("no key provided")
	ErrNoValue = errors.New("no value provided")
)

type notFoundError struct {
	key string
}

func (e notFoundError) Error() string {
	return fmt.Sprintf("item with key %s not found", e.key)
}

func NotFound(key string) error {
	return notFoundError{
		key: key,
	}
}

func IsNotFound(err error) bool {
	_, ok := err.(notFoundError)
	return ok
}
