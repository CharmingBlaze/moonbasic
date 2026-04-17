package mbpool

import (
	"fmt"
	"strings"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

type poolObj struct {
	name    string
	max     int
	factory string
	reset   string
	free    []heap.Handle
	busy    map[heap.Handle]struct{}
	h       *heap.Store
	mod     *Module
	release heap.ReleaseOnce
}

func (o *poolObj) TypeName() string { return "Pool" }

func (o *poolObj) TypeTag() uint16 { return heap.TagPool }

// Free returns busy and pooled handles to the heap (children) before the pool slot is cleared.
func (o *poolObj) Free() {
	o.release.Do(func() {
		if o.h == nil {
			return
		}
		for h := range o.busy {
			o.h.Free(h)
		}
		o.busy = nil
		for _, h := range o.free {
			o.h.Free(h)
		}
		o.free = nil
	})
}

func valToHandle(v value.Value) (heap.Handle, error) {
	if v.Kind != value.KindHandle {
		return 0, fmt.Errorf("expected handle return value")
	}
	return heap.Handle(v.IVal), nil
}

func (m *Module) requireHeap() error {
	if m.h == nil {
		return runtime.Errorf("POOL.*: heap not bound")
	}
	return nil
}

func (m *Module) getPool(args []value.Value, ix int, op string) (*poolObj, error) {
	if ix >= len(args) || args[ix].Kind != value.KindHandle {
		return nil, fmt.Errorf("%s: expected pool handle", op)
	}
	return heap.Cast[*poolObj](m.h, heap.Handle(args[ix].IVal))
}

// Register implements runtime.Module.
func (m *Module) Register(reg runtime.Registrar) {
	reg.Register("POOL.CREATE", "pool", runtime.AdaptLegacy(m.poolMake))
	reg.Register("POOL.MAKE", "pool", runtime.AdaptLegacy(m.poolMake))
	reg.Register("POOL.SETFACTORY", "pool", m.poolSetFactory)
	reg.Register("POOL.SETRESET", "pool", m.poolSetReset)
	reg.Register("POOL.PREWARM", "pool", runtime.AdaptLegacy(m.poolPrewarm))
	reg.Register("POOL.GET", "pool", runtime.AdaptLegacy(m.poolGet))
	reg.Register("POOL.RETURN", "pool", runtime.AdaptLegacy(m.poolReturn))
	reg.Register("POOL.FREE", "pool", runtime.AdaptLegacy(m.poolFreePool))
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}

func (m *Module) poolMake(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("POOL.MAKE expects (name$, capacity)")
	}
	name := strings.TrimSpace(runtime.ArgString(args[0]))
	n, ok := args[1].ToInt()
	if !ok {
		if f, okf := args[1].ToFloat(); okf {
			n = int64(f)
			ok = true
		}
	}
	if !ok || n < 1 {
		return value.Nil, fmt.Errorf("POOL.MAKE: capacity must be a positive integer")
	}
	o := &poolObj{
		name: name,
		max:  int(n),
		busy: make(map[heap.Handle]struct{}),
		h:    m.h,
		mod:  m,
	}
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) poolSetFactory(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("POOL.SETFACTORY expects (pool, functionName$)")
	}
	o, err := m.getPool(args, 0, "POOL.SETFACTORY")
	if err != nil {
		return value.Nil, err
	}
	fn, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	fn = strings.ToUpper(strings.TrimSpace(fn))
	if fn == "" {
		return value.Nil, fmt.Errorf("POOL.SETFACTORY: empty function name")
	}
	o.factory = fn
	return args[0], nil
}

func (m *Module) poolSetReset(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("POOL.SETRESET expects (pool, functionName$)")
	}
	o, err := m.getPool(args, 0, "POOL.SETRESET")
	if err != nil {
		return value.Nil, err
	}
	fn, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	fn = strings.ToUpper(strings.TrimSpace(fn))
	if fn == "" {
		return value.Nil, fmt.Errorf("POOL.SETRESET: empty function name")
	}
	o.reset = fn
	return args[0], nil
}

func (o *poolObj) runFactory() (heap.Handle, error) {
	if o.factory == "" {
		return 0, fmt.Errorf("POOL: factory not set (POOL.SETFACTORY)")
	}
	if o.mod == nil || o.mod.invoke == nil {
		return 0, fmt.Errorf("POOL: user function invoker not configured")
	}
	v, err := o.mod.invoke(o.factory, nil)
	if err != nil {
		return 0, err
	}
	return valToHandle(v)
}

func (o *poolObj) runReset(h heap.Handle) error {
	if o.reset == "" {
		return nil
	}
	if o.mod == nil || o.mod.invoke == nil {
		return fmt.Errorf("POOL: user function invoker not configured")
	}
	_, err := o.mod.invoke(o.reset, []value.Value{value.FromHandle(int32(h))})
	return err
}

func (m *Module) poolPrewarm(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("POOL.PREWARM expects (pool)")
	}
	o, err := m.getPool(args, 0, "POOL.PREWARM")
	if err != nil {
		return value.Nil, err
	}
	if o.factory == "" {
		return value.Nil, fmt.Errorf("POOL.PREWARM: set factory first")
	}
	for len(o.busy)+len(o.free) < o.max {
		h, err := o.runFactory()
		if err != nil {
			return value.Nil, err
		}
		if err := o.runReset(h); err != nil {
			o.h.Free(h)
			return value.Nil, err
		}
		o.free = append(o.free, h)
	}
	return args[0], nil
}

func (m *Module) poolGet(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("POOL.GET expects (pool)")
	}
	o, err := m.getPool(args, 0, "POOL.GET")
	if err != nil {
		return value.Nil, err
	}
	if o.factory == "" {
		return value.Nil, fmt.Errorf("POOL.GET: set factory first")
	}

	var h heap.Handle
	fromFree := false
	if n := len(o.free); n > 0 {
		h = o.free[n-1]
		o.free = o.free[:n-1]
		fromFree = true
	} else {
		if len(o.busy) >= o.max {
			return value.Nil, fmt.Errorf("POOL.GET: pool exhausted (capacity %d)", o.max)
		}
		h, err = o.runFactory()
		if err != nil {
			return value.Nil, err
		}
	}
	if err := o.runReset(h); err != nil {
		if fromFree {
			o.free = append(o.free, h)
		} else {
			o.h.Free(h)
		}
		return value.Nil, err
	}
	o.busy[h] = struct{}{}
	return value.FromHandle(int32(h)), nil
}

func (m *Module) poolReturn(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("POOL.RETURN expects (pool, handle)")
	}
	o, err := m.getPool(args, 0, "POOL.RETURN")
	if err != nil {
		return value.Nil, err
	}
	hv, err := valToHandle(args[1])
	if err != nil {
		return value.Nil, fmt.Errorf("POOL.RETURN: %w", err)
	}
	if _, ok := o.busy[hv]; !ok {
		return value.Nil, fmt.Errorf("POOL.RETURN: handle is not checked out from this pool")
	}
	delete(o.busy, hv)
	if err := o.runReset(hv); err != nil {
		o.busy[hv] = struct{}{}
		return value.Nil, err
	}
	o.free = append(o.free, hv)
	return args[0], nil
}

func (m *Module) poolFreePool(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("POOL.FREE expects (pool)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("POOL.FREE: pool must be a handle")
	}
	m.h.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}
