//go:build cgo || (windows && !cgo)

package mbdraw

import (
	"moonbasic/runtime"
	"moonbasic/vm/value"
)

// registerDrawNamespaceAliases registers DRAW.* names documented alongside DRAW3D.*
// (examples use Draw.Grid inside Camera.Begin; that compiles to DRAW.GRID).
func registerDrawNamespaceAliases(m *Module, r runtime.Registrar) {
	r.Register("DRAW.GRID", "draw", runtime.AdaptLegacy(m.drawGrid))
	r.Register("DRAW.LINE3D", "draw", runtime.AdaptLegacy(m.drawLine3D))
	r.Register("DRAW.POINT3D", "draw", runtime.AdaptLegacy(m.drawPoint3D))
	r.Register("DRAW.SPHERE", "draw", runtime.AdaptLegacy(m.drawSphere))
	r.Register("DRAW.SPHEREWIRES", "draw", runtime.AdaptLegacy(m.drawSphereWires))
	r.Register("DRAW.CUBE", "draw", runtime.AdaptLegacy(m.drawCube))
	r.Register("DRAW.CUBEWIRES", "draw", runtime.AdaptLegacy(m.drawCubeWires))
	r.Register("DRAW.CYLINDER", "draw", runtime.AdaptLegacy(m.drawCylinder))
	r.Register("DRAW.CYLINDERWIRES", "draw", runtime.AdaptLegacy(m.drawCylinderWires))
	r.Register("DRAW.CAPSULE", "draw", runtime.AdaptLegacy(m.drawCapsule))
	r.Register("DRAW.CAPSULEWIRES", "draw", runtime.AdaptLegacy(m.drawCapsuleWires))
	r.Register("DRAW.PLANE", "draw", runtime.AdaptLegacy(m.drawPlane))
	r.Register("DRAW.BOUNDINGBOX", "draw", runtime.AdaptLegacy(m.drawBBox))
	r.Register("DRAW.RAY", "draw", runtime.AdaptLegacy(m.drawRay))
	r.Register("DRAW.BILLBOARD", "draw", runtime.AdaptLegacy(m.drawBillboard))
	r.Register("DRAW.BILLBOARDREC", "draw", runtime.AdaptLegacy(m.drawBillboardRec))
	r.Register("DRAW.TEXT", "draw", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return m.text(rt, args)
	})
}
