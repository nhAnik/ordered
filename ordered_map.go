package ordered

import (
	"bytes"
	"container/list"
	"encoding"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/buger/jsonparser"
)

type valuePair[V any] struct {
	elem  *list.Element
	value V
}

// KeyValue represents a map elements as a key-value pair.
type KeyValue[K comparable, V any] struct {
	Key   K
	Value V
}

// Map represents an ordered map which is an extension of hashmap.
// Unlike hashmap, the ordered map maintains the insertion order
// i.e. the order in which the keys and their mapped values are
// inserted in the map. The insertion order is not changed if a key
// which already exists in the map is re-inserted.
type Map[K comparable, V any] struct {
	mp    map[K]*valuePair[V]
	items *list.List
}

// NewMap initializes an ordered map.
func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		mp:    make(map[K]*valuePair[V]),
		items: list.New(),
	}
}

// NewMapWithCapacity initializes an ordered map with the given
// initial capacity.
func NewMapWithCapacity[K comparable, V any](capacity int) *Map[K, V] {
	return &Map[K, V]{
		mp:    make(map[K]*valuePair[V], capacity),
		items: list.New(),
	}
}

// NewMapWithKVs initializes an ordered map and inserts the given key-value pair
// in the map.
func NewMapWithKVs[K comparable, V any](kvs ...KeyValue[K, V]) *Map[K, V] {
	om := NewMapWithCapacity[K, V](len(kvs))
	for _, kv := range kvs {
		om.Put(kv.Key, kv.Value)
	}
	return om
}

// Put inserts a key and its mapped value in the map. If the key already exists, the
// mapped value is replaced by the new value.
func (o *Map[K, V]) Put(key K, value V) {
	if _, ok := o.mp[key]; !ok {
		e := o.items.PushBack(key)
		o.mp[key] = &valuePair[V]{elem: e, value: value}
	} else {
		o.mp[key].value = value
	}
}

// Get returns the mapped value for the given key and a bool indicating
// whether the key exists or not.
func (o *Map[K, V]) Get(key K) (V, bool) {
	if val, ok := o.mp[key]; ok {
		return val.value, true
	}
	var dummy V
	return dummy, false
}

// GetOrDefault returns the mapped value for the given key if it exists.
// Otherwise, it returns the default value.
func (o *Map[K, V]) GetOrDefault(key K, defaultValue V) V {
	if val, ok := o.mp[key]; ok {
		return val.value
	}
	return defaultValue
}

// ContainsKey checks if the map contains a mapping for the given key.
func (o *Map[K, V]) ContainsKey(key K) bool {
	_, ok := o.mp[key]
	return ok
}

// Remove removes the key with its mapped value from the map and returns
// the value if the key exists.
func (o *Map[K, V]) Remove(key K) V {
	if vp, ok := o.mp[key]; ok {
		value := vp.value
		o.items.Remove(vp.elem)
		delete(o.mp, key)
		return value
	}
	var dummy V
	return dummy
}

// Len returns the number of elements in the map.
func (o *Map[K, V]) Len() int {
	return o.items.Len()
}

// Keys returns all the keys from the map according to their insertion order.
// The first element of the slice is the oldest key in the map.
func (o *Map[K, V]) Keys() []K {
	keys := make([]K, o.items.Len())
	idx := 0
	for e := o.items.Front(); e != nil; e = e.Next() {
		keys[idx] = e.Value.(K)
		idx++
	}
	return keys
}

// Values returns all the values from the map according to their insertion order.
// The first element of the slice is the oldest value in the map.
func (o *Map[K, V]) Values() []V {
	values := make([]V, o.items.Len())
	idx := 0
	for e := o.items.Front(); e != nil; e = e.Next() {
		key := e.Value.(K)
		if vp, ok := o.mp[key]; ok {
			values[idx] = vp.value
			idx++
		}
	}
	return values
}

// KeyValues returns all the keys and values from the map according to their
// insertion order. The first element of the slice is the oldest key and value
// in the map.
func (o *Map[K, V]) KeyValues() []KeyValue[K, V] {
	kvs := make([]KeyValue[K, V], o.items.Len())
	idx := 0
	for e := o.items.Front(); e != nil; e = e.Next() {
		key := e.Value.(K)
		if vp, ok := o.mp[key]; ok {
			kvs[idx] = KeyValue[K, V]{Key: key, Value: vp.value}
			idx++
		}
	}
	return kvs
}

// ForEach invokes the given function f for each element of the map.
func (o *Map[K, V]) ForEach(f func(K, V)) {
	for _, kv := range o.KeyValues() {
		f(kv.Key, kv.Value)
	}
}

// IsEmpty checks whether the map is empty or not.
func (o *Map[K, V]) IsEmpty() bool {
	return len(o.mp) == 0
}

// Clear removes all the keys and their mapped values from the map.
func (o *Map[K, V]) Clear() {
	for k := range o.mp {
		delete(o.mp, k)
	}
	var next *list.Element
	for e := o.items.Front(); e != nil; e = next {
		next = e.Next()
		o.items.Remove(e)
	}
}

// String returns the string representation of the map.
func (o *Map[K, V]) String() string {
	var sb strings.Builder
	sb.WriteString("map{")
	for idx, kv := range o.KeyValues() {
		if idx > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(kv.Key))
		sb.WriteByte(':')
		sb.WriteString(fmt.Sprint(kv.Value))
	}
	sb.WriteByte('}')
	return sb.String()
}

// MarshalJSON implements json.Marshaler interface.
func (o Map[K, V]) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for idx, kv := range o.KeyValues() {
		if idx > 0 {
			buf.WriteByte(',')
		}
		// key type must either be a string, an integer type, or implement encoding.TextMarshaler
		switch any(kv.Key).(type) {
		case string, encoding.TextMarshaler:
			keyBytes, err := json.Marshal(kv.Key)
			if err != nil {
				return nil, err
			}
			buf.Write(keyBytes)
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			var keyBytes bytes.Buffer
			b, _ := json.Marshal(kv.Key) // marshalling int/uint does not generate error
			keyBytes.WriteByte('"')
			keyBytes.Write(b)
			keyBytes.WriteByte('"')
			buf.Write(keyBytes.Bytes())
		default:
			return nil, errors.New("invalid key type")
		}

		buf.WriteByte(':')
		valBytes, err := json.Marshal(kv.Value)
		if err != nil {
			return nil, err
		}
		buf.Write(valBytes)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (o *Map[K, V]) UnmarshalJSON(b []byte) error {
	if o.items == nil || o.mp == nil {
		o.mp = make(map[K]*valuePair[V])
		o.items = list.New()
	}
	return jsonparser.ObjectEach(b, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		var k K
		// key type must either be a string, an integer type, or implement encoding.TextMarshaler
		switch any(k).(type) {
		case string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, encoding.TextMarshaler:
			if err := json.Unmarshal([]byte(fmt.Sprintf("\"%s\"", string(key))), &k); err != nil {
				return err
			}
		default:
			return errors.New("invalid key type")
		}
		var v V
		var valBytes []byte
		if dataType == jsonparser.String {
			valBytes = []byte(fmt.Sprintf("\"%s\"", string(value)))
		} else {
			valBytes = value
		}
		if err := json.Unmarshal(valBytes, &v); err != nil {
			return err
		}
		o.Put(k, v)
		return nil
	})
}

// GobEncode implements gob.GobEncoder interface.
func (o Map[K, V]) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(o.Len())
	for _, kv := range o.KeyValues() {
		if err := enc.Encode(kv.Key); err != nil {
			return nil, err
		}
		if err := enc.Encode(kv.Value); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

// GobDecode implements gob.GobDecoder interface.
func (o *Map[K, V]) GobDecode(b []byte) error {
	if o.items == nil || o.mp == nil {
		o.mp = make(map[K]*valuePair[V])
		o.items = list.New()
	}
	dec := gob.NewDecoder(bytes.NewBuffer(b))
	len := 0
	dec.Decode(&len)
	for i := 0; i < len; i++ {
		var k K
		var v V
		if err := dec.Decode(&k); err != nil {
			return err
		}
		if err := dec.Decode(&v); err != nil {
			return err
		}
		o.Put(k, v)
	}
	return nil
}
