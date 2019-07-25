package strato

import "fmt"

var (
	ErrNoAddress = ConfigError{"no Strato server address provided"}
	ErrNoPort         = ConfigError{"no server port supplied"}
	ErrPortOutOfRange = ConfigError{"port must be between 1024 and 49151"}
)

type (
	ConfigError struct {
		string
	}

	NotFoundError struct {
		location *Location
	}
)

func (e ConfigError) Error() string {
	return fmt.Sprintf("server config error: %s", e.string)
}

func IsConfigError(err error) bool {
	_, ok := err.(ConfigError)
	return ok
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
