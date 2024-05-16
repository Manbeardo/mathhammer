package core

type Weapon struct {
	tpl      *WeaponTemplate
	profiles []*WeaponProfile
}

func NewWeapon(tpl *WeaponTemplate) *Weapon {
	w := &Weapon{
		tpl:      tpl,
		profiles: make([]*WeaponProfile, len(tpl.Profiles)),
	}

	for i, ptpl := range tpl.Profiles {
		w.profiles[i] = NewWeaponProfile(ptpl)
	}

	return w
}

func (w *Weapon) WasActivated() bool {
	for _, p := range w.profiles {
		if p.wasActivated {
			return true
		}
	}
	return false
}

func (w *Weapon) MatchingProfile(tpl *WeaponProfileTemplate) *WeaponProfile {
	for _, p := range w.profiles {
		if p.tpl == tpl {
			return p
		}
	}
	return nil
}

type AttackKind string

const (
	RangedAttack AttackKind = "ranged"
	MeleeAttack  AttackKind = "melee"
)

type WeaponTemplate struct {
	Name     string
	Kind     AttackKind
	Profiles []*WeaponProfileTemplate
}
