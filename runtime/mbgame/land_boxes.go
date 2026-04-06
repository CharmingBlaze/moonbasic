package mbgame

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// landBoxes returns the best (maximum) snap Y among count axis-aligned boxes given as
// parallel float arrays (same layout as BOXTOPLAND per box). 0 means no landing on any box.
func landBoxes(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 12 {
		return value.Nil, fmt.Errorf("LANDBOXES expects 12 arguments (px#, py#, pz#, pvy#, pr#, plx#, ply#, plz#, plw#, plh#, pld#, count)")
	}
	h := rt.Heap
	if h == nil {
		return value.Nil, fmt.Errorf("LANDBOXES: heap not available")
	}
	px, ok1 := args[0].ToFloat()
	py, ok2 := args[1].ToFloat()
	pz, ok3 := args[2].ToFloat()
	pvy, ok4 := args[3].ToFloat()
	pr, ok5 := args[4].ToFloat()
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
		return value.Nil, fmt.Errorf("LANDBOXES: first five arguments must be numeric")
	}
	hx, err := arrayHandleArg(args[5], "plx#")
	if err != nil {
		return value.Nil, err
	}
	hy, err := arrayHandleArg(args[6], "ply#")
	if err != nil {
		return value.Nil, err
	}
	hz, err := arrayHandleArg(args[7], "plz#")
	if err != nil {
		return value.Nil, err
	}
	hw, err := arrayHandleArg(args[8], "plw#")
	if err != nil {
		return value.Nil, err
	}
	hh, err := arrayHandleArg(args[9], "plh#")
	if err != nil {
		return value.Nil, err
	}
	hd, err := arrayHandleArg(args[10], "pld#")
	if err != nil {
		return value.Nil, err
	}
	var n int
	if ni, ok := args[11].ToInt(); ok {
		n = int(ni)
	} else if f, okf := args[11].ToFloat(); okf {
		n = int(f)
	} else {
		return value.Nil, fmt.Errorf("LANDBOXES: count must be numeric")
	}
	if n < 0 {
		n = 0
	}
	best := 0.0
	for i := 0; i < n; i++ {
		bx, okx := h.ArrayGetFloat(hx, int64(i))
		by, oky := h.ArrayGetFloat(hy, int64(i))
		bz, okz := h.ArrayGetFloat(hz, int64(i))
		bw, okw := h.ArrayGetFloat(hw, int64(i))
		bh, okh := h.ArrayGetFloat(hh, int64(i))
		bd, okd := h.ArrayGetFloat(hd, int64(i))
		if !okx || !oky || !okz || !okw || !okh || !okd {
			return value.Nil, fmt.Errorf("LANDBOXES: could not read platform index %d", i)
		}
		snap := BoxTopLandSnap(px, py, pz, pvy, pr, bx, by, bz, bw, bh, bd)
		if snap > best {
			best = snap
		}
	}
	return value.FromFloat(best), nil
}

func arrayHandleArg(v value.Value, name string) (heap.Handle, error) {
	if v.Kind != value.KindHandle {
		return 0, fmt.Errorf("LANDBOXES: %s must be an array handle", name)
	}
	return heap.Handle(v.IVal), nil
}
