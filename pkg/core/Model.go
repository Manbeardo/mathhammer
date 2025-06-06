package core

type Model struct {
	tpl         *ModelTemplate
	weapons     []*Weapon
	woundsTaken int64
}

func NewModel(tpl *ModelTemplate) *Model {
	m := &Model{
		tpl: tpl,
	}
	for wtpl, count := range tpl.Weapons {
		for range count {
			m.weapons = append(m.weapons, NewWeapon(wtpl))
		}
	}
	return m
}

func (m *Model) IsDead() bool {
	return m.woundsTaken >= m.tpl.Wounds
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
	Weapons    map[*WeaponTemplate]int64
	Abilities  []Ability
}
