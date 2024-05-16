package effects

import "github.com/Manbeardo/mathhammer/pkg/core"

type Modifier struct {
	Mod  core.Modifier
	Kind core.ModifierKind
}

func (m Modifier) ApplyEffect(attack *core.Attack) {
	attack.Modifiers[m.Kind] = append(attack.Modifiers[m.Kind], m.Mod)
}
