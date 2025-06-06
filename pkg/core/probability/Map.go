package probability

import "math/big"

// maps a probability distribution from outcome type T to U
func Map[T comparable, U comparable](
	dist Distribution[T],
	mapper func(T) U,
	cmp func(U, U) int,
) Distribution[U] {
	out := map[U]*big.Rat{}
	for t, p := range dist.m {
		u := mapper(t)
		outP, ok := out[u]
		if !ok {
			outP = big.NewRat(0, 1)
			out[u] = outP
		}
		outP.Add(outP, p)
	}
	return NewDistributionFunc(out, cmp)
}
