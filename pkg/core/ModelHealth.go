package core

import (
	"encoding/json"
	"fmt"
)

type ModelHealthStr string

func EncodeModelHealth(modelHealth []int64) ModelHealthStr {
	b, err := json.Marshal(modelHealth)
	if err != nil {
		panic(fmt.Errorf("marshaling ModelHealth: %w", err))
	}
	return ModelHealthStr(b)
}

func (mh ModelHealthStr) ToSlice() []int64 {
	out := []int64{}
	err := json.Unmarshal([]byte(mh), &out)
	if err != nil {
		panic(fmt.Errorf("unmarshaling ModelHealth: %w", err))
	}
	return out
}
