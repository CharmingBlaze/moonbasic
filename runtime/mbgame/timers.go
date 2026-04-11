package mbgame

import (
	"fmt"
	"time"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

type gameTimerObj struct {
	end time.Time
}

func (o *gameTimerObj) TypeName() string { return "GameTimer" }

func (o *gameTimerObj) TypeTag() uint16 { return heap.TagGameTimer }

func (o *gameTimerObj) Free() {}

type stopwatchObj struct {
	t0 time.Time
}

func (o *stopwatchObj) TypeName() string { return "Stopwatch" }

func (o *stopwatchObj) TypeTag() uint16 { return heap.TagGameStopwatch }

func (o *stopwatchObj) Free() {}

func (m *Module) requireHeap(op string) error {
	if m.h == nil {
		return fmt.Errorf("%s: heap not bound", op)
	}
	return nil
}

func (m *Module) getTimer(args []value.Value, ix int, op string) (*gameTimerObj, error) {
	if err := m.requireHeap(op); err != nil {
		return nil, err
	}
	if ix >= len(args) || args[ix].Kind != value.KindHandle {
		return nil, fmt.Errorf("%s: argument %d must be timer handle", op, ix+1)
	}
	return heap.Cast[*gameTimerObj](m.h, heap.Handle(args[ix].IVal))
}

func (m *Module) getStopwatch(args []value.Value, ix int, op string) (*stopwatchObj, error) {
	if err := m.requireHeap(op); err != nil {
		return nil, err
	}
	if ix >= len(args) || args[ix].Kind != value.KindHandle {
		return nil, fmt.Errorf("%s: argument %d must be stopwatch handle", op, ix+1)
	}
	return heap.Cast[*stopwatchObj](m.h, heap.Handle(args[ix].IVal))
}

func (m *Module) registerTimers(r runtime.Registrar) {
	r.Register("TIMER.NEW", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if err := m.requireHeap("TIMER.NEW"); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("TIMER.NEW expects 1 argument (duration)")
		}
		sec, ok := argF(args[0])
		if !ok || sec < 0 {
			return value.Nil, fmt.Errorf("TIMER.NEW: duration must be a non-negative number")
		}
		o := &gameTimerObj{end: time.Now().Add(time.Duration(sec * float64(time.Second)))}
		id, err := m.h.Alloc(o)
		if err != nil {
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	})
	r.Register("TIMER.RESET", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		o, err := m.getTimer(args, 0, "TIMER.RESET")
		if err != nil {
			return value.Nil, err
		}
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("TIMER.RESET expects (timer, duration)")
		}
		sec, ok := argF(args[1])
		if !ok || sec < 0 {
			return value.Nil, fmt.Errorf("TIMER.RESET: duration must be a non-negative number")
		}
		o.end = time.Now().Add(time.Duration(sec * float64(time.Second)))
		return value.Nil, nil
	})
	r.Register("TIMER.FINISHED", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		o, err := m.getTimer(args, 0, "TIMER.FINISHED")
		if err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("TIMER.FINISHED expects (timer)")
		}
		return value.FromBool(time.Now().After(o.end)), nil
	})
	r.Register("TIMER.FREE", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if err := m.requireHeap("TIMER.FREE"); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 || args[0].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("TIMER.FREE expects timer handle")
		}
		if err := m.h.Free(heap.Handle(args[0].IVal)); err != nil {
			return value.Nil, err
		}
		return value.Nil, nil
	})

	r.Register("STOPWATCH.NEW", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if err := m.requireHeap("STOPWATCH.NEW"); err != nil {
			return value.Nil, err
		}
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("STOPWATCH.NEW expects 0 arguments")
		}
		o := &stopwatchObj{t0: time.Now()}
		id, err := m.h.Alloc(o)
		if err != nil {
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	})
	r.Register("STOPWATCH.RESET", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		o, err := m.getStopwatch(args, 0, "STOPWATCH.RESET")
		if err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("STOPWATCH.RESET expects (stopwatch)")
		}
		o.t0 = time.Now()
		return value.Nil, nil
	})
	r.Register("STOPWATCH.ELAPSED", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		o, err := m.getStopwatch(args, 0, "STOPWATCH.ELAPSED")
		if err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("STOPWATCH.ELAPSED expects (stopwatch)")
		}
		return value.FromFloat(time.Since(o.t0).Seconds()), nil
	})
	r.Register("STOPWATCH.FREE", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if err := m.requireHeap("STOPWATCH.FREE"); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 || args[0].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("STOPWATCH.FREE expects stopwatch handle")
		}
		if err := m.h.Free(heap.Handle(args[0].IVal)); err != nil {
			return value.Nil, err
		}
		return value.Nil, nil
	})

	m.registerTimerSim(r)
	m.registerTimerRemainingMerged(r)
}
