package prob

import "math/big"

// FlatMap maps each outcome of a probability distribution of outcome T
// to a new probability distribution with outcome U and flattens the
// results into a single distribution of outcome U
func FlatMap[T comparable, U comparable](
	dist Dist[T],
	mapper func(T) Dist[U],
	cmp func(U, U) int,
) Dist[U] {
	out := map[U]*big.Rat{}
	for t, tp := range dist.m {
		uDist := mapper(t)
		for u, up := range uDist.m {
			outP, ok := out[u]
			if !ok {
				outP = big.NewRat(0, 1)
				out[u] = outP
			}
			p := &big.Rat{}
			p.Mul(tp, up)
			outP.Add(outP, p)
		}
	}
	return NewDistributionFunc(out, cmp)
}
