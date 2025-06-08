package value

import (
	"cmp"

	"github.com/Manbeardo/mathhammer/pkg/core/prob"
	"github.com/Manbeardo/mathhammer/pkg/core/util"
)

type SumT []Interface

func Sum(values ...Interface) SumT {
	return SumT(values)
}

var _ Interface = (*SumT)(nil)

func (sum SumT) Distribution() prob.Dist[int64] {
	dists := []prob.Dist[int64]{}
	for _, i := range sum {
		dists = append(dists, i.Distribution())
	}
	return util.Must(prob.Reduce(
		dists,
		func(a, b int64) int64 { return a + b },
		cmp.Compare,
		0,
	))
}
