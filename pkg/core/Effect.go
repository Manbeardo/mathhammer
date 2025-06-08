package core

type Effect interface {
	ApplyEffect()
}

type EffectApplyModifier struct {
	// Mod  Modifier
	// Kind ModifierKind
}
