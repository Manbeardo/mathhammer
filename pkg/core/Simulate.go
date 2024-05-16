package core

func SimulateShooting(attacker *Unit, defender *Unit) {
	NewAttackSet(attacker, defender, RangedAttack).Run()
}

func SimulateFighting(attacker *Unit, defender *Unit) {
	NewAttackSet(attacker, defender, MeleeAttack).Run()
	NewAttackSet(defender, attacker, MeleeAttack).Run()
}
