package core

import (
	"github.com/Manbeardo/mathhammer/pkg/core/util"
)

func exampleUnitTpl_MEQ(count int) *UnitTemplate {
	return &UnitTemplate{
		Name:       "Marine Equivalent Squad",
		PointsCost: 100,
		Models: []util.EntryT[*ModelTemplate, int]{
			util.Entry(&ModelTemplate{
				Name:       "Jimmy Space",
				Toughness:  4,
				Save:       3,
				Wounds:     2,
				Leadership: 6,
			}, count),
		},
	}
}

func exampleUnitTpl_MEQWithRangedWeapon(count int, wep *WeaponProfileTemplate) *UnitTemplate {
	tpl := exampleUnitTpl_MEQ(count)
	tpl.Models[0].Key.Weapons = []util.EntryT[*WeaponTemplate, int]{
		util.Entry(&WeaponTemplate{
			Name:     "Bullet Gun",
			Kind:     RangedAttack,
			Profiles: []*WeaponProfileTemplate{wep},
		}, 1),
	}
	return tpl
}
