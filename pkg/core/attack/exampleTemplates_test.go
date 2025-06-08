package attack

import (
	"github.com/Manbeardo/mathhammer/pkg/core"
	"github.com/Manbeardo/mathhammer/pkg/core/util"
)

func exampleUnitTpl_MEQ(count int) *core.UnitTemplate {
	return &core.UnitTemplate{
		Name:       "Marine Equivalent Squad",
		PointsCost: 100,
		Models: []util.Entry[*core.ModelTemplate, int]{
			{
				Key: &core.ModelTemplate{
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

func exampleUnitTpl_MEQWithWeaponProfile(count int, wep *core.WeaponProfileTemplate) *core.UnitTemplate {
	tpl := exampleUnitTpl_MEQ(count)
	tpl.Models[0].Key.Weapons = []util.Entry[*core.WeaponTemplate, int]{
		{
			Key: &core.WeaponTemplate{
				Name:     "Bullet Gun",
				Profiles: []*core.WeaponProfileTemplate{wep},
			},
			Value: 1,
		},
	}
	return tpl
}
