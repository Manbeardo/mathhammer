package check

type Opts struct {
	Count                    int64
	SuccessTarget            int64
	CriticalSuccessThreshold int64
	CriticalFailureThreshold int64
	ModifierFn               func(int64) int64
}

func (opts Opts) eval(v int64) Outcome {
	if v <= opts.CriticalFailureThreshold {
		return Outcome{CriticalFailures: 1}
	}
	if opts.CriticalSuccessThreshold > 0 && v >= opts.CriticalSuccessThreshold {
		return Outcome{CriticalSuccesses: 1}
	}
	if opts.ModifierFn != nil {
		v = opts.ModifierFn(v)
	}
	if v >= opts.SuccessTarget {
		return Outcome{NormalSuccesses: 1}
	}
	return Outcome{NormalFailures: 1}
}
