package mbrand

import (
	"fmt"
	"math/rand/v2"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

type rngObj struct {
	r *rand.Rand
}

func (o *rngObj) TypeName() string { return "Rng" }

func (o *rngObj) TypeTag() uint16 { return heap.TagRng }

func (o *rngObj) Free() { o.r = nil }

// Register implements runtime.Module.
func (m *Module) Register(r runtime.Registrar) {
	r.Register("RAND.MAKE", "rand", runtime.AdaptLegacy(m.randMake))
	r.Register("RAND.NEXT", "rand", runtime.AdaptLegacy(m.randNext))
	r.Register("RAND.NEXTF", "rand", runtime.AdaptLegacy(m.randNextF))
	r.Register("RAND.FREE", "rand", runtime.AdaptLegacy(m.randFree))
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}

func (m *Module) requireHeap() error {
	if m.h == nil {
		return runtime.Errorf("RAND.* builtins: heap not bound")
	}
	return nil
}

func (m *Module) getRng(args []value.Value, ix int, op string) (*rngObj, error) {
	if err := m.requireHeap(); err != nil {
		return nil, err
	}
	if ix >= len(args) || args[ix].Kind != value.KindHandle {
		return nil, fmt.Errorf("%s: argument %d must be rng handle", op, ix+1)
	}
	return heap.Cast[*rngObj](m.h, heap.Handle(args[ix].IVal))
}

func (m *Module) randMake(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, runtime.Errorf("RAND.MAKE expects 1 argument (seed)")
	}
	seed, ok := args[0].ToInt()
	if !ok {
		if f, okf := args[0].ToFloat(); okf {
			seed = int64(f)
		} else {
			return value.Nil, runtime.Errorf("RAND.MAKE: seed must be numeric")
		}
	}
	s := uint64(seed)
	// Second PCG stream word derived from seed so different seeds differ in both words.
	seq := s*0x9e3779b97f4a7c15 + 0x1234567890abcdef
	src := rand.NewPCG(s, seq)
	r := rand.New(src)
	id, err := m.h.Alloc(&rngObj{r: r})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) randNext(args []value.Value) (value.Value, error) {
	o, err := m.getRng(args, 0, "RAND.NEXT")
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, runtime.Errorf("RAND.NEXT expects 3 arguments (rng, min, max)")
	}
	min, ok1 := args[1].ToInt()
	max, ok2 := args[2].ToInt()
	if !ok1 || !ok2 {
		return value.Nil, runtime.Errorf("RAND.NEXT: min and max must be integers")
	}
	if max < min {
		min, max = max, min
	}
	span := max - min + 1
	if span <= 0 {
		return value.Nil, runtime.Errorf("RAND.NEXT: integer range overflow (min/max too far apart)")
	}
	n := o.r.Int64N(span) + min
	return value.FromInt(n), nil
}

func (m *Module) randNextF(args []value.Value) (value.Value, error) {
	o, err := m.getRng(args, 0, "RAND.NEXTF")
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, runtime.Errorf("RAND.NEXTF expects 1 argument (rng)")
	}
	return value.FromFloat(o.r.Float64()), nil
}

func (m *Module) randFree(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, runtime.Errorf("RAND.FREE expects rng handle")
	}
	if err := m.h.Free(heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}
