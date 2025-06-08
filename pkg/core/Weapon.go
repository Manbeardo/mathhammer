package core

type WeaponTemplate struct {
	Name     string
	Profiles []*WeaponProfileTemplate
}

func (wtpl *WeaponTemplate) AvailableCount(unit *Unit, health UnitHealth) int64 {
	sum := 0
	for i, m := range unit.models {
		if health[i] == 0 {
			continue
		}
		for _, e := range m.tpl.Weapons {
			tpl, count := e.Key, e.Value
			if tpl == wtpl {
				sum += count
			}
		}
	}
	return int64(sum)
}
