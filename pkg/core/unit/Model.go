package unit

import "github.com/Manbeardo/mathhammer/pkg/core/util"

type Model struct {
	tpl    *ModelTemplate
	points float64
}

func NewModel(unit *UnitTemplate, model *ModelTemplate) *Model {
	return &Model{
		tpl:    model,
		points: float64(unit.PointsCost) / float64(unit.CoreModelCount()),
	}
}

func (m Model) Save() int64 {
	return m.tpl.Save
}

type ModelTemplate struct {
	Name       string
	Toughness  int64
	Save       int64
	Wounds     int64
	Leadership int64
	Weapons    []util.Entry[*WeaponTemplate, int]
}
