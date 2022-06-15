package safemap

import (
	"bytes"
	"encoding/gob"
	"hash/crc32"
	"hash/fnv"
	"sync"
)

const (
	defaultShard = 64
)

//New
func New[K comparable, V any](shard ...uint32) *SafeMap[K, V] {
	var shardNum uint32 = defaultShard
	if len(shard) != 0 && shard[0] > 0 {
		shardNum = shard[0]
	}
	var sm = &SafeMap[K, V]{
		locks: make([]sync.RWMutex, shardNum),
		maps:  make([]map[K]V, shardNum),
		shard: shardNum,
	}
	for i := range sm.maps {
		sm.maps[i] = make(map[K]V)
	}
	return sm
}

type SafeMap[K comparable, V any] struct {
	locks []sync.RWMutex
	maps  []map[K]V
	shard uint32
}

//Get returns the value associated with key.
func (s *SafeMap[K, V]) Get(key K) (V, bool) {
	idx := s.idx(key)
	s.locks[idx].RLock()
	v, ok := s.maps[idx][key]
	s.locks[idx].RUnlock()
	if ok {
		return v, ok
	}
	return *new(V), false
}

//GetOrSet returns the value for the key if present, otherwise sets and returns the default value
func (s *SafeMap[K, V]) GetOrSet(key K, value V) V {
	idx := s.idx(key)
	s.locks[idx].Lock()
	v, ok := s.maps[idx][key]
	if !ok {
		s.maps[idx][key] = value
	}
	s.locks[idx].Unlock()
	return v
}

//Set sets the value for key.
func (s *SafeMap[K, V]) Set(key K, value V) {
	idx := s.idx(key)
	s.locks[idx].Lock()
	s.maps[idx][key] = value
	s.locks[idx].Unlock()
}

//Del deletes the value for key.
func (s *SafeMap[K, V]) Del(key K) {
	idx := s.idx(key)
	s.locks[idx].Lock()
	delete(s.maps[idx], key)
	s.locks[idx].Unlock()
}

//Reset deletes all values.
func (s *SafeMap[K, V]) Reset() {
	for i := uint32(0); i < s.shard; i++ {
		s.locks[i].Lock()
		s.maps[i] = make(map[K]V)
		s.locks[i].Unlock()
	}
}

//Len
func (s *SafeMap[K, V]) Len() int {
	var n int
	for i := uint32(0); i < s.shard; i++ {
		s.locks[i].RLock()
		n += len(s.maps[i])
		s.locks[i].RUnlock()
	}
	return n
}

//Range calls f sequentially for each key and value present in the map. If f returns false, range stops the iteration.
func (s *SafeMap[K, V]) Range(f func(K, V) bool) {
	for i := uint32(0); i < s.shard; i++ {
		s.locks[i].RLock()
		for k, v := range s.maps[i] {
			if !f(k, v) {
				s.locks[i].RUnlock()
				return
			}
		}
		s.locks[i].RUnlock()
	}
}

//idx
func (s SafeMap[K, V]) idx(key K) uint32 {
	return keyid(key) % s.shard
}

var (
	encpool = sync.Pool{
		New: func() interface{} {
			b := new(bytes.Buffer)
			return &gobEncoder{
				buf: b,
				enc: gob.NewEncoder(b),
			}
		},
	}
)

func keyid[T comparable](o T) uint32 {
	switch val := interface{}(o).(type) {
	case int:
		return uint32(val)
	case int8:
		return uint32(val)
	case int16:
		return uint32(val)
	case int32:
		return uint32(val)
	case int64:
		return uint32(val)
	case uint:
		return uint32(val)
	case uint8:
		return uint32(val)
	case uint16:
		return uint32(val)
	case uint32:
		return uint32(val)
	case uint64:
		return uint32(val)
	case float32:
		return uint32(val)
	case float64:
		return uint32(val)
	case string:
		return uint32(crc32.ChecksumIEEE([]byte(val)))
	case []byte:
		return uint32(crc32.ChecksumIEEE(val))
	case bool:
		if val {
			return 1
		}
		return 0
	default:
		if v, ok := val.(KeyID); ok {
			return v.ToUint32()
		}
	}
	enc := encpool.Get().(*gobEncoder)
	defer enc.Free()
	return enc.Hash32(o)
}

type KeyID interface {
	ToUint32() uint32
}

type gobEncoder struct {
	enc *gob.Encoder
	buf *bytes.Buffer
}

func (e *gobEncoder) Hash32(o any) uint32 {
	e.enc.Encode(o)
	h := fnv.New32()
	h.Write(e.buf.Bytes())
	return h.Sum32()
}

//Free
func (e *gobEncoder) Free() {
	e.buf.Reset()
	encpool.Put(e)
}
