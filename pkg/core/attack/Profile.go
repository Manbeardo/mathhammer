package attack

import (
	"slices"

	"github.com/Manbeardo/mathhammer/pkg/core/check"
	"github.com/Manbeardo/mathhammer/pkg/core/modifier"
	"github.com/Manbeardo/mathhammer/pkg/core/prob"
	"github.com/Manbeardo/mathhammer/pkg/core/unit"
	"github.com/Manbeardo/mathhammer/pkg/core/util"
	"github.com/Manbeardo/mathhammer/pkg/core/value"
)

type Profile struct {
	Attack
	AttackerWeaponProfile *unit.WeaponProfileKind
	AttackerWeaponCount   int
	DefenderHealth        prob.Dist[unit.Health]
}

func (a Profile) attacks() prob.Dist[int64] {
	if a.DistanceInches > a.AttackerWeaponProfile.Datasheet().RangeInches ||
		(a.DistanceInches == 0 && a.AttackerWeaponProfile.Datasheet().RangeInches > 0) {
		return value.Int(0).Distribution()
	}

	return value.Sum(
		slices.Repeat([]value.Interface{
			a.AttackerWeaponProfile.Datasheet().Attacks,
		}, int(a.AttackerWeaponCount))...,
	).Distribution()
}

func (a Profile) hits(attacks prob.Dist[int64]) prob.Dist[check.Outcome] {
	skill := a.AttackerWeaponProfile.Datasheet().Skill
	return check.Calculate(value.Roll(6), check.Opts{
		Count:                    attacks,
		SuccessTarget:            value.Int(skill),
		CriticalSuccessThreshold: 6,
		CriticalFailureThreshold: 1,
	})
}

func (a Profile) wounds(hits prob.Dist[int64]) prob.Dist[check.Outcome] {
	strengthDist := a.AttackerWeaponProfile.Datasheet().Strength.Distribution()
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
	))

	return check.Calculate(value.Roll(6), check.Opts{
		Count:                    hits,
		SuccessTarget:            targetDist,
		CriticalSuccessThreshold: 6,
		CriticalFailureThreshold: 1,
	})
}

func (a Profile) allocateWound(health unit.Health) (m *unit.Model, idx int) {
	// TODO: [PRECISION]
	for i, m := range a.AttackerUnit.Models() {
		if m.IsAlive(health) {
			return m, i
		}
	}
	return nil, -1
}

func (a Profile) resolveNormalWounds(woundDist prob.Dist[int64]) prob.Dist[unit.Health] {
	ap := a.AttackerWeaponProfile.Datasheet().ArmorPenetration
	saveModifiers := modifier.Set{
		modifier.Add(ap),
	}

	return util.Must(prob.FlatMap(
		woundDist,
		func(wounds int64) prob.Dist[unit.Health] {
			healthDist := a.DefenderHealth
			for range wounds {
				// TODO: memoize this
				healthDist = util.Must(prob.FlatMap(
					healthDist,
					func(health unit.Health) prob.Dist[unit.Health] {
						model, idx := a.allocateWound(health)
						if model == nil {
							return util.Must(prob.FromConst(health))
						}

						save := saveModifiers.Apply(modifier.ModelArmourSave, model.Datasheet().Save)

						checkDist := check.Calculate(value.Roll(6), check.Opts{
							SuccessTarget:            value.Int(save),
							CriticalFailureThreshold: 1,
						})

						return util.Must(prob.Map(
							checkDist,
							func(outcome check.Outcome) unit.Health {
								healthCopy := slices.Clone(health)
								damage := a.AttackerWeaponProfile.Datasheet().Damage
								for range outcome.Failures() {
									health := healthCopy[idx]
									if damage > health {
										healthCopy[idx] = 0
									} else {
										healthCopy[idx] -= damage
									}
								}
								return healthCopy
							},
						))
					},
				))
			}
			return healthDist
		},
	))
}

func (a Profile) ResolveProfile() prob.Dist[unit.Health] {
	attacks := a.attacks()

	hitOutcomes := a.hits(attacks)
	hits := util.Must(prob.Map(
		hitOutcomes,
		func(outcome check.Outcome) int64 {
			return outcome.Successes()
		},
	))
	// TODO: [LETHAL HITS]
	// TODO: [SUSTAINED HITS]

	woundOutcomes := a.wounds(hits)
	wounds := util.Must(prob.Map(
		woundOutcomes,
		func(outcome check.Outcome) int64 {
			return outcome.Successes()
		},
	))
	// TODO: [DEVASTATING WOUNDS]
	// TODO: mortal wounds

	health := a.resolveNormalWounds(wounds)

	return health
}
