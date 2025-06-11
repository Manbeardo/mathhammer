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
			battle := unit.NewBattle()
			defender := battle.NewUnit(exampleUnitTpl_MEQ(2))
			attackerTpl := exampleUnitTpl_MEQ(1)

			strongProfile := unit.NewWeaponProfileTemplate(unit.WeaponProfileDatasheet{
				Name:             "killy",
				RangeInches:      12,
				Attacks:          value.Int(2),
				Skill:            2,
				Strength:         value.Int(8),
				ArmorPenetration: 4,
				Damage:           2,
			})
			weakProfile1 := unit.NewWeaponProfileTemplate(unit.WeaponProfileDatasheet{
				Name:             "spray",
				RangeInches:      12,
				Attacks:          value.Int(10),
				Skill:            5,
				Strength:         value.Int(1),
				ArmorPenetration: 0,
				Damage:           1,
			})
			weakProfile2 := unit.NewWeaponProfileTemplate(unit.WeaponProfileDatasheet{
				Name:             "bayonetto",
				RangeInches:      2,
				Attacks:          value.Int(1),
				Skill:            3,
				Strength:         value.Int(8),
				ArmorPenetration: 1,
				Damage:           2,
			})

			attackerTpl.Models[0].Weapons = []*unit.WeaponTemplate{
				unit.NewWeaponTemplate(
					unit.WeaponDatasheet{Name: "BIG GUN"},
					weakProfile1,
					strongProfile,
					weakProfile2,
				),
			}

			attacker := battle.NewUnit(attackerTpl)

			result := NewAttack(AttackOpts{
				AttackerUnit:   attacker,
				DefenderUnit:   defender,
				DistanceInches: 6,
			}).ResolveAttack()

			assert.Equal(t,
				[]util.Entry[*unit.WeaponProfileKind, int]{
					{K: attacker.WeaponProfiles()[0][1].Kind(), V: 1},
				},
				util.EntriesFromSeq2(result.SelectedProfiles.Iter()),
			)
		})
	})
}
