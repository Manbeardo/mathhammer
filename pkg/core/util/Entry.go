package util

import "slices"

type Entry[T any, U any] struct {
	Key   T
	Value U
}

func Entries[T comparable, U any](m map[T]U) []Entry[T, U] {
	es := []Entry[T, U]{}
	for k, v := range m {
		es = append(es, Entry[T, U]{k, v})
	}
	return es
}

func OrderedEntries[T comparable, U any](m map[T]U, fn func(T, T) int) []Entry[T, U] {
	es := Entries(m)
	slices.SortFunc(es, func(a, b Entry[T, U]) int {
		return fn(a.Key, b.Key)
	})
	return es
}
