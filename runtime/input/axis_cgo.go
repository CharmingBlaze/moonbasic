//go:build cgo || (windows && !cgo)

package input

import (
	"fmt"
	"math"

	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func axisValueFromKeys(negVal, posVal value.Value) (float64, error) {
	kn, err := KeyCodeFromValue(negVal)
	if err != nil {
		return 0, err
	}
	kp, err := KeyCodeFromValue(posVal)
	if err != nil {
		return 0, err
	}
	neg := rl.IsKeyDown(kn)
	pos := rl.IsKeyDown(kp)
	if pos && !neg {
		return 1.0, nil
	}
	if neg && !pos {
		return -1.0, nil
	}
	return 0.0, nil
}

func (m *Module) inputAxis(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("INPUT.AXIS expects 2 arguments (negKey, posKey)")
	}
	ax, err := axisValueFromKeys(args[0], args[1])
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(ax), nil
}

func (m *Module) inputAxisDeg(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("INPUT.AXISDEG expects 4 arguments (negKey, posKey, degreesPerSec, dt)")
	}
	ax, err := axisValueFromKeys(args[0], args[1])
	if err != nil {
		return value.Nil, err
	}
	degPerSec, ok1 := args[2].ToFloat()
	dt, ok2 := args[3].ToFloat()
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("INPUT.AXISDEG: degreesPerSec and dt must be numeric")
	}
	delta := ax * (degPerSec * math.Pi / 180.0) * dt
	return value.FromFloat(delta), nil
}

func (m *Module) axisX(args []value.Value) (value.Value, error) {
	return value.FromFloat(float64(rl.GetGamepadAxisMovement(0, rl.GamepadAxisLeftX))), nil
}

func (m *Module) axisY(args []value.Value) (value.Value, error) {
	return value.FromFloat(float64(rl.GetGamepadAxisMovement(0, rl.GamepadAxisLeftY))), nil
}

func (m *Module) axisDX(args []value.Value) (value.Value, error) {
	return value.FromFloat(float64(rl.GetGamepadAxisMovement(0, rl.GamepadAxisRightX))), nil
}

func (m *Module) axisDY(args []value.Value) (value.Value, error) {
	return value.FromFloat(float64(rl.GetGamepadAxisMovement(0, rl.GamepadAxisRightY))), nil
}

func (m *Module) axisDPadX(args []value.Value) (value.Value, error) {
	v := 0.0
	if rl.IsGamepadButtonDown(0, rl.GamepadButtonLeftFaceLeft) {
		v -= 1.0
	}
	if rl.IsGamepadButtonDown(0, rl.GamepadButtonLeftFaceRight) {
		v += 1.0
	}
	return value.FromFloat(v), nil
}

func (m *Module) axisDPadY(args []value.Value) (value.Value, error) {
	v := 0.0
	if rl.IsGamepadButtonDown(0, rl.GamepadButtonLeftFaceUp) {
		v -= 1.0
	}
	if rl.IsGamepadButtonDown(0, rl.GamepadButtonLeftFaceDown) {
		v += 1.0
	}
	return value.FromFloat(v), nil
}
