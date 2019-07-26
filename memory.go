package strato

import "strings"

type Memory struct {
	values map[Location]*Value
	docs   []*Document
}

var (
	//_ Cache  = (*Memory)(nil)
	_ KV     = (*Memory)(nil)
	_ Search = (*Memory)(nil)
)

func New() *Memory {
	values := make(map[Location]*Value)

	docs := make([]*Document, 0)

	return &Memory{
		values: values,
		docs:   docs,
	}
}

func (m *Memory) KVGet(location *Location) (*Value, error) {
	val, ok := m.values[*location]

	if !ok {
		return nil, NotFound(location)
	}

	return val, nil
}

func (m *Memory) KVPut(location *Location, value *Value) {
	m.values[*location] = value
}

func (m *Memory) KVDelete(location *Location) {
	if _, ok := m.values[*location]; ok {
		delete(m.values, *location)
	}
}

func (m *Memory) Index(doc *Document) {
	doc = doc.prepare()

	m.docs = append(m.docs, doc)
}

func (m *Memory) Query(q string) []*Document {
	strings.ToLower(q)

	docs := make([]*Document, 0)

	if len(m.docs) == 0 {
		return nil
	}

	for _, doc := range m.docs {
		if strings.Contains(doc.Content, q) {
			docs = append(docs, doc)
		}
	}

	return docs
}
