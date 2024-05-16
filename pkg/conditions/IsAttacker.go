package conditions

import "github.com/Manbeardo/mathhammer/pkg/core"

func IsAttacker() Condition {
	return func(ac core.AbilityContext) bool {
		return ac.IsAttacker()
	}
}
