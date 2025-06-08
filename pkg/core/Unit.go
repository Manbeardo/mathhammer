package core

import (
	"slices"

	"github.com/Manbeardo/mathhammer/pkg/core/util"
)

type Unit struct {
	tpl             *UnitTemplate
	models          []*Model
	weaponTemplates []*WeaponTemplate
	leaders         []*Unit
	startingHealth  UnitHealth
}

func NewUnit(tpl *UnitTemplate) *Unit {
	u := &Unit{
		tpl: tpl,
	}

	foundWeapons := map[*WeaponTemplate]struct{}{}

	for _, ltpl := range tpl.Leaders {
		u.leaders = append(u.leaders, NewUnit(ltpl))
	}

	unitTemplates := []*UnitTemplate{tpl}
	unitTemplates = append(unitTemplates, tpl.Leaders...)

	for _, utpl := range unitTemplates {
		for _, e := range utpl.Models {
			mtpl, count := e.Key, e.Value
			for range count {
				u.models = append(u.models, NewModel(utpl, mtpl))
				u.startingHealth = append(u.startingHealth, mtpl.Wounds)
			}
			for _, e := range mtpl.Weapons {
				wtpl := e.Key
				if _, exists := foundWeapons[wtpl]; !exists {
					u.weaponTemplates = append(u.weaponTemplates, wtpl)
					foundWeapons[wtpl] = struct{}{}
				}
			}
		}
	}

	return u
}

func (u *Unit) Abilities() []Ability {
	return u.tpl.Abilities
}

func (u *Unit) Toughness(health UnitHealth) int64 {
	// a unit's toughness is equal to the highest toughness
	// among its bodyguard models
	t := int64(0)
	for i := range u.tpl.CoreModelCount() {
		if health[i] == 0 {
			continue
		}
		mtpl := u.models[i].tpl
		if mtpl.Toughness > t {
			t = mtpl.Toughness
		}
	}
	return t
}

func (u *Unit) StartingHealth() UnitHealth {
	return slices.Clone(u.startingHealth)
}

func (u *Unit) Model(i int) *Model {
	return u.models[i]
}

func (u *Unit) Models() []*Model {
	return slices.Clone(u.models)
}

func (u *Unit) SurvivingModels(health UnitHealth) []*ModelTemplate {
	out := []*ModelTemplate{}
	for i, m := range u.models {
		if health[i] > 0 {
			out = append(out, m.tpl)
		}
	}
	return out
}

func (u *Unit) PointsLost(health UnitHealth) float64 {
	sum := 0.0
	for i, m := range u.models {
		remainingRatio := 1 - (float64(health[i]) / float64(m.tpl.Wounds))
		sum += remainingRatio * m.points
	}
	return sum
}

func (u *Unit) WeaponTemplates() []*WeaponTemplate {
	return slices.Clone(u.weaponTemplates)
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
