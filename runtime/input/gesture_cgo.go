//go:build cgo || (windows && !cgo)

package input

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/vm/value"
)

func (m *Module) gestureEnable(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("GESTURE.ENABLE expects 1 numeric flags bitmask")
	}
	flags, ok := gestureUIntArg(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("GESTURE.ENABLE: flags must be numeric")
	}
	rl.SetGesturesEnabled(flags)
	return value.Nil, nil
}

func (m *Module) gestureIsDetected(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("GESTURE.ISDETECTED expects 1 gesture id (GESTURE_* constant)")
	}
	g, ok := gestureIntArg(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("GESTURE.ISDETECTED: gesture must be numeric")
	}
	return value.FromBool(rl.IsGestureDetected(rl.Gestures(g))), nil
}

func (m *Module) gestureGetDetected(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("GESTURE.GETDETECTED expects no arguments")
	}
	return value.FromInt(int64(rl.GetGestureDetected())), nil
}

func (m *Module) gestureGetHoldDuration(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("GESTURE.GETHOLDDURATION expects no arguments")
	}
	return value.FromFloat(float64(rl.GetGestureHoldDuration())), nil
}

func (m *Module) gestureGetDragVectorX(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("GESTURE.GETDRAGVECTORX expects no arguments")
	}
	v := rl.GetGestureDragVector()
	return value.FromFloat(float64(v.X)), nil
}

func (m *Module) gestureGetDragVectorY(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("GESTURE.GETDRAGVECTORY expects no arguments")
	}
	v := rl.GetGestureDragVector()
	return value.FromFloat(float64(v.Y)), nil
}

func (m *Module) gestureGetDragAngle(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("GESTURE.GETDRAGANGLE expects no arguments")
	}
	return value.FromFloat(float64(rl.GetGestureDragAngle())), nil
}

func (m *Module) gestureGetPinchVectorX(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("GESTURE.GETPINCHVECTORX expects no arguments")
	}
	v := rl.GetGesturePinchVector()
	return value.FromFloat(float64(v.X)), nil
}

func (m *Module) gestureGetPinchVectorY(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("GESTURE.GETPINCHVECTORY expects no arguments")
	}
	v := rl.GetGesturePinchVector()
	return value.FromFloat(float64(v.Y)), nil
}

func (m *Module) gestureGetPinchAngle(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("GESTURE.GETPINCHANGLE expects no arguments")
	}
	return value.FromFloat(float64(rl.GetGesturePinchAngle())), nil
}

func gestureIntArg(v value.Value) (int32, bool) {
	if i, ok := v.ToInt(); ok {
		return int32(i), true
	}
	if f, ok := v.ToFloat(); ok {
		return int32(f), true
	}
	return 0, false
}

func gestureUIntArg(v value.Value) (uint32, bool) {
	if i, ok := v.ToInt(); ok {
		return uint32(i), true
	}
	if f, ok := v.ToFloat(); ok {
		return uint32(f), true
	}
	return 0, false
}
