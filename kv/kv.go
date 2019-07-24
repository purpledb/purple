package kv

type KV interface {
	Get(location *Location) (*Value, error)
	Put(location *Location, value *Value) error
	Delete(location *Location) error
}
