package core

import (
	"cmp"
	"math"
	"slices"
)

type ModifierKind string

const (
	ModWeaponAttacks  ModifierKind = "c_weapon_attacks"
	ModWeaponSkill    ModifierKind = "c_weapon_skill"
	ModWeaponStrength ModifierKind = "c_weapon_strength"
	ModWeaponArmorPen ModifierKind = "c_weapon_armor_penetration"
	ModWeaponDamage   ModifierKind = "c_weapon_damage"

	ModWeaponRollHit   ModifierKind = "r_weapon_hit"
	ModWeaponRollWound ModifierKind = "r_weapon_wound"

	ModModelToughness  ModifierKind = "c_model_toughness"
	ModModelArmourSave ModifierKind = "c_model_armour_save"
	ModModelLeadership ModifierKind = "c_model_leadership"
	ModModelOC         ModifierKind = "c_model_objective_control"

	ModModelRollSave ModifierKind = "r_model_save"
)

type Modifier interface {
	Apply(float64) float64
	Priority() int
}

type Modifiers []Modifier

func (ms Modifiers) Apply(kind ModifierKind, in int) int {
	slices.SortFunc(ms, func(a, b Modifier) int {
		return cmp.Compare(a.Priority(), b.Priority())
	})
	runningValue := 0.0
	for _, m := range ms {
		runningValue = m.Apply(runningValue)
	}
	result := int(math.Ceil(runningValue))

	switch kind {
	case ModWeaponRollHit, ModWeaponRollWound:
		if result < in-1 {
			result = in - 1
		}
		if result > in+1 {
			result = in + 1
		}
	case ModModelRollSave:
		if result > in+1 {
			result = in + 1
		}
	case ModWeaponSkill, ModModelArmourSave:
		if result < 2 {
			result = 2
		}
	case ModModelLeadership:
		if result < 4 {
			result = 4
		}
		if result > 9 {
			result = 9
		}
	case ModWeaponAttacks, ModWeaponStrength, ModWeaponDamage, ModModelToughness:
		if result < 1 {
			result = 1
		}
	case ModWeaponArmorPen:
		if result > 0 {
			result = 0
		}
	case ModModelOC:
		if result < 0 {
			result = 0
		}
	}

	return result
}

type ReplacementModifier struct {
	N int
}

func (m ReplacementModifier) Apply(in float64) float64 {
	return float64(m.N)
}

func (m ReplacementModifier) Priority() int {
	return 0
}

type DivisionModifier struct {
	N int
}

func (m DivisionModifier) Apply(in float64) float64 {
	return in / float64(m.N)
}

func (m DivisionModifier) Priority() int {
	return 1
}

type MultiplicationModifier struct {
	N int
}

func (m MultiplicationModifier) Apply(in float64) float64 {
	return in * float64(m.N)
}

func (m MultiplicationModifier) Priority() int {
	return 2
}

type AdditionModifier struct {
	N int
}

func (m AdditionModifier) Apply(in float64) float64 {
	return in + float64(m.N)
}

func (m AdditionModifier) Priority() int {
	return 3
}

type SubtractionModifier struct {
	N int
}

func (m SubtractionModifier) Apply(in float64) float64 {
	return in - float64(m.N)
}

func (m SubtractionModifier) Priority() int {
	return 4
}
