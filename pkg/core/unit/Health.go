package unit

import (
	"cmp"
	"fmt"
	"math/big"

	"github.com/Manbeardo/mathhammer/pkg/core/prob"
	"github.com/Manbeardo/mathhammer/pkg/core/util"
)

type Health []int64

func (mh Health) WoundsRemaining() int64 {
	sum := int64(0)
	for _, h := range mh {
		sum += h
	}
	return sum
}

func (mh Health) ModelsRemaining() int64 {
	sum := int64(0)
	for _, h := range mh {
		if h > 0 {
			sum += 1
		}
	}
	return sum
}

func (mh Health) StringKey() string {
	return fmt.Sprintf("%v", mh)
}

func (mh Health) ToDist() prob.Dist[Health] {
	return util.Must(prob.FromConst(mh))
}

func MeanWoundsRemaining(dist prob.Dist[Health]) *big.Rat {
	avg := big.NewRat(0, 1)
	for health, p := range dist.Iter() {
		w := big.NewRat(health.WoundsRemaining(), 1)
		var partial big.Rat
		partial.Mul(w, p)
		avg.Add(avg, &partial)
	}
	return avg
}

func CompareHealth(a, b Health) int {
	if c := cmp.Compare(a.WoundsRemaining(), b.WoundsRemaining()); c != 0 {
		return c
	}
	for i := 1; i <= len(a) && i <= len(b); i++ {
		if c := cmp.Compare(a[len(a)-i], b[len(b)-i]); c != 0 {
			return c
		}
	}
	return 0
}
