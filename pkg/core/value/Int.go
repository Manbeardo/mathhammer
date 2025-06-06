package value

import (
	"math/big"

	"github.com/Manbeardo/mathhammer/pkg/core/probability"
)

type Int struct {
	N int64
}

var _ Interface = (*Int)(nil)

func (v Int) Distribution() probability.Distribution[int64] {
	return probability.NewDistribution(map[int64]*big.Rat{
		(v.N): big.NewRat(1, 1),
	})
}
