package unit

import (
	"slices"
)

type WeaponDatasheet struct {
	Name string
}

type WeaponTemplate struct {
	baseTemplate[WeaponDatasheet]
	Profiles []*WeaponProfileTemplate
}

func NewWeaponTemplate(stats WeaponDatasheet, profiles ...*WeaponProfileTemplate) *WeaponTemplate {
	return &WeaponTemplate{
		baseTemplate: createBaseTemplate(stats),
		Profiles:     slices.Clone(profiles),
	}
}

type WeaponID struct {
	childID[*Model, ModelID]
}

func (id WeaponID) getInstance() *Weapon {
	return id.getBattle().weapons[id]
}

type Weapon struct {
	instance[*Weapon, WeaponID, WeaponDatasheet, *WeaponTemplate]
	profiles []*WeaponProfile
	kind     *WeaponKind
}

func (b *Battle) newWeapon(id WeaponID, tpl *WeaponTemplate) *Weapon {
	w := &Weapon{
		instance: createInstance(b, id, tpl),
	}
	b.weapons[w.id] = w

	if typeID, ok := b.weaponKindIDs[tpl]; ok {
		w.kind = b.weaponKinds[typeID]
	} else {
		w.kind = b.newWeaponKind(tpl)
	}

	for i, ptpl := range tpl.Profiles {
		pid := WeaponProfileID{
			childID: createChildID(w.id, i),
		}
		w.profiles = append(w.profiles, b.newWeaponProfile(pid, ptpl))
	}

	return w
}

func (w *Weapon) Kind() *WeaponKind {
	return w.kind
}

func (w *Weapon) Profiles() []*WeaponProfile {
	return slices.Clone(w.profiles)
}

type WeaponKindID struct {
	childID[*Battle, BattleID]
}

func (id WeaponKindID) getInstance() *WeaponKind {
	return id.getBattle().weaponKinds[id]
}

type WeaponKind struct {
	instance[*WeaponKind, WeaponKindID, WeaponDatasheet, *WeaponTemplate]
	profiles []*WeaponProfileKind
}

func (b *Battle) newWeaponKind(tpl *WeaponTemplate) *WeaponKind {
	id := WeaponKindID{
		childID: createChildID(b.id, len(b.weaponKinds)),
	}
	wk := &WeaponKind{
		instance: createInstance(b, id, tpl),
	}
	b.weaponKinds[wk.id] = wk
	b.weaponKindIDs[tpl] = wk.id

	for i, ptpl := range tpl.Profiles {
		pid := WeaponProfileKindID{
			childID: createChildID(wk.id, i),
		}
		wk.profiles = append(wk.profiles, b.newWeaponProfileKind(pid, ptpl))
	}

	return wk
}

func (wk *WeaponKind) Profiles() []*WeaponProfileKind {
	return slices.Clone(wk.profiles)
}
