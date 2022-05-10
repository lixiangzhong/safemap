package safemap

import "sync"

type SafeMap[K comparable, V any] struct {
	m *sync.Map
}

func New[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{m: new(sync.Map)}
}

// Get returns the value for the given key.
func (m *SafeMap[K, V]) Get(key K) (V, bool) {
	val, ok := m.m.Load(key)
	if !ok {
		return *new(V), false
	}
	return val.(V), true
}

// Set sets the value for the given key.
func (m *SafeMap[K, V]) Set(key K, value V) {
	m.m.Store(key, value)
}

// Del deletes the value for the given key.
func (m *SafeMap[K, V]) Del(key K) {
	m.m.Delete(key)
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
func (m *SafeMap[K, V]) Range(f func(K, V) bool) {
	m.m.Range(func(key, value interface{}) bool {
		return f(key.(K), value.(V))
	})
}
