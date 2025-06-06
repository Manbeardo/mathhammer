package probability

import (
	"cmp"
	"fmt"
	"iter"
	"maps"
	"math/big"
	"slices"

	"github.com/Manbeardo/mathhammer/pkg/core/util"
)

type Distribution[T comparable] struct {
	m        map[T]*big.Rat
	outcomes []T
}

func NewDistribution[T cmp.Ordered](m map[T]*big.Rat) Distribution[T] {
	return NewDistributionFunc(m, cmp.Compare)
}

func NewDistributionFunc[T comparable](m map[T]*big.Rat, cmp func(T, T) int) Distribution[T] {
	var psum big.Rat
	for _, p := range m {
		psum.Add(&psum, p)
	}
	if f, ok := psum.Float64(); !ok || f != 1 {
		panic("sum of all probabilities must be 1")
	}

	d := Distribution[T]{
		m:        maps.Clone(m),
		outcomes: slices.Collect(maps.Keys(m)),
	}
	slices.SortFunc(d.outcomes, cmp)
	return d
}

func (d Distribution[T]) Format(w fmt.State, v rune) {
	util.PrettyFormat(w, v, d)
}

func (d Distribution[T]) Iter() iter.Seq2[T, *big.Rat] {
	return func(yield func(T, *big.Rat) bool) {
		for _, k := range d.outcomes {
			if !yield(k, d.m[k]) {
				return
			}
		}
	}
}
