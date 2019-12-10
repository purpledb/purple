package data

import "encoding/json"

type Set struct {
	items []string
}

func NewSet(items ...string) *Set {
	is := make([]string, 0)

	for _, i := range items {
		is = append(is, i)
	}

	return &Set{
		items: is,
	}
}

func (s *Set) contains(item string) bool {
	for _, i := range s.items {
		if i == item {
			return true
		}
	}

	return false
}

func (s *Set) Add(item string) {
	if !s.contains(item) {
		s.items = append(s.items, item)
	}
}

func (s *Set) Remove(item string) {
	if s.contains(item) {
		for idx, i := range s.items {
			if i == item {
				s.items = append(s.items[:idx], s.items[idx+1:]...)
			}
		}
	}
}

func (s *Set) Get() []string {
	if s.items == nil {
		return []string{}
	} else {
		return s.items
	}
}

func (s *Set) AsBytes() ([]byte, error) {
	return json.Marshal(s.items)
}
