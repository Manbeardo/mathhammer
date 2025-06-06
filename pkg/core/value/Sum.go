package value

import (
	"cmp"

	"github.com/Manbeardo/mathhammer/pkg/core/probability"
)

type Sum []Interface

var _ Interface = (*Sum)(nil)

func (sum Sum) Distribution() probability.Distribution[int64] {
	dists := []probability.Distribution[int64]{}
	for _, i := range sum {
		dists = append(dists, i.Distribution())
	}
	return probability.Reduce(
		dists,
		func(a, b int64) int64 { return a + b },
		cmp.Compare,
		0,
	)
}
