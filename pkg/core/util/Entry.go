package util

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
