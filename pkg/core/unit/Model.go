package unit

import "slices"

type ModelDatasheet struct {
	Name       string
	Toughness  int64
	Save       int64
	Wounds     int64
	Leadership int64
}

type ModelTemplate struct {
	baseTemplate[ModelDatasheet]
	Weapons []*WeaponTemplate
}

func NewModelTemplate(stats ModelDatasheet, weapons ...*WeaponTemplate) *ModelTemplate {
	return &ModelTemplate{
		baseTemplate: createBaseTemplate(stats),
		Weapons:      slices.Clone(weapons),
	}
}

type ModelID struct {
	childID[*Unit, UnitID]
}

func (id ModelID) getInstance() *Model {
	return id.getBattle().models[id]
}

type Model struct {
	instance[*Model, ModelID, ModelDatasheet, *ModelTemplate]
	weapons []*Weapon
}

func (b *Battle) newModel(id ModelID, tpl *ModelTemplate) *Model {
	m := &Model{
		instance: createInstance(b, id, tpl),
	}
	b.models[m.id] = m

	for i, wtpl := range tpl.Weapons {
		wid := WeaponID{
			childID: createChildID(m.id, i),
		}
		m.weapons = append(m.weapons, b.newWeapon(wid, wtpl))
	}

	return m
}

func (m *Model) Weapons() []*Weapon {
	return slices.Clone(m.weapons)
}

func (m *Model) WeaponProfiles() [][]*WeaponProfile {
	out := [][]*WeaponProfile{}
	for _, w := range m.weapons {
		out = append(out, w.Profiles())
	}
	return out
}

func (m *Model) IsAlive(health Health) bool {
	return health[m.id.index] > 0
}

func (m *Model) Points() float64 {
	unit := m.id.parent.getInstance()
	return float64(unit.Datasheet().PointsCost) / float64(len(unit.models))
}

func (m *Model) PointsLost(health Health) float64 {
	remainingRatio := 1.0 - (float64(health[m.id.index]) / float64(m.Datasheet().Wounds))
	return remainingRatio * m.Points()
}
