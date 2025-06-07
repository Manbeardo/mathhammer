package effects

import "github.com/Manbeardo/mathhammer/pkg/core"

type Blast struct{}

func (b Blast) ApplyEffect(attack *core.Attack) {
	// Modifier{
	// 	Kind: core.ModWeaponAttacks,
	// 	Mod:  core.AdditionModifier{N: attack.DefenderUnit.ModelCount() / 5},
	// }.ApplyEffect(attack)
}
