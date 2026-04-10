//go:build cgo || (windows && !cgo)

package mbentity

import (
	texmod "moonbasic/runtime/texture"

	"moonbasic/runtime/mbmodel3d"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// drawSpriteBillboard draws a billboard entity using either a runtime texture atlas (TEXTURE.*)
// or a legacy internal textureObj (LOADSPRITE path).
func (m *Module) drawSpriteBillboard(e *ent) {
	if e == nil || e.texHandle == 0 {
		return
	}
	obj, ok := m.h.Get(e.texHandle)
	if !ok {
		return
	}
	cam, okCam := mbmodel3d.ActiveCamera3D()
	if !okCam {
		return
	}
	wp := m.worldPos(e)
	col := m.entTintResolved(e)
	size := rl.Vector2{X: e.w * e.scale.X, Y: e.h * e.scale.Y}

	_, _, r := e.getRot()
	origin := rl.Vector2{X: size.X * 0.5, Y: size.Y * 0.5}

	switch t := obj.(type) {
	case *texmod.TextureObject:
		src := t.FrameSourceRect()
		rl.DrawBillboardPro(cam, t.Tex, src, wp, cam.Up, size, origin, r, col)
	case *textureObj:
		src := rl.Rectangle{X: 0, Y: 0, Width: float32(t.tex.Width), Height: float32(t.tex.Height)}
		rl.DrawBillboardPro(cam, t.tex, src, wp, cam.Up, size, origin, r, col)
	default:
		return
	}
}
