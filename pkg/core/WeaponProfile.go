package core

import "github.com/Manbeardo/mathhammer/pkg/core/value"

type WeaponProfile struct {
	tpl *WeaponProfileTemplate
}

func NewWeaponProfile(tpl *WeaponProfileTemplate) *WeaponProfile {
	return &WeaponProfile{
		tpl: tpl,
	}
}

type WeaponProfileTemplate struct {
	Name             string
	RangeInches      int64
	Attacks          value.Interface
	Skill            int64
	Strength         value.Interface
	ArmorPenetration int64
	Damage           int64
}
