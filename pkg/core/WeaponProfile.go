package core

import "github.com/Manbeardo/mathhammer/pkg/core/value"

type WeaponProfile struct {
	tpl          *WeaponProfileTemplate
	wasActivated bool
}

func NewWeaponProfile(tpl *WeaponProfileTemplate) *WeaponProfile {
	return &WeaponProfile{
		tpl: tpl,
	}
}

func (w *WeaponProfile) Abilities() []Ability {
	return w.tpl.Abilities
}

type WeaponProfileTemplate struct {
	Attacks          value.Interface
	Skill            int64
	Strength         value.Interface
	ArmorPenetration int64
	Damage           int64
	Abilities        []Ability
}
