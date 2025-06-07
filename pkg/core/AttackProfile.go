package core

import (
	"cmp"
	"slices"

	"github.com/Manbeardo/mathhammer/pkg/core/check"
	"github.com/Manbeardo/mathhammer/pkg/core/modifier"
	"github.com/Manbeardo/mathhammer/pkg/core/prob"
	"github.com/Manbeardo/mathhammer/pkg/core/value"
)

type AttackProfile struct {
	Attack
	AttackerWeaponProfile  *WeaponProfileTemplate
	AttackerWeaponCount    int64
	DefenderStartingHealth prob.Dist[UnitHealthStr]
}

func (a AttackProfile) attacks() prob.Dist[int64] {
	if a.DistanceInches > a.AttackerWeaponProfile.RangeInches ||
		(a.DistanceInches == 0 && a.AttackerWeaponProfile.RangeInches > 0) {
		return value.Int(0).Distribution()
	}

	return value.Sum(
		slices.Repeat([]value.Interface{
			a.AttackerWeaponProfile.Attacks,
		}, int(a.AttackerWeaponCount))...,
	).Distribution()
}

func (a AttackProfile) hits(attacks prob.Dist[int64]) prob.Dist[check.Outcome] {
	skill := a.AttackerWeaponProfile.Skill
	return check.Calculate(value.Roll(6), check.Opts{
		Count:                    attacks,
		SuccessTarget:            value.Int(skill),
		CriticalSuccessThreshold: 6,
		CriticalFailureThreshold: 1,
	})
}

func (a AttackProfile) wounds(hits prob.Dist[int64]) prob.Dist[check.Outcome] {
	strengthDist := a.AttackerWeaponProfile.Strength.Distribution()
	toughness := a.DefenderToughness
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
	})
}

func (a AttackProfile) allocateWound(healthSlice []int64) (m *Model, idx int) {
	// TODO: [PRECISION]
	for idx, health := range healthSlice {
		if health > 0 {
			return a.DefenderUnit.Model(idx), idx
		}
	}
	return nil, -1
}

func (a AttackProfile) resolveNormalWounds(woundDist prob.Dist[int64]) prob.Dist[UnitHealthStr] {
	ap := a.AttackerWeaponProfile.ArmorPenetration
	saveModifiers := modifier.Set{
		modifier.Add(ap),
	}

	return prob.FlatMap(
		woundDist,
		func(wounds int64) prob.Dist[UnitHealthStr] {
			healthDist := a.DefenderStartingHealth
			for range wounds {
				// TODO: memoize this
				healthDist = prob.FlatMap(
					healthDist,
					func(healthStr UnitHealthStr) prob.Dist[UnitHealthStr] {
						healthSlice := healthStr.ToSlice()
						model, idx := a.allocateWound(healthSlice)
						if model == nil {
							return prob.NewConstDist(healthStr)
						}

						save := saveModifiers.Apply(modifier.ModelArmourSave, model.tpl.Save)

						checkDist := check.Calculate(value.Roll(6), check.Opts{
							SuccessTarget:            value.Int(save),
							CriticalFailureThreshold: 1,
						})

						return prob.Map(
							checkDist,
							func(outcome check.Outcome) UnitHealthStr {
								healthSliceCopy := slices.Clone(healthSlice)
								damage := a.AttackerWeaponProfile.Damage
								for range outcome.Failures() {
									health := healthSliceCopy[idx]
									if damage > health {
										healthSliceCopy[idx] = 0
									} else {
										healthSliceCopy[idx] -= damage
									}
								}
								return healthSliceCopy.ToKey()
							},
							cmp.Compare,
						)
					},
					cmp.Compare,
				)
			}
			return healthDist
		},
		cmp.Compare,
	)
}

func (a AttackProfile) ResolveProfile() prob.Dist[UnitHealthStr] {
	attacks := a.attacks()

	hitOutcomes := a.hits(attacks)
	hits := prob.Map(
		hitOutcomes,
		func(outcome check.Outcome) int64 {
			return outcome.Successes()
		},
		cmp.Compare,
	)
	// TODO: [LETHAL HITS]
	// TODO: [SUSTAINED HITS]

	woundOutcomes := a.wounds(hits)
	wounds := prob.Map(
		woundOutcomes,
		func(outcome check.Outcome) int64 {
			return outcome.Successes()
		},
		cmp.Compare,
	)
	// TODO: [DEVASTATING WOUNDS]
	// TODO: mortal wounds

	health := a.resolveNormalWounds(wounds)

	return health
}
