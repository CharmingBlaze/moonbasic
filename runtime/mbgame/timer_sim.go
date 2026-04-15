package mbgame

import (
	"fmt"
	"time"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// gameTimerSimObj is a delta-time driven timer (game/sim time), distinct from
// wall-clock TIMER.NEW handles.
type gameTimerSimObj struct {
	duration float64
	elapsed  float64
	running  bool
	loop     bool
	doneEdge bool
}

func (o *gameTimerSimObj) TypeName() string { return "GameTimerSim" }

func (o *gameTimerSimObj) TypeTag() uint16 { return heap.TagGameTimerSim }

func (o *gameTimerSimObj) Free() {}

func (o *gameTimerSimObj) remaining() float64 {
	if o.duration <= 0 {
		return 0
	}
	left := o.duration - o.elapsed
	if left < 0 {
		return 0
	}
	return left
}

func (o *gameTimerSimObj) fraction() float64 {
	if o.duration <= 0 {
		return 1
	}
	f := o.elapsed / o.duration
	if f > 1 {
		return 1
	}
	if f < 0 {
		return 0
	}
	return f
}

func (o *gameTimerSimObj) update(dt float64) {
	if !o.running || o.duration <= 0 {
		return
	}
	o.elapsed += dt
	if o.elapsed < o.duration {
		return
	}
	if o.loop {
		for o.elapsed >= o.duration {
			o.elapsed -= o.duration
		}
		o.doneEdge = true
	} else {
		o.elapsed = o.duration
		o.running = false
		o.doneEdge = true
	}
}

func (m *Module) getTimerSim(args []value.Value, ix int, op string) (*gameTimerSimObj, error) {
	if err := m.requireHeap(op); err != nil {
		return nil, err
	}
	if ix >= len(args) || args[ix].Kind != value.KindHandle {
		return nil, fmt.Errorf("%s: argument %d must be timer handle", op, ix+1)
	}
	return heap.Cast[*gameTimerSimObj](m.h, heap.Handle(args[ix].IVal))
}

func (m *Module) registerTimerSim(r runtime.Registrar) {
	simTimerNew := func(op string) func(*runtime.Runtime, ...value.Value) (value.Value, error) {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			_ = rt
			if err := m.requireHeap(op); err != nil {
				return value.Nil, err
			}
			if len(args) != 1 {
				return value.Nil, fmt.Errorf("%s expects 1 argument (duration)", op)
			}
			sec, ok := argF(args[0])
			if !ok || sec < 0 {
				return value.Nil, fmt.Errorf("%s: duration must be a non-negative number", op)
			}
			o := &gameTimerSimObj{duration: sec}
			id, err := m.h.Alloc(o)
			if err != nil {
				return value.Nil, err
			}
			return value.FromHandle(id), nil
		}
	}
	r.Register("TIMER.MAKE", "game", simTimerNew("TIMER.MAKE"))
	r.Register("TIMER.CREATE", "game", simTimerNew("TIMER.CREATE"))
	r.Register("TIMER.START", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		o, err := m.getTimerSim(args, 0, "TIMER.START")
		if err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("TIMER.START expects (timer)")
		}
		o.running = true
		o.elapsed = 0
		o.doneEdge = false
		return value.Nil, nil
	})
	r.Register("TIMER.STOP", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		o, err := m.getTimerSim(args, 0, "TIMER.STOP")
		if err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("TIMER.STOP expects (timer)")
		}
		o.running = false
		return value.Nil, nil
	})
	r.Register("TIMER.REWIND", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		o, err := m.getTimerSim(args, 0, "TIMER.REWIND")
		if err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("TIMER.REWIND expects (timer)")
		}
		o.elapsed = 0
		o.doneEdge = false
		return value.Nil, nil
	})
	r.Register("TIMER.SETLOOP", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		o, err := m.getTimerSim(args, 0, "TIMER.SETLOOP")
		if err != nil {
			return value.Nil, err
		}
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("TIMER.SETLOOP expects (timer, loop)")
		}
		var pool []string
		if rt.Prog != nil {
			pool = rt.Prog.StringTable
		}
		o.loop = value.Truthy(args[1], pool, rt.Heap)
		return value.Nil, nil
	})
	r.Register("TIMER.UPDATE", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		o, err := m.getTimerSim(args, 0, "TIMER.UPDATE")
		if err != nil {
			return value.Nil, err
		}
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("TIMER.UPDATE expects (timer, dt)")
		}
		dt, ok := argF(args[1])
		if !ok || dt < 0 {
			return value.Nil, fmt.Errorf("TIMER.UPDATE: dt must be a non-negative number")
		}
		o.update(dt)
		return value.Nil, nil
	})
	r.Register("TIMER.DONE", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		o, err := m.getTimerSim(args, 0, "TIMER.DONE")
		if err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("TIMER.DONE expects (timer)")
		}
		if o.doneEdge {
			o.doneEdge = false
			return value.FromBool(true), nil
		}
		return value.FromBool(false), nil
	})
	r.Register("TIMER.FRACTION", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		o, err := m.getTimerSim(args, 0, "TIMER.FRACTION")
		if err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("TIMER.FRACTION expects (timer)")
		}
		return value.FromFloat(o.fraction()), nil
	})
}

// registerTimerRemainingMerged wires TIMER.REMAINING to wall-clock and sim timers.
func (m *Module) registerTimerRemainingMerged(r runtime.Registrar) {
	r.Register("TIMER.REMAINING", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if err := m.requireHeap("TIMER.REMAINING"); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 || args[0].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("TIMER.REMAINING expects (timer)")
		}
		h := heap.Handle(args[0].IVal)
		if sim, err := heap.Cast[*gameTimerSimObj](m.h, h); err == nil {
			return value.FromFloat(sim.remaining()), nil
		}
		if wall, err := heap.Cast[*gameTimerObj](m.h, h); err == nil {
			rem := time.Until(wall.end).Seconds()
			if rem < 0 {
				rem = 0
			}
			return value.FromFloat(rem), nil
		}
		return value.Nil, fmt.Errorf("TIMER.REMAINING: not a timer handle")
	})
}
