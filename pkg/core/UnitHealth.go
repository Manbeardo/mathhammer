package core

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/Manbeardo/mathhammer/pkg/core/prob"
	"github.com/Manbeardo/mathhammer/pkg/core/util"
)

type UnitHealth []int64

func (mh UnitHealth) WoundsRemaining() int64 {
	sum := int64(0)
	for _, h := range mh {
		sum += h
	}
	return sum
}

func (mh UnitHealth) ModelsRemaining() int64 {
	sum := int64(0)
	for _, h := range mh {
		if h > 0 {
			sum += 1
		}
	}
	return sum
}

func (mh UnitHealth) ToKey() UnitHealthStr {
	b, err := json.Marshal(mh)
	if err != nil {
		panic(fmt.Errorf("marshaling ModelHealth: %w", err))
	}
	return UnitHealthStr(b)
}

func (mh UnitHealth) ToDist() prob.Dist[UnitHealthStr] {
	return util.Must(prob.NewConstDist(mh.ToKey()))
}

func (mh UnitHealth) String() string {
	return string(mh.ToKey())
}

type UnitHealthStr string

func (mh UnitHealthStr) ToSlice() UnitHealth {
	out := []int64{}
	err := json.Unmarshal([]byte(mh), &out)
	if err != nil {
		panic(fmt.Errorf("unmarshaling ModelHealth: %w", err))
	}
	return out
}

func MeanWoundsRemaining(dist prob.Dist[UnitHealthStr]) *big.Rat {
	avg := big.NewRat(0, 1)
	for healthStr, p := range dist.Iter() {
		w := big.NewRat(healthStr.ToSlice().WoundsRemaining(), 1)
		var partial big.Rat
		partial.Mul(w, p)
		avg.Add(avg, &partial)
	}
	return avg
}
