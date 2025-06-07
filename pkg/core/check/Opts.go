package check

import "github.com/Manbeardo/mathhammer/pkg/core/value"

type Opts struct {
	Count                    value.Interface
	SuccessTarget            value.Interface
	CriticalSuccessThreshold int64
	CriticalFailureThreshold int64
	ModifierFn               func(int64) int64
	// TODO: rerolls
}

func (opts Opts) eval(v int64, target int64) Outcome {
	if v <= opts.CriticalFailureThreshold {
		return Outcome{CriticalFailures: 1}
	}
	if opts.CriticalSuccessThreshold > 0 && v >= opts.CriticalSuccessThreshold {
		return Outcome{CriticalSuccesses: 1}
	}
	if opts.ModifierFn != nil {
		v = opts.ModifierFn(v)
	}
	if v >= target {
		return Outcome{NormalSuccesses: 1}
	}
	return Outcome{NormalFailures: 1}
}
