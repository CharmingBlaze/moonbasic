//go:build !cgo && !windows

package input

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerAxis(r runtime.Registrar) {
	r.Register("INPUT.AXIS", "input", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("INPUT.AXIS expects 2 arguments (negKey, posKey)")
		}
		return value.FromFloat(0.0), nil
	}))
	axisDegStub := func(args []value.Value) (value.Value, error) {
		if len(args) != 4 {
			return value.Nil, fmt.Errorf("INPUT.AXISDEG expects 4 arguments (negKey, posKey, degreesPerSec, dt)")
		}
		return value.FromFloat(0.0), nil
	}
	r.Register("INPUT.AXISDEG", "input", runtime.AdaptLegacy(axisDegStub))
	r.Register("INPUT.ORBIT", "input", runtime.AdaptLegacy(axisDegStub))
}
