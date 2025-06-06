package core

import (
	"cmp"
	"math"
	"slices"

	"github.com/Manbeardo/mathhammer/pkg/core/check"
	"github.com/Manbeardo/mathhammer/pkg/core/prob"
	"github.com/Manbeardo/mathhammer/pkg/core/value"
)

type Attack struct {
	AttackTargets
	Modifiers map[ModifierKind]Modifiers
}

func NewAttack(targets AttackTargets) *Attack {
	return &Attack{
		AttackTargets: targets,
		Modifiers:     map[ModifierKind]Modifiers{},
	}
}

func (a *Attack) Modify(kind ModifierKind, in int64) int64 {
	return a.Modifiers[kind].Apply(kind, in)
}

func (a *Attack) ModifyDist(kind ModifierKind, dist prob.Dist[int64]) prob.Dist[int64] {
	return a.Modifiers[kind].ApplyDist(kind, dist)
}

func (a *Attack) Clone() *Attack {
	clone := NewAttack(a.AttackTargets)
	for k, v := range a.Modifiers {
		clone.Modifiers[k] = slices.Clone(v)
	}
	return clone
}

func (a *Attack) AllocateWound() {
	m := a.findModelToWound(a.DefenderUnit.models)
	for _, l := range a.DefenderUnit.leaders {
		if m != nil {
			continue
		}
		m = a.findModelToWound(l.models)
	}
	a.DefenderModel = m
}

func (a *Attack) findModelToWound(models []*Model) *Model {
	for _, m := range models {
		if !m.IsDead() && m.woundsTaken > 0 {
			return m
		}
	}
	for _, m := range models {
		if !m.IsDead() {
			return m
		}
	}
	return nil
}

func (a *Attack) ApplyTriggerEffects(trigger AbilityTrigger) {
	effects := getEffects(a.AttackTargets, trigger)
	for _, e := range effects {
		e.ApplyEffect(a)
	}
}

func (a *Attack) Attacks() prob.Dist[int64] {
	return a.ModifyDist(
		ModWeaponAttacks,
		a.AttackerWeaponProfile.tpl.Attacks.Distribution(),
	)
}

func (a *Attack) Hits(attacks prob.Dist[int64]) prob.Dist[check.Outcome] {
	skill := a.Modify(
		ModWeaponSkill,
		a.AttackerWeaponProfile.tpl.Skill,
	)
	return check.Calculate(value.Roll(6), check.Opts{
		Count:                    attacks,
		SuccessTarget:            value.Int(skill),
		CriticalSuccessThreshold: 6,
		CriticalFailureThreshold: 1,
		ModifierFn: func(in int64) int64 {
			return a.Modify(ModWeaponRollHit, in)
		},
	})
}

func (a *Attack) Wounds(hits prob.Dist[int64]) prob.Dist[check.Outcome] {
	strengthDist := a.ModifyDist(
		ModWeaponStrength,
		a.AttackerWeaponProfile.tpl.Strength.Distribution(),
	)
	toughness := a.Modify(
		ModModelToughness,
		a.DefenderModel.tpl.Toughness,
	)
	targetDist := prob.Map(
		strengthDist,
		func(strength int64) int64 {
			switch {
			case strength >= toughness*2:
				return 2
			case strength > toughness:
				return 3
			case strength == toughness:
				return 4
			case strength*2 <= toughness:
				return 6
			default:
				return 5
			}
		},
		cmp.Compare,
	)

	return check.Calculate(value.Roll(6), check.Opts{
		Count:                    hits,
		SuccessTarget:            targetDist,
		CriticalSuccessThreshold: 6,
		CriticalFailureThreshold: 1,
		ModifierFn: func(in int64) int64 {
			return a.Modify(ModWeaponRollWound, in)
		},
	})
}

func (a *Attack) Saves() RollResult {
	ap := a.Modify(
		ModWeaponArmorPen,
		a.AttackerWeaponProfile.tpl.ArmorPenetration,
	)
	a.Modifiers[ModModelRollSave] = append(
		a.Modifiers[ModModelRollSave],
		AdditionModifier{N: ap},
	)
	saveTarget := a.Modify(
		ModModelArmourSave,
		a.DefenderModel.tpl.Save,
	)

	return Roller{
		Value:                    RollValue{N: 6},
		SuccessTarget:            saveTarget,
		CriticalSuccessThreshold: math.MaxInt,
		CriticalFailureThreshold: 1,
		ModifyFn: func(in int) int {
			return a.Modify(ModModelRollSave, in)
		},
	}.Roll()
}

func (a *Attack) EvalDamage() int {
	return a.Modify(
		ModWeaponDamage,
		a.AttackerWeaponProfile.tpl.Damage.Eval(),
	)
}

type AttackTargets struct {
	AttackerUnit          *Unit
	AttackerModel         *Model
	AttackerWeaponProfile *WeaponProfile

	DefenderUnit  *Unit
	DefenderModel *Model
}

func (ctx AttackTargets) AllBearers() []AbilityBearer {
	return slices.DeleteFunc([]AbilityBearer{
		ctx.AttackerUnit,
		ctx.AttackerModel,
		ctx.AttackerWeaponProfile,
		ctx.DefenderUnit,
		ctx.DefenderModel,
	}, func(ab AbilityBearer) bool {
		return ab == nil
	})
}
