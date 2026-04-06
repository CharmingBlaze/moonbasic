//go:build cgo || (windows && !cgo)

package mbsprite

import (
	"fmt"
	"strconv"
	"strings"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

type animStateDef struct {
	first, last int
	fps         float32
	loop        bool
}

type animTransition struct {
	from, to, cond string
}

type animMachine struct {
	states      map[string]animStateDef
	trans       []animTransition
	floatParams map[string]float64
	boolParams  map[string]bool
	current     string
	accum       float32
	frameIdx    int // 0..(last-first)
}

func ensureAnim(s *spriteObj) *animMachine {
	if s.anim == nil {
		s.anim = &animMachine{
			states:      make(map[string]animStateDef),
			floatParams: make(map[string]float64),
			boolParams:  make(map[string]bool),
		}
	}
	return s.anim
}

func (s *spriteObj) recomputeAnimCellWidth() {
	if s.anim == nil || s.atlasRegionW <= 0 {
		return
	}
	max := 0
	for _, st := range s.anim.states {
		if st.last > max {
			max = st.last
		}
	}
	if max <= 0 {
		return
	}
	denom := int32(max + 1)
	if denom > 0 && s.atlasRegionW%denom == 0 {
		s.frameW = s.atlasRegionW / denom
		s.frameH = s.atlasRegionH
	}
}

func (s *spriteObj) syncAnimFrame() {
	if s.anim == nil || s.anim.current == "" {
		return
	}
	st, ok := s.anim.states[s.anim.current]
	if !ok {
		return
	}
	am := s.anim
	if am.frameIdx > st.last-st.first {
		am.frameIdx = st.last - st.first
	}
	if am.frameIdx < 0 {
		am.frameIdx = 0
	}
	s.curFrame = st.first + am.frameIdx
}

func (m *Module) registerAnim(reg runtime.Registrar) {
	reg.Register("ANIM.DEFINE", "sprite", m.animDefine)
	reg.Register("ANIM.ADDTRANSITION", "sprite", m.animAddTransition)
	reg.Register("ANIM.UPDATE", "sprite", m.animUpdate)
	reg.Register("ANIM.SETPARAM", "sprite", m.animSetParam)
}

func argInt(v value.Value) (int, bool) {
	if i, ok := v.ToInt(); ok {
		return int(i), true
	}
	if f, ok := v.ToFloat(); ok {
		return int(f), true
	}
	return 0, false
}

func (m *Module) animDefine(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 6 || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("ANIM.DEFINE expects (sprite, name$, first, last, fps, loop)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("ANIM.DEFINE: invalid sprite handle")
	}
	s, err := heap.Cast[*spriteObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	name, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	fi, ok1 := argInt(args[2])
	la, ok2 := argInt(args[3])
	fps, ok3 := argF(args[4])
	loop := args[5].Kind == value.KindBool && args[5].IVal != 0
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("ANIM.DEFINE: first, last, fps must be numeric")
	}
	if fi > la {
		fi, la = la, fi
	}
	if fps <= 0 {
		fps = 1
	}
	am := ensureAnim(s)
	am.states[strings.TrimSpace(name)] = animStateDef{first: fi, last: la, fps: fps, loop: loop}
	s.recomputeAnimCellWidth()
	if am.current == "" {
		am.current = strings.TrimSpace(name)
		am.frameIdx = 0
		am.accum = 0
	}
	s.syncAnimFrame()
	return value.Nil, nil
}

func (m *Module) animAddTransition(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 4 || args[1].Kind != value.KindString || args[2].Kind != value.KindString || args[3].Kind != value.KindString {
		return value.Nil, fmt.Errorf("ANIM.ADDTRANSITION expects (sprite, from$, to$, condition$)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("ANIM.ADDTRANSITION: invalid sprite handle")
	}
	s, err := heap.Cast[*spriteObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	from, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	to, err := rt.ArgString(args, 2)
	if err != nil {
		return value.Nil, err
	}
	cond, err := rt.ArgString(args, 3)
	if err != nil {
		return value.Nil, err
	}
	am := ensureAnim(s)
	am.trans = append(am.trans, animTransition{from: strings.TrimSpace(from), to: strings.TrimSpace(to), cond: strings.TrimSpace(cond)})
	return value.Nil, nil
}

func (m *Module) animUpdate(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ANIM.UPDATE expects (sprite, dt)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("ANIM.UPDATE: invalid sprite handle")
	}
	s, err := heap.Cast[*spriteObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	dt, ok := argF(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("ANIM.UPDATE: dt must be numeric")
	}
	if dt < 0 {
		dt = 0
	}
	if s.anim == nil || s.anim.current == "" {
		return value.Nil, nil
	}
	s.recomputeAnimCellWidth()
	am := s.anim
	// Transitions from current state
	for _, tr := range am.trans {
		if tr.from != am.current {
			continue
		}
		ok, err := evalAnimCond(tr.cond, am.floatParams, am.boolParams)
		if err != nil || !ok {
			continue
		}
		if _, exists := am.states[tr.to]; exists {
			am.current = tr.to
			am.frameIdx = 0
			am.accum = 0
			break
		}
	}
	st, okSt := am.states[am.current]
	if !okSt {
		return value.Nil, nil
	}
	span := st.last - st.first + 1
	if span < 1 {
		return value.Nil, nil
	}
	am.accum += dt * st.fps
	for am.accum >= 1 {
		am.accum--
		if am.frameIdx >= span-1 {
			if st.loop {
				am.frameIdx = 0
			}
		} else {
			am.frameIdx++
		}
	}
	s.syncAnimFrame()
	return value.Nil, nil
}

func (m *Module) animSetParam(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("ANIM.SETPARAM expects (sprite, name$, value)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("ANIM.SETPARAM: invalid sprite handle")
	}
	s, err := heap.Cast[*spriteObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	name, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	name = strings.TrimSpace(strings.ToLower(name))
	am := ensureAnim(s)
	switch args[2].Kind {
	case value.KindBool:
		am.boolParams[name] = args[2].IVal != 0
	default:
		if f, okf := argF(args[2]); okf {
			am.floatParams[name] = float64(f)
		} else if i, oki := args[2].ToInt(); oki {
			am.floatParams[name] = float64(i)
		} else {
			return value.Nil, fmt.Errorf("ANIM.SETPARAM: value must be numeric or bool")
		}
	}
	return value.Nil, nil
}

func evalAnimCond(expr string, fp map[string]float64, bp map[string]bool) (bool, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return false, nil
	}
	ops := []string{">=", "<=", "==", "!=", ">", "<"}
	for _, op := range ops {
		idx := strings.Index(expr, op)
		if idx < 0 {
			continue
		}
		left := strings.TrimSpace(expr[:idx])
		right := strings.TrimSpace(expr[idx+len(op):])
		if left == "" || right == "" {
			return false, nil
		}
		left = strings.ToLower(left)
		v, ok := fp[left]
		if !ok {
			if b, okb := bp[left]; okb {
				v = 0
				if b {
					v = 1
				}
			} else {
				return false, nil
			}
		}
		rhs, err := strconv.ParseFloat(right, 64)
		if err != nil {
			return false, nil
		}
		switch op {
		case ">=":
			return v >= rhs, nil
		case "<=":
			return v <= rhs, nil
		case "==":
			return v == rhs, nil
		case "!=":
			return v != rhs, nil
		case ">":
			return v > rhs, nil
		case "<":
			return v < rhs, nil
		}
	}
	// single identifier → bool param or numeric truthy
	id := strings.TrimSpace(strings.ToLower(expr))
	if b, ok := bp[id]; ok {
		return b, nil
	}
	if v, ok := fp[id]; ok {
		return v != 0, nil
	}
	// allow "true" / "false"
	if id == "true" {
		return true, nil
	}
	if id == "false" {
		return false, nil
	}
	return false, nil
}
