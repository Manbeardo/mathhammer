package core

import (
	"testing"

	"github.com/Manbeardo/mathhammer/pkg/core/util"
	"github.com/Manbeardo/mathhammer/pkg/core/value"
	"github.com/stretchr/testify/assert"
)

func TestAttack(t *testing.T) {
	t.Run("ResolveAttack", func(t *testing.T) {
		t.Run("selects the best profile from weapons with multiple profiles", func(t *testing.T) {
			defender := NewUnit(exampleUnitTpl_MEQ(2))
			attackerTpl := exampleUnitTpl_MEQ(1)

			strongProfile := &WeaponProfileTemplate{
				Name:             "killy",
				RangeInches:      12,
				Attacks:          value.Int(2),
				Skill:            2,
				Strength:         value.Int(8),
				ArmorPenetration: 4,
				Damage:           2,
			}
			weakProfile1 := &WeaponProfileTemplate{
				Name:             "spray",
				RangeInches:      12,
				Attacks:          value.Int(10),
				Skill:            5,
				Strength:         value.Int(1),
				ArmorPenetration: 0,
				Damage:           1,
			}
			weakProfile2 := &WeaponProfileTemplate{
				Name:             "bayonetto",
				RangeInches:      2,
				Attacks:          value.Int(1),
				Skill:            3,
				Strength:         value.Int(8),
				ArmorPenetration: 1,
				Damage:           2,
			}

			attackerTpl.Models[0].Key.Weapons = []util.Entry[*WeaponTemplate, int]{
				{
					Key: &WeaponTemplate{
						Name: "BIG GUN",
						Profiles: []*WeaponProfileTemplate{
							weakProfile1,
							strongProfile,
							weakProfile2,
						},
					},
					Value: 1,
				},
			}
			attacker := NewUnit(attackerTpl)

			result := NewAttack(AttackOpts{
				AttackerUnit:   attacker,
				DefenderUnit:   defender,
				DistanceInches: 6,
			}).ResolveAttack()

			assert.Equal(t,
				[]util.Entry[*WeaponProfileTemplate, int64]{
					{Key: strongProfile, Value: 1},
				},
				result.SelectedProfiles,
			)
		})
	})
}
