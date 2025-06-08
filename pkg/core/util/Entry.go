package util

type Entry[K any, V any] struct {
	K K
	V V
}

func Entries[K comparable, V any](m map[K]V) []Entry[K, V] {
	es := []Entry[K, V]{}
	for k, v := range m {
		es = append(es, Entry[K, V]{k, v})
	}
	return es
}
