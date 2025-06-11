package util

import "iter"

type Entry[K, V any] struct {
	K K
	V V
}

func EntriesFromSeq2[K, V any](seq iter.Seq2[K, V]) []Entry[K, V] {
	es := []Entry[K, V]{}
	for k, v := range seq {
		es = append(es, Entry[K, V]{k, v})
	}
	return es
}

func EntriesFromMap[K comparable, V any](m map[K]V) []Entry[K, V] {
	es := []Entry[K, V]{}
	for k, v := range m {
		es = append(es, Entry[K, V]{k, v})
	}
	return es
}
