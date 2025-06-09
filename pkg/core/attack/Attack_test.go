package attack

import (
	"testing"

	"github.com/Manbeardo/mathhammer/pkg/core/unit"
	"github.com/Manbeardo/mathhammer/pkg/core/util"
	"github.com/Manbeardo/mathhammer/pkg/core/value"
	"github.com/stretchr/testify/assert"
)

func TestAttack(t *testing.T) {
	t.Run("ResolveAttack", func(t *testing.T) {
		t.Run("selects the best profile from weapons with multiple profiles", func(t *testing.T) {
			defender := unit.NewUnit(exampleUnitTpl_MEQ(2))
			attackerTpl := exampleUnitTpl_MEQ(1)

			strongProfile := &unit.WeaponProfileTemplate{
				Name:             "killy",
				RangeInches:      12,
				Attacks:          value.Int(2),
				Skill:            2,
				Strength:         value.Int(8),
				ArmorPenetration: 4,
				Damage:           2,
			}
			weakProfile1 := &unit.WeaponProfileTemplate{
				Name:             "spray",
				RangeInches:      12,
				Attacks:          value.Int(10),
				Skill:            5,
				Strength:         value.Int(1),
				ArmorPenetration: 0,
				Damage:           1,
			}
			weakProfile2 := &unit.WeaponProfileTemplate{
				Name:             "bayonetto",
				RangeInches:      2,
				Attacks:          value.Int(1),
				Skill:            3,
				Strength:         value.Int(8),
				ArmorPenetration: 1,
				Damage:           2,
			}

			attackerTpl.Models[0].K.Weapons = []util.Entry[*unit.WeaponTemplate, int]{
				{
					K: &unit.WeaponTemplate{
						Name: "BIG GUN",
						Profiles: []*unit.WeaponProfileTemplate{
							weakProfile1,
							strongProfile,
							weakProfile2,
						},
					},
					V: 1,
				},
			}
			attacker := unit.NewUnit(attackerTpl)

			result := NewAttack(AttackOpts{
				AttackerUnit:   attacker,
				DefenderUnit:   defender,
				DistanceInches: 6,
			}).ResolveAttack()

			assert.Equal(t,
				[]util.Entry[*unit.WeaponProfileTemplate, int64]{
					{K: strongProfile, V: 1},
				},
				result.SelectedProfiles,
			)
		})
	})
}
