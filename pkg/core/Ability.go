package core

type AbilityTrigger string

const (
	TriggerSelectTargetUnit AbilityTrigger = "select target unit"
)

type Ability interface {
	ID() string
	Trigger() AbilityTrigger
	ShouldApply(AbilityContext) bool
	Effects(AbilityContext) []Effect
	CompareTo(Ability) int
}

type AbilityBearer interface {
	Abilities() []Ability
}

type AbilityContext struct {
	// AttackTargets
	Bearer AbilityBearer
}

// func (ctx AbilityContext) IsAttacker() bool {
// 	return ctx.Bearer == ctx.AttackerUnit ||
// 		ctx.Bearer == ctx.AttackerModel ||
// 		ctx.Bearer == ctx.AttackerWeaponProfile
// }

// func (ctx AbilityContext) IsDefender() bool {
// 	return ctx.Bearer == ctx.DefenderUnit ||
// 		ctx.Bearer == ctx.DefenderModel
// }
