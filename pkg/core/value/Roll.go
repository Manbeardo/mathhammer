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
	out := []prob.EntryT[int64]{}
	for i := range v.N {
		out = append(out, prob.EntryT[int64]{
			Key: i + 1, Value: big.NewRat(1, v.N),
		})
	}
	return util.Must(prob.FromEntries(out))
}
