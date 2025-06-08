package prob

import (
	"math/big"

	"github.com/Manbeardo/mathhammer/pkg/core/util"
)

// reduces probability distributions of multiple independent events into an aggregate
// probability distribution
func Reduce[T any, U any](
	dists []Dist[T],
	collector func(U, T) U,
	initialValue U,
) (Dist[U], error) {
	out, err := FromEntries([]util.Entry[U, *big.Rat]{
		{Key: initialValue, Value: big.NewRat(1, 1)},
	})
	if err != nil {
		return out, err
	}
	for _, dist := range dists {
		nextPmap := map[Key]*big.Rat{}
		nextVmap := map[Key]U{}
		for tk, tp := range dist.pmap {
			tv := dist.vmap[tk]
			for prevUk, prevUp := range out.pmap {
				prevUv := out.vmap[prevUk]
				partialP := big.NewRat(0, 1)
				partialP.Mul(tp, prevUp)
				nextUv := collector(prevUv, tv)
				nextUk := out.key(nextUv)
				nextVmap[nextUk] = nextUv
				nextUp, ok := nextPmap[nextUk]
				if !ok {
					nextUp = big.NewRat(0, 1)
					nextPmap[nextUk] = nextUp
				}
				nextUp.Add(nextUp, partialP)
			}
		}
		out.pmap = nextPmap
		out.vmap = nextVmap
	}
	return out.validate()
}
