//go:build cgo || (windows && !cgo)

package input

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// mouseFacade, keyFacade, gamepadFacade are stateless facades (no native resources).
type mouseFacade struct{ release heap.ReleaseOnce }

func (o *mouseFacade) Free() { o.release.Do(func() {}) }
func (o *mouseFacade) TypeTag() uint16  { return heap.TagInputFacade }
func (o *mouseFacade) TypeName() string { return "MOUSE" }

type keyFacade struct{ release heap.ReleaseOnce }

func (o *keyFacade) Free() { o.release.Do(func() {}) }
func (o *keyFacade) TypeTag() uint16  { return heap.TagInputFacade }
func (o *keyFacade) TypeName() string { return "KEY" }

type gamepadFacade struct{ release heap.ReleaseOnce }

func (o *gamepadFacade) Free() { o.release.Do(func() {}) }
func (o *gamepadFacade) TypeTag() uint16  { return heap.TagInputFacade }
func (o *gamepadFacade) TypeName() string { return "GAMEPAD" }

// facade implementation below

func (m *Module) makeMouse(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("MOUSE expects 0 arguments")
	}
	if m.h == nil {
		return value.Nil, fmt.Errorf("MOUSE: heap not bound")
	}
	if m.mouseH != 0 {
		return value.FromHandle(m.mouseH), nil
	}
	id, err := m.h.Alloc(&mouseFacade{})
	if err != nil {
		return value.Nil, err
	}
	m.mouseH = id
	return value.FromHandle(id), nil
}

func (m *Module) makeKey(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("KEY expects 0 arguments")
	}
	if m.h == nil {
		return value.Nil, fmt.Errorf("KEY: heap not bound")
	}
	if m.keyH != 0 {
		return value.FromHandle(m.keyH), nil
	}
	id, err := m.h.Alloc(&keyFacade{})
	if err != nil {
		return value.Nil, err
	}
	m.keyH = id
	return value.FromHandle(id), nil
}

func (m *Module) makeGamepad(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("GAMEPAD expects 0 arguments")
	}
	if m.h == nil {
		return value.Nil, fmt.Errorf("GAMEPAD: heap not bound")
	}
	if m.gamepadH != 0 {
		return value.FromHandle(m.gamepadH), nil
	}
	id, err := m.h.Alloc(&gamepadFacade{})
	if err != nil {
		return value.Nil, err
	}
	m.gamepadH = id
	return value.FromHandle(id), nil
}

func (m *Module) mouseDX(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("MOUSE.DX expects (mouseHandle)")
	}
	d := rl.GetMouseDelta()
	return value.FromFloat(float64(d.X)), nil
}

func (m *Module) mouseDY(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("MOUSE.DY expects (mouseHandle)")
	}
	d := rl.GetMouseDelta()
	return value.FromFloat(float64(d.Y)), nil
}

func (m *Module) mouseWheel(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("MOUSE.WHEEL expects (mouseHandle)")
	}
	return value.FromFloat(float64(rl.GetMouseWheelMove())), nil
}

func (m *Module) mouseDown(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("MOUSE.DOWN expects (mouseHandle, button)")
	}
	b, err := argMouseButton(args[1])
	if err != nil {
		return value.Nil, err
	}
	return value.FromBool(rl.IsMouseButtonDown(b)), nil
}

func (m *Module) mousePressed(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("MOUSE.PRESSED expects (mouseHandle, button)")
	}
	b, err := argMouseButton(args[1])
	if err != nil {
		return value.Nil, err
	}
	return value.FromBool(rl.IsMouseButtonPressed(b)), nil
}

func (m *Module) mouseReleased(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("MOUSE.RELEASED expects (mouseHandle, button)")
	}
	b, err := argMouseButton(args[1])
	if err != nil {
		return value.Nil, err
	}
	return value.FromBool(rl.IsMouseButtonReleased(b)), nil
}

func (m *Module) keyDown(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("KEY.DOWN expects (keyHandle, key)")
	}
	kc, err := KeyCodeFromValue(args[1])
	if err != nil {
		return value.Nil, err
	}
	return value.FromBool(rl.IsKeyDown(kc)), nil
}

func (m *Module) keyHit(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("KEY.HIT expects (keyHandle, key)")
	}
	kc, err := KeyCodeFromValue(args[1])
	if err != nil {
		return value.Nil, err
	}
	return value.FromBool(rl.IsKeyPressed(kc)), nil
}

func (m *Module) keyUp(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("KEY.UP expects (keyHandle, key)")
	}
	kc, err := KeyCodeFromValue(args[1])
	if err != nil {
		return value.Nil, err
	}
	return value.FromBool(rl.IsKeyReleased(kc)), nil
}

func (m *Module) gpAxis(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("GAMEPAD.AXIS expects (padHandle, gamepad, axis)")
	}
	gp, ok1 := args[1].ToInt()
	ax, ok2 := args[2].ToInt()
	if !ok1 || !ok2 || gp < 0 || ax < 0 {
		return value.Nil, fmt.Errorf("GAMEPAD.AXIS: non-negative int required")
	}
	if !rl.IsGamepadAvailable(int32(gp)) {
		return value.FromFloat(0), nil
	}
	v := rl.GetGamepadAxisMovement(int32(gp), int32(ax))
	return value.FromFloat(float64(v)), nil
}

func (m *Module) gpButton(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("GAMEPAD.BUTTON expects (padHandle, gamepad, button)")
	}
	gp, ok1 := args[1].ToInt()
	btn, ok2 := args[2].ToInt()
	if !ok1 || !ok2 || gp < 0 || btn < 0 {
		return value.Nil, fmt.Errorf("GAMEPAD.BUTTON: non-negative int required")
	}
	if !rl.IsGamepadAvailable(int32(gp)) {
		return value.FromBool(false), nil
	}
	return value.FromBool(rl.IsGamepadButtonDown(int32(gp), int32(btn))), nil
}
func (m *Module) mouseHitGlobal(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("MOUSEHIT expects (button)")
	}
	b, err := argMouseButton(args[0])
	if err != nil {
		return value.Nil, err
	}
	return value.FromInt(boolToInt(rl.IsMouseButtonPressed(b))), nil
}

func (m *Module) mouseXGlobal(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("MOUSEX expects 0 arguments")
	}
	return value.FromFloat(float64(rl.GetMousePosition().X)), nil
}

func (m *Module) mouseYGlobal(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("MOUSEY expects 0 arguments")
	}
	return value.FromFloat(float64(rl.GetMousePosition().Y)), nil
}

func boolToInt(b bool) int64 {
	if b {
		return 1
	}
	return 0
}
