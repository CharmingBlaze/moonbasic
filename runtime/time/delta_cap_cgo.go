//go:build cgo || (windows && !cgo)

package mbtime

import (
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

var (
	deltaMu  sync.Mutex
	noClamp  bool
	maxDelta = 0.05 // seconds; used when noClamp is false
)

func applyDeltaCap(dt float64) float64 {
	deltaMu.Lock()
	defer deltaMu.Unlock()
	if noClamp {
		return dt
	}
	if dt > maxDelta {
		return maxDelta
	}
	return dt
}

// DeltaSeconds returns Raylib frame time with optional cap (see TIME.SETMAXDELTA). Used by GAME.DT.
func DeltaSeconds(rt *runtime.Runtime) float64 {
	if rt != nil && rt.GamePaused {
		return 0
	}
	dt := applyDeltaCap(float64(rl.GetFrameTime()))
	if rt != nil {
		s := rt.TimeScale
		if s != 0 && s != 1 {
			dt *= s
		}
	}
	return dt
}

func registerDeltaCapCommands(reg runtime.Registrar) {
	reg.Register("TIME.SETMAXDELTA", "time", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if len(args) != 1 {
			return value.Nil, errArgs(1, len(args))
		}
		var sec float64
		if f, ok := args[0].ToFloat(); ok {
			sec = f
		} else if i, ok := args[0].ToInt(); ok {
			sec = float64(i)
		} else {
			return value.Nil, runtime.Errorf("TIME.SETMAXDELTA expects 1 numeric argument (seconds; 0 = no cap)")
		}
		deltaMu.Lock()
		defer deltaMu.Unlock()
		if sec <= 0 {
			noClamp = true
		} else {
			noClamp = false
			maxDelta = sec
		}
		return value.Nil, nil
	})
}
