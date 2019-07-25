package strato

import "strings"

type Memory struct {
	values map[Location]*Value
	docs   []*Document
}

var (
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

func (m *Memory) Get(location *Location) (*Value, error) {
	val, ok := m.values[*location]

	if !ok {
		return nil, NotFound(location)
	}

	return val, nil
}

func (m *Memory) Put(location *Location, value *Value) {
	m.values[*location] = value
}

func (m *Memory) Delete(location *Location) {
	if _, ok := m.values[*location]; ok {
		delete(m.values, *location)
	}
}

func (m *Memory) Index(doc *Document) {
	cnt := strings.ToLower(doc.Content)

	m.docs = append(m.docs, &Document{ID: doc.ID, Content: cnt})
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
