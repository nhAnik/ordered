package ordered

import (
	"container/list"
)

type valuePair[V any] struct {
	elem  *list.Element
	value V
}

type KeyValue[K comparable, V any] struct {
	Key   K
	Value V
}

type Map[K comparable, V any] struct {
	mp    map[K]*valuePair[V]
	items *list.List
}

func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		mp:    make(map[K]*valuePair[V], 0),
		items: list.New(),
	}
}

func NewMapWithElems[K comparable, V any](kvs ...KeyValue[K, V]) *Map[K, V] {
	om := NewMap[K, V]()
	for _, kv := range kvs {
		om.Put(kv.Key, kv.Value)
	}
	return om
}

func (o *Map[K, V]) Put(key K, value V) {
	if _, ok := o.mp[key]; !ok {
		e := o.items.PushBack(key)
		o.mp[key] = &valuePair[V]{elem: e, value: value}
	} else {
		o.mp[key].value = value
	}
}

func (o *Map[K, V]) Get(key K) (V, bool) {
	val, ok := o.mp[key]
	if ok {
		return val.value, true
	}
	var dummy V
	return dummy, false
}

func (o *Map[K, V]) ContainsKey(key K) bool {
	_, ok := o.mp[key]
	return ok
}

func (o *Map[K, V]) Remove(key K) {
	if vp, ok := o.mp[key]; ok {
		o.items.Remove(vp.elem)
		delete(o.mp, key)
	}
}

func (o *Map[K, V]) Len() int {
	return o.items.Len()
}

func (o *Map[K, V]) Keys() []K {
	keys := make([]K, o.items.Len())
	idx := 0
	for e := o.items.Front(); e != nil; e = e.Next() {
		keys[idx] = e.Value.(K)
		idx++
	}
	return keys
}

func (o *Map[K, V]) Values() []V {
	values := make([]V, o.items.Len())
	idx := 0
	for e := o.items.Front(); e != nil; e = e.Next() {
		key := e.Value.(K)
		values[idx] = o.mp[key].value
		idx++
	}
	return values
}

func (o *Map[K, V]) KeyValues() []KeyValue[K, V] {
	kvs := make([]KeyValue[K, V], o.items.Len())
	idx := 0
	for e := o.items.Front(); e != nil; e = e.Next() {
		key := e.Value.(K)
		value := o.mp[key].value
		kvs[idx] = KeyValue[K, V]{Key: key, Value: value}
		idx++
	}
	return kvs
}

func (o *Map[K, V]) IsEmpty() bool {
	return len(o.mp) == 0
}

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
