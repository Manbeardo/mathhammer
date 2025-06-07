package core

import "github.com/Manbeardo/mathhammer/pkg/core/util"

type Model struct {
	tpl     *ModelTemplate
	weapons []*Weapon
}

func NewModel(tpl *ModelTemplate) *Model {
	m := &Model{
		tpl: tpl,
	}
	for _, entry := range tpl.Weapons {
		wtpl, count := entry.Key, entry.Value
		for range count {
			m.weapons = append(m.weapons, NewWeapon(wtpl))
		}
	}
	return m
}

func (m *Model) MatchingWeaponProfiles(tpl *WeaponProfileTemplate) []*WeaponProfile {
	out := []*WeaponProfile{}
	for _, weapon := range m.weapons {
		p := weapon.MatchingProfile(tpl)
		if p != nil {
			out = append(out, p)
		}
	}
	return out
}

func (m *Model) Abilities() []Ability {
	return m.tpl.Abilities
}

type ModelTemplate struct {
	Name       string
	Toughness  int64
	Save       int64
	Wounds     int64
	Leadership int64
	Weapons    []util.Entry[*WeaponTemplate, int]
	Abilities  []Ability
}
