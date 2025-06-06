package check

import "cmp"

type Outcome struct {
	NormalSuccesses   int64
	NormalFailures    int64
	CriticalSuccesses int64
	CriticalFailures  int64
}

func (o Outcome) Successes() int64 {
	return o.NormalSuccesses + o.CriticalSuccesses
}

func (o Outcome) Failures() int64 {
	return o.NormalFailures + o.CriticalFailures
}

func CompareOutcomes(a, b Outcome) int {
	if c := cmp.Compare(b.Successes(), a.Successes()); c != 0 {
		return c
	}
	if c := cmp.Compare(b.CriticalSuccesses, a.CriticalSuccesses); c != 0 {
		return c
	}
	if c := cmp.Compare(a.Failures(), b.Failures()); c != 0 {
		return c
	}
	if c := cmp.Compare(a.CriticalFailures, b.CriticalFailures); c != 0 {
		return c
	}
	return 0
}

func SumOutcomes(outcomes ...Outcome) Outcome {
	out := Outcome{}
	for _, o := range outcomes {
		out.NormalSuccesses += o.NormalSuccesses
		out.NormalFailures += o.NormalFailures
		out.CriticalSuccesses += o.CriticalSuccesses
		out.CriticalFailures += o.CriticalFailures
	}
	return out
}
