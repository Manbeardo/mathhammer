package core

type Unit struct {
	tpl     *UnitTemplate
	models  []*Model
	leaders []*Unit
}

func NewUnit(tpl *UnitTemplate) *Unit {
	u := &Unit{
		tpl:     tpl,
		leaders: make([]*Unit, len(tpl.Leaders)),
	}

	for mtpl, count := range tpl.Models {
		for range count {
			u.models = append(u.models, NewModel(mtpl))
		}
	}

	for i, ltpl := range tpl.Leaders {
		u.leaders[i] = NewUnit(ltpl)
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
	for _, m := range u.models {
		if !m.IsDead() {
			out = append(out, m)
		}
	}
	for _, l := range u.leaders {
		out = append(out, l.SurvivingModels()...)
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

func (u *Unit) PointsLost() float64 {
	ppm := float64(u.tpl.PointsCost) / float64(len(u.models))
	lost := 0.0
	for _, m := range u.models {
		if m.IsDead() {
			lost += ppm
		} else {
			lost += ppm * (float64(m.woundsTaken) / float64(m.tpl.Wounds))
		}
	}
	for _, l := range u.leaders {
		lost += l.PointsLost()
	}
	return lost
}

func (u *Unit) ProfileTemplatesForAttack(kind AttackKind) map[*WeaponProfileTemplate]struct{} {
	out := map[*WeaponProfileTemplate]struct{}{}
	for _, model := range u.SurvivingModels() {
		for _, weapon := range model.weapons {
			if weapon.tpl.Kind != RangedAttack {
				continue
			}
			for _, profile := range weapon.tpl.Profiles {
				out[profile] = struct{}{}
			}
		}
	}
	return out
}

type UnitTemplate struct {
	Name       string
	Models     map[*ModelTemplate]int
	PointsCost int
	Leaders    []*UnitTemplate
	Abilities  []Ability
}

func (u *UnitTemplate) CoreModelCount() int {
	total := 0
	for _, count := range u.Models {
		total += count
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
