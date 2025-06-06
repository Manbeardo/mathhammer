package prob

import "math/big"

// reduces probability distributions of multiple independent events into an aggregate
// probability distribution
func Reduce[T comparable, U comparable](
	dists []Dist[T],
	collector func(U, T) U,
	cmp func(U, U) int,
	initialValue U,
) Dist[U] {
	prev := map[U]*big.Rat{
		(initialValue): big.NewRat(1, 1),
	}
	for _, dist := range dists {
		next := map[U]*big.Rat{}
		for t, p := range dist.m {
			for prevU, prevP := range prev {
				partialP := big.NewRat(0, 1)
				partialP.Mul(p, prevP)
				nextU := collector(prevU, t)
				nextP, ok := next[nextU]
				if !ok {
					nextP = big.NewRat(0, 1)
					next[nextU] = nextP
				}
				nextP.Add(nextP, partialP)
			}
		}
		prev = next
	}
	return NewDistributionFunc(prev, cmp)
}
