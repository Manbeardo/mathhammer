package unit

import (
	"slices"
)

type Unit struct {
	tpl             *Template
	models          []*Model
	weaponTemplates []*WeaponTemplate
	leaders         []*Unit
	startingHealth  Health
}

func NewUnit(tpl *Template) *Unit {
	u := &Unit{
		tpl: tpl,
	}

	foundWeapons := map[*WeaponTemplate]struct{}{}

	for _, ltpl := range tpl.Leaders {
		u.leaders = append(u.leaders, NewUnit(ltpl))
	}

	unitTemplates := []*Template{tpl}
	unitTemplates = append(unitTemplates, tpl.Leaders...)

	for _, utpl := range unitTemplates {
		for _, e := range utpl.Models {
			mtpl, count := e.K, e.V
			for range count {
				u.models = append(u.models, NewModel(utpl, mtpl))
				u.startingHealth = append(u.startingHealth, mtpl.Wounds)
			}
			for _, e := range mtpl.Weapons {
				wtpl := e.K
				if _, exists := foundWeapons[wtpl]; !exists {
					u.weaponTemplates = append(u.weaponTemplates, wtpl)
					foundWeapons[wtpl] = struct{}{}
				}
			}
		}
	}

	return u
}

func (u *Unit) Toughness(health Health) int64 {
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

func (u *Unit) StartingHealth() Health {
	return slices.Clone(u.startingHealth)
}

func (u *Unit) Model(i int) *Model {
	return u.models[i]
}

func (u *Unit) Models() []*Model {
	return slices.Clone(u.models)
}

func (u *Unit) SurvivingModels(health Health) []*ModelTemplate {
	out := []*ModelTemplate{}
	for i, m := range u.models {
		if health[i] > 0 {
			out = append(out, m.tpl)
		}
	}
	return out
}

func (u *Unit) PointsLost(health Health) float64 {
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
