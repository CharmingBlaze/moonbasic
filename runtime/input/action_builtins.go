package input

import (
	"fmt"

	"moonbasic/vm/value"
)

// action mapping implementation below

func numToI32(v value.Value) (int32, bool) {
	if i, ok := v.ToInt(); ok {
		return int32(i), true
	}
	if f, ok := v.ToFloat(); ok {
		return int32(f), true
	}
	return 0, false
}

func (m *Module) inMapKey(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("INPUT.MAPKEY expects (action, keyCode)")
	}
	act := args[0].String()
	kc, ok := numToI32(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("INPUT.MAPKEY: keyCode must be numeric")
	}
	sign := defaultKeyAxisSign(kc)
	appendBind(act, inputBind{Kind: bkKey, Code: kc, KeyAxisSign: sign})
	return value.Nil, nil
}

func (m *Module) inMapGamepadButton(args []value.Value) (value.Value, error) {
	if len(args) != 3 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("INPUT.MAPGAMEPADBUTTON expects (action, gamepadIndex, buttonCode)")
	}
	act := args[0].String()
	pad, ok1 := numToI32(args[1])
	btn, ok2 := numToI32(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("INPUT.MAPGAMEPADBUTTON: gamepad and button must be numeric")
	}
	appendBind(act, inputBind{Kind: bkGamepadBtn, Pad: pad, Code: btn})
	return value.Nil, nil
}

func (m *Module) inMapGamepadAxis(args []value.Value) (value.Value, error) {
	if len(args) != 3 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("INPUT.MAPGAMEPADAXIS expects (action, gamepadIndex, axisCode)")
	}
	act := args[0].String()
	pad, ok1 := numToI32(args[1])
	ax, ok2 := numToI32(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("INPUT.MAPGAMEPADAXIS: gamepad and axis must be numeric")
	}
	appendBind(act, inputBind{Kind: bkGamepadAxis, Pad: pad, Code: ax})
	return value.Nil, nil
}

func (m *Module) inActionPressed(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("INPUT.ACTIONPRESSED expects (action)")
	}
	act := args[0].String()
	return value.FromBool(actionPressedAny(act)), nil
}

func (m *Module) inActionDown(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("INPUT.ACTIONDOWN expects (action)")
	}
	act := args[0].String()
	return value.FromBool(actionDownAny(act)), nil
}

func (m *Module) inActionReleased(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("INPUT.ACTIONRELEASED expects (action)")
	}
	act := args[0].String()
	return value.FromBool(actionReleasedAny(act)), nil
}

func (m *Module) inActionAxis(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("INPUT.ACTIONAXIS expects (action)")
	}
	act := args[0].String()
	return value.FromFloat(actionAxisSum(act)), nil
}

func (m *Module) inSaveMappings(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("INPUT.SAVEMAPPINGS expects (path)")
	}
	path := args[0].String()
	if err := saveMappingsFile(path); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) inLoadMappings(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("INPUT.LOADMAPPINGS expects (path)")
	}
	path := args[0].String()
	if err := loadMappingsFile(path); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) inMapMouse(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("ACTION.MAPMOUSE expects (action, buttonCode)")
	}
	act := args[0].String()
	btn, ok := numToI32(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("ACTION.MAPMOUSE: buttonCode must be numeric")
	}
	appendBind(act, inputBind{Kind: bkMouseBtn, Code: btn})
	return value.Nil, nil
}

func (m *Module) inMapJoy(args []value.Value) (value.Value, error) {
	// Alias for gamepad button
	return m.inMapGamepadButton(args)
}

func (m *Module) inActionReset(args []value.Value) (value.Value, error) {
	clearAllBindings()
	return value.Nil, nil
}
