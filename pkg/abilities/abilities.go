package abilities

import (
	"cmp"

	"github.com/Manbeardo/mathhammer/pkg/conditions"
	"github.com/Manbeardo/mathhammer/pkg/core"
)

type EffectsFn func(ctx core.AbilityContext) []core.Effect

type ability struct {
	id           string
	trigger      core.AbilityTrigger
	condition    conditions.Condition
	effectsFn    EffectsFn
	compareValue int
}

func (a ability) ID() string {
	return a.id
}

func (a ability) Trigger() core.AbilityTrigger {
	return a.trigger
}

func (a ability) ShouldApply(ctx core.AbilityContext) bool {
	return a.condition(ctx)
}

func (a ability) Effects(ctx core.AbilityContext) []core.Effect {
	return a.effectsFn(ctx)
}

func (a ability) CompareTo(other core.Ability) int {
	asStruct, ok := other.(ability)
	if !ok {
		return 0
	}
	return cmp.Compare(a.compareValue, asStruct.compareValue)
}
