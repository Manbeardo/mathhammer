package check

import (
	"math/big"
	"testing"

	"github.com/Manbeardo/mathhammer/pkg/core/prob"
	"github.com/Manbeardo/mathhammer/pkg/core/util"
	"github.com/Manbeardo/mathhammer/pkg/core/value"
	"github.com/stretchr/testify/assert"
)

func TestCalculate(t *testing.T) {
	t.Run("1D6 4+ has a 1/2 success rate", func(t *testing.T) {
		r := Calculate(value.Roll(6), Opts{
			SuccessTarget: value.Int(4),
		})
		assert.Equal(t, util.Must(prob.FromMapFunc(map[Outcome]*big.Rat{
			{NormalSuccesses: 1}: big.NewRat(1, 2),
			{NormalFailures: 1}:  big.NewRat(1, 2),
		}, CompareOutcomes)), r)
	})
	t.Run("1D6 6+ with +5 modifier has a 1/1 success rate", func(t *testing.T) {
		r := Calculate(value.Roll(6), Opts{
			SuccessTarget: value.Int(6),
			ModifierFn:    func(i int64) int64 { return i + 5 },
		})
		assert.Equal(t, util.Must(prob.FromConst(Outcome{NormalSuccesses: 1})), r)
	})
	t.Run("2D6 7+ has a 21/36 success rate", func(t *testing.T) {
		r := Calculate(value.Sum(value.Roll(6), value.Roll(6)), Opts{
			SuccessTarget: value.Int(7),
		})
		assert.Equal(t, util.Must(prob.FromMapFunc(map[Outcome]*big.Rat{
			{NormalSuccesses: 1}: big.NewRat(21, 36),
			{NormalFailures: 1}:  big.NewRat(15, 36),
		}, CompareOutcomes)), r)
	})
	t.Run("Modifiers are applied after crits", func(t *testing.T) {
		r := Calculate(value.Roll(6), Opts{
			SuccessTarget:            value.Int(4),
			CriticalSuccessThreshold: 6,
			CriticalFailureThreshold: 1,
			ModifierFn:               func(i int64) int64 { return i + 3 },
		})
		assert.Equal(t, util.Must(prob.FromMapFunc(map[Outcome]*big.Rat{
			{NormalSuccesses: 1}:   big.NewRat(4, 6),
			{CriticalSuccesses: 1}: big.NewRat(1, 6),
			{CriticalFailures: 1}:  big.NewRat(1, 6),
		}, CompareOutcomes)), r)
	})
	t.Run("3x1D6 4+ has a correct distribution", func(t *testing.T) {
		r := Calculate(value.Roll(6), Opts{
			Count:         value.Int(3),
			SuccessTarget: value.Int(4),
		})
		assert.Equal(t, util.Must(prob.FromMapFunc(map[Outcome]*big.Rat{
			{NormalSuccesses: 3, NormalFailures: 0}: big.NewRat(1, 8),
			{NormalSuccesses: 2, NormalFailures: 1}: big.NewRat(3, 8),
			{NormalSuccesses: 1, NormalFailures: 2}: big.NewRat(3, 8),
			{NormalSuccesses: 0, NormalFailures: 3}: big.NewRat(1, 8),
		}, CompareOutcomes)), r)
	})

	t.Run("1D2x1D6 4+ has a correct distribution", func(t *testing.T) {
		r := Calculate(value.Roll(6), Opts{
			Count:         value.Roll(2),
			SuccessTarget: value.Int(4),
		})
		assert.Equal(t, util.Must(prob.FromMapFunc(map[Outcome]*big.Rat{
			{NormalSuccesses: 1, NormalFailures: 0}: big.NewRat(1, 4),
			{NormalSuccesses: 0, NormalFailures: 1}: big.NewRat(1, 4),
			{NormalSuccesses: 2, NormalFailures: 0}: big.NewRat(1, 8),
			{NormalSuccesses: 1, NormalFailures: 1}: big.NewRat(2, 8),
			{NormalSuccesses: 0, NormalFailures: 2}: big.NewRat(1, 8),
		}, CompareOutcomes)), r)
	})

	t.Run("2x1D4 D2+1+ has a correct distribution", func(t *testing.T) {
		r := Calculate(value.Roll(4), Opts{
			Count:         value.Int(2),
			SuccessTarget: value.Sum(value.Roll(2), value.Int(1)),
		})
		assert.Equal(t, util.Must(prob.FromMapFunc(map[Outcome]*big.Rat{
			{NormalSuccesses: 2, NormalFailures: 0}: big.NewRat(13, 32),
			{NormalSuccesses: 1, NormalFailures: 1}: big.NewRat(14, 32),
			{NormalSuccesses: 0, NormalFailures: 2}: big.NewRat(5, 32),
		}, CompareOutcomes)), r)
	})
}
