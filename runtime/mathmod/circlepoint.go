package mathmod

import (
	"fmt"
	"math"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) registerCirclePoint(r runtime.Registrar) {
	r.Register("MATH.CIRCLEPOINT", "math", m.mathCirclePoint)
	r.Register("CIRCLEPOINT", "math", m.mathCirclePoint)
}

// mathCirclePoint returns (x, z) on a circle: center (cx, cz), radius, index i (1..count), count points.
func (m *Module) mathCirclePoint(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.h == nil {
		return value.Nil, fmt.Errorf("MATH.CIRCLEPOINT: heap not bound")
	}
	if len(args) != 5 {
		return value.Nil, errNArgs(5, len(args))
	}
	cx, _ := args[0].ToFloat()
	cz, _ := args[1].ToFloat()
	rad, _ := args[2].ToFloat()
	idx, ok := args[3].ToInt()
	if !ok {
		f, _ := args[3].ToFloat()
		idx = int64(f)
	}
	cnt, ok := args[4].ToInt()
	if !ok {
		f, _ := args[4].ToFloat()
		cnt = int64(f)
	}
	if cnt < 1 {
		return value.Nil, fmt.Errorf("MATH.CIRCLEPOINT: count must be >= 1")
	}
	if idx < 1 {
		return value.Nil, fmt.Errorf("MATH.CIRCLEPOINT: index i must be >= 1")
	}
	ang := float64(idx-1) * 2 * math.Pi / float64(cnt)
	ex := cx + math.Cos(ang)*rad
	ez := cz + math.Sin(ang)*rad
	arr, err := heap.NewArrayOfKind([]int64{2}, heap.ArrayKindFloat, 0)
	if err != nil {
		return value.Nil, err
	}
	arr.Floats[0] = ex
	arr.Floats[1] = ez
	h, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(h), nil
}
