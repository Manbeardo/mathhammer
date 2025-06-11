package unit

type BattleID struct {
	battle *Battle
}

func (id BattleID) getBattle() *Battle {
	return id.battle
}

func (id BattleID) getInstance() *Battle {
	return id.battle
}

func (id BattleID) StringKey() string {
	return ""
}

type Battle struct {
	id                 BattleID
	units              map[UnitID]*Unit
	models             map[ModelID]*Model
	weapons            map[WeaponID]*Weapon
	weaponProfiles     map[WeaponProfileID]*WeaponProfile
	weaponKinds        map[WeaponKindID]*WeaponKind
	weaponKindIDs      map[*WeaponTemplate]WeaponKindID
	weaponProfileKinds map[WeaponProfileKindID]*WeaponProfileKind
}

func NewBattle() *Battle {
	b := &Battle{
		units:              map[UnitID]*Unit{},
		models:             map[ModelID]*Model{},
		weapons:            map[WeaponID]*Weapon{},
		weaponProfiles:     map[WeaponProfileID]*WeaponProfile{},
		weaponKinds:        map[WeaponKindID]*WeaponKind{},
		weaponKindIDs:      map[*WeaponTemplate]WeaponKindID{},
		weaponProfileKinds: map[WeaponProfileKindID]*WeaponProfileKind{},
	}
	b.id = BattleID{
		battle: b,
	}
	return b
}

func (b *Battle) StringKey() string {
	return b.id.StringKey()
}
