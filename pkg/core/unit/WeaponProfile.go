package unit

import (
	"github.com/Manbeardo/mathhammer/pkg/core/value"
)

type WeaponProfileDatasheet struct {
	Name             string
	RangeInches      int64
	Attacks          value.Interface
	Skill            int64
	Strength         value.Interface
	ArmorPenetration int64
	Damage           int64
}

type WeaponProfileTemplate struct {
	baseTemplate[WeaponProfileDatasheet]
}

func NewWeaponProfileTemplate(stats WeaponProfileDatasheet) *WeaponProfileTemplate {
	return &WeaponProfileTemplate{
		baseTemplate: createBaseTemplate(stats),
	}
}

type WeaponProfileID struct {
	childID[*Weapon, WeaponID]
}

func (id WeaponProfileID) getInstance() *WeaponProfile {
	return id.getBattle().weaponProfiles[id]
}

type WeaponProfile struct {
	instance[*WeaponProfile, WeaponProfileID, WeaponProfileDatasheet, *WeaponProfileTemplate]
}

func (b *Battle) newWeaponProfile(id WeaponProfileID, tpl *WeaponProfileTemplate) *WeaponProfile {
	p := &WeaponProfile{
		instance: createInstance(b, id, tpl),
	}
	b.weaponProfiles[p.id] = p

	return p
}

func (wp *WeaponProfile) Kind() *WeaponProfileKind {
	weapon := wp.battle.weapons[wp.id.parent]
	return weapon.Kind().profiles[wp.id.index]
}

type WeaponProfileKindID struct {
	childID[*WeaponKind, WeaponKindID]
}

func (id WeaponProfileKindID) getInstance() *WeaponProfileKind {
	return id.getBattle().weaponProfileKinds[id]
}

type WeaponProfileKind struct {
	instance[*WeaponProfileKind, WeaponProfileKindID, WeaponProfileDatasheet, *WeaponProfileTemplate]
}

func (b *Battle) newWeaponProfileKind(id WeaponProfileKindID, tpl *WeaponProfileTemplate) *WeaponProfileKind {
	pk := &WeaponProfileKind{
		instance: createInstance(b, id, tpl),
	}
	b.weaponProfileKinds[pk.id] = pk

	return pk
}
