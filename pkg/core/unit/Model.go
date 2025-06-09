package unit

type Model struct {
	tpl    *ModelTemplate
	points float64
}

func NewModel(unit *Template, model *ModelTemplate) *Model {
	return &Model{
		tpl:    model,
		points: float64(unit.PointsCost) / float64(unit.CoreModelCount()),
	}
}

func (m Model) Save() int64 {
	return m.tpl.Save
}
