package core

import (
	"github.com/Manbeardo/mathhammer/pkg/core/util"
)

func exampleUnitTpl_MEQ(count int) *UnitTemplate {
	return &UnitTemplate{
		Name:       "Marine Equivalent Squad",
		PointsCost: 100,
		Models: []util.Entry[*ModelTemplate, int]{
			{
				Key: &ModelTemplate{
					Name:       "Jimmy Space",
					Toughness:  4,
					Save:       3,
					Wounds:     2,
					Leadership: 6,
				},
				Value: count,
			},
		},
	}
}

func exampleUnitTpl_MEQWithWeaponProfile(count int, wep *WeaponProfileTemplate) *UnitTemplate {
	tpl := exampleUnitTpl_MEQ(count)
	tpl.Models[0].Key.Weapons = []util.Entry[*WeaponTemplate, int]{
		{
			Key: &WeaponTemplate{
				Name:     "Bullet Gun",
				Profiles: []*WeaponProfileTemplate{wep},
			},
			Value: 1,
		},
	}
	return tpl
}
