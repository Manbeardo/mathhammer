package check

import (
	"slices"

	"github.com/Manbeardo/mathhammer/pkg/core/prob"
	"github.com/Manbeardo/mathhammer/pkg/core/value"
)

func Calculate(v value.Interface, opts Opts) prob.Dist[Outcome] {
	countV := opts.Count
	if countV == nil {
		countV = value.Int(1)
	}
	countDist := countV.Distribution()

	targetV := opts.SuccessTarget
	if targetV == nil {
		targetV = value.Int(0)
	}
	targetDist := opts.SuccessTarget.Distribution()

	valueDist := v.Distribution()

	return prob.FlatMap(
		targetDist,
		func(target int64) prob.Dist[Outcome] {
			rollDist := prob.Map(
				valueDist,
				func(v int64) Outcome {
					return opts.eval(v, target)
				},
				CompareOutcomes,
			)
			return prob.FlatMap(
				countDist,
				func(count int64) prob.Dist[Outcome] {
					return prob.Reduce(
						slices.Repeat([]prob.Dist[Outcome]{rollDist}, int(count)),
						func(a, b Outcome) Outcome { return SumOutcomes(a, b) },
						CompareOutcomes,
						Outcome{},
					)
				},
				CompareOutcomes,
			)
		},
		CompareOutcomes,
	)
}
