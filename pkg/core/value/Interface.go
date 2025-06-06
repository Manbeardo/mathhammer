package value

import "github.com/Manbeardo/mathhammer/pkg/core/probability"

type Interface interface {
	Distribution() probability.Distribution[int64]
}
