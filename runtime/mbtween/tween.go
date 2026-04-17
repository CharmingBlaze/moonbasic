package mbtween

import (
	"fmt"
	"strings"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

type tweenStep struct {
	VarName string
	Target  float64
	Dur     float64
	Ease    string
	fwdSnap float64 // global value when forward leg of this step began (for yoyo reverse)
}

type tweenObj struct {
	steps []tweenStep

	idx          int
	phaseForward bool // false = reverse pass (yoyo)
	tStep        float64
	entered      bool
	segOrigin    float64
	segDest      float64

	running bool
	yoyo    bool
	loopMax int // -1 = infinite, >0 = max cycles, 0 = unset -> 1

	loopsDone int
	onDone    string

	h   *heap.Store
	mod *Module
}

func (o *tweenObj) TypeName() string { return "Tween" }

func (o *tweenObj) TypeTag() uint16 { return heap.TagTween }

func (o *tweenObj) Free() {}

func normGlobal(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}

func (m *Module) readGlobal(name string) (float64, error) {
	k := normGlobal(name)
	if m.getG == nil {
		return 0, fmt.Errorf("TWEEN: global accessor not configured")
	}
	v, ok := m.getG(k)
	if !ok {
		return 0, nil
	}
	if f, ok := v.ToFloat(); ok {
		return f, nil
	}
	if i, ok := v.ToInt(); ok {
		return float64(i), nil
	}
	return 0, fmt.Errorf("TWEEN: global %q is not numeric", k)
}

func (m *Module) writeGlobal(name string, f float64) error {
	k := normGlobal(name)
	if m.setG == nil {
		return fmt.Errorf("TWEEN: global accessor not configured")
	}
	m.setG(k, value.FromFloat(f))
	return nil
}

func (m *Module) requireHeap() error {
	if m.h == nil {
		return runtime.Errorf("TWEEN.*: heap not bound")
	}
	return nil
}

func (m *Module) getTween(args []value.Value, ix int, op string) (*tweenObj, error) {
	if ix >= len(args) || args[ix].Kind != value.KindHandle {
		return nil, fmt.Errorf("%s: expected tween handle", op)
	}
	return heap.Cast[*tweenObj](m.h, heap.Handle(args[ix].IVal))
}

// Register implements runtime.Module.
func (m *Module) Register(reg runtime.Registrar) {
	reg.Register("TWEEN.CREATE", "tween", runtime.AdaptLegacy(m.twMake))
	reg.Register("TWEEN.MAKE", "tween", runtime.AdaptLegacy(m.twMake))
	reg.Register("TWEEN.TO", "tween", m.twTo)
	reg.Register("TWEEN.THEN", "tween", m.twThen)
	reg.Register("TWEEN.ONCOMPLETE", "tween", m.twOnComplete)
	reg.Register("TWEEN.START", "tween", runtime.AdaptLegacy(m.twStart))
	reg.Register("TWEEN.UPDATE", "tween", runtime.AdaptLegacy(m.twUpdate))
	reg.Register("TWEEN.LOOP", "tween", runtime.AdaptLegacy(m.twLoop))
	reg.Register("TWEEN.YOYO", "tween", runtime.AdaptLegacy(m.twYoyo))
	reg.Register("TWEEN.STOP", "tween", runtime.AdaptLegacy(m.twStop))

	reg.Register("TWEEN.ISPLAYING", "tween", runtime.AdaptLegacy(m.twIsPlaying))
	reg.Register("TWEEN.ISFINISHED", "tween", runtime.AdaptLegacy(m.twIsFinished))
	reg.Register("TWEEN.PROGRESS", "tween", runtime.AdaptLegacy(m.twProgress))
	reg.Register("TWEEN.GETLOOP", "tween", runtime.AdaptLegacy(m.twGetLoop))
	reg.Register("TWEEN.GETYOYO", "tween", runtime.AdaptLegacy(m.twGetYoyo))
	reg.Register("TWEEN.FREE", "tween", runtime.AdaptLegacy(m.twFree))
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}

func (m *Module) twMake(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("TWEEN.MAKE expects 0 arguments")
	}
	o := &tweenObj{
		phaseForward: true,
		loopMax:      1,
		h:            m.h,
		mod:          m,
	}
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) twTo(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 5 || args[1].Kind != value.KindString || args[4].Kind != value.KindString {
		return value.Nil, fmt.Errorf("TWEEN.TO expects (tween, varName$, target, seconds, easing$)")
	}
	o, err := m.getTween(args, 0, "TWEEN.TO")
	if err != nil {
		return value.Nil, err
	}
	if o.running {
		return value.Nil, fmt.Errorf("TWEEN.TO: cannot modify while running (STOP first)")
	}
	vn, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	to, ok1 := args[2].ToFloat()
	if !ok1 {
		if i, ok := args[2].ToInt(); ok {
			to = float64(i)
			ok1 = true
		}
	}
	if !ok1 {
		return value.Nil, fmt.Errorf("TWEEN.TO: target must be numeric")
	}
	dur, ok2 := args[3].ToFloat()
	if !ok2 {
		if i, ok := args[3].ToInt(); ok {
			dur = float64(i)
			ok2 = true
		}
	}
	if !ok2 || dur <= 0 {
		return value.Nil, fmt.Errorf("TWEEN.TO: duration must be positive")
	}
	easeName, err := rt.ArgString(args, 4)
	if err != nil {
		return value.Nil, err
	}
	o.steps = append(o.steps, tweenStep{
		VarName: normGlobal(vn),
		Target:  to,
		Dur:     dur,
		Ease:    easeName,
	})
	return args[0], nil
}

func (m *Module) twThen(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return m.twTo(rt, args...)
}

func (m *Module) twOnComplete(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("TWEEN.ONCOMPLETE expects (tween, functionName$)")
	}
	o, err := m.getTween(args, 0, "TWEEN.ONCOMPLETE")
	if err != nil {
		return value.Nil, err
	}
	if o.running {
		return value.Nil, fmt.Errorf("TWEEN.ONCOMPLETE: cannot change while running")
	}
	fn, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	fn = strings.ToUpper(strings.TrimSpace(fn))
	if fn == "" {
		return value.Nil, fmt.Errorf("TWEEN.ONCOMPLETE: empty function name")
	}
	o.onDone = fn
	return args[0], nil
}

func (m *Module) twStart(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("TWEEN.START expects (tween)")
	}
	o, err := m.getTween(args, 0, "TWEEN.START")
	if err != nil {
		return value.Nil, err
	}
	if len(o.steps) == 0 {
		return value.Nil, fmt.Errorf("TWEEN.START: no steps (use TWEEN.TO first)")
	}
	o.running = true
	o.idx = 0
	o.phaseForward = true
	o.tStep = 0
	o.entered = false
	o.loopsDone = 0
	return args[0], nil
}

func (m *Module) twLoop(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("TWEEN.LOOP expects (tween, count)")
	}
	o, err := m.getTween(args, 0, "TWEEN.LOOP")
	if err != nil {
		return value.Nil, err
	}
	if o.running {
		return value.Nil, fmt.Errorf("TWEEN.LOOP: cannot change while running")
	}
	n, ok := args[1].ToInt()
	if !ok {
		if f, okf := args[1].ToFloat(); okf {
			n = int64(f)
			ok = true
		}
	}
	if !ok {
		return value.Nil, fmt.Errorf("TWEEN.LOOP: count must be numeric")
	}
	if n <= 0 {
		o.loopMax = -1
	} else {
		o.loopMax = int(n)
	}
	return args[0], nil
}

func (m *Module) twYoyo(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("TWEEN.YOYO expects (tween)")
	}
	o, err := m.getTween(args, 0, "TWEEN.YOYO")
	if err != nil {
		return value.Nil, err
	}
	if o.running {
		return value.Nil, fmt.Errorf("TWEEN.YOYO: cannot change while running")
	}
	o.yoyo = true
	return args[0], nil
}

func (m *Module) twGetLoop(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("TWEEN.GETLOOP expects (tween)")
	}
	o, err := m.getTween(args, 0, "TWEEN.GETLOOP")
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(o.loopMax)), nil
}

func (m *Module) twGetYoyo(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("TWEEN.GETYOYO expects (tween)")
	}
	o, err := m.getTween(args, 0, "TWEEN.GETYOYO")
	if err != nil {
		return value.Nil, err
	}
	return value.FromBool(o.yoyo), nil
}

func (m *Module) twFree(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("TWEEN.FREE expects (tween)")
	}
	m.h.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}

func (m *Module) twStop(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("TWEEN.STOP expects (tween)")
	}
	o, err := m.getTween(args, 0, "TWEEN.STOP")
	if err != nil {
		return value.Nil, err
	}
	o.running = false
	o.entered = false
	o.tStep = 0
	return args[0], nil
}

func (m *Module) twIsPlaying(args []value.Value) (value.Value, error) {
	o, err := m.getTween(args, 0, "TWEEN.ISPLAYING")
	if err != nil {
		return value.FromBool(false), nil
	}
	return value.FromBool(o.running), nil
}

func (m *Module) twIsFinished(args []value.Value) (value.Value, error) {
	o, err := m.getTween(args, 0, "TWEEN.ISFINISHED")
	if err != nil {
		return value.FromBool(true), nil
	}
	return value.FromBool(!o.running && o.loopsDone >= o.loopMax && o.loopMax > 0), nil
}

func (m *Module) twProgress(args []value.Value) (value.Value, error) {
	o, err := m.getTween(args, 0, "TWEEN.PROGRESS")
	if err != nil {
		return value.FromFloat(0), nil
	}
	if len(o.steps) == 0 {
		return value.FromFloat(0), nil
	}
	s := o.steps[o.idx]
	if s.Dur <= 0 {
		return value.FromFloat(1), nil
	}
	return value.FromFloat(o.tStep / s.Dur), nil
}

func (o *tweenObj) beginSegment(m *Module) error {
	s := &o.steps[o.idx]
	origin, err := m.readGlobal(s.VarName)
	if err != nil {
		return err
	}
	o.segOrigin = origin
	if o.phaseForward {
		s.fwdSnap = origin
		o.segDest = s.Target
	} else {
		o.segDest = s.fwdSnap
	}
	return nil
}

func (o *tweenObj) finishTween(m *Module) {
	o.running = false
	o.entered = false
	o.tStep = 0
	if o.onDone == "" || m.invoke == nil {
		return
	}
	_, _ = m.invoke(o.onDone, nil)
}

func (o *tweenObj) completeCycle(m *Module) {
	if o.loopMax < 0 {
		// Infinite loop — restart forward without incrementing toward a cap.
		o.phaseForward = true
		o.idx = 0
		o.tStep = 0
		o.entered = false
		return
	}
	o.loopsDone++
	if o.loopsDone >= o.loopMax {
		o.finishTween(m)
		return
	}
	o.phaseForward = true
	o.idx = 0
	o.tStep = 0
	o.entered = false
}

func (o *tweenObj) advanceAfterStep(m *Module) error {
	n := len(o.steps)
	if o.phaseForward {
		next := o.idx + 1
		if next >= n {
			if o.yoyo {
				o.phaseForward = false
				o.idx = n - 1
				o.tStep = 0
				o.entered = false
				return nil
			}
			o.completeCycle(m)
			if o.running {
				o.tStep = 0
				o.entered = false
			}
			return nil
		}
		o.idx = next
	} else {
		next := o.idx - 1
		if next < 0 {
			o.phaseForward = true
			o.idx = 0
			o.completeCycle(m)
			if o.running {
				o.tStep = 0
				o.entered = false
			}
			return nil
		}
		o.idx = next
	}
	o.tStep = 0
	o.entered = false
	return nil
}

func (m *Module) twUpdate(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("TWEEN.UPDATE expects (tween, dt)")
	}
	o, err := m.getTween(args, 0, "TWEEN.UPDATE")
	if err != nil {
		return value.Nil, err
	}
	dt, ok := args[1].ToFloat()
	if !ok {
		if i, ok2 := args[1].ToInt(); ok2 {
			dt = float64(i)
			ok = true
		}
	}
	if !ok {
		return value.Nil, fmt.Errorf("TWEEN.UPDATE: dt must be numeric")
	}
	if dt < 0 {
		dt = 0
	}
	if !o.running || len(o.steps) == 0 {
		return value.Nil, nil
	}

	// Consume time — may complete multiple steps in one frame
	for dt > 0 && o.running {
		if !o.entered {
			if err := o.beginSegment(m); err != nil {
				o.running = false
				return value.Nil, err
			}
			o.entered = true
		}
		s := &o.steps[o.idx]
		dur := s.Dur
		if dur <= 0 {
			dur = 1e-6
		}
		if o.tStep >= dur {
			if err := m.writeGlobal(s.VarName, o.segDest); err != nil {
				o.running = false
				return value.Nil, err
			}
			if err := o.advanceAfterStep(m); err != nil {
				return value.Nil, err
			}
			continue
		}
		remain := dur - o.tStep
		if remain <= 0 {
			remain = 0
		}
		use := dt
		if use > remain {
			use = remain
		}
		o.tStep += use
		dt -= use

		p := o.tStep / dur
		if p >= 1 {
			if err := m.writeGlobal(s.VarName, o.segDest); err != nil {
				o.running = false
				return value.Nil, err
			}
			if err := o.advanceAfterStep(m); err != nil {
				return value.Nil, err
			}
			continue
		}
		cur := o.segOrigin + (o.segDest-o.segOrigin)*Ease(p, s.Ease)
		if err := m.writeGlobal(s.VarName, cur); err != nil {
			o.running = false
			return value.Nil, err
		}
		break
	}
	return value.Nil, nil
}
