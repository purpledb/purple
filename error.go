package strato

import "fmt"

var (
	ErrNoLocation     = KVError{"no location specified"}
	ErrNoValue        = KVError{"no value specified"}
	ErrNoAddress      = ConfigError{"no server address provided"}
	ErrNoPort         = ConfigError{"no server port supplied"}
	ErrPortOutOfRange = ConfigError{"port must be between 1024 and 49151"}
)

type (
	KVError struct {
		string
	}

	ConfigError struct {
		string
	}

	NotFoundError struct {
		location *Location
	}
)

func (e KVError) Error() string {
	return fmt.Sprintf("KV error: %s", e.string)
}

func IsKVError(err error) bool {
	_, ok := err.(KVError)
	return ok
}

func (e ConfigError) Error() string {
	return fmt.Sprintf("config error: %s", e.string)
}

func IsConfigError(err error) bool {
	_, ok := err.(ConfigError)
	return ok
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf(`no value found for %s`, e.location.String())
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
