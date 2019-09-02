package strato

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrExpired              = CacheError{"item has expired"}
	ErrNoCacheItem          = CacheError{"no item found"}
	ErrNoCacheKey           = CacheError{"no cache key specified"}
	ErrNoCacheValue         = CacheError{"no cache value specified"}
	ErrNoBucket             = KVError{"no bucket specified"}
	ErrNoKey                = KVError{"no key specified"}
	ErrNoLocation           = KVError{"no location specified"}
	ErrNoValue              = KVError{"no value specified"}
	ErrNoSet                = SetError{"set does not exist"}
	ErrNoAddress            = ConfigError{"no server address provided"}
	ErrNoPort               = ConfigError{"no server port supplied"}
	ErrPortOutOfRange       = ConfigError{"port must be between 1024 and 49151"}
	ErrBackendNotRecognized = ConfigError{"backend key not recognized"}
	ErrNoBackend            = ConfigError{"no backend specified"}
)

type (
	CacheError struct {
		string
	}

	KVError struct {
		string
	}

	ConfigError struct {
		string
	}

	NotFoundError struct {
		location *Location
	}

	SetError struct {
		string
	}
)

func (e CacheError) Error() string {
	return fmt.Sprintf("cache error: %s", e.string)
}

func IsNoCacheKey(err error) bool {
	return err == ErrNoCacheKey
}

func IsNoCacheValue(err error) bool {
	return err == ErrNoCacheValue
}

func IsExpired(err error) bool {
	return err == ErrExpired
}

func IsNoItemFound(err error) bool {
	return err == ErrNoCacheItem
}

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

func (e SetError) Error() string {
	return fmt.Sprintf("set error: %s", e.string)
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf(`no value found for %s`, e.location.String())
}

func (e NotFoundError) AsProtoStatus() error {
	return status.Error(codes.NotFound, e.Error())
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
