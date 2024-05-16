package conditions

import "github.com/Manbeardo/mathhammer/pkg/core"

func IsDefender() Condition {
	return func(ac core.AbilityContext) bool {
		return ac.IsDefender()
	}
}
