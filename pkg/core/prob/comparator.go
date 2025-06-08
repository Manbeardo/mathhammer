package prob

import (
	"cmp"
	"reflect"
)

type Comparator[T any] = func(a, b T) int

func convertingComparatorFor[Underlying cmp.Ordered](t reflect.Type) Comparator[any] {
	u := reflect.TypeFor[Underlying]()
	var convert func(any) Underlying
	if t == u {
		convert = func(a any) Underlying {
			return a.(Underlying)
		}
	} else {
		convert = func(a any) Underlying {
			return reflect.ValueOf(a).Convert(u).Interface().(Underlying)
		}
	}
	return func(a, b any) int {
		at, bt := convert(a), convert(b)
		return cmp.Compare(at, bt)
	}
}

func unsafeComparatorFor[T any]() Comparator[any] {
	t := reflect.TypeFor[T]()
	switch t.Kind() {
	case reflect.Int:
		return convertingComparatorFor[int](t)
	case reflect.Int8:
		return convertingComparatorFor[int8](t)
	case reflect.Int16:
		return convertingComparatorFor[int16](t)
	case reflect.Int32:
		return convertingComparatorFor[int32](t)
	case reflect.Int64:
		return convertingComparatorFor[int64](t)
	case reflect.Uint:
		return convertingComparatorFor[uint](t)
	case reflect.Uint8:
		return convertingComparatorFor[uint8](t)
	case reflect.Uint16:
		return convertingComparatorFor[uint16](t)
	case reflect.Uint32:
		return convertingComparatorFor[uint32](t)
	case reflect.Uint64:
		return convertingComparatorFor[uint64](t)
	case reflect.Uintptr:
		return convertingComparatorFor[uintptr](t)
	case reflect.Float32:
		return convertingComparatorFor[float32](t)
	case reflect.Float64:
		return convertingComparatorFor[float64](t)
	case reflect.String:
		return convertingComparatorFor[string](t)
	}
	return nil
}
