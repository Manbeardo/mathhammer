package prob

import (
	"maps"
	"math/big"
)

// FlatMap maps each outcome of a probability distribution of outcome T
// to a new probability distribution with outcome U and flattens the
// results into a single distribution of outcome U
func FlatMap[T any, U any](
	dist Dist[T],
	mapper func(T) Dist[U],
) (Dist[U], error) {
	out, err := empty[U]()
	if err != nil {
		return out, err
	}

	for tk, tv := range dist.vmap {
		tp := dist.pmap[tk]
		uDist := mapper(tv)
		maps.Copy(out.vmap, uDist.vmap)
		for uk, up := range uDist.pmap {
			outP, ok := out.pmap[uk]
			if !ok {
				outP = big.NewRat(0, 1)
				out.pmap[uk] = outP
			}
			p := &big.Rat{}
			p.Mul(tp, up)
			outP.Add(outP, p)
		}
	}

	return out.validate()
}
