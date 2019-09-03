package set

type Set interface {
	GetSet(set string) ([]string, error)
	AddToSet(set, item string) ([]string, error)
	RemoveFromSet(set, item string) ([]string, error)
}
