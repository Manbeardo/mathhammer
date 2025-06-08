package core

import (
	"math/big"
	"testing"

	"github.com/Manbeardo/mathhammer/pkg/core/prob"
	"github.com/Manbeardo/mathhammer/pkg/core/util"
	"github.com/stretchr/testify/assert"
)

func TestMeanWoundsRemaining(t *testing.T) {
	t.Run("weights results correctly", func(t *testing.T) {
		dist := util.Must(prob.FromEntries(
			[]util.Entry[UnitHealth, *big.Rat]{
				{Key: UnitHealth{0, 10}, Value: big.NewRat(1, 3)},
				{Key: UnitHealth{2, 3}, Value: big.NewRat(2, 3)},
			},
		))
		mean := MeanWoundsRemaining(dist)
		assert.Equal(t, big.NewRat(20, 3), mean)
	})
}
