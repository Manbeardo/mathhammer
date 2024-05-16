package core

import (
	"math"
	"slices"
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

func (a *Attack) Modify(kind ModifierKind, in int) int {
	return a.Modifiers[kind].Apply(kind, in)
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

func (a *Attack) EvalAttacks() int {
	unmodified := a.AttackerWeaponProfile.tpl.Attacks.Eval()
	final := a.Modifiers[ModWeaponAttacks].Apply(ModWeaponAttacks, unmodified)
	return final
}

func (a *Attack) RollHits(attacks int) RollResult {
	skill := a.Modify(
		ModWeaponSkill,
		a.AttackerWeaponProfile.tpl.Skill,
	)
	return Roller{
		Value:                    RollValue{N: 6},
		SuccessTarget:            skill,
		CriticalSuccessThreshold: 6,
		CriticalFailureThreshold: 1,
		ModifyFn: func(in int) int {
			return a.Modify(ModWeaponRollHit, in)
		},
	}.RollN(attacks)
}

func (a *Attack) RollWounds(hits int) RollResult {
	// TODO: figure out whether this is rolled per hit or per profile
	strength := a.Modify(
		ModWeaponStrength,
		a.AttackerWeaponProfile.tpl.Strength.Eval(),
	)
	toughness := a.Modify(
		ModModelToughness,
		a.DefenderModel.tpl.Toughness,
	)
	var target int
	if strength >= toughness*2 {
		target = 2
	} else if strength > toughness {
		target = 3
	}
	if strength == toughness {
		target = 4
	}
	if strength*2 <= toughness {
		target = 6
	} else if strength < toughness {
		target = 5
	}

	return Roller{
		Value:                    RollValue{N: 6},
		SuccessTarget:            target,
		CriticalSuccessThreshold: 6,
		CriticalFailureThreshold: 1,
		ModifyFn: func(in int) int {
			return a.Modify(ModWeaponRollWound, in)
		},
	}.RollN(hits)
}

func (a *Attack) RollSave() RollResult {
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
