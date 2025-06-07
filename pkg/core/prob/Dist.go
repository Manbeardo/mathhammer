// package prob defines various types and functions for manipulating discrete
// probability distributions
package prob

import (
	"cmp"
	"fmt"
	"iter"
	"maps"
	"math/big"
	"slices"

	"github.com/Manbeardo/mathhammer/pkg/core/util"
)

// a discrete probability distribution
type Dist[T comparable] struct {
	m        map[T]*big.Rat
	outcomes []T
}

func NewDistribution[T cmp.Ordered](m map[T]*big.Rat) Dist[T] {
	return NewDistributionFunc(m, cmp.Compare)
}

func NewDistributionFunc[T comparable](m map[T]*big.Rat, cmp func(T, T) int) Dist[T] {
	var psum big.Rat
	for _, p := range m {
		psum.Add(&psum, p)
	}
	if f, ok := psum.Float64(); !ok || f != 1 {
		panic("sum of all probabilities must be 1")
	}

	d := Dist[T]{
		m:        maps.Clone(m),
		outcomes: slices.Collect(maps.Keys(m)),
	}
	slices.SortFunc(d.outcomes, cmp)
	return d
}

func (d Dist[T]) Format(w fmt.State, v rune) {
	util.PrettyFormat(w, v, d)
}

func (d Dist[T]) Iter() iter.Seq2[T, *big.Rat] {
	return func(yield func(T, *big.Rat) bool) {
		for _, k := range d.outcomes {
			if !yield(k, d.m[k]) {
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
	for _, t := range d.outcomes {
		sum = sum.Add(sum, d.m[t])
		if sum.Cmp(p) >= 0 {
			return t
		}
	}
	return d.outcomes[len(d.outcomes)-1]
}

func (d Dist[T]) Median() T {
	return d.Percentile(big.NewRat(1, 2))
}
