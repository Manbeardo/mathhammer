package check

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResult(t *testing.T) {
	t.Run("getters work correctly", func(t *testing.T) {
		r := Result{
			ns: *big.NewRat(1, 1),
			cs: *big.NewRat(2, 1),
			nf: *big.NewRat(3, 1),
			cf: *big.NewRat(5, 1),
		}
		assert.Equal(t, 3.0, r.Successes())
		assert.Equal(t, 1.0, r.NormalSuccesses())
		assert.Equal(t, 2.0, r.CriticalSuccesses())
		assert.Equal(t, 8.0, r.Failures())
		assert.Equal(t, 3.0, r.NormalFailures())
		assert.Equal(t, 5.0, r.CriticalFailures())
	})
	t.Run("Add", func(t *testing.T) {
		a := Result{
			ns: *big.NewRat(1, 1),
			cs: *big.NewRat(2, 1),
			nf: *big.NewRat(3, 1),
			cf: *big.NewRat(5, 1),
		}
		b := Result{
			ns: *big.NewRat(7, 1),
			cs: *big.NewRat(11, 1),
			nf: *big.NewRat(13, 1),
			cf: *big.NewRat(17, 1),
		}
		var sum Result
		sum.Add(&a, &b)

		// sum is correct
		assert.Equal(t, Result{
			ns: *big.NewRat(8, 1),
			cs: *big.NewRat(13, 1),
			nf: *big.NewRat(16, 1),
			cf: *big.NewRat(22, 1),
		}, sum)
		// a is unmodified
		assert.Equal(t, Result{
			ns: *big.NewRat(1, 1),
			cs: *big.NewRat(2, 1),
			nf: *big.NewRat(3, 1),
			cf: *big.NewRat(5, 1),
		}, a)
		// b is unmodified
		assert.Equal(t, Result{
			ns: *big.NewRat(7, 1),
			cs: *big.NewRat(11, 1),
			nf: *big.NewRat(13, 1),
			cf: *big.NewRat(17, 1),
		}, b)
	})

	t.Run("Mul", func(t *testing.T) {
		a := Result{
			ns: *big.NewRat(1, 1),
			cs: *big.NewRat(2, 1),
			nf: *big.NewRat(3, 1),
			cf: *big.NewRat(5, 1),
		}
		b := Result{
			ns: *big.NewRat(7, 1),
			cs: *big.NewRat(11, 1),
			nf: *big.NewRat(13, 1),
			cf: *big.NewRat(17, 1),
		}
		var product Result
		product.Mul(&a, &b)

		// product is correct
		assert.Equal(t, Result{
			ns: *big.NewRat(7, 1),
			cs: *big.NewRat(22, 1),
			nf: *big.NewRat(39, 1),
			cf: *big.NewRat(85, 1),
		}, product)
		// a is unmodified
		assert.Equal(t, Result{
			ns: *big.NewRat(1, 1),
			cs: *big.NewRat(2, 1),
			nf: *big.NewRat(3, 1),
			cf: *big.NewRat(5, 1),
		}, a)
		// b is unmodified
		assert.Equal(t, Result{
			ns: *big.NewRat(7, 1),
			cs: *big.NewRat(11, 1),
			nf: *big.NewRat(13, 1),
			cf: *big.NewRat(17, 1),
		}, b)
	})
}
