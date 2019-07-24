package memory

import "strato/kv"

type Memory struct {
	values map[*kv.Location]*kv.Value
}

var _ kv.KV = (*Memory)(nil)

func (m *Memory) Get(location *kv.Location) (*kv.Value, error) {
	val, ok := m.values[location]

	if !ok {
		return nil, kv.NotFound(location)
	}

	return val, nil
}

func (m *Memory) Put(location *kv.Location, value *kv.Value) error {
	m.values[location] = value

	return nil
}

func (m *Memory) Delete(location *kv.Location) error {
	delete(m.values, location)

	return nil
}
