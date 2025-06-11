package unit

import (
	"slices"

	"github.com/Manbeardo/mathhammer/pkg/core/util"
)

type UnitDatasheet struct {
	Name       string
	PointsCost int
}

type UnitTemplate struct {
	baseTemplate[UnitDatasheet]
	Models []*ModelTemplate
}

func NewUnitTemplate(stats UnitDatasheet, models ...*ModelTemplate) *UnitTemplate {
	return &UnitTemplate{
		baseTemplate: createBaseTemplate(stats),
		Models:       slices.Clone(models),
	}
}

type UnitID struct {
	childID[*Battle, BattleID]
}

func (id UnitID) getInstance() *Unit {
	return id.getBattle().units[id]
}

type Unit struct {
	instance[*Unit, UnitID, UnitDatasheet, *UnitTemplate]
	models []*Model
}

func (b *Battle) NewUnit(tpl *UnitTemplate) *Unit {
	id := UnitID{
		childID: createChildID(b.id, len(b.units)),
	}
	u := &Unit{
		instance: createInstance(b, id, tpl),
	}
	b.units[u.id] = u

	for i, mtpl := range tpl.Models {
		mid := ModelID{
			childID: createChildID(u.id, i),
		}
		u.models = append(u.models, b.newModel(mid, mtpl))
	}

	return u
}

func (u *Unit) Toughness(health Health) int64 {
	// a unit's toughness is equal to the highest toughness
	// among its bodyguard models
	ts := []int64{}
	for i := range u.models {
		if health[i] == 0 {
			continue
		}
		ts = append(ts, u.models[i].Datasheet().Toughness)
	}
	if len(ts) == 0 {
		return 0
	}
	return slices.Max(ts)
}

func (u *Unit) StartingHealth() Health {
	health := Health{}
	for _, model := range u.models {
		health = append(health, model.Datasheet().Wounds)
	}
	return health
}

func (u *Unit) Models() []*Model {
	return slices.Clone(u.models)
}

func (u *Unit) WeaponProfiles() [][]*WeaponProfile {
	out := [][]*WeaponProfile{}
	for _, m := range u.models {
		out = append(out, m.WeaponProfiles()...)
	}
	return out
}

func (u *Unit) SurvivingModels(health Health) []*Model {
	out := []*Model{}
	for _, m := range u.models {
		if m.IsAlive(health) {
			out = append(out, m)
		}
	}
	return out
}

func (u *Unit) PointsLost(health Health) float64 {
	sum := 0.0
	for _, m := range u.models {
		sum += m.PointsLost(health)
	}
	return sum
}

func (u *Unit) SurvivingWeapons(health Health) *util.OrderedMap[*WeaponKind, []*Weapon] {
	out := util.NewOrderedMap[*WeaponKind, []*Weapon]()
	for _, m := range u.SurvivingModels(health) {
		for _, w := range m.weapons {
			ws, _ := out.Get(w.kind)
			out.Put(w.kind, append(ws, w))
		}
	}
	return out
}
