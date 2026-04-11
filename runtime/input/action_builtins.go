package input

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerActionMapping(r runtime.Registrar) {
	r.Register("INPUT.MAPKEY", "input", m.inMapKey)
	r.Register("INPUT.MAPGAMEPADBUTTON", "input", m.inMapGamepadButton)
	r.Register("INPUT.MAPGAMEPADAXIS", "input", m.inMapGamepadAxis)
	r.Register("INPUT.ACTIONPRESSED", "input", m.inActionPressed)
	r.Register("INPUT.ACTIONDOWN", "input", m.inActionDown)
	r.Register("INPUT.ACTIONRELEASED", "input", m.inActionReleased)
	r.Register("INPUT.ACTIONAXIS", "input", m.inActionAxis)
	r.Register("INPUT.SAVEMAPPINGS", "input", m.inSaveMappings)
	r.Register("INPUT.LOADMAPPINGS", "input", m.inLoadMappings)
}

func numToI32(v value.Value) (int32, bool) {
	if i, ok := v.ToInt(); ok {
		return int32(i), true
	}
	if f, ok := v.ToFloat(); ok {
		return int32(f), true
	}
	return 0, false
}

func (m *Module) inMapKey(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = m
	if len(args) != 2 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("INPUT.MAPKEY expects (action, keyCode)")
	}
	act, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	kc, ok := numToI32(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("INPUT.MAPKEY: keyCode must be numeric")
	}
	sign := defaultKeyAxisSign(kc)
	appendBind(act, inputBind{Kind: bkKey, Code: kc, KeyAxisSign: sign})
	return value.Nil, nil
}

func (m *Module) inMapGamepadButton(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = m
	if len(args) != 3 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("INPUT.MAPGAMEPADBUTTON expects (action, gamepadIndex, buttonCode)")
	}
	act, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	pad, ok1 := numToI32(args[1])
	btn, ok2 := numToI32(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("INPUT.MAPGAMEPADBUTTON: gamepad and button must be numeric")
	}
	appendBind(act, inputBind{Kind: bkGamepadBtn, Pad: pad, Code: btn})
	return value.Nil, nil
}

func (m *Module) inMapGamepadAxis(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = m
	if len(args) != 3 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("INPUT.MAPGAMEPADAXIS expects (action, gamepadIndex, axisCode)")
	}
	act, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	pad, ok1 := numToI32(args[1])
	ax, ok2 := numToI32(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("INPUT.MAPGAMEPADAXIS: gamepad and axis must be numeric")
	}
	appendBind(act, inputBind{Kind: bkGamepadAxis, Pad: pad, Code: ax})
	return value.Nil, nil
}

func (m *Module) inActionPressed(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = m
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("INPUT.ACTIONPRESSED expects (action)")
	}
	act, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	return value.FromBool(actionPressedAny(act)), nil
}

func (m *Module) inActionDown(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = m
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("INPUT.ACTIONDOWN expects (action)")
	}
	act, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	return value.FromBool(actionDownAny(act)), nil
}

func (m *Module) inActionReleased(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = m
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("INPUT.ACTIONRELEASED expects (action)")
	}
	act, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	return value.FromBool(actionReleasedAny(act)), nil
}

func (m *Module) inActionAxis(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = m
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("INPUT.ACTIONAXIS expects (action)")
	}
	act, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(actionAxisSum(act)), nil
}

func (m *Module) inSaveMappings(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = m
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("INPUT.SAVEMAPPINGS expects (path)")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	if err := saveMappingsFile(path); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) inLoadMappings(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = m
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("INPUT.LOADMAPPINGS expects (path)")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	if err := loadMappingsFile(path); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}
