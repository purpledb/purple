package strato

type Cache interface {
	Get(key string) (string, error)
	Put(key, value string) error
}
