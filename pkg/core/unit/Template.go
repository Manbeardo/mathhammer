package unit

import (
	"github.com/Manbeardo/mathhammer/pkg/core/util"
	"github.com/Manbeardo/mathhammer/pkg/core/value"
)

type Template struct {
	Name       string
	Models     []util.Entry[*ModelTemplate, int]
	PointsCost int
	Leaders    []*Template
}

func (u *Template) CoreModelCount() int {
	total := 0
	for _, entry := range u.Models {
		total += entry.V
	}
	return total
}

func (u *Template) ModelCount() int {
	total := u.CoreModelCount()
	for _, l := range u.Leaders {
		total += l.ModelCount()
	}
	return total
}

type ModelTemplate struct {
	Name       string
	Toughness  int64
	Save       int64
	Wounds     int64
	Leadership int64
	Weapons    []util.Entry[*WeaponTemplate, int]
}

type WeaponTemplate struct {
	Name     string
	Profiles []*WeaponProfileTemplate
}

func (wtpl *WeaponTemplate) AvailableCount(unit *Unit, health Health) int64 {
	sum := 0
	for i, m := range unit.models {
		if health[i] == 0 {
			continue
		}
		for _, e := range m.tpl.Weapons {
			tpl, count := e.K, e.V
			if tpl == wtpl {
				sum += count
			}
		}
	}
	return int64(sum)
}

type WeaponProfileTemplate struct {
	Name             string
	RangeInches      int64
	Attacks          value.Interface
	Skill            int64
	Strength         value.Interface
	ArmorPenetration int64
	Damage           int64
}
