package blitzengine

import (
	"fmt"
	"math"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func ival(x float64) value.Value {
	return value.FromInt(int64(math.Round(x)))
}

func (m *Module) setColor(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("SETCOLOR expects 3 arguments (r, g, b) 0–255")
	}
	r, _ := args[0].ToInt()
	g, _ := args[1].ToInt()
	b, _ := args[2].ToInt()
	m.pen.setColor(int(r), int(g), int(b))
	return value.Nil, nil
}

func (m *Module) setAlpha(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("SETALPHA expects 1 argument (alpha# 0..1)")
	}
	a, ok := args[0].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("SETALPHA: alpha must be numeric")
	}
	m.pen.setAlpha(a)
	return value.Nil, nil
}

func (m *Module) setOrigin(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("SETORIGIN expects 2 arguments (x#, y#)")
	}
	x, ok1 := args[0].ToFloat()
	y, ok2 := args[1].ToFloat()
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("SETORIGIN: x and y must be numeric")
	}
	m.pen.setOrigin(x, y)
	return value.Nil, nil
}

func (m *Module) setViewport(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("SETVIEWPORT expects 4 arguments (x, y, w, h)")
	}
	x, _ := args[0].ToInt()
	y, _ := args[1].ToInt()
	w, _ := args[2].ToInt()
	h, _ := args[3].ToInt()
	m.pen.setViewport(int32(x), int32(y), int32(w), int32(h))
	_, err := call(rt,"RENDER.SETSCISSOR", args[0], args[1], args[2], args[3])
	return value.Nil, err
}

func (m *Module) plot(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PLOT expects 2 arguments (x, y)")
	}
	ox, oy := m.pen.offset()
	x := valFloat(args[0]) + ox
	y := valFloat(args[1]) + oy
	r, g, b := m.pen.rgb()
	a := m.pen.rgbaA()
	return call(rt,"DRAW.PLOT", ival(x), ival(y), value.FromInt(int64(r)), value.FromInt(int64(g)), value.FromInt(int64(b)), value.FromInt(int64(a)))
}

func (m *Module) line(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("LINE expects 4 arguments (x1, y1, x2, y2)")
	}
	ox, oy := m.pen.offset()
	x1 := valFloat(args[0]) + ox
	y1 := valFloat(args[1]) + oy
	x2 := valFloat(args[2]) + ox
	y2 := valFloat(args[3]) + oy
	r, g, b := m.pen.rgb()
	a := m.pen.rgbaA()
	return call(rt,"DRAW.LINE",
		ival(x1), ival(y1), ival(x2), ival(y2),
		value.FromInt(int64(r)), value.FromInt(int64(g)), value.FromInt(int64(b)), value.FromInt(int64(a)),
	)
}

func (m *Module) rect(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("RECT expects 5 arguments (x, y, w, h, filled)")
	}
	ox, oy := m.pen.offset()
	x := valFloat(args[0]) + ox
	y := valFloat(args[1]) + oy
	w := valFloat(args[2])
	h := valFloat(args[3])
	filled, _ := args[4].ToInt()
	r, g, b := m.pen.rgb()
	a := m.pen.rgbaA()
	ri, gi, bi, ai := int64(r), int64(g), int64(b), int64(a)
	if filled != 0 {
		return call(rt,"DRAW.RECTANGLE",
			ival(x), ival(y), ival(w), ival(h),
			value.FromInt(ri), value.FromInt(gi), value.FromInt(bi), value.FromInt(ai),
		)
	}
	return call(rt,"DRAW.RECTLINES",
		ival(x), ival(y), ival(w), ival(h),
		value.FromFloat(1),
		value.FromInt(ri), value.FromInt(gi), value.FromInt(bi), value.FromInt(ai),
	)
}

func (m *Module) oval(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("OVAL expects 5 arguments (x, y, w, h, filled)")
	}
	ox, oy := m.pen.offset()
	x := valFloat(args[0]) + ox
	y := valFloat(args[1]) + oy
	w := valFloat(args[2])
	h := valFloat(args[3])
	filled, _ := args[4].ToInt()
	cx := x + w*0.5
	cy := y + h*0.5
	rx := w * 0.5
	ry := h * 0.5
	r, g, b := m.pen.rgb()
	a := m.pen.rgbaA()
	ri, gi, bi, ai := int64(r), int64(g), int64(b), int64(a)
	if filled != 0 {
		return call(rt,"DRAW.OVAL",
			ival(cx), ival(cy), value.FromFloat(rx), value.FromFloat(ry),
			value.FromInt(ri), value.FromInt(gi), value.FromInt(bi), value.FromInt(ai),
		)
	}
	return call(rt,"DRAW.OVALLINES",
		ival(cx), ival(cy), value.FromFloat(rx), value.FromFloat(ry),
		value.FromInt(ri), value.FromInt(gi), value.FromInt(bi), value.FromInt(ai),
	)
}

func (m *Module) textDraw(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("TEXT expects 3 arguments (x, y, text$)")
	}
	ox, oy := m.pen.offset()
	x := valFloat(args[0]) + ox
	y := valFloat(args[1]) + oy
	r, g, b := m.pen.rgb()
	a := m.pen.rgbaA()
	const size = 20
	return call(rt,"DRAW.TEXT", args[2],
		ival(x), ival(y), value.FromInt(size),
		value.FromInt(int64(r)), value.FromInt(int64(g)), value.FromInt(int64(b)), value.FromInt(int64(a)),
	)
}

func (m *Module) setFog(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("SETFOG expects 5 arguments (r, g, b, near#, far#)")
	}
	if _, err := call(rt,"FOG.ENABLE", value.FromBool(true)); err != nil {
		return value.Nil, err
	}
	if _, err := call(rt,"FOG.SETCOLOR", args[0], args[1], args[2], value.FromInt(255)); err != nil {
		return value.Nil, err
	}
	if _, err := call(rt,"FOG.SETRANGE", args[3], args[4]); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) createLight(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CREATELIGHT expects 2 arguments (type, parent) — parent ignored")
	}
	var kindStr string
	switch args[0].Kind {
	case value.KindString:
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		kindStr = s
	case value.KindInt, value.KindFloat:
		n, _ := args[0].ToInt()
		kindStr = lightKindFromInt(n)
	default:
		return value.Nil, fmt.Errorf("CREATELIGHT: type must be int (1=dir,2=point,3=spot) or kind string")
	}
	v, err := call(rt,"LIGHT.MAKE", value.FromStringIndex(rt.Heap.Intern(kindStr)))
	if err != nil {
		return value.Nil, err
	}
	pid, _ := args[1].ToInt()
	if pid != 0 {
		// No light–entity parenting in core API; ignore.
		_ = pid
	}
	return v, nil
}

func (m *Module) createCube(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CREATECUBE expects 1 argument (parent entity id)")
	}
	v, err := call(rt,"ENTITY.CREATECUBE", value.FromFloat(1), value.FromFloat(1), value.FromFloat(1))
	if err != nil {
		return value.Nil, err
	}
	pid, _ := args[0].ToInt()
	if pid != 0 {
		return call(rt,"ENTITY.PARENT", v, args[0])
	}
	return v, nil
}

func (m *Module) createSphereBlitz(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CREATESPHERE expects 2 arguments (segments, parent)")
	}
	seg, _ := args[0].ToInt()
	if seg < 3 {
		seg = 12
	}
	v, err := call(rt,"ENTITY.CREATESPHERE", value.FromFloat(1), value.FromInt(seg))
	if err != nil {
		return value.Nil, err
	}
	pid, _ := args[1].ToInt()
	if pid != 0 {
		return call(rt,"ENTITY.PARENT", v, args[1])
	}
	return v, nil
}

func (m *Module) createPlaneBlitz(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CREATEPLANE expects 2 arguments (divisions, parent)")
	}
	div, _ := args[0].ToInt()
	sz := float64(div)
	if sz <= 0 {
		sz = 10
	}
	v, err := call(rt,"ENTITY.CREATEPLANE", value.FromFloat(sz))
	if err != nil {
		return value.Nil, err
	}
	pid, _ := args[1].ToInt()
	if pid != 0 {
		return call(rt,"ENTITY.PARENT", v, args[1])
	}
	return v, nil
}

func (m *Module) createMeshBlitz(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CREATEMESH expects 1 argument (parent)")
	}
	v, err := call(rt,"ENTITY.CREATEMESH")
	if err != nil {
		return value.Nil, err
	}
	pid, _ := args[0].ToInt()
	if pid != 0 {
		return call(rt,"ENTITY.PARENT", v, args[0])
	}
	return v, nil
}

func (m *Module) loadMeshParent(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("LOADMESH expects 2 arguments (file$, parent)")
	}
	v, err := call(rt,"ENTITY.LOADMESH", args[0])
	if err != nil {
		return value.Nil, err
	}
	pid, _ := args[1].ToInt()
	if pid != 0 {
		return call(rt,"ENTITY.PARENT", v, args[1])
	}
	return v, nil
}

func (m *Module) loadAnimMeshParent(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("LOADANIMMESH expects 2 arguments (file$, parent)")
	}
	v, err := call(rt,"ENTITY.LOADANIMATEDMESH", args[0])
	if err != nil {
		return value.Nil, err
	}
	pid, _ := args[1].ToInt()
	if pid != 0 {
		return call(rt,"ENTITY.PARENT", v, args[1])
	}
	return v, nil
}

func (m *Module) meshWidth(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("MESHWIDTH expects 1 argument (mesh handle)")
	}
	minV, err := call(rt,"MESH.GETBBOXMINX", args[0])
	if err != nil {
		return value.Nil, err
	}
	maxV, err := call(rt,"MESH.GETBBOXMAXX", args[0])
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(math.Abs(valFloat(maxV) - valFloat(minV))), nil
}

func (m *Module) meshHeight(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("MESHHEIGHT expects 1 argument (mesh handle)")
	}
	minV, err := call(rt,"MESH.GETBBOXMINY", args[0])
	if err != nil {
		return value.Nil, err
	}
	maxV, err := call(rt,"MESH.GETBBOXMAXY", args[0])
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(math.Abs(valFloat(maxV) - valFloat(minV))), nil
}

func (m *Module) meshDepth(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("MESHDEPTH expects 1 argument (mesh handle)")
	}
	minV, err := call(rt,"MESH.GETBBOXMINZ", args[0])
	if err != nil {
		return value.Nil, err
	}
	maxV, err := call(rt,"MESH.GETBBOXMAXZ", args[0])
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(math.Abs(valFloat(maxV) - valFloat(minV))), nil
}

func (m *Module) createTexture(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CREATETEXTURE expects 2 arguments (w, h)")
	}
	img, err := call(rt,"IMAGE.MAKE", args[0], args[1])
	if err != nil {
		return value.Nil, err
	}
	return call(rt,"TEXTURE.FROMIMAGE", img)
}

func (m *Module) spriteAt(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("SPRITE expects 3 arguments (id, x, y)")
	}
	if _, err := call(rt,"SPRITE.SETPOS", args[0], args[1], args[2]); err != nil {
		return value.Nil, err
	}
	return call(rt,"SPRITE.DRAW", args[0], args[1], args[2])
}

func (m *Module) spriteNoOpTint(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	_ = args
	return value.Nil, nil
}

func (m *Module) entityScaleX(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITYSCALEX expects 1 argument (entity)")
	}
	return value.FromFloat(1), nil
}

func (m *Module) entityScaleY(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITYSCALEY expects 1 argument (entity)")
	}
	return value.FromFloat(1), nil
}

func (m *Module) entityScaleZ(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITYSCALEZ expects 1 argument (entity)")
	}
	return value.FromFloat(1), nil
}
