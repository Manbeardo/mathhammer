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
				K: &core.ModelTemplate{
					Name:       "Jimmy Space",
					Toughness:  4,
					Save:       3,
					Wounds:     2,
					Leadership: 6,
				},
				V: count,
			},
		},
	}
}

func exampleUnitTpl_MEQWithWeaponProfile(count int, wep *core.WeaponProfileTemplate) *core.UnitTemplate {
	tpl := exampleUnitTpl_MEQ(count)
	tpl.Models[0].K.Weapons = []util.Entry[*core.WeaponTemplate, int]{
		{
			K: &core.WeaponTemplate{
				Name:     "Bullet Gun",
				Profiles: []*core.WeaponProfileTemplate{wep},
			},
			V: 1,
		},
	}
	return tpl
}
