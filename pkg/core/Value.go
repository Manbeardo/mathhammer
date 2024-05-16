package core

import "math/rand"

type Value interface {
	Eval() int
}

type RollValue struct {
	N int
}

func (v RollValue) Eval() int {
	return rand.Intn(v.N) + 1
}

type IntValue struct {
	N int
}

func (v IntValue) Eval() int {
	return v.N
}

type SumValue []Value

func (v SumValue) Eval() int {
	sum := 0
	for _, inner := range v {
		sum += inner.Eval()
	}
	return sum
}
