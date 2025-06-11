package util

import (
	"iter"
	"slices"
)

type OrderedMap[K comparable, V any] struct {
	m    map[K]V
	keys []K
}

func NewOrderedMap[K comparable, V any](entries ...Entry[K, V]) *OrderedMap[K, V] {
	om := &OrderedMap[K, V]{
		m: map[K]V{},
	}
	for _, e := range entries {
		om.Put(e.K, e.V)
	}
	return om
}

func (om *OrderedMap[K, V]) Put(k K, v V) {
	if _, exists := om.m[k]; !exists {
		om.keys = append(om.keys, k)
	}
	om.m[k] = v
}

func (om *OrderedMap[K, V]) Del(k K) {
	if _, exists := om.m[k]; exists {
		i := slices.Index(om.keys, k)
		om.keys = slices.Delete(om.keys, i, i+1)
	}
	delete(om.m, k)
}

func (om *OrderedMap[K, V]) Get(k K) (V, bool) {
	v, ok := om.m[k]
	return v, ok
}

func (om *OrderedMap[K, V]) Idx(k K) (int, bool) {
	if _, exists := om.m[k]; !exists {
		return 0, false
	}
	return slices.Index(om.keys, k), true
}

func (om *OrderedMap[K, V]) Len() int {
	return len(om.keys)
}

func (om *OrderedMap[K, V]) MustGet(k K) V {
	return om.m[k]
}

func (om *OrderedMap[K, V]) Iter() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, k := range om.keys {
			if !yield(k, om.m[k]) {
				return
			}
		}
	}
}
