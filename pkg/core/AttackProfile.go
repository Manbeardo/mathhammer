package core

import (
	"cmp"
	"math/big"
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
	DefenderStartingHealth ModelHealth
}

func (a *AttackProfile) attacks() prob.Dist[int64] {
	return value.Sum(
		slices.Repeat([]value.Interface{
			a.AttackerWeaponProfile.Attacks,
		}, int(a.AttackerWeaponCount))...,
	).Distribution()
}

func (a *AttackProfile) hits(attacks prob.Dist[int64]) prob.Dist[check.Outcome] {
	skill := a.AttackerWeaponProfile.Skill
	return check.Calculate(value.Roll(6), check.Opts{
		Count:                    attacks,
		SuccessTarget:            value.Int(skill),
		CriticalSuccessThreshold: 6,
		CriticalFailureThreshold: 1,
	})
}

func (a *AttackProfile) wounds(hits prob.Dist[int64]) prob.Dist[check.Outcome] {
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

func (a *AttackProfile) allocateWound(healthSlice []int64) (m *Model, idx int) {
	// TODO: [PRECISION]
	for idx, health := range healthSlice {
		if health > 0 {
			return a.DefenderUnit.Model(idx), idx
		}
	}
	return nil, -1
}

func (a *AttackProfile) resolveNormalWounds(woundDist prob.Dist[int64]) prob.Dist[ModelHealthStr] {
	ap := a.AttackerWeaponProfile.ArmorPenetration
	saveModifiers := modifier.Set{
		modifier.Add(ap),
	}

	return prob.FlatMap(
		woundDist,
		func(wounds int64) prob.Dist[ModelHealthStr] {
			healthDist := prob.NewDistribution(map[ModelHealthStr]*big.Rat{
				(a.DefenderStartingHealth.ToKey()): big.NewRat(1, 1),
			})
			for range wounds {
				// TODO: memoize this
				healthDist = prob.FlatMap(
					healthDist,
					func(healthStr ModelHealthStr) prob.Dist[ModelHealthStr] {
						healthSlice := healthStr.ToSlice()
						model, idx := a.allocateWound(healthSlice)

						save := saveModifiers.Apply(modifier.ModelArmourSave, model.tpl.Save)

						checkDist := check.Calculate(value.Roll(6), check.Opts{
							SuccessTarget:            value.Int(save),
							CriticalFailureThreshold: 1,
						})

						return prob.Map(
							checkDist,
							func(outcome check.Outcome) ModelHealthStr {
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

func (a *AttackProfile) Resolve() prob.Dist[ModelHealthStr] {
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
