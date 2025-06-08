// package prob defines various types and functions for manipulating discrete
// probability distributions
package prob

import (
	"cmp"
	"fmt"
	"iter"
	"maps"
	"math/big"
	"reflect"
	"slices"

	"github.com/Manbeardo/mathhammer/pkg/core/util"
)

type StringKeyer interface {
	StringKey() string
}

type MapT[T comparable] = map[T]*big.Rat
type EntryT[T any] = util.Entry[T, *big.Rat]

type keyFuncKind int

type Key string

const (
	keyFuncInvalidFirst keyFuncKind = iota // invalid value bounding the start of the valid range
	keyFuncInterface
	keyFuncComparable
	keyFuncInvalidLast // invalid value bounding the end of the valid range
)

// a discrete probability distribution
type Dist[T any] struct {
	vmap        map[Key]T
	pmap        map[Key]*big.Rat
	keyFuncKind keyFuncKind
}

func empty[T any]() (Dist[T], error) {
	keyFunc, err := keyFuncKindFor[T]()
	if err != nil {
		return Dist[T]{}, err
	}

	return Dist[T]{
		vmap:        map[Key]T{},
		pmap:        map[Key]*big.Rat{},
		keyFuncKind: keyFunc,
	}, nil
}

func FromEntries[T any](m []EntryT[T]) (Dist[T], error) {
	d, err := empty[T]()
	if err != nil {
		return d, err
	}

	for _, e := range m {
		v, p := e.Key, e.Value
		k := d.key(v)
		if _, exists := d.vmap[k]; exists {
			return d, fmt.Errorf("multiple values with same key: %s", k)
		}
		d.vmap[k] = v
		d.pmap[k] = p
	}

	return d.validate()
}

func FromMap[T comparable](m MapT[T]) (Dist[T], error) {
	return FromEntries(util.Entries(m))
}

// FromConst returns a distribution whose sole outcome is v
func FromConst[T any](v T) (Dist[T], error) {
	return FromEntries(
		[]EntryT[T]{
			{Key: v, Value: big.NewRat(1, 1)},
		},
	)
}

func keyFuncKindFor[T any]() (keyFuncKind, error) {
	t := reflect.TypeFor[T]()
	if t.AssignableTo(reflect.TypeFor[StringKeyer]()) {
		return keyFuncInterface, nil
	}
	if t.Comparable() {
		return keyFuncComparable, nil
	}
	return keyFuncInvalidFirst, fmt.Errorf("%s does not satisfy comparable or StringKeyer", t.String())
}

func (d Dist[T]) validate() (Dist[T], error) {
	var psum big.Rat

	if d.keyFuncKind <= keyFuncInvalidFirst || d.keyFuncKind >= keyFuncInvalidLast {
		return d, fmt.Errorf("invalid keyFuncKind: %d", d.keyFuncKind)
	}

	for _, p := range d.pmap {
		psum.Add(&psum, p)
	}

	if f, ok := psum.Float64(); !ok || f != 1 {
		return d, fmt.Errorf("sum of all probabilities must be 1 (is: %f)", f)
	}

	return d, nil
}

func (d Dist[T]) key(v any) Key {
	switch d.keyFuncKind {
	case keyFuncInterface:
		return Key(v.(StringKeyer).StringKey())
	case keyFuncComparable:
		return Key(fmt.Sprintf("%v", v))
	default:
		panic(fmt.Errorf("invalid keyFuncKind: %d", d.keyFuncKind))
	}
}

func (d Dist[T]) vmapComparator() Comparator[util.Entry[Key, T]] {
	if d.keyFuncKind == keyFuncComparable {
		vcmp := unsafeComparatorFor[T]()
		if vcmp != nil {
			return func(a, b util.Entry[Key, T]) int {
				return vcmp(a.Value, b.Value)
			}
		}
	}
	return func(a, b util.Entry[Key, T]) int {
		return cmp.Compare(a.Key, b.Key)
	}
}

func (d Dist[T]) Keys() []Key {
	entries := util.Entries(d.vmap)
	slices.SortFunc(entries, d.vmapComparator())
	keys := []Key{}
	for _, e := range entries {
		keys = append(keys, e.Key)
	}
	return keys
}

func (d Dist[T]) Lookup(k Key) (EntryT[T], bool) {
	p, ok := d.pmap[k]
	if !ok {
		return EntryT[T]{}, false
	}
	v := d.vmap[k]
	return EntryT[T]{Key: v, Value: p}, true
}

func (d Dist[T]) Format(w fmt.State, v rune) {
	util.PrettyFormat(w, v, d)
}

func (d Dist[T]) Iter() iter.Seq2[T, *big.Rat] {
	return func(yield func(T, *big.Rat) bool) {
		for _, k := range d.Keys() {
			if !yield(d.vmap[k], d.pmap[k]) {
				return
			}
		}
	}
}

func (d Dist[T]) Distribution() Dist[T] {
	return d
}

func (d Dist[T]) Percentile(p float64, cmp func(T, T) int) T {
	values := slices.Collect(maps.Values(d.vmap))
	slices.SortFunc(values, cmp)

	sum := big.NewRat(0, 1)
	for _, v := range values {
		sum = sum.Add(sum, d.pmap[d.key(v)])
		sumf, _ := sum.Float64()
		if sumf > p {
			return v
		}
	}
	return values[len(values)-1]
}

func (d Dist[T]) Median(cmp func(T, T) int) T {
	return d.Percentile(0.5, cmp)
}

func (d Dist[T]) StringKey() string {
	// vmap and keyFuncKind are implementation details that don't matter
	// for equality checks when the value type's key function works correctly
	return fmt.Sprintf("%#v", util.PrettyMap[Key, *big.Rat](d.pmap))
}
