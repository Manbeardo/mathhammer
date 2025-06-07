package core

import (
	"cmp"
	"math/big"
	"testing"

	"github.com/Manbeardo/mathhammer/pkg/core/check"
	"github.com/Manbeardo/mathhammer/pkg/core/prob"
	"github.com/Manbeardo/mathhammer/pkg/core/value"
	"github.com/stretchr/testify/assert"
)

func TestAttackProfile(t *testing.T) {
	t.Run("attacks", func(t *testing.T) {
		t.Run("applies the weapon profile the correct number of times", func(t *testing.T) {
			defender := NewUnit(exampleUnitTpl_MEQ(10))
			wep := &WeaponProfileTemplate{
				Attacks:          value.Int(1),
				Skill:            3,
				Strength:         value.Int(4),
				ArmorPenetration: 1,
				Damage:           1,
			}
			attacker := NewUnit(exampleUnitTpl_MEQWithRangedWeapon(10, wep))

			a := &AttackProfile{
				Attack: Attack{
					AttackerUnit:      attacker,
					DefenderUnit:      defender,
					DefenderToughness: defender.Toughness(),
				},
				AttackerWeaponProfile:  wep,
				AttackerWeaponCount:    10,
				DefenderStartingHealth: defender.ModelHealth(),
			}

			attackDist := a.attacks()

			assert.Equal(t, prob.NewDistribution(map[int64]*big.Rat{
				10: big.NewRat(1, 1),
			}), attackDist)
		})

		t.Run("handles random values correctly", func(t *testing.T) {
			defender := NewUnit(exampleUnitTpl_MEQ(10))
			wep := &WeaponProfileTemplate{
				Attacks:          value.Roll(2),
				Skill:            0,
				Strength:         value.Int(5),
				ArmorPenetration: 1,
				Damage:           1,
			}
			attacker := NewUnit(exampleUnitTpl_MEQWithRangedWeapon(4, wep))

			a := &AttackProfile{
				Attack: Attack{
					AttackerUnit:      attacker,
					DefenderUnit:      defender,
					DefenderToughness: defender.Toughness(),
				},
				AttackerWeaponProfile:  wep,
				AttackerWeaponCount:    4,
				DefenderStartingHealth: defender.ModelHealth(),
			}

			attackDist := a.attacks()

			assert.Equal(t, prob.NewDistribution(map[int64]*big.Rat{
				8: big.NewRat(1, 16),
				7: big.NewRat(4, 16),
				6: big.NewRat(6, 16),
				5: big.NewRat(4, 16),
				4: big.NewRat(1, 16),
			}), attackDist)
		})
	})

	t.Run("hits", func(t *testing.T) {
		t.Run("calculates hits correctly for a simple example", func(t *testing.T) {
			defender := NewUnit(exampleUnitTpl_MEQ(10))
			wep := &WeaponProfileTemplate{
				Attacks:          value.Int(1),
				Skill:            4,
				Strength:         value.Int(4),
				ArmorPenetration: 1,
				Damage:           1,
			}
			attacker := NewUnit(exampleUnitTpl_MEQWithRangedWeapon(3, wep))

			a := &AttackProfile{
				Attack: Attack{
					AttackerUnit:      attacker,
					DefenderUnit:      defender,
					DefenderToughness: defender.Toughness(),
				},
				AttackerWeaponProfile:  wep,
				AttackerWeaponCount:    3,
				DefenderStartingHealth: defender.ModelHealth(),
			}

			hitDist := prob.Map(
				a.hits(value.Int(3).Distribution()),
				func(o check.Outcome) int64 { return o.Successes() },
				cmp.Compare,
			)

			assert.Equal(t, prob.NewDistribution(map[int64]*big.Rat{
				3: big.NewRat(1, 8),
				2: big.NewRat(3, 8),
				1: big.NewRat(3, 8),
				0: big.NewRat(1, 8),
			}), hitDist)
		})

		t.Run("calculates hits correctly for a simple example with random attacks", func(t *testing.T) {
			defender := NewUnit(exampleUnitTpl_MEQ(10))
			wep := &WeaponProfileTemplate{
				Attacks:          value.Int(1),
				Skill:            4,
				Strength:         value.Int(4),
				ArmorPenetration: 1,
				Damage:           1,
			}
			attacker := NewUnit(exampleUnitTpl_MEQWithRangedWeapon(3, wep))

			a := &AttackProfile{
				Attack: Attack{
					AttackerUnit:      attacker,
					DefenderUnit:      defender,
					DefenderToughness: defender.Toughness(),
				},
				AttackerWeaponProfile:  wep,
				AttackerWeaponCount:    3,
				DefenderStartingHealth: defender.ModelHealth(),
			}

			hitDist := prob.Map(
				a.hits(value.Roll(2).Distribution()),
				func(o check.Outcome) int64 { return o.Successes() },
				cmp.Compare,
			)

			assert.Equal(t, prob.NewDistribution(map[int64]*big.Rat{
				2: big.NewRat(1, 8),
				1: big.NewRat(4, 8),
				0: big.NewRat(3, 8),
			}), hitDist)
		})
	})

	t.Run("wounds", func(t *testing.T) {
		t.Run("calculates wounds correctly for a simple example", func(t *testing.T) {
			defender := NewUnit(exampleUnitTpl_MEQ(10))
			wep := &WeaponProfileTemplate{
				Attacks:          value.Int(1),
				Skill:            4,
				Strength:         value.Int(3),
				ArmorPenetration: 1,
				Damage:           1,
			}
			attacker := NewUnit(exampleUnitTpl_MEQWithRangedWeapon(3, wep))

			a := &AttackProfile{
				Attack: Attack{
					AttackerUnit:      attacker,
					DefenderUnit:      defender,
					DefenderToughness: defender.Toughness(),
				},
				AttackerWeaponProfile:  wep,
				AttackerWeaponCount:    3,
				DefenderStartingHealth: defender.ModelHealth(),
			}

			woundDist := prob.Map(
				a.wounds(value.Int(3).Distribution()),
				func(o check.Outcome) int64 { return o.Successes() },
				cmp.Compare,
			)

			assert.Equal(t, prob.NewDistribution(map[int64]*big.Rat{
				3: big.NewRat(1, 27),
				2: big.NewRat(6, 27),
				1: big.NewRat(12, 27),
				0: big.NewRat(8, 27),
			}), woundDist)
		})

		t.Run("calculates wounds correctly for a simple example with random attacks", func(t *testing.T) {
			defender := NewUnit(exampleUnitTpl_MEQ(10))
			wep := &WeaponProfileTemplate{
				Attacks:          value.Int(1),
				Skill:            4,
				Strength:         value.Int(3),
				ArmorPenetration: 1,
				Damage:           1,
			}
			attacker := NewUnit(exampleUnitTpl_MEQWithRangedWeapon(3, wep))

			a := &AttackProfile{
				Attack: Attack{
					AttackerUnit:      attacker,
					DefenderUnit:      defender,
					DefenderToughness: defender.Toughness(),
				},
				AttackerWeaponProfile:  wep,
				AttackerWeaponCount:    3,
				DefenderStartingHealth: defender.ModelHealth(),
			}

			woundDist := prob.Map(
				a.wounds(value.Roll(2).Distribution()),
				func(o check.Outcome) int64 { return o.Successes() },
				cmp.Compare,
			)

			assert.Equal(t, prob.NewDistribution(map[int64]*big.Rat{
				2: big.NewRat(1, 18),
				1: big.NewRat(7, 18),
				0: big.NewRat(10, 18),
			}), woundDist)
		})
	})

	t.Run("resolveNormalWounds", func(t *testing.T) {
		t.Run("calculates health correctly for a simple example", func(t *testing.T) {
			defender := NewUnit(exampleUnitTpl_MEQ(2))
			wep := &WeaponProfileTemplate{
				Attacks:          value.Int(1),
				Skill:            4,
				Strength:         value.Int(3),
				ArmorPenetration: 1,
				Damage:           1,
			}
			attacker := NewUnit(exampleUnitTpl_MEQWithRangedWeapon(2, wep))

			a := &AttackProfile{
				Attack: Attack{
					AttackerUnit:      attacker,
					DefenderUnit:      defender,
					DefenderToughness: defender.Toughness(),
				},
				AttackerWeaponProfile:  wep,
				AttackerWeaponCount:    3,
				DefenderStartingHealth: defender.ModelHealth(),
			}

			woundDist := prob.Map(
				a.resolveNormalWounds(value.Int(2).Distribution()),
				func(s ModelHealthStr) int64 { return s.ToSlice().WoundsRemaining() },
				cmp.Compare,
			)

			assert.Equal(t, prob.NewDistribution(map[int64]*big.Rat{
				4: big.NewRat(1, 4),
				3: big.NewRat(2, 4),
				2: big.NewRat(1, 4),
			}), woundDist)
		})
	})

	t.Run("Resolve", func(t *testing.T) {
		t.Run("works correctly in a basic MEQ example", func(t *testing.T) {
			defender := NewUnit(exampleUnitTpl_MEQ(2))
			wep := &WeaponProfileTemplate{
				Attacks:          value.Int(2),
				Skill:            3,
				Strength:         value.Int(4),
				ArmorPenetration: 1,
				Damage:           1,
			}
			attacker := NewUnit(exampleUnitTpl_MEQWithRangedWeapon(2, wep))

			a := &AttackProfile{
				Attack: Attack{
					AttackerUnit:      attacker,
					DefenderUnit:      defender,
					DefenderToughness: defender.Toughness(),
				},
				AttackerWeaponProfile:  wep,
				AttackerWeaponCount:    2,
				DefenderStartingHealth: defender.ModelHealth(),
			}

			healthDist := prob.Map(
				a.Resolve(),
				func(s ModelHealthStr) int64 { return s.ToSlice().WoundsRemaining() },
				cmp.Compare,
			)

			assert.Equal(t, prob.NewDistribution(map[int64]*big.Rat{
				0: big.NewRat(1, 1296),   // unitcrunch: <0.5%
				1: big.NewRat(5, 324),    // unitcrunch: 1.5%
				2: big.NewRat(25, 216),   // unitcrunch: 11.3%
				3: big.NewRat(125, 324),  // unitcrunch: 38.2%
				4: big.NewRat(625, 1296), // unitcrunch: 49%
			}), healthDist)
		})
	})
}
