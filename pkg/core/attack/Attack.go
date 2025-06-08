package attack

import (
	"github.com/Manbeardo/mathhammer/pkg/core"
	"github.com/Manbeardo/mathhammer/pkg/core/prob"
	"github.com/Manbeardo/mathhammer/pkg/core/util"
)

type Attack struct {
	AttackOpts
	DefenderToughness int64
}

type AttackOpts struct {
	AttackerUnit          *core.Unit
	DefenderUnit          *core.Unit
	InitialAttackerHealth core.UnitHealth
	InitialDefenderHealth core.UnitHealth
	DistanceInches        int64
}

type Result struct {
	AttackerHealth   prob.Dist[core.UnitHealth]
	DefenderHealth   prob.Dist[core.UnitHealth]
	SelectedProfiles []util.Entry[*core.WeaponProfileTemplate, int64]
}

func NewAttack(opts AttackOpts) Attack {
	if opts.InitialAttackerHealth == nil {
		opts.InitialAttackerHealth = opts.AttackerUnit.StartingHealth()
	}
	if opts.InitialDefenderHealth == nil {
		opts.InitialDefenderHealth = opts.DefenderUnit.StartingHealth()
	}
	return Attack{
		AttackOpts:        opts,
		DefenderToughness: opts.DefenderUnit.Toughness(opts.InitialDefenderHealth),
	}
}

func (a Attack) ResolveAttack() Result {
	selectedProfiles := []util.Entry[*core.WeaponProfileTemplate, int64]{}
	defenderHealth := a.InitialDefenderHealth.ToDist()
	for _, wtpl := range a.AttackerUnit.WeaponTemplates() {
		count := wtpl.AvailableCount(a.AttackerUnit, a.InitialAttackerHealth)
		bestResult := defenderHealth
		bestWoundsRemaining := core.MeanWoundsRemaining(defenderHealth)
		var bestProfile *core.WeaponProfileTemplate
		for _, profile := range wtpl.Profiles {
			result := Profile{
				Attack:                a,
				AttackerWeaponProfile: profile,
				AttackerWeaponCount:   count,
				DefenderHealth:        defenderHealth,
			}.ResolveProfile()
			woundsRemaining := core.MeanWoundsRemaining(result)
			if bestWoundsRemaining.Cmp(woundsRemaining) == 1 {
				bestResult = result
				bestWoundsRemaining = woundsRemaining
				bestProfile = profile
			}
		}
		defenderHealth = bestResult
		if bestProfile != nil {
			selectedProfiles = append(selectedProfiles, util.Entry[*core.WeaponProfileTemplate, int64]{
				K: bestProfile,
				V: count,
			})
		}
	}
	return Result{
		AttackerHealth:   a.InitialAttackerHealth.ToDist(),
		DefenderHealth:   defenderHealth,
		SelectedProfiles: selectedProfiles,
	}
}
