//go:build cgo || (windows && !cgo)

package input

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	mbcamera "moonbasic/runtime/camera"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) registerInputAdvanced(r runtime.Registrar) {
	r.Register("INPUT.MOUSEX", "input", runtime.AdaptLegacy(m.inMouseX))
	r.Register("INPUT.MOUSEY", "input", runtime.AdaptLegacy(m.inMouseY))
	r.Register("INPUT.MOUSEDOWN", "input", runtime.AdaptLegacy(m.inMouseDown))
	r.Register("INPUT.MOUSEHIT", "input", runtime.AdaptLegacy(m.inMouseHit))
	r.Register("INPUT.SETMOUSESCALE", "input", runtime.AdaptLegacy(m.inSetMouseScale))
	r.Register("INPUT.SETMOUSEOFFSET", "input", runtime.AdaptLegacy(m.inSetMouseOffset))
	r.Register("INPUT.GETMOUSEWORLDPOS", "input", runtime.AdaptLegacy(m.inGetMouseWorldPos))
	r.Register("INPUT.TOUCHCOUNT", "input", runtime.AdaptLegacy(m.inTouchCount))
	r.Register("INPUT.TOUCHX", "input", runtime.AdaptLegacy(m.inTouchX))
	r.Register("INPUT.TOUCHY", "input", runtime.AdaptLegacy(m.inTouchY))
	r.Register("INPUT.TOUCHPRESSED", "input", runtime.AdaptLegacy(m.inTouchPressed))
	r.Register("INPUT.GETTOUCHPOINTID", "input", runtime.AdaptLegacy(m.inGetTouchPointID))
	r.Register("INPUT.GAMEPADBUTTONCOUNT", "input", runtime.AdaptLegacy(m.inGamepadButtonCount))
	r.Register("INPUT.GAMEPADAXISCOUNT", "input", runtime.AdaptLegacy(m.inGamepadAxisCount))
	r.Register("INPUT.SETGAMEPADMAPPINGS", "input", m.inSetGamepadMappings)
}

func argMouseButton(v value.Value) (rl.MouseButton, error) {
	i, ok := v.ToInt()
	if !ok {
		if f, okf := v.ToFloat(); okf {
			i = int64(f)
		} else {
			return 0, fmt.Errorf("expected numeric mouse button id")
		}
	}
	return rl.MouseButton(i), nil
}

func (m *Module) inMouseX(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("INPUT.MOUSEX expects 0 arguments")
	}
	return value.FromInt(int64(rl.GetMouseX())), nil
}

func (m *Module) inMouseY(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("INPUT.MOUSEY expects 0 arguments")
	}
	return value.FromInt(int64(rl.GetMouseY())), nil
}

func (m *Module) inMouseDown(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("INPUT.MOUSEDOWN expects 1 argument (button int)")
	}
	b, err := argMouseButton(args[0])
	if err != nil {
		return value.Nil, fmt.Errorf("INPUT.MOUSEDOWN: %w", err)
	}
	return value.FromBool(rl.IsMouseButtonDown(b)), nil
}

func (m *Module) inMouseHit(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("INPUT.MOUSEHIT expects 1 argument (button int)")
	}
	b, err := argMouseButton(args[0])
	if err != nil {
		return value.Nil, fmt.Errorf("INPUT.MOUSEHIT: %w", err)
	}
	return value.FromBool(rl.IsMouseButtonPressed(b)), nil
}

func (m *Module) inSetMouseScale(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("INPUT.SETMOUSESCALE expects 2 arguments (sx, sy)")
	}
	sx, ok1 := args[0].ToFloat()
	if !ok1 {
		if i, ok := args[0].ToInt(); ok {
			sx = float64(i)
		} else {
			return value.Nil, fmt.Errorf("INPUT.SETMOUSESCALE: sx must be numeric")
		}
	}
	sy, ok2 := args[1].ToFloat()
	if !ok2 {
		if i, ok := args[1].ToInt(); ok {
			sy = float64(i)
		} else {
			return value.Nil, fmt.Errorf("INPUT.SETMOUSESCALE: sy must be numeric")
		}
	}
	rl.SetMouseScale(float32(sx), float32(sy))
	return value.Nil, nil
}

func (m *Module) inSetMouseOffset(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("INPUT.SETMOUSEOFFSET expects 2 arguments (x, y)")
	}
	x, ok1 := args[0].ToInt()
	if !ok1 {
		if f, ok := args[0].ToFloat(); ok {
			x = int64(f)
		} else {
			return value.Nil, fmt.Errorf("INPUT.SETMOUSEOFFSET: x must be numeric")
		}
	}
	y, ok2 := args[1].ToInt()
	if !ok2 {
		if f, ok := args[1].ToFloat(); ok {
			y = int64(f)
		} else {
			return value.Nil, fmt.Errorf("INPUT.SETMOUSEOFFSET: y must be numeric")
		}
	}
	setMouseOffsetCompat(int(x), int(y))
	return value.Nil, nil
}

// inGetMouseWorldPos returns a 1D array [x,y,z] where the screen ray hits the horizontal plane y=0.
func (m *Module) inGetMouseWorldPos(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("INPUT.GETMOUSEWORLDPOS: heap not bound")
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("INPUT.GETMOUSEWORLDPOS expects 3 arguments (cam, sx, sy)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("INPUT.GETMOUSEWORLDPOS: cam must be a handle")
	}
	cam, err := mbcamera.RayCamera3D(m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	sx, ok1 := args[1].ToFloat()
	if !ok1 {
		if i, ok := args[1].ToInt(); ok {
			sx = float64(i)
		} else {
			return value.Nil, fmt.Errorf("INPUT.GETMOUSEWORLDPOS: sx must be numeric")
		}
	}
	sy, ok2 := args[2].ToFloat()
	if !ok2 {
		if i, ok := args[2].ToInt(); ok {
			sy = float64(i)
		} else {
			return value.Nil, fmt.Errorf("INPUT.GETMOUSEWORLDPOS: sy must be numeric")
		}
	}
	ray := rl.GetScreenToWorldRay(rl.Vector2{X: float32(sx), Y: float32(sy)}, cam)
	p := intersectRayHorizontalPlane(ray, 0)
	arr, err := heap.NewArray([]int64{3})
	if err != nil {
		return value.Nil, err
	}
	_ = arr.Set([]int64{0}, float64(p.X))
	_ = arr.Set([]int64{1}, float64(p.Y))
	_ = arr.Set([]int64{2}, float64(p.Z))
	id, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func intersectRayHorizontalPlane(ray rl.Ray, planeY float32) rl.Vector3 {
	dy := ray.Direction.Y
	const eps = 1e-5
	if math.Abs(float64(dy)) < eps {
		return ray.Position
	}
	t := (planeY - ray.Position.Y) / dy
	if t < 0 {
		t = 0
	}
	return rl.Vector3{
		X: ray.Position.X + ray.Direction.X*t,
		Y: planeY,
		Z: ray.Position.Z + ray.Direction.Z*t,
	}
}

func (m *Module) inTouchCount(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("INPUT.TOUCHCOUNT expects 0 arguments")
	}
	return value.FromInt(int64(rl.GetTouchPointCount())), nil
}

func (m *Module) inTouchX(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("INPUT.TOUCHX expects 1 argument (idx)")
	}
	idx, ok := args[0].ToInt()
	if !ok {
		if f, okf := args[0].ToFloat(); okf {
			idx = int64(f)
		} else {
			return value.Nil, fmt.Errorf("INPUT.TOUCHX: idx must be numeric")
		}
	}
	p := rl.GetTouchPosition(int32(idx))
	return value.FromInt(int64(p.X)), nil
}

func (m *Module) inTouchY(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("INPUT.TOUCHY expects 1 argument (idx)")
	}
	idx, ok := args[0].ToInt()
	if !ok {
		if f, okf := args[0].ToFloat(); okf {
			idx = int64(f)
		} else {
			return value.Nil, fmt.Errorf("INPUT.TOUCHY: idx must be numeric")
		}
	}
	p := rl.GetTouchPosition(int32(idx))
	return value.FromInt(int64(p.Y)), nil
}

func (m *Module) inTouchPressed(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("INPUT.TOUCHPRESSED expects 1 argument (idx)")
	}
	idx, ok := args[0].ToInt()
	if !ok {
		if f, okf := args[0].ToFloat(); okf {
			idx = int64(f)
		} else {
			return value.Nil, fmt.Errorf("INPUT.TOUCHPRESSED: idx must be numeric")
		}
	}
	if idx == 0 {
		return value.FromBool(rl.IsMouseButtonPressed(rl.MouseLeftButton)), nil
	}
	return value.FromBool(false), nil
}

func (m *Module) inGetTouchPointID(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("INPUT.GETTOUCHPOINTID expects 1 argument (idx)")
	}
	idx, ok := args[0].ToInt()
	if !ok {
		if f, okf := args[0].ToFloat(); okf {
			idx = int64(f)
		} else {
			return value.Nil, fmt.Errorf("INPUT.GETTOUCHPOINTID: idx must be numeric")
		}
	}
	return value.FromInt(int64(rl.GetTouchPointId(int32(idx)))), nil
}

func (m *Module) inGamepadButtonCount(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("INPUT.GAMEPADBUTTONCOUNT expects 1 argument (pad)")
	}
	pad, ok := args[0].ToInt()
	if !ok {
		if f, okf := args[0].ToFloat(); okf {
			pad = int64(f)
		} else {
			return value.Nil, fmt.Errorf("INPUT.GAMEPADBUTTONCOUNT: pad must be numeric")
		}
	}
	if !rl.IsGamepadAvailable(int32(pad)) {
		return value.FromInt(0), nil
	}
	return value.FromInt(int64(rl.GamepadButtonRightThumb + 1)), nil
}

func (m *Module) inGamepadAxisCount(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("INPUT.GAMEPADAXISCOUNT expects 1 argument (pad)")
	}
	pad, ok := args[0].ToInt()
	if !ok {
		if f, okf := args[0].ToFloat(); okf {
			pad = int64(f)
		} else {
			return value.Nil, fmt.Errorf("INPUT.GAMEPADAXISCOUNT: pad must be numeric")
		}
	}
	return value.FromInt(int64(rl.GetGamepadAxisCount(int32(pad)))), nil
}

func (m *Module) inSetGamepadMappings(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("INPUT.SETGAMEPADMAPPINGS expects 1 string argument (mappings)")
	}
	mappings, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	n := rl.SetGamepadMappings(mappings)
	return value.FromInt(int64(n)), nil
}
