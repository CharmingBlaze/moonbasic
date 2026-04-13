//go:build cgo

// UI.BUTTON / UI.PROGRESSBAR / UI.LABEL3D use raygui; raygui-go has no Windows purego build, so these
// register only when CGO_ENABLED=1 (see purego_register_windows.go for the minimal GUI path).

package mbgui

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/runtime/mbentity"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
	raygui "github.com/gen2brain/raylib-go/raygui"
)

func registerGuiQoLAPI(m *Module, r runtime.Registrar) {
	r.Register("UI.BUTTON", "gui", runtime.AdaptLegacy(m.uiButton))
	r.Register("UI.PROGRESSBAR", "gui", runtime.AdaptLegacy(m.uiProgressBar))
	r.Register("UI.INVENTORYICON", "gui", runtime.AdaptLegacy(m.uiInventoryIcon))
	r.Register("UI.LABEL3D", "gui", m.uiLabel3D)
}

func (m *Module) uiButton(args []value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("UI.BUTTON expects (label$, x, y, w, h)")
	}
	label := args[0].String()
	x, _ := args[1].ToFloat()
	y, _ := args[2].ToFloat()
	w, _ := args[3].ToFloat()
	h, _ := args[4].ToFloat()

	bounds := rl.Rectangle{X: float32(x), Y: float32(y), Width: float32(w), Height: float32(h)}
	clicked := raygui.Button(bounds, label)
	
	if clicked { return value.FromInt(1), nil }
	return value.FromInt(0), nil
}

func (m *Module) uiProgressBar(args []value.Value) (value.Value, error) {
	if len(args) != 6 {
		return value.Nil, fmt.Errorf("UI.PROGRESSBAR expects (x, y, w, h, percent, colorInt)")
	}
	x, _ := args[0].ToFloat()
	y, _ := args[1].ToFloat()
	w, _ := args[2].ToFloat()
	h, _ := args[3].ToFloat()
	val, _ := args[4].ToFloat()

	bounds := rl.Rectangle{X: float32(x), Y: float32(y), Width: float32(w), Height: float32(h)}
	raygui.ProgressBar(bounds, "", "", float32(val), 0, 100)
	return value.Nil, nil
}

func (m *Module) uiInventoryIcon(args []value.Value) (value.Value, error) {
	// Simple stub demonstrating projecting a 3D overlay or 2D slice via texture
	return value.Nil, nil
}

func (m *Module) uiLabel3D(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("UI.LABEL3D expects (text$, entity#, offsetVec#)")
	}
	txt, _ := rt.ArgString(args, 0)
	eid, _ := args[1].ToInt()
	
	if m.cam == nil { return value.Nil, nil }
	cam, ok := m.cam.ActiveCamera()
	if !ok { return value.Nil, nil }
	
	entMod := mbentity.ModulesByStore[m.h]
	if entMod == nil { return value.Nil, nil }
	
	wp, ok := entMod.GetWorldPosByID(int(eid))
	if !ok { return value.Nil, nil }
	
	// Add offset (simplified: assuming handle logic or fixed float)
	// For now, just project the entity base position.
	pos2d := rl.GetWorldToScreenEx(wp, cam, int32(rl.GetRenderWidth()), int32(rl.GetRenderHeight()))
	
	rl.DrawText(txt, int32(pos2d.X), int32(pos2d.Y), 20, rl.RayWhite)
	return value.Nil, nil
}
