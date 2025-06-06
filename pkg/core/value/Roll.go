package value

import (
	"math/big"

	"github.com/Manbeardo/mathhammer/pkg/core/prob"
)

type RollT struct {
	N int64
}

func Roll(n int64) RollT {
	return RollT{N: n}
}

var _ Interface = (*RollT)(nil)

func (v RollT) Distribution() prob.Dist[int64] {
	out := map[int64]*big.Rat{}
	for i := range v.N {
		out[i+1] = big.NewRat(1, v.N)
	}
	return prob.NewDistribution(out)
}
