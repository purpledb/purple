package kv

import "fmt"

type NotFoundError struct {
	location *Location
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("no value found for %s", e.location.String())
}

func NotFound(location *Location) NotFoundError {
	return NotFoundError{
		location: location,
	}
}

func IsNotFound(err error) bool {
	_, ok := err.(NotFoundError)
	return ok
}