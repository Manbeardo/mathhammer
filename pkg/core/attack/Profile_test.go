package attack

import (
	"math/big"
	"testing"

	"github.com/Manbeardo/mathhammer/pkg/core/check"
	"github.com/Manbeardo/mathhammer/pkg/core/prob"
	"github.com/Manbeardo/mathhammer/pkg/core/unit"
	"github.com/Manbeardo/mathhammer/pkg/core/util"
	"github.com/Manbeardo/mathhammer/pkg/core/value"
	"github.com/stretchr/testify/assert"
)

func TestProfile(t *testing.T) {
	t.Run("attacks", func(t *testing.T) {
		t.Run("applies the weapon profile the correct number of times", func(t *testing.T) {
			battle := unit.NewBattle()
			defender := battle.NewUnit(exampleUnitTpl_MEQ(10))
			wep := unit.NewWeaponProfileTemplate(unit.WeaponProfileDatasheet{
				RangeInches:      12,
				Attacks:          value.Int(1),
				Skill:            3,
				Strength:         value.Int(4),
				ArmorPenetration: 1,
				Damage:           1,
			})
			attacker := battle.NewUnit(exampleUnitTpl_MEQWithWeaponProfile(10, wep))

			a := Profile{
				Attack: NewAttack(AttackOpts{
					AttackerUnit:   attacker,
					DefenderUnit:   defender,
					DistanceInches: 6,
				}),
				AttackerWeaponProfile: attacker.WeaponProfiles()[0][0].Kind(),
				AttackerWeaponCount:   10,
				DefenderHealth:        defender.StartingHealth().ToDist(),
			}

			attackDist := a.attacks()

			assert.Equal(t, util.Must(prob.FromConst(int64(10))), attackDist)
		})

		t.Run("0 attacks when attack is outside of weapon range", func(t *testing.T) {
			battle := unit.NewBattle()
			defender := battle.NewUnit(exampleUnitTpl_MEQ(10))
			wep := unit.NewWeaponProfileTemplate(unit.WeaponProfileDatasheet{
				RangeInches:      10,
				Attacks:          value.Int(1),
				Skill:            3,
				Strength:         value.Int(4),
				ArmorPenetration: 1,
				Damage:           1,
			})
			attacker := battle.NewUnit(exampleUnitTpl_MEQWithWeaponProfile(10, wep))

			a := Profile{
				Attack: NewAttack(AttackOpts{
					AttackerUnit:   attacker,
					DefenderUnit:   defender,
					DistanceInches: 11,
				}),
				AttackerWeaponProfile: attacker.WeaponProfiles()[0][0].Kind(),
				AttackerWeaponCount:   10,
				DefenderHealth:        defender.StartingHealth().ToDist(),
			}

			attackDist := a.attacks()

			assert.Equal(t, util.Must(prob.FromConst(int64(0))), attackDist)
		})

		t.Run("0 attacks when using melee weapon in ranged attack", func(t *testing.T) {
			battle := unit.NewBattle()
			defender := battle.NewUnit(exampleUnitTpl_MEQ(10))
			wep := unit.NewWeaponProfileTemplate(unit.WeaponProfileDatasheet{
				RangeInches:      12,
				Attacks:          value.Int(1),
				Skill:            3,
				Strength:         value.Int(4),
				ArmorPenetration: 1,
				Damage:           1,
			})
			attacker := battle.NewUnit(exampleUnitTpl_MEQWithWeaponProfile(10, wep))

			a := Profile{
				Attack: NewAttack(AttackOpts{
					AttackerUnit:   attacker,
					DefenderUnit:   defender,
					DistanceInches: 0,
				}),
				AttackerWeaponProfile: attacker.WeaponProfiles()[0][0].Kind(),
				AttackerWeaponCount:   10,
				DefenderHealth:        defender.StartingHealth().ToDist(),
			}

			attackDist := a.attacks()

			assert.Equal(t, util.Must(prob.FromConst(int64(0))), attackDist)
		})

		t.Run("handles random values correctly", func(t *testing.T) {
			battle := unit.NewBattle()
			defender := battle.NewUnit(exampleUnitTpl_MEQ(10))
			wep := unit.NewWeaponProfileTemplate(unit.WeaponProfileDatasheet{
				RangeInches:      12,
				Attacks:          value.Roll(2),
				Skill:            0,
				Strength:         value.Int(5),
				ArmorPenetration: 1,
				Damage:           1,
			})
			attacker := battle.NewUnit(exampleUnitTpl_MEQWithWeaponProfile(4, wep))

			a := Profile{
				Attack: NewAttack(AttackOpts{
					AttackerUnit:   attacker,
					DefenderUnit:   defender,
					DistanceInches: 6,
				}),
				AttackerWeaponProfile: attacker.WeaponProfiles()[0][0].Kind(),
				AttackerWeaponCount:   4,
				DefenderHealth:        defender.StartingHealth().ToDist(),
			}

			attackDist := a.attacks()

			assert.Equal(t, util.Must(prob.FromMap(prob.MapT[int64]{
				8: big.NewRat(1, 16),
				7: big.NewRat(4, 16),
				6: big.NewRat(6, 16),
				5: big.NewRat(4, 16),
				4: big.NewRat(1, 16),
			})), attackDist)
		})
	})

	t.Run("hits", func(t *testing.T) {
		t.Run("calculates hits correctly for a simple example", func(t *testing.T) {
			battle := unit.NewBattle()
			defender := battle.NewUnit(exampleUnitTpl_MEQ(10))
			wep := unit.NewWeaponProfileTemplate(unit.WeaponProfileDatasheet{
				RangeInches:      12,
				Attacks:          value.Int(1),
				Skill:            4,
				Strength:         value.Int(4),
				ArmorPenetration: 1,
				Damage:           1,
			})
			attacker := battle.NewUnit(exampleUnitTpl_MEQWithWeaponProfile(3, wep))

			a := Profile{
				Attack: NewAttack(AttackOpts{
					AttackerUnit:   attacker,
					DefenderUnit:   defender,
					DistanceInches: 6,
				}),
				AttackerWeaponProfile: attacker.WeaponProfiles()[0][0].Kind(),
				AttackerWeaponCount:   3,
				DefenderHealth:        defender.StartingHealth().ToDist(),
			}

			hitDist := util.Must(prob.Map(
				a.hits(value.Int(3).Distribution()),
				func(o check.Outcome) int64 { return o.Successes() },
			))

			assert.Equal(t, util.Must(prob.FromMap(prob.MapT[int64]{
				3: big.NewRat(1, 8),
				2: big.NewRat(3, 8),
				1: big.NewRat(3, 8),
				0: big.NewRat(1, 8),
			})), hitDist)
		})

		t.Run("calculates hits correctly for a simple example with random attacks", func(t *testing.T) {
			battle := unit.NewBattle()
			defender := battle.NewUnit(exampleUnitTpl_MEQ(10))
			wep := unit.NewWeaponProfileTemplate(unit.WeaponProfileDatasheet{
				RangeInches:      12,
				Attacks:          value.Int(1),
				Skill:            4,
				Strength:         value.Int(4),
				ArmorPenetration: 1,
				Damage:           1,
			})
			attacker := battle.NewUnit(exampleUnitTpl_MEQWithWeaponProfile(3, wep))

			a := Profile{
				Attack: NewAttack(AttackOpts{
					AttackerUnit:   attacker,
					DefenderUnit:   defender,
					DistanceInches: 6,
				}),
				AttackerWeaponProfile: attacker.WeaponProfiles()[0][0].Kind(),
				AttackerWeaponCount:   3,
				DefenderHealth:        defender.StartingHealth().ToDist(),
			}

			hitDist := util.Must(prob.Map(
				a.hits(value.Roll(2).Distribution()),
				func(o check.Outcome) int64 { return o.Successes() },
			))

			assert.Equal(t, util.Must(prob.FromMap(prob.MapT[int64]{
				2: big.NewRat(1, 8),
				1: big.NewRat(4, 8),
				0: big.NewRat(3, 8),
			})), hitDist)
		})
	})

	t.Run("wounds", func(t *testing.T) {
		t.Run("calculates wounds correctly for a simple example", func(t *testing.T) {
			battle := unit.NewBattle()
			defender := battle.NewUnit(exampleUnitTpl_MEQ(10))
			wep := unit.NewWeaponProfileTemplate(unit.WeaponProfileDatasheet{
				RangeInches:      12,
				Attacks:          value.Int(1),
				Skill:            4,
				Strength:         value.Int(3),
				ArmorPenetration: 1,
				Damage:           1,
			})
			attacker := battle.NewUnit(exampleUnitTpl_MEQWithWeaponProfile(3, wep))

			a := Profile{
				Attack: NewAttack(AttackOpts{
					AttackerUnit:   attacker,
					DefenderUnit:   defender,
					DistanceInches: 6,
				}),
				AttackerWeaponProfile: attacker.WeaponProfiles()[0][0].Kind(),
				AttackerWeaponCount:   3,
				DefenderHealth:        defender.StartingHealth().ToDist(),
			}

			woundDist := util.Must(prob.Map(
				a.wounds(value.Int(3).Distribution()),
				func(o check.Outcome) int64 { return o.Successes() },
			))

			assert.Equal(t, util.Must(prob.FromMap(prob.MapT[int64]{
				3: big.NewRat(1, 27),
				2: big.NewRat(6, 27),
				1: big.NewRat(12, 27),
				0: big.NewRat(8, 27),
			})), woundDist)
		})

		t.Run("calculates wounds correctly for a simple example with random attacks", func(t *testing.T) {
			battle := unit.NewBattle()
			defender := battle.NewUnit(exampleUnitTpl_MEQ(10))
			wep := unit.NewWeaponProfileTemplate(unit.WeaponProfileDatasheet{
				RangeInches:      12,
				Attacks:          value.Int(1),
				Skill:            4,
				Strength:         value.Int(3),
				ArmorPenetration: 1,
				Damage:           1,
			})
			attacker := battle.NewUnit(exampleUnitTpl_MEQWithWeaponProfile(3, wep))

			a := Profile{
				Attack: NewAttack(AttackOpts{
					AttackerUnit:   attacker,
					DefenderUnit:   defender,
					DistanceInches: 6,
				}),
				AttackerWeaponProfile: attacker.WeaponProfiles()[0][0].Kind(),
				AttackerWeaponCount:   3,
				DefenderHealth:        defender.StartingHealth().ToDist(),
			}

			woundDist := util.Must(prob.Map(
				a.wounds(value.Roll(2).Distribution()),
				func(o check.Outcome) int64 { return o.Successes() },
			))

			assert.Equal(t, util.Must(prob.FromMap(prob.MapT[int64]{
				2: big.NewRat(1, 18),
				1: big.NewRat(7, 18),
				0: big.NewRat(10, 18),
			})), woundDist)
		})
	})

	t.Run("resolveNormalWounds", func(t *testing.T) {
		t.Run("calculates health correctly for a simple example", func(t *testing.T) {
			battle := unit.NewBattle()
			defender := battle.NewUnit(exampleUnitTpl_MEQ(2))
			wep := unit.NewWeaponProfileTemplate(unit.WeaponProfileDatasheet{
				RangeInches:      12,
				Attacks:          value.Int(1),
				Skill:            4,
				Strength:         value.Int(3),
				ArmorPenetration: 1,
				Damage:           1,
			})
			attacker := battle.NewUnit(exampleUnitTpl_MEQWithWeaponProfile(2, wep))

			a := Profile{
				Attack: NewAttack(AttackOpts{
					AttackerUnit:   attacker,
					DefenderUnit:   defender,
					DistanceInches: 6,
				}),
				AttackerWeaponProfile: attacker.WeaponProfiles()[0][0].Kind(),
				AttackerWeaponCount:   3,
				DefenderHealth:        defender.StartingHealth().ToDist(),
			}

			woundDist := util.Must(prob.Map(
				a.resolveNormalWounds(value.Int(2).Distribution()),
				func(s unit.Health) int64 { return s.WoundsRemaining() },
			))

			assert.Equal(t, util.Must(prob.FromMap(prob.MapT[int64]{
				4: big.NewRat(1, 4),
				3: big.NewRat(2, 4),
				2: big.NewRat(1, 4),
			})), woundDist)
		})
	})

	t.Run("ResolveProfile", func(t *testing.T) {
		t.Run("works correctly in a basic MEQ example", func(t *testing.T) {
			battle := unit.NewBattle()
			defender := battle.NewUnit(exampleUnitTpl_MEQ(2))
			wep := unit.NewWeaponProfileTemplate(unit.WeaponProfileDatasheet{
				RangeInches:      12,
				Attacks:          value.Int(2),
				Skill:            3,
				Strength:         value.Int(4),
				ArmorPenetration: 1,
				Damage:           1,
			})
			attacker := battle.NewUnit(exampleUnitTpl_MEQWithWeaponProfile(2, wep))

			a := Profile{
				Attack: NewAttack(AttackOpts{
					AttackerUnit:   attacker,
					DefenderUnit:   defender,
					DistanceInches: 6,
				}),
				AttackerWeaponProfile: attacker.WeaponProfiles()[0][0].Kind(),
				AttackerWeaponCount:   2,
				DefenderHealth:        defender.StartingHealth().ToDist(),
			}

			healthDist := util.Must(prob.Map(
				a.ResolveProfile(),
				func(s unit.Health) int64 { return s.WoundsRemaining() },
			))

			assert.Equal(t, util.Must(prob.FromMap(prob.MapT[int64]{
				0: big.NewRat(1, 1296),   // unitcrunch: <0.5%
				1: big.NewRat(5, 324),    // unitcrunch: 1.5%
				2: big.NewRat(25, 216),   // unitcrunch: 11.3%
				3: big.NewRat(125, 324),  // unitcrunch: 38.2%
				4: big.NewRat(625, 1296), // unitcrunch: 49%
			})), healthDist)
		})
	})
}
