package flag

type Flag interface {
	FlagGet(key string) (bool, error)
	FlagSet(key string) error
	FlagUnset(key string) error
}
