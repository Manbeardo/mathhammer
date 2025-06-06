package check

import (
	"slices"

	"github.com/Manbeardo/mathhammer/pkg/core/probability"
	"github.com/Manbeardo/mathhammer/pkg/core/value"
)

// TODO: convert to a distribution of outcomes

func Calculate(v value.Interface, opts Opts) probability.Distribution[Outcome] {
	count := opts.Count
	if count == 0 {
		count = 1
	}
	dist := probability.Map(
		v.Distribution(),
		func(i int64) Outcome {
			return opts.eval(i)
		},
		CompareOutcomes,
	)
	return probability.Reduce(
		slices.Repeat([]probability.Distribution[Outcome]{dist}, int(count)),
		func(a, b Outcome) Outcome { return SumOutcomes(a, b) },
		CompareOutcomes,
		Outcome{},
	)
}
