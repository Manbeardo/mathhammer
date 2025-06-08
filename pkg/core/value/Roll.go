package value

import (
	"math/big"

	"github.com/Manbeardo/mathhammer/pkg/core/prob"
	"github.com/Manbeardo/mathhammer/pkg/core/util"
)

type RollT struct {
	N int64
}

func Roll(n int64) RollT {
	return RollT{N: n}
}

var _ Interface = (*RollT)(nil)

func (v RollT) Distribution() prob.Dist[int64] {
	out := []util.Entry[int64, *big.Rat]{}
	for i := range v.N {
		out = append(out, util.Entry[int64, *big.Rat]{
			Key: i + 1, Value: big.NewRat(1, v.N),
		})
	}
	return util.Must(prob.NewDist(out))
}
