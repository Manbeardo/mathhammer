package core

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
	Attacks          Value
	Skill            int
	Strength         Value
	ArmorPenetration int
	Damage           Value
	Abilities        []Ability
}
