package core

import (
	"github.com/Manbeardo/mathhammer/pkg/core/prob"
	"github.com/Manbeardo/mathhammer/pkg/core/util"
)

type Attack struct {
	AttackOpts
	DefenderToughness int64
}

type AttackOpts struct {
	AttackerUnit          *Unit
	DefenderUnit          *Unit
	InitialAttackerHealth UnitHealth
	InitialDefenderHealth UnitHealth
	DistanceInches        int64
}

type AttackResult struct {
	AttackerHealth   prob.Dist[UnitHealthStr]
	DefenderHealth   prob.Dist[UnitHealthStr]
	SelectedProfiles []util.Entry[*WeaponProfileTemplate, int64]
}

func NewAttack(opts AttackOpts) Attack {
	if opts.InitialAttackerHealth == nil {
		opts.InitialAttackerHealth = opts.AttackerUnit.startingHealth
	}
	if opts.InitialDefenderHealth == nil {
		opts.InitialDefenderHealth = opts.DefenderUnit.startingHealth
	}
	return Attack{
		AttackOpts:        opts,
		DefenderToughness: opts.DefenderUnit.Toughness(opts.InitialDefenderHealth),
	}
}

func (a Attack) ResolveAttack() AttackResult {
	selectedProfiles := []util.Entry[*WeaponProfileTemplate, int64]{}
	defenderHealth := a.InitialDefenderHealth.ToDist()
	for _, wtpl := range a.AttackerUnit.WeaponTemplates() {
		count := wtpl.AvailableCount(a.AttackerUnit, a.InitialAttackerHealth)
		bestResult := defenderHealth
		bestWoundsRemaining := MeanWoundsRemaining(defenderHealth)
		var bestProfile *WeaponProfileTemplate
		for _, profile := range wtpl.Profiles {
			result := AttackProfile{
				Attack:                a,
				AttackerWeaponProfile: profile,
				AttackerWeaponCount:   count,
				DefenderHealth:        defenderHealth,
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
		AttackerHealth:   a.InitialAttackerHealth.ToDist(),
		DefenderHealth:   defenderHealth,
		SelectedProfiles: selectedProfiles,
	}
}
