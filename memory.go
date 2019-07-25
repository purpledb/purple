package strato

type Memory struct {
	values map[Location]*Value
}

var _ KV = (*Memory)(nil)

func New() *Memory {
	values := make(map[Location]*Value)

	return &Memory{
		values: values,
	}
}

func (m *Memory) Get(location *Location) (*Value, error) {
	val, ok := m.values[*location]

	if !ok {
		return nil, NotFound(location)
	}

	return val, nil
}

func (m *Memory) Put(location *Location, value *Value) error {
	m.values[*location] = value

	return nil
}

func (m *Memory) Delete(location *Location) error {
	if _, ok := m.values[*location]; ok {
		delete(m.values, *location)
	}

	return nil
}
