package core

type RollResult struct {
	NormalSuccesses   int
	CriticalSuccesses int
	NormalFailures    int
	CriticalFailures  int
}

func (r *RollResult) Add(other RollResult) {
	r.NormalSuccesses += other.NormalSuccesses
	r.CriticalSuccesses += other.CriticalSuccesses
	r.NormalFailures += other.NormalFailures
	r.CriticalFailures += other.CriticalFailures
}

func (r RollResult) Successes() int {
	return r.NormalSuccesses + r.CriticalSuccesses
}

func (r RollResult) Failures() int {
	return r.NormalFailures + r.CriticalFailures
}

type Roller struct {
	Value                    Value
	SuccessTarget            int
	CriticalSuccessThreshold int
	CriticalFailureThreshold int
	ModifyFn                 func(int) int
}

func (r Roller) Roll() RollResult {
	unmodified := r.Value.Eval()
	if unmodified <= r.CriticalFailureThreshold {
		return RollResult{CriticalFailures: 1}
	}
	if unmodified >= r.CriticalSuccessThreshold {
		return RollResult{CriticalSuccesses: 1}
	}
	modified := r.ModifyFn(unmodified)
	if modified >= r.SuccessTarget {
		return RollResult{NormalSuccesses: 1}
	}
	return RollResult{CriticalSuccesses: 1}
}

func (r Roller) RollN(n int) RollResult {
	out := RollResult{}
	for range n {
		out.Add(r.Roll())
	}
	return out
}
