package attack

import (
	"cmp"
	"slices"

	"github.com/Manbeardo/mathhammer/pkg/core"
	"github.com/Manbeardo/mathhammer/pkg/core/check"
	"github.com/Manbeardo/mathhammer/pkg/core/modifier"
	"github.com/Manbeardo/mathhammer/pkg/core/prob"
	"github.com/Manbeardo/mathhammer/pkg/core/util"
	"github.com/Manbeardo/mathhammer/pkg/core/value"
)

type AttackProfile struct {
	Attack
	AttackerWeaponProfile *core.WeaponProfileTemplate
	AttackerWeaponCount   int64
	DefenderHealth        prob.Dist[core.UnitHealthStr]
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
	targetDist := util.Must(prob.Map(
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
	))

	return check.Calculate(value.Roll(6), check.Opts{
		Count:                    hits,
		SuccessTarget:            targetDist,
		CriticalSuccessThreshold: 6,
		CriticalFailureThreshold: 1,
	})
}

func (a AttackProfile) allocateWound(healthSlice []int64) (m *core.Model, idx int) {
	// TODO: [PRECISION]
	for idx, health := range healthSlice {
		if health > 0 {
			return a.DefenderUnit.Model(idx), idx
		}
	}
	return nil, -1
}

func (a AttackProfile) resolveNormalWounds(woundDist prob.Dist[int64]) prob.Dist[core.UnitHealthStr] {
	ap := a.AttackerWeaponProfile.ArmorPenetration
	saveModifiers := modifier.Set{
		modifier.Add(ap),
	}

	return util.Must(prob.FlatMap(
		woundDist,
		func(wounds int64) prob.Dist[core.UnitHealthStr] {
			healthDist := a.DefenderHealth
			for range wounds {
				// TODO: memoize this
				healthDist = util.Must(prob.FlatMap(
					healthDist,
					func(healthStr core.UnitHealthStr) prob.Dist[core.UnitHealthStr] {
						healthSlice := healthStr.ToSlice()
						model, idx := a.allocateWound(healthSlice)
						if model == nil {
							return util.Must(prob.NewConstDist(healthStr))
						}

						save := saveModifiers.Apply(modifier.ModelArmourSave, model.Save())

						checkDist := check.Calculate(value.Roll(6), check.Opts{
							SuccessTarget:            value.Int(save),
							CriticalFailureThreshold: 1,
						})

						return util.Must(prob.Map(
							checkDist,
							func(outcome check.Outcome) core.UnitHealthStr {
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
						))
					},
					cmp.Compare,
				))
			}
			return healthDist
		},
		cmp.Compare,
	))
}

func (a AttackProfile) ResolveProfile() prob.Dist[core.UnitHealthStr] {
	attacks := a.attacks()

	hitOutcomes := a.hits(attacks)
	hits := util.Must(prob.Map(
		hitOutcomes,
		func(outcome check.Outcome) int64 {
			return outcome.Successes()
		},
		cmp.Compare,
	))
	// TODO: [LETHAL HITS]
	// TODO: [SUSTAINED HITS]

	woundOutcomes := a.wounds(hits)
	wounds := util.Must(prob.Map(
		woundOutcomes,
		func(outcome check.Outcome) int64 {
			return outcome.Successes()
		},
		cmp.Compare,
	))
	// TODO: [DEVASTATING WOUNDS]
	// TODO: mortal wounds

	health := a.resolveNormalWounds(wounds)

	return health
}
