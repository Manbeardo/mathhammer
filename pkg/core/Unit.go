package core

import (
	"slices"

	"github.com/Manbeardo/mathhammer/pkg/core/util"
)

type Unit struct {
	tpl     *UnitTemplate
	models  []*Model
	leaders []*Unit
	health  UnitHealth
}

func NewUnit(tpl *UnitTemplate) *Unit {
	u := &Unit{
		tpl:     tpl,
		leaders: make([]*Unit, len(tpl.Leaders)),
	}

	for _, e := range tpl.Models {
		mtpl, count := e.Key, e.Value
		for range count {
			u.models = append(u.models, NewModel(mtpl))
			u.health = append(u.health, mtpl.Wounds)
		}
	}

	for i, ltpl := range tpl.Leaders {
		u.leaders[i] = NewUnit(ltpl)
		u.health = append(u.health, u.leaders[i].health...)
	}

	return u
}

func (u *Unit) IsDead() bool {
	return len(u.SurvivingModels()) == 0
}

func (u *Unit) Abilities() []Ability {
	return u.tpl.Abilities
}

func (u *Unit) SurvivingModels() []*Model {
	out := []*Model{}
	for i, m := range u.Models() {
		if u.health[i] > 0 {
			out = append(out, m)
		}
	}
	return out
}

func (u *Unit) ModelCount() int {
	return u.tpl.ModelCount()
}

func (u *Unit) PointsCost() int {
	p := u.tpl.PointsCost
	for _, l := range u.leaders {
		p += l.PointsCost()
	}
	return p
}

func (u *Unit) Toughness() int64 {
	// a unit's toughness is equal to the highest toughness
	// among its bodyguard models
	t := int64(0)
	for i, m := range u.models {
		if u.health[i] == 0 {
			continue
		}
		if m.tpl.Toughness > t {
			t = m.tpl.Toughness
		}
	}
	return t
}

func (u *Unit) Model(idx int) *Model {
	if idx < len(u.models) {
		return u.models[idx]
	}
	idx -= len(u.models)
	for _, leader := range u.leaders {
		if idx < len(leader.models) {
			return leader.models[idx]
		}
		idx -= len(leader.models)
	}
	return nil
}

func (u *Unit) Models() []*Model {
	ms := slices.Clone(u.models)
	for _, l := range u.leaders {
		ms = append(ms, l.Models()...)
	}
	return ms
}

func (u *Unit) PointsLost() float64 {
	ppm := float64(u.tpl.PointsCost) / float64(len(u.models))
	lost := 0.0
	for i, m := range u.models {
		lost += ppm * (1.0 - (float64(u.health[i]) - float64(m.tpl.Wounds)))
	}
	for _, l := range u.leaders {
		lost += l.PointsLost()
	}
	return lost
}

func (u *Unit) Weapons() map[*WeaponTemplate]int64 {
	out := map[*WeaponTemplate]int64{}
	for i, m := range u.Models() {
		if u.health[i] == 0 {
			continue
		}
		for _, w := range m.weapons {
			out[w.tpl] += 1
		}
	}
	return out
}

type UnitTemplate struct {
	Name       string
	Models     []util.Entry[*ModelTemplate, int]
	PointsCost int
	Leaders    []*UnitTemplate
	Abilities  []Ability
}

func (u *UnitTemplate) CoreModelCount() int {
	total := 0
	for _, entry := range u.Models {
		total += entry.Value
	}
	return total
}

func (u *UnitTemplate) ModelCount() int {
	total := u.CoreModelCount()
	for _, l := range u.Leaders {
		total += l.ModelCount()
	}
	return total
}
