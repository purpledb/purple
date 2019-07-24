package kv

import "fmt"

type NotFoundError struct {
	Location Location
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("no value found for %s", e.Location.String())
}

func NotFound(location Location) NotFoundError {
	return NotFoundError{
		Location: location,
	}
}

func IsNotFound(err error) bool {
	_, ok := err.(NotFoundError)
	return ok
}