// package prob defines various types and functions for manipulating discrete
// probability distributions
package prob

import (
	"fmt"
	"iter"
	"maps"
	"math/big"
	"reflect"
	"slices"
	"strings"

	"github.com/Manbeardo/mathhammer/pkg/core/util"
)

type StringKeyer interface {
	StringKey() string
}

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

func FromEntries[T any](m []util.Entry[T, *big.Rat]) (Dist[T], error) {
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

func FromMap[T comparable](m map[T]*big.Rat) (Dist[T], error) {
	return FromEntries(util.Entries(m))
}

// FromConst returns a distribution whose sole outcome is v
func FromConst[T any](v T) (Dist[T], error) {
	return FromEntries(
		[]util.Entry[T, *big.Rat]{
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

func (d Dist[T]) Keys() []Key {
	keys := slices.Collect(maps.Keys(d.pmap))
	slices.Sort(keys)
	return keys
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

func (d Dist[T]) Percentile(p *big.Rat, cmp func(T, T) int) T {
	values := slices.Collect(maps.Values(d.vmap))
	slices.SortFunc(values, cmp)

	sum := big.NewRat(0, 1)
	for _, v := range values {
		sum = sum.Add(sum, d.pmap[d.key(v)])
		if sum.Cmp(p) >= 0 {
			return v
		}
	}
	return values[len(values)-1]
}

func (d Dist[T]) Median(cmp func(T, T) int) T {
	return d.Percentile(big.NewRat(1, 2), cmp)
}

func (d Dist[T]) StringKey() string {
	keys := slices.Collect(maps.Keys(d.pmap))
	slices.Sort(keys)

	b := strings.Builder{}
	b.WriteString("{")
	for i, k := range keys {
		if i > 0 {
			b.WriteString(", ")
		}
		p := d.pmap[k]
		fmt.Fprintf(&b, "%s: %s", k, p.RatString())
	}
	b.WriteString("}")
	return b.String()
}
