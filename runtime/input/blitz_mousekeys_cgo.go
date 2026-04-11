//go:build cgo || (windows && !cgo)

package input

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	mbruntime "moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerBlitzMouseKeys(r mbruntime.Registrar) {
	r.Register("MOUSEDOWN", "input", mbruntime.AdaptLegacy(m.inMouseDown))
	r.Register("MOUSEHIT", "input", mbruntime.AdaptLegacy(m.inMouseHit))
	r.Register("MOUSEXSPEED", "input", mbruntime.AdaptLegacy(m.inMouseXSpeed))
	r.Register("MOUSEYSPEED", "input", mbruntime.AdaptLegacy(m.inMouseYSpeed))
	r.Register("FlushMouse", "input", mbruntime.AdaptLegacy(m.inFlushMouse))
	r.Register("WaitMouse", "input", mbruntime.AdaptLegacy(m.inWaitMouse))
	r.Register("MoveMouse", "input", mbruntime.AdaptLegacy(m.inMoveMouse))
	r.Register("HidePointer", "input", mbruntime.AdaptLegacy(m.inHidePointer))
	r.Register("ShowPointer", "input", mbruntime.AdaptLegacy(m.inShowPointer))
	r.Register("GetKey", "input", mbruntime.AdaptLegacy(m.inGetKey))
	r.Register("WaitKey", "input", mbruntime.AdaptLegacy(m.inWaitKey))
	r.Register("FlushKeys", "input", mbruntime.AdaptLegacy(m.inFlushKeys))
}

func (m *Module) inFlushMouse(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("FlushMouse expects 0 arguments")
	}
	return value.Nil, nil
}

func (m *Module) inFlushKeys(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("FlushKeys expects 0 arguments")
	}
	return value.Nil, nil
}

func (m *Module) inWaitMouse(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("WaitMouse expects 0 arguments")
	}
	for !rl.WindowShouldClose() {
		for b := rl.MouseButtonLeft; b <= rl.MouseButtonBack; b++ {
			if rl.IsMouseButtonPressed(b) {
				return value.FromInt(int64(b)), nil
			}
		}
		rl.WaitTime(1.0 / 240.0)
	}
	return value.FromInt(0), nil
}

func (m *Module) inWaitKey(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("WaitKey expects 0 arguments")
	}
	for !rl.WindowShouldClose() {
		if k, ok := pollAnyKeyPressed(); ok {
			return value.FromInt(int64(k)), nil
		}
		rl.WaitTime(1.0 / 240.0)
	}
	return value.FromInt(0), nil
}

func (m *Module) inGetKey(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("GetKey expects 0 arguments")
	}
	if k, ok := pollAnyKeyPressed(); ok {
		return value.FromInt(int64(k)), nil
	}
	return value.FromInt(0), nil
}

func pollAnyKeyPressed() (int32, bool) {
	for k := int32(32); k <= 400; k++ {
		if rl.IsKeyPressed(k) {
			return k, true
		}
	}
	return 0, false
}

func (m *Module) inMoveMouse(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("MoveMouse expects (x, y)")
	}
	x, ok1 := args[0].ToInt()
	y, ok2 := args[1].ToInt()
	if !ok1 {
		if xf, ok := args[0].ToFloat(); ok {
			x = int64(xf)
			ok1 = true
		}
	}
	if !ok2 {
		if yf, ok := args[1].ToFloat(); ok {
			y = int64(yf)
			ok2 = true
		}
	}
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("MoveMouse: x,y must be numeric")
	}
	setMousePositionCompat(int(x), int(y))
	return value.Nil, nil
}

func (m *Module) inHidePointer(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("HidePointer expects 0 arguments")
	}
	rl.HideCursor()
	return value.Nil, nil
}

func (m *Module) inShowPointer(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("ShowPointer expects 0 arguments")
	}
	rl.ShowCursor()
	return value.Nil, nil
}
