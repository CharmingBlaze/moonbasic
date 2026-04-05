package mbgame

import (
	"fmt"

	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) rndRange(minV, maxV int64) int64 {
	if maxV < minV {
		minV, maxV = maxV, minV
	}
	span := maxV - minV + 1
	if span <= 0 {
		return minV
	}
	return minV + m.rng.Int63n(span)
}

func (m *Module) weightedRndIndex(h *heap.Store, arrH heap.Handle) (int, error) {
	n := h.ArrayFlatLen(arrH)
	if n < 1 {
		return -1, fmt.Errorf("weights array empty")
	}
	sum := 0.0
	for i := 0; i < n; i++ {
		w, ok := h.ArrayGetFloat(arrH, int64(i))
		if !ok {
			return -1, fmt.Errorf("invalid weight at %d", i)
		}
		if w > 0 {
			sum += w
		}
	}
	if sum <= 0 {
		return 0, nil
	}
	r := m.rng.Float64() * sum
	acc := 0.0
	for i := 0; i < n; i++ {
		w, _ := h.ArrayGetFloat(arrH, int64(i))
		if w <= 0 {
			continue
		}
		acc += w
		if r <= acc {
			return int(i), nil
		}
	}
	return int(n - 1), nil
}

func (m *Module) shuffleArray(st *heap.Store, arrH heap.Handle) error {
	a, err := heap.Cast[*heap.Array](st, arrH)
	if err != nil {
		return err
	}
	n := st.ArrayFlatLen(arrH)
	if n < 2 {
		return nil
	}
	for i := n - 1; i > 0; i-- {
		j := int(m.rng.Int63n(int64(i + 1)))
		vi, e1 := a.Get([]int64{int64(i)})
		vj, e2 := a.Get([]int64{int64(j)})
		if e1 != nil || e2 != nil {
			return fmt.Errorf("SHUFFLE: numeric array required")
		}
		_ = a.Set([]int64{int64(i)}, vj)
		_ = a.Set([]int64{int64(j)}, vi)
	}
	return nil
}

func (m *Module) randomElement(st *heap.Store, arrH heap.Handle) (value.Value, error) {
	n := st.ArrayFlatLen(arrH)
	if n < 1 {
		return value.Nil, fmt.Errorf("array empty")
	}
	ix := m.rng.Int63n(int64(n))
	v, ok := st.ArrayGetFloat(arrH, ix)
	if ok {
		return value.FromFloat(v), nil
	}
	return value.Nil, fmt.Errorf("RANDOMELEMENT: numeric array required")
}
