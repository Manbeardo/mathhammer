package core

type Effect interface {
	ApplyEffect(attack *Attack)
}

type EffectApplyModifier struct {
	// Mod  Modifier
	// Kind ModifierKind
}
