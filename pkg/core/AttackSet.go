package core

type AttackSet struct {
	profileTplToAttacks map[*WeaponProfileTemplate][]*Attack
}

func NewAttackSet(attacker *Unit, defender *Unit, kind AttackKind) *AttackSet {
	profileTpls := attacker.ProfileTemplatesForAttack(kind)
	attack := NewAttack(AttackTargets{
		AttackerUnit: attacker,
		DefenderUnit: defender,
	})

	tplToAttacks := map[*WeaponProfileTemplate][]*Attack{}
	// select target for every weapon profile
	for profileTpl := range profileTpls {
		for _, model := range attacker.SurvivingModels() {
			attack := attack.Clone()
			attack.AttackerModel = model
			// we start with an arbitrary defender model in order to expose the toughness value
			attack.DefenderModel = defender.models[0]
			profiles := model.MatchingWeaponProfiles(profileTpl)
			for _, profile := range profiles {
				attack := attack.Clone()
				attack.AttackerWeaponProfile = profile
				attack.ApplyTriggerEffects(TriggerSelectTargetUnit)
				tplToAttacks[profileTpl] = append(tplToAttacks[profileTpl], attack)
			}
		}
	}

	return &AttackSet{
		profileTplToAttacks: tplToAttacks,
	}
}

func (as *AttackSet) Run() {
	// roll attacks per profile
	// for _, attacks := range as.profileTplToAttacks {
	// 	for _, attack := range attacks {
	// 		attackCount := attack.EvalAttacks()
	// 		hitResults := attack.RollHits(attackCount)
	// 		woundResults := attack.RollWounds(hitResults.Successes())

	// 		for range woundResults.Successes() {
	// 			if attack.DefenderUnit.IsDead() {
	// 				return
	// 			}
	// 			attack := attack.Clone()
	// 			attack.AllocateWound()
	// 			saveResult := attack.RollSave()
	// 			if saveResult.Successes() > 0 {
	// 				continue
	// 			}
	// 			dmg := attack.EvalDamage()
	// 			attack.DefenderModel.woundsTaken += dmg
	// 		}
	// 	}
	// }
}
