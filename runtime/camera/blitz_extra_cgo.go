//go:build cgo || (windows && !cgo)

package mbcamera

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerBlitzCameraExtras(m *Module, r runtime.Registrar) {
	r.Register("CameraRange", "camera", runtime.AdaptLegacy(m.camSetRange))
	r.Register("CameraZoom", "camera", runtime.AdaptLegacy(m.camZoom))
	r.Register("CameraProject", "camera", runtime.AdaptLegacy(m.camWorldToScreen))
	r.Register("CameraFogMode", "camera", m.cameraFogModeBlitz)
	r.Register("CameraFogRange", "camera", m.cameraFogRangeBlitz)
	r.Register("CameraFogColor", "camera", m.cameraFogColorBlitz)
}

// CameraFogMode(camera, mode#) — mode 0 = disable fog, non-zero = enable (linear fog via FOG.* pipeline).
func (m *Module) cameraFogModeBlitz(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CameraFogMode expects (camera, mode#)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CameraFogMode: camera handle required")
	}
	mode, ok := args[1].ToInt()
	if !ok {
		return value.Nil, fmt.Errorf("CameraFogMode: mode must be numeric")
	}
	reg := runtime.ActiveRegistry()
	if reg == nil {
		return value.Nil, fmt.Errorf("CameraFogMode: registry not active")
	}
	on := mode != 0
	_, err := reg.Call("FOG.ENABLE", []value.Value{value.FromBool(on)})
	_ = rt
	return value.Nil, err
}

func (m *Module) cameraFogRangeBlitz(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("CameraFogRange expects (camera, near#, far#)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CameraFogRange: camera handle required")
	}
	near, ok1 := argF(args[1])
	far, ok2 := argF(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("CameraFogRange: near/far must be numeric")
	}
	reg := runtime.ActiveRegistry()
	if reg == nil {
		return value.Nil, fmt.Errorf("CameraFogRange: registry not active")
	}
	_, err := reg.Call("FOG.SETRANGE", []value.Value{value.FromFloat(float64(near)), value.FromFloat(float64(far))})
	_ = rt
	return value.Nil, err
}

func (m *Module) cameraFogColorBlitz(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("CameraFogColor expects (camera, r#, g#, b#)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CameraFogColor: camera handle required")
	}
	rf, ok1 := argF(args[1])
	gf, ok2 := argF(args[2])
	bf, ok3 := argF(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("CameraFogColor: r,g,b must be numeric")
	}
	reg := runtime.ActiveRegistry()
	if reg == nil {
		return value.Nil, fmt.Errorf("CameraFogColor: registry not active")
	}
	ri := int64(rf)
	gi := int64(gf)
	bi := int64(bf)
	if rf > 1 || gf > 1 || bf > 1 {
		// assume 0..255
	} else {
		ri = int64(rf * 255)
		gi = int64(gf * 255)
		bi = int64(bf * 255)
	}
	_, err := reg.Call("FOG.SETCOLOR", []value.Value{
		value.FromInt(ri), value.FromInt(gi), value.FromInt(bi), value.FromInt(255),
	})
	_ = rt
	return value.Nil, err
}
