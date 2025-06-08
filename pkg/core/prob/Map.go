package prob

import "math/big"

// Map maps a probability distribution from outcome type T to U
func Map[T any, U any](
	dist Dist[T],
	mapper func(T) U,
) (Dist[U], error) {
	out, err := empty[U]()
	if err != nil {
		return out, err
	}
	for tk, p := range dist.pmap {
		tv := dist.vmap[tk]
		uv := mapper(tv)
		uk := out.key(uv)
		out.vmap[uk] = uv
		outP, ok := out.pmap[uk]
		if !ok {
			outP = big.NewRat(0, 1)
			out.pmap[uk] = outP
		}
		outP.Add(outP, p)
	}
	return out.finalize()
}
