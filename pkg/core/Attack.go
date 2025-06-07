package core

import (
	"cmp"

	"github.com/Manbeardo/mathhammer/pkg/core/prob"
	"github.com/Manbeardo/mathhammer/pkg/core/util"
)

type Attack struct {
	AttackOpts
	DefenderToughness int64
}

type AttackOpts struct {
	AttackerUnit   *Unit
	DefenderUnit   *Unit
	DistanceInches int64
}

type AttackResult struct {
	AttackerHealth   prob.Dist[UnitHealthStr]
	DefenderHealth   prob.Dist[UnitHealthStr]
	SelectedProfiles []util.Entry[*WeaponProfileTemplate, int64]
}

func NewAttack(opts AttackOpts) Attack {
	return Attack{
		AttackOpts:        opts,
		DefenderToughness: opts.DefenderUnit.Toughness(),
	}
}

func (a Attack) ResolveAttack() AttackResult {
	weapons := util.OrderedEntries(
		a.AttackerUnit.Weapons(),
		func(a, b *WeaponTemplate) int {
			return cmp.Compare(a.Name, b.Name)
		},
	)
	defenderHealth := a.DefenderUnit.health.ToDist()
	selectedProfiles := []util.Entry[*WeaponProfileTemplate, int64]{}
	for _, e := range weapons {
		wtpl, count := e.Key, e.Value
		bestResult := defenderHealth
		bestWoundsRemaining := MeanWoundsRemaining(defenderHealth)
		var bestProfile *WeaponProfileTemplate
		for _, profile := range wtpl.Profiles {
			result := AttackProfile{
				Attack:                 a,
				AttackerWeaponProfile:  profile,
				AttackerWeaponCount:    count,
				DefenderStartingHealth: defenderHealth,
			}.ResolveProfile()
			woundsRemaining := MeanWoundsRemaining(result)
			if bestWoundsRemaining.Cmp(woundsRemaining) == 1 {
				bestResult = result
				bestWoundsRemaining = woundsRemaining
				bestProfile = profile
			}
		}
		defenderHealth = bestResult
		if bestProfile != nil {
			selectedProfiles = append(selectedProfiles, util.Entry[*WeaponProfileTemplate, int64]{
				Key:   bestProfile,
				Value: count,
			})
		}
	}
	return AttackResult{
		AttackerHealth:   a.AttackerUnit.health.ToDist(),
		DefenderHealth:   defenderHealth,
		SelectedProfiles: selectedProfiles,
	}
}
