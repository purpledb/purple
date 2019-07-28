package strato

type Set interface {
	GetSet(set string) []string
	AddToSet(set, item string)
	RemoveFromSet(set, item string)
}
