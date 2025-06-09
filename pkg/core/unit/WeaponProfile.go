package unit

type WeaponProfile struct {
	tpl *WeaponProfileTemplate
}

func NewWeaponProfile(tpl *WeaponProfileTemplate) *WeaponProfile {
	return &WeaponProfile{
		tpl: tpl,
	}
}
