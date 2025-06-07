package effects

import (
	"github.com/Manbeardo/mathhammer/pkg/core"
	"github.com/Manbeardo/mathhammer/pkg/core/modifier"
)

type Modifier struct {
	Mod  modifier.Interface
	Kind modifier.Kind
}

func (m Modifier) ApplyEffect(attack *core.Attack) {
	// attack.Modifiers[m.Kind] = append(attack.Modifiers[m.Kind], m.Mod)
}
