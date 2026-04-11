// Package mbscene implements SCENE.* — register loader FUNCTION names, run update/draw hooks, and optional transitions.
package mbscene

import (
	"fmt"
	"strings"
	"sync"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

// Module holds scene registry and active lifecycle hooks.
type Module struct {
	mu     sync.Mutex
	invoke func(string, []value.Value) (value.Value, error)

	loaders map[string]string // scene id (upper) -> load FUNCTION name (upper)

	currentID string
	updateFn  string
	drawFn    string

	asyncLoad string

	trPending     string
	trDur         float64
	trPhase       int // 0 none, 1 waiting transition out, 2 waiting fade in
	trNeedsFadeIn bool
}

// NewModule creates the scene module.
func NewModule() *Module {
	return &Module{
		loaders: make(map[string]string),
	}
}

// SetUserInvoker wires VM.CallUserFunction for scene loaders and per-frame hooks.
func (m *Module) SetUserInvoker(fn func(string, []value.Value) (value.Value, error)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.invoke = fn
}

// Register implements runtime.Module.
func (m *Module) Register(reg runtime.Registrar) {
	reg.Register("SCENE.REGISTER", "scene", m.sceneRegister)
	reg.Register("SCENE.SETHANDLERS", "scene", m.sceneSetHandlers)
	reg.Register("SCENE.LOAD", "scene", m.sceneLoad)
	reg.Register("SCENE.LOADASYNC", "scene", m.sceneLoadAsync)
	reg.Register("SCENE.LOADWITHTRANSITION", "scene", m.sceneLoadWithTransition)
	reg.Register("SCENE.UPDATE", "scene", m.sceneUpdate)
	reg.Register("SCENE.DRAW", "scene", m.sceneDraw)
	reg.Register("SCENE.CURRENT", "scene", m.sceneCurrent)
	reg.Register("SCENE.SWITCH", "scene", m.sceneSwitch)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.loaders = make(map[string]string)
	m.currentID, m.updateFn, m.drawFn = "", "", ""
	m.asyncLoad, m.trPending = "", ""
	m.trPhase = 0
}

// Reset implements runtime.Module.
func (m *Module) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.currentID, m.updateFn, m.drawFn = "", "", ""
	m.asyncLoad, m.trPending = "", ""
	m.trPhase = 0
}

func up(s string) string { return strings.ToUpper(strings.TrimSpace(s)) }

func truthy(reg *runtime.Runtime, v value.Value) bool {
	if reg == nil {
		return value.Truthy(v, nil, nil)
	}
	var pool []string
	if reg.Prog != nil {
		pool = reg.Prog.StringTable
	}
	return value.Truthy(v, pool, reg.Heap)
}

func (m *Module) sceneRegister(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindString || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("SCENE.REGISTER expects (sceneId, loadFunctionName)")
	}
	id, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	fn, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	idU, fnU := up(id), up(fn)
	if idU == "" || fnU == "" {
		return value.Nil, fmt.Errorf("SCENE.REGISTER: scene id and function name must be non-empty")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.loaders[idU] = fnU
	return value.Nil, nil
}

func (m *Module) sceneSetHandlers(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindString || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("SCENE.SETHANDLERS expects (updateFunctionName, drawFunctionName)")
	}
	u, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	d, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.updateFn, m.drawFn = up(u), up(d)
	return value.Nil, nil
}

func (m *Module) sceneLoad(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("SCENE.LOAD expects sceneId")
	}
	id, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	m.mu.Lock()
	m.asyncLoad = ""
	m.trPending, m.trPhase = "", 0
	m.mu.Unlock()
	return value.Nil, m.runLoad(rt, up(id))
}

func (m *Module) sceneLoadAsync(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("SCENE.LOADASYNC expects sceneId")
	}
	id, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	idU := up(id)
	if idU == "" {
		return value.Nil, fmt.Errorf("SCENE.LOADASYNC: scene id must be non-empty")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.asyncLoad = idU
	return value.Nil, nil
}

func (m *Module) sceneLoadWithTransition(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 || args[0].Kind != value.KindString || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("SCENE.LOADWITHTRANSITION expects (sceneId, kind, duration)")
	}
	id, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	kindStr, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	dur, ok := args[2].ToFloat()
	if !ok {
		if i, ok2 := args[2].ToInt(); ok2 {
			dur = float64(i)
			ok = true
		}
	}
	if !ok || dur <= 0 {
		return value.Nil, fmt.Errorf("SCENE.LOADWITHTRANSITION: duration must be positive")
	}
	idU, kindU := up(id), strings.ToLower(strings.TrimSpace(kindStr))
	if idU == "" {
		return value.Nil, fmt.Errorf("SCENE.LOADWITHTRANSITION: scene id must be non-empty")
	}
	switch kindU {
	case "fade", "wipe":
	default:
		return value.Nil, fmt.Errorf("SCENE.LOADWITHTRANSITION: kind must be \"fade\" or \"wipe\"")
	}

	reg := runtime.ActiveRegistry()
	if reg == nil {
		return value.Nil, runtime.Errorf("SCENE.LOADWITHTRANSITION: no active runtime")
	}

	m.mu.Lock()
	m.asyncLoad = ""
	m.trPending = idU
	m.trDur = dur
	m.trPhase = 1
	m.trNeedsFadeIn = true
	m.mu.Unlock()

	if err := startTransitionOut(rt, reg, kindU, dur); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func startTransitionOut(rt *runtime.Runtime, reg *runtime.Runtime, kind string, dur float64) error {
	if kind == "fade" {
		_, err := reg.Call("TRANSITION.FADEOUT", []value.Value{value.FromFloat(dur)})
		return err
	}
	left := rt.RetString("left")
	_, err := reg.Call("TRANSITION.WIPE", []value.Value{left, value.FromFloat(dur)})
	return err
}

func (m *Module) transitionPoll(rt *runtime.Runtime) error {
	reg := runtime.ActiveRegistry()
	if reg == nil {
		return nil
	}

	m.mu.Lock()
	phase := m.trPhase
	pending := m.trPending
	dur := m.trDur
	needsIn := m.trNeedsFadeIn
	m.mu.Unlock()

	if phase == 0 {
		return nil
	}

	doneVal, err := reg.Call("TRANSITION.ISDONE", nil)
	if err != nil {
		return err
	}
	if !truthy(reg, doneVal) {
		return nil
	}

	switch phase {
	case 1:
		if err := m.runLoad(rt, pending); err != nil {
			m.mu.Lock()
			m.trPhase, m.trPending = 0, ""
			m.mu.Unlock()
			return err
		}
		if !needsIn {
			m.mu.Lock()
			m.trPhase, m.trPending = 0, ""
			m.mu.Unlock()
			return nil
		}
		if _, err := reg.Call("TRANSITION.FADEIN", []value.Value{value.FromFloat(dur)}); err != nil {
			m.mu.Lock()
			m.trPhase, m.trPending = 0, ""
			m.mu.Unlock()
			return err
		}
		m.mu.Lock()
		m.trPhase = 2
		m.mu.Unlock()
	case 2:
		m.mu.Lock()
		m.trPhase, m.trPending = 0, ""
		m.mu.Unlock()
	}
	return nil
}

func (m *Module) runLoad(rt *runtime.Runtime, sceneID string) error {
	if sceneID == "" {
		return fmt.Errorf("SCENE: empty scene id")
	}
	m.mu.Lock()
	loader := m.loaders[sceneID]
	invoke := m.invoke
	m.mu.Unlock()

	if loader == "" {
		return fmt.Errorf("SCENE: unknown scene %q (use SCENE.REGISTER first)", sceneID)
	}
	if invoke == nil {
		return runtime.Errorf("SCENE: user function invoker not configured")
	}

	m.mu.Lock()
	m.updateFn, m.drawFn = "", ""
	m.mu.Unlock()

	if _, err := invoke(loader, nil); err != nil {
		return err
	}

	m.mu.Lock()
	m.currentID = sceneID
	m.mu.Unlock()
	return nil
}

func (m *Module) sceneUpdate(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("SCENE.UPDATE expects (dt)")
	}
	dt, ok := args[0].ToFloat()
	if !ok {
		if i, ok2 := args[0].ToInt(); ok2 {
			dt = float64(i)
			ok = true
		}
	}
	if !ok {
		return value.Nil, fmt.Errorf("SCENE.UPDATE: dt must be numeric")
	}

	m.mu.Lock()
	async := m.asyncLoad
	m.asyncLoad = ""
	m.mu.Unlock()
	if async != "" {
		if err := m.runLoad(rt, async); err != nil {
			return value.Nil, err
		}
	}

	if err := m.transitionPoll(rt); err != nil {
		return value.Nil, err
	}

	m.mu.Lock()
	fn := m.updateFn
	invoke := m.invoke
	m.mu.Unlock()

	if fn == "" || invoke == nil {
		return value.Nil, nil
	}
	_, err := invoke(fn, []value.Value{value.FromFloat(dt)})
	return value.Nil, err
}

func (m *Module) sceneDraw(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = args
	if err := m.transitionPoll(rt); err != nil {
		return value.Nil, err
	}

	m.mu.Lock()
	fn := m.drawFn
	invoke := m.invoke
	m.mu.Unlock()

	if fn == "" || invoke == nil {
		return value.Nil, nil
	}
	_, err := invoke(fn, nil)
	return value.Nil, err
}

func (m *Module) sceneCurrent(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = args
	m.mu.Lock()
	id := m.currentID
	m.mu.Unlock()
	return rt.RetString(id), nil
}

func (m *Module) sceneSwitch(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("SCENE.SWITCH expects (sceneId, fadeDuration)")
	}
	// Use default "fade" transition
	return m.sceneLoadWithTransition(rt, args[0], value.FromStringIndex(rt.Heap.Intern("fade")), args[1])
}
