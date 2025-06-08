package value

import (
	"github.com/Manbeardo/mathhammer/pkg/core/prob"
	"github.com/Manbeardo/mathhammer/pkg/core/util"
)

type IntT struct {
	N int64
}

func Int(n int64) IntT {
	return IntT{N: n}
}

var _ Interface = (*IntT)(nil)

func (v IntT) Distribution() prob.Dist[int64] {
	return util.Must(prob.FromConst(v.N))
}
