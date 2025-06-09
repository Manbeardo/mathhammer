package attack

import (
	"github.com/Manbeardo/mathhammer/pkg/core/unit"
	"github.com/Manbeardo/mathhammer/pkg/core/util"
)

func exampleUnitTpl_MEQ(count int) *unit.UnitTemplate {
	return &unit.UnitTemplate{
		Name:       "Marine Equivalent Squad",
		PointsCost: 100,
		Models: []util.Entry[*unit.ModelTemplate, int]{
			{
				K: &unit.ModelTemplate{
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

func exampleUnitTpl_MEQWithWeaponProfile(count int, wep *unit.WeaponProfileTemplate) *unit.UnitTemplate {
	tpl := exampleUnitTpl_MEQ(count)
	tpl.Models[0].K.Weapons = []util.Entry[*unit.WeaponTemplate, int]{
		{
			K: &unit.WeaponTemplate{
				Name:     "Bullet Gun",
				Profiles: []*unit.WeaponProfileTemplate{wep},
			},
			V: 1,
		},
	}
	return tpl
}
