package strato

type Set interface {
	GetSet(set string) ([]string, error)
	AddToSet(set, item string) error
	RemoveFromSet(set, item string) error
}
