//go:build !cgo && !windows

package mbdraw

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func stub(name string) runtime.BuiltinFn {
	return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return value.Nil, fmt.Errorf("%s requires CGO_ENABLED=1", name)
	}
}

// Register implements runtime.Module.
func (m *Module) Register(r runtime.Registrar) {
	// This is a partial list. Add all commands from docs/audit/archives/IMPL_AUDIT.txt here.
	r.Register("DRAW.RECTANGLE", "draw", stub("DRAW.RECTANGLE"))
	r.Register("DRAW.RECTANGLE_ROUNDED", "draw", stub("DRAW.RECTANGLE_ROUNDED"))
	r.Register("DRAW.TEXTURE", "draw", stub("DRAW.TEXTURE"))
	r.Register("DRAW.TEXTURENPATCH", "draw", stub("DRAW.TEXTURENPATCH"))
	r.Register("TEXTURE.LOAD", "draw", stub("TEXTURE.LOAD"))
	r.Register("TEXTURE.FROMIMAGE", "draw", stub("TEXTURE.FROMIMAGE"))
	r.Register("TEXTURE.FREE", "draw", stub("TEXTURE.FREE"))
	r.Register("DRAW.CIRCLE", "draw", stub("DRAW.CIRCLE"))
	r.Register("DRAW.CIRCLELINES", "draw", stub("DRAW.CIRCLELINES"))
	r.Register("DRAW.ELLIPSE", "draw", stub("DRAW.ELLIPSE"))
	r.Register("DRAW.ELLIPSELINES", "draw", stub("DRAW.ELLIPSELINES"))
	r.Register("DRAW.OVAL", "draw", stub("DRAW.OVAL"))
	r.Register("DRAW.OVALLINES", "draw", stub("DRAW.OVALLINES"))
	r.Register("DRAW.RING", "draw", stub("DRAW.RING"))
	r.Register("DRAW.RINGLINES", "draw", stub("DRAW.RINGLINES"))
	r.Register("DRAW.TRIANGLE", "draw", stub("DRAW.TRIANGLE"))
	r.Register("DRAW.TRIANGLELINES", "draw", stub("DRAW.TRIANGLELINES"))
	r.Register("DRAW.POLY", "draw", stub("DRAW.POLY"))
	r.Register("DRAW.POLYLINES", "draw", stub("DRAW.POLYLINES"))
	r.Register("DRAW.RECTLINES", "draw", stub("DRAW.RECTLINES"))
	r.Register("DRAW.RECTPRO", "draw", stub("DRAW.RECTPRO"))
	r.Register("DRAW.RECTGRADV", "draw", stub("DRAW.RECTGRADV"))
	r.Register("DRAW.RECTGRADH", "draw", stub("DRAW.RECTGRADH"))
	r.Register("DRAW.RECTGRAD", "draw", stub("DRAW.RECTGRAD"))
	r.Register("DRAW.LINE", "draw", stub("DRAW.LINE"))
	r.Register("DRAW.LINEEX", "draw", stub("DRAW.LINEEX"))
	r.Register("DRAW.LINEBEZIER", "draw", stub("DRAW.LINEBEZIER"))
	r.Register("DRAW.LINEBEZIERQUAD", "draw", stub("DRAW.LINEBEZIERQUAD"))
	r.Register("DRAW.LINEBEZIERCUBIC", "draw", stub("DRAW.LINEBEZIERCUBIC"))
	r.Register("DRAW.SPLINELINEAR", "draw", stub("DRAW.SPLINELINEAR"))
	r.Register("DRAW.SPLINEBASIS", "draw", stub("DRAW.SPLINEBASIS"))
	r.Register("DRAW.SPLINECATMULLROM", "draw", stub("DRAW.SPLINECATMULLROM"))
	r.Register("DRAW.SPLINEBEZIERQUAD", "draw", stub("DRAW.SPLINEBEZIERQUAD"))
	r.Register("DRAW.SPLINEBEZIERCUBIC", "draw", stub("DRAW.SPLINEBEZIERCUBIC"))
	r.Register("DRAW.TEXT", "draw", stub("DRAW.TEXT"))
	r.Register("RENDER.DRAWFPS", "render", stub("RENDER.DRAWFPS"))
	r.Register("DRAW.TEXTEX", "draw", stub("DRAW.TEXTEX"))
	r.Register("DRAW.TEXTFONT", "draw", stub("DRAW.TEXTFONT"))
	r.Register("DRAW.TEXTPRO", "draw", stub("DRAW.TEXTPRO"))
	r.Register("DRAW.TEXTWIDTH", "draw", stub("DRAW.TEXTWIDTH"))
	r.Register("DRAW.TEXTFONTWIDTH", "draw", stub("DRAW.TEXTFONTWIDTH"))
	r.Register("MEASURETEXT", "draw", stub("MEASURETEXT"))
	r.Register("MEASURETEXTEX", "draw", stub("MEASURETEXTEX"))
	r.Register("GETTEXTCODEPOINTCOUNT", "draw", stub("GETTEXTCODEPOINTCOUNT"))
	r.Register("DRAW.TEXTUREV", "draw", stub("DRAW.TEXTUREV"))
	r.Register("DRAW.TEXTUREEX", "draw", stub("DRAW.TEXTUREEX"))
	r.Register("DRAW.TEXTUREREC", "draw", stub("DRAW.TEXTUREREC"))
	r.Register("DRAW.TEXTUREPRO", "draw", stub("DRAW.TEXTUREPRO"))
	r.Register("DRAW.TEXTUREFULL", "draw", stub("DRAW.TEXTUREFULL"))
	r.Register("DRAW.TEXTUREFLIPPED", "draw", stub("DRAW.TEXTUREFLIPPED"))
	r.Register("DRAW.TEXTURETILED", "draw", stub("DRAW.TEXTURETILED"))
	r.Register("DRAW.ARC", "draw", stub("DRAW.ARC"))
	r.Register("DRAW.DOT", "draw", stub("DRAW.DOT"))
	r.Register("DRAW.PIXEL", "draw", stub("DRAW.PIXEL"))
	r.Register("DRAW.PLOT", "draw", stub("DRAW.PLOT"))
	r.Register("DRAW.PIXELV", "draw", stub("DRAW.PIXELV"))
	r.Register("DRAW.SETPIXELCOLOR", "draw", stub("DRAW.SETPIXELCOLOR"))
	r.Register("DRAW.GRID2D", "draw", stub("DRAW.GRID2D"))
	r.Register("DRAW.GETPIXELCOLOR", "draw", stub("DRAW.GETPIXELCOLOR"))
	r.Register("DRAW3D.GRID", "draw", stub("DRAW3D.GRID"))
	r.Register("DRAW3D.LINE", "draw", stub("DRAW3D.LINE"))
	r.Register("DRAW3D.POINT", "draw", stub("DRAW3D.POINT"))
	r.Register("DRAW3D.SPHERE", "draw", stub("DRAW3D.SPHERE"))
	r.Register("DRAW3D.SPHEREWIRES", "draw", stub("DRAW3D.SPHEREWIRES"))
	r.Register("DRAW3D.CUBE", "draw", stub("DRAW3D.CUBE"))
	r.Register("DRAW3D.CUBEWIRES", "draw", stub("DRAW3D.CUBEWIRES"))
	r.Register("DRAW3D.CYLINDER", "draw", stub("DRAW3D.CYLINDER"))
	r.Register("DRAW3D.CYLINDERWIRES", "draw", stub("DRAW3D.CYLINDERWIRES"))
	r.Register("DRAW3D.CAPSULE", "draw", stub("DRAW3D.CAPSULE"))
	r.Register("DRAW3D.CAPSULEWIRES", "draw", stub("DRAW3D.CAPSULEWIRES"))
	r.Register("DRAW3D.PLANE", "draw", stub("DRAW3D.PLANE"))
	r.Register("DRAW3D.BBOX", "draw", stub("DRAW3D.BBOX"))
	r.Register("DRAW3D.RAY", "draw", stub("DRAW3D.RAY"))
	r.Register("DRAW3D.BILLBOARD", "draw", stub("DRAW3D.BILLBOARD"))
	r.Register("DRAW3D.BILLBOARDREC", "draw", stub("DRAW3D.BILLBOARDREC"))
	r.Register("BOX", "draw", stub("BOX"))
	r.Register("BOXW", "draw", stub("BOXW"))
	r.Register("WIRECUBE", "draw", stub("WIRECUBE"))
	r.Register("BALL", "draw", stub("BALL"))
	r.Register("BALLW", "draw", stub("BALLW"))
	r.Register("GRID3", "draw", stub("GRID3"))
	r.Register("FLAT", "draw", stub("FLAT"))
	r.Register("CAP", "draw", stub("CAP"))
	r.Register("CAPW", "draw", stub("CAPW"))
	r.Register("DRAW.GRID", "draw", stub("DRAW.GRID"))
	r.Register("DRAW.LINE3D", "draw", stub("DRAW.LINE3D"))
	r.Register("DRAW.POINT3D", "draw", stub("DRAW.POINT3D"))
	r.Register("DRAW.SPHERE", "draw", stub("DRAW.SPHERE"))
	r.Register("DRAW.SPHEREWIRES", "draw", stub("DRAW.SPHEREWIRES"))
	r.Register("DRAW.CUBE", "draw", stub("DRAW.CUBE"))
	r.Register("DRAW.CUBEWIRES", "draw", stub("DRAW.CUBEWIRES"))
	r.Register("DRAW.CYLINDER", "draw", stub("DRAW.CYLINDER"))
	r.Register("DRAW.CYLINDERWIRES", "draw", stub("DRAW.CYLINDERWIRES"))
	r.Register("DRAW.CAPSULE", "draw", stub("DRAW.CAPSULE"))
	r.Register("DRAW.CAPSULEWIRES", "draw", stub("DRAW.CAPSULEWIRES"))
	r.Register("DRAW.PLANE", "draw", stub("DRAW.PLANE"))
	r.Register("DRAW.BOUNDINGBOX", "draw", stub("DRAW.BOUNDINGBOX"))
	r.Register("DRAW.RAY", "draw", stub("DRAW.RAY"))
	r.Register("DRAW.BILLBOARD", "draw", stub("DRAW.BILLBOARD"))
	r.Register("DRAW.BILLBOARDREC", "draw", stub("DRAW.BILLBOARDREC"))
	r.Register("DRAW.CIRCLESECTOR", "draw", stub("DRAW.CIRCLESECTOR"))
	r.Register("DRAW.CIRCLEGRADIENT", "draw", stub("DRAW.CIRCLEGRADIENT"))
	r.Register("DRAW.PROGRESSBAR", "draw", stub("DRAW.PROGRESSBAR"))
	r.Register("DRAW.HEALTHBAR", "draw", stub("DRAW.HEALTHBAR"))
	r.Register("DRAW.CENTERTEXT", "draw", stub("DRAW.CENTERTEXT"))
	r.Register("DRAW.RIGHTTEXT", "draw", stub("DRAW.RIGHTTEXT"))
	r.Register("DRAW.SHADOWTEXT", "draw", stub("DRAW.SHADOWTEXT"))
	r.Register("DRAW.OUTLINETEXT", "draw", stub("DRAW.OUTLINETEXT"))
	r.Register("DRAW.CROSSHAIR", "draw", stub("DRAW.CROSSHAIR"))
	r.Register("DRAW.RECTGRID", "draw", stub("DRAW.RECTGRID"))

	// Object-style draw wrappers (CGO implementation in prim*_wrapper_cgo.go, obj_text_texture_cgo.go).
	for _, name := range []string{
		"DRAWPRIM3D.POS", "DRAWPRIM3D.SIZE", "DRAWPRIM3D.COLOR", "DRAWPRIM3D.COL", "DRAWPRIM3D.WIRE", "DRAWPRIM3D.RADIUS",
		"DRAWPRIM3D.ENDPOINT", "DRAWPRIM3D.CYL", "DRAWPRIM3D.BBOX", "DRAWPRIM3D.SLICES", "DRAWPRIM3D.RINGS", "DRAWPRIM3D.GRID",
		"DRAWPRIM3D.SETRAY", "DRAWPRIM3D.SETTEXTURE", "DRAWPRIM3D.SRCTEX", "DRAWPRIM3D.DRAW", "DRAWPRIM3D.FREE",
		"DRAWCUBE", "DRAWCUBEWIRES", "DRAWSPHERE", "DRAWSPHEREW", "DRAWCYLINDER", "DRAWCYLINDERW", "DRAWCAP", "DRAWCAPW",
		"DRAWPLANE", "DRAWBBOX", "DRAWRAY", "DRAWLINE3D", "DRAWPOINT3D", "DRAWGRID3D", "DRAWBILLBOARD", "DRAWBILLBOARDREC",
		"DRAWPRIM2D.POS", "DRAWPRIM2D.SIZE", "DRAWPRIM2D.COLOR", "DRAWPRIM2D.COL", "DRAWPRIM2D.OUTLINE", "DRAWPRIM2D.P2", "DRAWPRIM2D.P3",
		"DRAWPRIM2D.RING", "DRAWPRIM2D.SEGS", "DRAWPRIM2D.SIDES", "DRAWPRIM2D.ROT", "DRAWPRIM2D.THICK", "DRAWPRIM2D.DRAW", "DRAWPRIM2D.FREE",
		"DRAWCIRCLE2", "DRAWCIRCLE2W", "DRAWELLIPSE2", "DRAWELLIPSE2W", "DRAWRECT2", "DRAWRECT2W", "DRAWLINE2", "DRAWTRI2", "DRAWTRI2W",
		"DRAWRING2", "DRAWRING2W", "DRAWPOLY2", "DRAWPOLY2W",
		"TEXTDRAW.POS", "TEXTDRAW.SIZE", "TEXTDRAW.COLOR", "TEXTDRAW.COL", "TEXTDRAW.SETTEXT", "TEXTDRAW.DRAW", "TEXTDRAW.FREE", "TEXTOBJ",
		"TEXTEXOBJ.POS", "TEXTEXOBJ.SIZE", "TEXTEXOBJ.SPACING", "TEXTEXOBJ.COLOR", "TEXTEXOBJ.SETTEXT", "TEXTEXOBJ.DRAW", "TEXTEXOBJ.FREE", "TEXTOBJEX",
		"DRAWTEX2.POS", "DRAWTEX2.COLOR", "DRAWTEX2.COL", "DRAWTEX2.SETTEXTURE", "DRAWTEX2.DRAW", "DRAWTEX2.FREE", "DRAWTEX2",
		"DRAWTEXREC", "DRAWTEXREC.SRC", "DRAWTEXREC.POS", "DRAWTEXREC.COLOR", "DRAWTEXREC.COL", "DRAWTEXREC.SETTEXTURE", "DRAWTEXREC.DRAW", "DRAWTEXREC.FREE",
		"DRAWTEXPRO", "DRAWTEXPRO.SRC", "DRAWTEXPRO.DST", "DRAWTEXPRO.ORIGIN", "DRAWTEXPRO.ROT", "DRAWTEXPRO.COLOR", "DRAWTEXPRO.COL", "DRAWTEXPRO.SETTEXTURE", "DRAWTEXPRO.DRAW", "DRAWTEXPRO.FREE",
	} {
		r.Register(name, "draw", stub(name))
	}
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
