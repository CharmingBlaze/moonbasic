//go:build !cgo && !windows

package input

import (
	"fmt"
	"moonbasic/vm/value"
)

const gestureHint = "GESTURE.* natives require CGO: set CGO_ENABLED=1 and install a C compiler, then rebuild"

func (m *Module) gestureEnable(args []value.Value) (value.Value, error) {
	return value.Nil, fmt.Errorf("GESTURE.ENABLE: %s", gestureHint)
}

func (m *Module) gestureIsDetected(args []value.Value) (value.Value, error) {
	return value.False, nil
}

func (m *Module) gestureGetDetected(args []value.Value) (value.Value, error) {
	return value.FromInt(0), nil
}

func (m *Module) gestureGetHoldDuration(args []value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}

func (m *Module) gestureGetDragVectorX(args []value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}

func (m *Module) gestureGetDragVectorY(args []value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}

func (m *Module) gestureGetDragAngle(args []value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}

func (m *Module) gestureGetPinchVectorX(args []value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}

func (m *Module) gestureGetPinchVectorY(args []value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}

func (m *Module) gestureGetPinchAngle(args []value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}
