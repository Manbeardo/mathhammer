package value

import (
	"math/big"

	"github.com/Manbeardo/mathhammer/pkg/core/probability"
)

type Roll struct {
	N int64
}

var _ Interface = (*Roll)(nil)

func (v Roll) Distribution() probability.Distribution[int64] {
	out := map[int64]*big.Rat{}
	for i := range v.N {
		out[i+1] = big.NewRat(1, v.N)
	}
	return probability.NewDistribution(out)
}
