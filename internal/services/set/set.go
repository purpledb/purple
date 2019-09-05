package set

type Set interface {
	SetGet(set string) ([]string, error)
	SetAdd(set, item string) ([]string, error)
	SetRemove(set, item string) ([]string, error)
}
