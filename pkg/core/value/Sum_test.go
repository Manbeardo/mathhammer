package value

import (
	"math/big"
	"testing"

	"github.com/Manbeardo/mathhammer/pkg/core/prob"
	"github.com/Manbeardo/mathhammer/pkg/core/util"
	"github.com/stretchr/testify/assert"
)

func TestSum(t *testing.T) {
	t.Run("has a correct distribution for 2D6", func(t *testing.T) {
		v := SumT{
			RollT{N: 6},
			RollT{N: 6},
		}
		dist := v.Distribution()
		assert.Equal(t, util.Must(prob.FromMap(prob.MapT[int64]{
			2:  big.NewRat(1, 36),
			3:  big.NewRat(2, 36),
			4:  big.NewRat(3, 36),
			5:  big.NewRat(4, 36),
			6:  big.NewRat(5, 36),
			7:  big.NewRat(6, 36),
			8:  big.NewRat(5, 36),
			9:  big.NewRat(4, 36),
			10: big.NewRat(3, 36),
			11: big.NewRat(2, 36),
			12: big.NewRat(1, 36),
		})), dist)
	})
}
