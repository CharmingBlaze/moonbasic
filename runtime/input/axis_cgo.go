//go:build cgo || (windows && !cgo)

package input

import (
	"fmt"
	"math"

	"moonbasic/runtime"
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

func registerAxis(r runtime.Registrar) {
	r.Register("INPUT.AXIS", "input", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("INPUT.AXIS expects 2 arguments (negKey, posKey)")
		}
		ax, err := axisValueFromKeys(args[0], args[1])
		if err != nil {
			return value.Nil, err
		}
		return value.FromFloat(ax), nil
	}))
	axisDegLegacy := func(args []value.Value) (value.Value, error) {
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
	// Same as Input.Axis(neg,pos) * DEGPERSEC(degreesPerSec#, dt#) — radians to add e.g. to cam yaw.
	r.Register("INPUT.AXISDEG", "input", runtime.AdaptLegacy(axisDegLegacy))
	// Alias of INPUT.AXISDEG (camera orbit / yaw).
	r.Register("INPUT.ORBIT", "input", runtime.AdaptLegacy(axisDegLegacy))
}
