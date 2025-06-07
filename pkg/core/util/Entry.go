package util

type EntryT[T comparable, U any] struct {
	Key   T
	Value U
}

func Entry[T comparable, U any](key T, value U) EntryT[T, U] {
	return EntryT[T, U]{
		Key:   key,
		Value: value,
	}
}
