package modifier

import (
	"cmp"
	"math"
	"slices"

	"github.com/Manbeardo/mathhammer/pkg/core/prob"
)

type Kind string

const (
	WeaponAttacks  Kind = "c_weapon_attacks"
	WeaponSkill    Kind = "c_weapon_skill"
	WeaponStrength Kind = "c_weapon_strength"
	WeaponArmorPen Kind = "c_weapon_armor_penetration"
	WeaponDamage   Kind = "c_weapon_damage"

	WeaponRollHit   Kind = "r_weapon_hit"
	WeaponRollWound Kind = "r_weapon_wound"

	ModelToughness  Kind = "c_model_toughness"
	ModelArmourSave Kind = "c_model_armour_save"
	ModelLeadership Kind = "c_model_leadership"
	ModelOC         Kind = "c_model_objective_control"

	ModelRollSave Kind = "r_model_save"
)

type Interface interface {
	Apply(float64) float64
	Priority() int64
}

type Set []Interface

func (ms Set) ApplyDist(
	kind Kind,
	dist prob.Dist[int64],
) prob.Dist[int64] {
	return prob.Map(
		dist,
		func(in int64) int64 {
			return ms.Apply(kind, in)
		},
		cmp.Compare,
	)
}

func (ms Set) Apply(kind Kind, in int64) int64 {
	slices.SortFunc(ms, func(a, b Interface) int {
		return cmp.Compare(a.Priority(), b.Priority())
	})
	runningValue := float64(in)
	for _, m := range ms {
		runningValue = m.Apply(runningValue)
	}
	result := int64(math.Ceil(runningValue))

	switch kind {
	case WeaponRollHit, WeaponRollWound:
		if result < in-1 {
			result = in - 1
		}
		if result > in+1 {
			result = in + 1
		}
	case ModelRollSave:
		if result > in+1 {
			result = in + 1
		}
	case WeaponSkill, ModelArmourSave:
		if result < 2 {
			result = 2
		}
	case ModelLeadership:
		if result < 4 {
			result = 4
		}
		if result > 9 {
			result = 9
		}
	case WeaponAttacks, WeaponStrength, WeaponDamage, ModelToughness:
		if result < 1 {
			result = 1
		}
	case WeaponArmorPen:
		if result > 0 {
			result = 0
		}
	case ModelOC:
		if result < 0 {
			result = 0
		}
	}

	return result
}

type ReplacementModifier struct {
	N int64
}

func Replace(n int64) ReplacementModifier {
	return ReplacementModifier{N: n}
}
func (m ReplacementModifier) Apply(in float64) float64 {
	return float64(m.N)
}

func (m ReplacementModifier) Priority() int64 {
	return 0
}

type DivisionModifier struct {
	N int64
}

func Divide(n int64) DivisionModifier {
	return DivisionModifier{N: n}
}

func (m DivisionModifier) Apply(in float64) float64 {
	return in / float64(m.N)
}

func (m DivisionModifier) Priority() int64 {
	return 1
}

type MultiplicationModifier struct {
	N int64
}

func Multiply(n int64) MultiplicationModifier {
	return MultiplicationModifier{N: n}
}

func (m MultiplicationModifier) Apply(in float64) float64 {
	return in * float64(m.N)
}

func (m MultiplicationModifier) Priority() int64 {
	return 2
}

type AdditionModifier struct {
	N int64
}

func Add(n int64) AdditionModifier {
	return AdditionModifier{N: n}
}

func (m AdditionModifier) Apply(in float64) float64 {
	return in + float64(m.N)
}

func (m AdditionModifier) Priority() int64 {
	return 3
}

type SubtractionModifier struct {
	N int64
}

func Subtract(n int64) SubtractionModifier {
	return SubtractionModifier{N: n}
}

func (m SubtractionModifier) Apply(in float64) float64 {
	return in - float64(m.N)
}

func (m SubtractionModifier) Priority() int64 {
	return 4
}
