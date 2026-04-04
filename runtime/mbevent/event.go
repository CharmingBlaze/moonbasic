package mbevent

import (
	"fmt"
	"strings"
	"sync"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

type listenerRec struct {
	fn   string
	once bool
}

var (
	evMu   sync.RWMutex
	evSubs = make(map[string][]listenerRec)
)

func normEvent(s string) string { return strings.ToLower(strings.TrimSpace(s)) }

// offRemoveFirst removes the first listener on name matching fnU.
func offRemoveFirst(name, fnU string) {
	evMu.Lock()
	defer evMu.Unlock()
	ls := evSubs[name]
	for i, r := range ls {
		if r.fn == fnU {
			evSubs[name] = append(ls[:i], ls[i+1:]...)
			return
		}
	}
}

func removeAllNamed(name, fnU string) {
	evMu.Lock()
	defer evMu.Unlock()
	ls := evSubs[name]
	out := ls[:0]
	for _, r := range ls {
		if r.fn != fnU {
			out = append(out, r)
		}
	}
	if len(out) == 0 {
		delete(evSubs, name)
	} else {
		evSubs[name] = out
	}
}

// Register implements runtime.Module.
func (m *Module) Register(reg runtime.Registrar) {
	reg.Register("EVENT.ON", "event", m.evOn)
	reg.Register("EVENT.ONCE", "event", m.evOnce)
	reg.Register("EVENT.OFF", "event", m.evOff)
	reg.Register("EVENT.FIRE", "event", m.evFire)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {
	evMu.Lock()
	defer evMu.Unlock()
	evSubs = make(map[string][]listenerRec)
}

func (m *Module) evOn(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindString || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("EVENT.ON expects (eventName$, functionName$)")
	}
	en, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	fn, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	name := normEvent(en)
	fnU := strings.ToUpper(strings.TrimSpace(fn))
	if name == "" || fnU == "" {
		return value.Nil, fmt.Errorf("EVENT.ON: names must be non-empty")
	}
	evMu.Lock()
	defer evMu.Unlock()
	evSubs[name] = append(evSubs[name], listenerRec{fn: fnU, once: false})
	return value.Nil, nil
}

func (m *Module) evOnce(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindString || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("EVENT.ONCE expects (eventName$, functionName$)")
	}
	en, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	fn, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	name := normEvent(en)
	fnU := strings.ToUpper(strings.TrimSpace(fn))
	if name == "" || fnU == "" {
		return value.Nil, fmt.Errorf("EVENT.ONCE: names must be non-empty")
	}
	evMu.Lock()
	defer evMu.Unlock()
	evSubs[name] = append(evSubs[name], listenerRec{fn: fnU, once: true})
	return value.Nil, nil
}

func (m *Module) evOff(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindString || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("EVENT.OFF expects (eventName$, functionName$)")
	}
	en, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	fn, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	name := normEvent(en)
	fnU := strings.ToUpper(strings.TrimSpace(fn))
	if name == "" || fnU == "" {
		return value.Nil, fmt.Errorf("EVENT.OFF: names must be non-empty")
	}
	removeAllNamed(name, fnU)
	return value.Nil, nil
}

func (m *Module) evFire(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) < 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("EVENT.FIRE expects (eventName$, ...payload)")
	}
	en, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	name := normEvent(en)
	if name == "" {
		return value.Nil, fmt.Errorf("EVENT.FIRE: empty event name")
	}
	if m.invoke == nil {
		return value.Nil, runtime.Errorf("EVENT.FIRE: user function invoker not configured")
	}

	evMu.RLock()
	ls := evSubs[name]
	copyLs := append([]listenerRec(nil), ls...)
	evMu.RUnlock()

	payload := args[1:]
	var onceToRemove []string

	for _, rec := range copyLs {
		_, err := m.invoke(rec.fn, payload)
		if err != nil {
			return value.Nil, err
		}
		if rec.once {
			onceToRemove = append(onceToRemove, rec.fn)
		}
	}
	for _, fnU := range onceToRemove {
		offRemoveFirst(name, fnU)
	}
	return value.Nil, nil
}
