package core

type Event interface {
	eventMarker()
}

type baseEvent struct {
	Attacker *UnitTemplate
	Defender *UnitTemplate
}

func (_ baseEvent) eventMarker() {}

type EventSelectRangedWeapon struct {
	baseEvent
}

type EventCalculateAttacks struct {
	baseEvent
	Weapon *WeaponTemplate
}

type EventUnitDoneAttacking struct {
	baseEvent
	Weapon    *WeaponTemplate
	DidAttack bool
}
