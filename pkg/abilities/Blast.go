package abilities

import (
	"github.com/Manbeardo/mathhammer/pkg/core"
	"github.com/Manbeardo/mathhammer/pkg/effects"
)

func Blast() core.Ability {
	return ability{
		id:      "blast",
		trigger: core.TriggerSelectTargetUnit,
		// condition: conditions.IsAttacker(),
		effectsFn: func(ctx core.AbilityContext) []core.Effect {
			return []core.Effect{effects.Blast{}}
		},
	}
}
