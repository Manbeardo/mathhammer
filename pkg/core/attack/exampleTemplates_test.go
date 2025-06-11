package attack

import (
	"slices"

	"github.com/Manbeardo/mathhammer/pkg/core/unit"
)

func exampleUnitTpl_MEQ(count int) *unit.UnitTemplate {
	return unit.NewUnitTemplate(
		unit.UnitDatasheet{
			Name:       "Marine Equivalent Squad",
			PointsCost: 100,
		},
		slices.Repeat([]*unit.ModelTemplate{
			unit.NewModelTemplate(unit.ModelDatasheet{
				Name:       "Jimmy Space",
				Toughness:  4,
				Save:       3,
				Wounds:     2,
				Leadership: 6,
			}),
		}, count)...,
	)
}

func exampleUnitTpl_MEQWithWeaponProfile(count int, wep *unit.WeaponProfileTemplate) *unit.UnitTemplate {
	tpl := exampleUnitTpl_MEQ(count)
	tpl.Models[0].Weapons = slices.Repeat(
		[]*unit.WeaponTemplate{unit.NewWeaponTemplate(
			unit.WeaponDatasheet{
				Name: "Bullet Gun",
			},
			wep,
		)},
		count,
	)
	return tpl
}
