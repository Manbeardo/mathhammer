package core

import (
	"encoding/json"
	"fmt"
)

type ModelHealth []int64

func (mh ModelHealth) WoundsRemaining() int64 {
	sum := int64(0)
	for _, h := range mh {
		sum += h
	}
	return sum
}

func (mh ModelHealth) ModelsRemaining() int64 {
	sum := int64(0)
	for _, h := range mh {
		if h > 0 {
			sum += 1
		}
	}
	return sum
}

func (mh ModelHealth) ToKey() ModelHealthStr {
	b, err := json.Marshal(mh)
	if err != nil {
		panic(fmt.Errorf("marshaling ModelHealth: %w", err))
	}
	return ModelHealthStr(b)
}

func (mh ModelHealth) String() string {
	return string(mh.ToKey())
}

type ModelHealthStr string

func (mh ModelHealthStr) ToSlice() ModelHealth {
	out := []int64{}
	err := json.Unmarshal([]byte(mh), &out)
	if err != nil {
		panic(fmt.Errorf("unmarshaling ModelHealth: %w", err))
	}
	return out
}
