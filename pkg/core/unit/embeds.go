package unit

import (
	"fmt"
)

type id[Instance any] interface {
	getBattle() *Battle
	getInstance() Instance
	StringKey() string
}

type childID[Parent any, ParentID id[Parent]] struct {
	parent ParentID
	index  int
}

func (id childID[P, PID]) getBattle() *Battle {
	return id.parent.getBattle()
}

func (id childID[P, PID]) StringKey() string {
	return fmt.Sprintf("%s %d", id.parent.StringKey(), id.index)
}

func createChildID[P any, PID id[P]](parent PID, idx int) childID[P, PID] {
	return childID[P, PID]{
		parent: parent,
		index:  idx,
	}
}

type template[Datasheet any] interface {
	Datasheet() Datasheet
}

type baseTemplate[Datasheet any] struct {
	datasheet Datasheet
}

func createBaseTemplate[D any](Datasheet D) baseTemplate[D] {
	return baseTemplate[D]{
		datasheet: Datasheet,
	}
}

func (tpl baseTemplate[D]) Datasheet() D {
	return tpl.datasheet
}

type instance[T any, ID id[T], Stats any, Tpl template[Stats]] struct {
	id     ID
	battle *Battle
	tpl    Tpl
}

func createInstance[T any, ID id[T], Datasheet any, Tpl template[Datasheet]](
	battle *Battle,
	id ID,
	tpl Tpl,
) instance[T, ID, Datasheet, Tpl] {
	return instance[T, ID, Datasheet, Tpl]{
		id:     id,
		battle: battle,
		tpl:    tpl,
	}
}

func (i *instance[T, ID, Datasheet, Tpl]) Datasheet() Datasheet {
	return i.tpl.Datasheet()
}

func (i *instance[T, ID, Details, Tpl]) StringKey() string {
	return i.id.StringKey()
}
