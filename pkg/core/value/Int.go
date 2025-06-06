package value

import (
	"math/big"

	"github.com/Manbeardo/mathhammer/pkg/core/prob"
)

type IntT struct {
	N int64
}

func Int(n int64) IntT {
	return IntT{N: n}
}

var _ Interface = (*IntT)(nil)

func (v IntT) Distribution() prob.Dist[int64] {
	return prob.NewDistribution(map[int64]*big.Rat{
		(v.N): big.NewRat(1, 1),
	})
}
