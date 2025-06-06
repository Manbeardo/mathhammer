package check

import (
	"fmt"
	"math/big"

	"github.com/Manbeardo/mathhammer/pkg/core/util"
)

type Result struct {
	ns big.Rat // normal successes
	cs big.Rat // crit successes
	nf big.Rat // normal failures
	cf big.Rat // crit failures
}

type bareResult Result

// creates a normalized copy of the result, ensuring that its denominators are initialized
func (r Result) normalize() *Result {
	r.ns.Set(&r.ns)
	r.cs.Set(&r.cs)
	r.nf.Set(&r.nf)
	r.cf.Set(&r.cf)
	return &r
}

func (r Result) Format(w fmt.State, v rune) {
	util.PrettyFormat(w, v, bareResult(r))
}

func ConstResultFromRat(rat *big.Rat) *Result {
	return &Result{
		ns: *rat,
		cs: *rat,
		nf: *rat,
		cf: *rat,
	}
}

func ConstResultFromInt(v int64) *Result {
	return ConstResultFromRat(big.NewRat(v, 1))
}

// Add sets r to the sum a+b and returns r
func (r *Result) Add(a, b *Result) *Result {
	r.ns.Add(&a.ns, &b.ns)
	r.cs.Add(&a.cs, &b.cs)
	r.nf.Add(&a.nf, &b.nf)
	r.cf.Add(&a.cf, &b.cf)
	return r
}

// Mul sets r to the product a*b and returns r
func (r *Result) Mul(a, b *Result) *Result {
	r.ns.Mul(&a.ns, &b.ns)
	r.cs.Mul(&a.cs, &b.cs)
	r.nf.Mul(&a.nf, &b.nf)
	r.cf.Mul(&a.cf, &b.cf)
	return r
}

func (r Result) Successes() float64 {
	var sum big.Rat
	sum.Add(&r.ns, &r.cs)
	return ratToFloat(sum)
}

func (r Result) NormalSuccesses() float64 {
	return ratToFloat(r.ns)
}

func (r Result) CriticalSuccesses() float64 {
	return ratToFloat(r.cs)
}

func (r Result) Failures() float64 {
	var sum big.Rat
	sum.Add(&r.nf, &r.cf)
	return ratToFloat(sum)
}

func (r Result) NormalFailures() float64 {
	return ratToFloat(r.nf)
}

func (r Result) CriticalFailures() float64 {
	return ratToFloat(r.cf)
}

func ratToFloat(r big.Rat) float64 {
	out, ok := r.Float64()
	if !ok {
		panic("numeric overflow")
	}
	return out
}
