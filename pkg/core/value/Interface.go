package value

import "github.com/Manbeardo/mathhammer/pkg/core/prob"

type Interface interface {
	Distribution() prob.Dist[int64]
}
