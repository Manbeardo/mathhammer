// package prob defines various types and functions for manipulating discrete
// probability distributions
package prob

import (
	"cmp"
	"fmt"
	"iter"
	"math/big"
	"reflect"
	"slices"

	"github.com/Manbeardo/mathhammer/pkg/core/util"
)

type StringKeyer interface {
	StringKey() string
}

type keyFuncKind int

const (
	keyFuncInvalid keyFuncKind = iota
	keyFuncInterface
	keyFuncComparable
)

// a discrete probability distribution
type Dist[T any] struct {
	vmap        map[string]T
	pmap        map[string]*big.Rat
	sorted      []T
	keyFuncKind keyFuncKind
}

func empty[T any]() (Dist[T], error) {
	keyFunc, err := keyFuncKindFor[T]()
	if err != nil {
		return Dist[T]{}, err
	}

	return Dist[T]{
		vmap:        map[string]T{},
		pmap:        map[string]*big.Rat{},
		keyFuncKind: keyFunc,
	}, nil
}

func NewDist[T cmp.Ordered](m []util.Entry[T, *big.Rat]) (Dist[T], error) {
	return NewDistFunc(m, cmp.Compare)
}

func NewDistFunc[T any](m []util.Entry[T, *big.Rat], cmp func(T, T) int) (Dist[T], error) {
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

	return d.finalize(cmp)
}

func FromMap[T cmp.Ordered](m map[T]*big.Rat) (Dist[T], error) {
	return FromMapFunc(m, cmp.Compare)
}

func FromMapFunc[T comparable](m map[T]*big.Rat, cmp func(T, T) int) (Dist[T], error) {
	return NewDistFunc(util.OrderedEntries(m, cmp), cmp)
}

// NewConstDist returns a distribution whose sole outcome is v
func NewConstDist[T any](v T) (Dist[T], error) {
	return NewDistFunc(
		[]util.Entry[T, *big.Rat]{
			{Key: v, Value: big.NewRat(1, 1)},
		},
		func(T, T) int { return 0 },
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
	return keyFuncInvalid, fmt.Errorf("%s does not satisfy comparable or StringKeyer", t.String())
}

func (d Dist[T]) finalize(cmp func(T, T) int) (Dist[T], error) {
	d.sorted = nil
	var psum big.Rat
	for k, v := range d.vmap {
		p := d.pmap[k]
		d.sorted = append(d.sorted, v)
		psum.Add(&psum, p)
	}
	slices.SortFunc(d.sorted, cmp)

	if f, ok := psum.Float64(); !ok || f != 1 {
		return d, fmt.Errorf("sum of all probabilities must be 1 (is: %f)", f)
	}

	return d, nil
}

func (d Dist[T]) key(v any) string {
	switch d.keyFuncKind {
	case keyFuncInterface:
		return v.(StringKeyer).StringKey()
	case keyFuncComparable:
		return fmt.Sprintf("%v", v)
	default:
		panic(fmt.Errorf("invalid keyFuncKind: %d", d.keyFuncKind))
	}
}

func (d Dist[T]) Format(w fmt.State, v rune) {
	util.PrettyFormat(w, v, d)
}

func (d Dist[T]) Iter() iter.Seq2[T, *big.Rat] {
	return func(yield func(T, *big.Rat) bool) {
		for _, v := range d.sorted {
			if !yield(v, d.pmap[d.key(v)]) {
				return
			}
		}
	}
}

func (d Dist[T]) Distribution() Dist[T] {
	return d
}

func (d Dist[T]) Percentile(p *big.Rat) T {
	// d.outcomes is already sorted, so we just pick the first
	// value whose cumulative probability is >= p
	sum := big.NewRat(0, 1)
	for _, v := range d.sorted {
		sum = sum.Add(sum, d.pmap[d.key(v)])
		if sum.Cmp(p) >= 0 {
			return v
		}
	}
	return d.sorted[len(d.sorted)-1]
}

func (d Dist[T]) Median() T {
	return d.Percentile(big.NewRat(1, 2))
}
