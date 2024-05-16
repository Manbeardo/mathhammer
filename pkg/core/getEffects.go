package core

type abilityAndBearer struct {
	ability Ability
	bearer  AbilityBearer
}

func getEffects(targets AttackTargets, trigger AbilityTrigger) []Effect {
	idToAbilityAndBearer := map[string]abilityAndBearer{}
	for _, bearer := range targets.AllBearers() {
		for _, ability := range bearer.Abilities() {
			if ability.Trigger() != trigger ||
				!ability.ShouldApply(AbilityContext{
					Bearer:        bearer,
					AttackTargets: targets,
				}) {
				continue
			}
			conflict, hasConflict := idToAbilityAndBearer[ability.ID()]
			if hasConflict && ability.CompareTo(conflict.ability) < 1 {
				continue
			}
			idToAbilityAndBearer[ability.ID()] = abilityAndBearer{
				ability: ability,
				bearer:  bearer,
			}
		}
	}
	effects := []Effect{}
	for _, aAndB := range idToAbilityAndBearer {
		effects = append(effects, aAndB.ability.Effects(AbilityContext{
			Bearer:        aAndB.bearer,
			AttackTargets: targets,
		})...)
	}
	return effects
}
