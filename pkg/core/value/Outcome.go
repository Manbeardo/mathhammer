package value

import (
	"fmt"
	"math/big"

	"github.com/Manbeardo/mathhammer/pkg/core/util"
)

type Outcome struct {
	V int64   // value
	P big.Rat // probability
}

func (o Outcome) Format(w fmt.State, v rune) {
	util.PrettyFormat(w, v, o)
}
