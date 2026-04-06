//go:build cgo || (windows && !cgo)

package mbmodel3d

import (
	"fmt"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/runtime/convert"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerModelRenderHierarchy(m *Module, reg runtime.Registrar) {
	reg.Register("MODEL.SETALPHA", "model", runtime.AdaptLegacy(m.modelSetAlpha))
	reg.Register("MODEL.SETWIREFRAME", "model", runtime.AdaptLegacy(m.modelSetWireframe))
	reg.Register("MODEL.SETCULL", "model", runtime.AdaptLegacy(m.modelSetCull))
	reg.Register("MODEL.SETLIGHTING", "model", runtime.AdaptLegacy(m.modelSetLighting))
	reg.Register("MODEL.SETFOG", "model", runtime.AdaptLegacy(m.modelSetFog))
	reg.Register("MODEL.SETBLEND", "model", runtime.AdaptLegacy(m.modelSetBlend))
	reg.Register("MODEL.SETDEPTH", "model", runtime.AdaptLegacy(m.modelSetDepth))
	reg.Register("MODEL.SETDIFFUSE", "model", runtime.AdaptLegacy(m.modelSetDiffuse))
	reg.Register("MODEL.SETSPECULAR", "model", runtime.AdaptLegacy(m.modelSetSpecular))
	reg.Register("MODEL.SETSPECULARPOW", "model", runtime.AdaptLegacy(m.modelSetSpecularPow))
	reg.Register("MODEL.SETEMISSIVE", "model", runtime.AdaptLegacy(m.modelSetEmissive))
	reg.Register("MODEL.SETAMBIENTCOLOR", "model", runtime.AdaptLegacy(m.modelSetAmbientColor))

	reg.Register("MODEL.CLONE", "model", runtime.AdaptLegacy(m.modelClone))
	reg.Register("MODEL.INSTANCE", "model", runtime.AdaptLegacy(m.modelInstance))
	reg.Register("MODEL.ATTACHTO", "model", runtime.AdaptLegacy(m.modelAttachTo))
	reg.Register("MODEL.DETACH", "model", runtime.AdaptLegacy(m.modelDetach))
	reg.Register("MODEL.EXISTS", "model", runtime.AdaptLegacy(m.modelExists))
	reg.Register("MODEL.SETGPUSKINNING", "model", runtime.AdaptLegacy(m.modelSetGPUSkinning))
}

func clampU8(v int64) uint8 {
	switch {
	case v < 0:
		return 0
	case v > 255:
		return 255
	default:
		return uint8(v)
	}
}

func (m *Module) modelSetAlpha(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("MODEL.SETALPHA expects (model, a)")
	}
	o, err := m.getModel(args, 0, "MODEL.SETALPHA")
	if err != nil {
		return value.Nil, err
	}
	a, ok := args[1].ToInt()
	if !ok {
		if f, okf := args[1].ToFloat(); okf {
			a = int64(f)
		} else {
			return value.Nil, fmt.Errorf("MODEL.SETALPHA: a must be numeric")
		}
	}
	au := clampU8(a)
	mats := o.model.GetMaterials()
	for i := range mats {
		mp := mats[i].GetMap(rl.MapAlbedo)
		c := mp.Color
		c.A = au
		mp.Color = c
	}
	return value.Nil, nil
}

func (m *Module) modelSetWireframe(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	o, err := m.getModel(args, 0, "MODEL.SETWIREFRAME")
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("MODEL.SETWIREFRAME expects (model, enable?)")
	}
	en, ok := argBool(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("MODEL.SETWIREFRAME: enable must be bool or numeric")
	}
	o.wireframe = en
	return value.Nil, nil
}

func (m *Module) modelSetCull(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	o, err := m.getModel(args, 0, "MODEL.SETCULL")
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("MODEL.SETCULL expects (model, enable?)")
	}
	en, ok := argBool(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("MODEL.SETCULL: enable must be bool or numeric")
	}
	o.cull = en
	return value.Nil, nil
}

func (m *Module) modelSetLighting(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	o, err := m.getModel(args, 0, "MODEL.SETLIGHTING")
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("MODEL.SETLIGHTING expects (model, enable?)")
	}
	en, ok := argBool(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("MODEL.SETLIGHTING: enable must be bool or numeric")
	}
	o.lighting = en
	return value.Nil, nil
}

func (m *Module) modelSetFog(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	o, err := m.getModel(args, 0, "MODEL.SETFOG")
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("MODEL.SETFOG expects (model, enable?)")
	}
	en, ok := argBool(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("MODEL.SETFOG: enable must be bool or numeric")
	}
	o.fog = en
	return value.Nil, nil
}

func (m *Module) modelSetBlend(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("MODEL.SETBLEND expects (model, mode)")
	}
	o, err := m.getModel(args, 0, "MODEL.SETBLEND")
	if err != nil {
		return value.Nil, err
	}
	mode, ok := argInt(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("MODEL.SETBLEND: mode must be numeric (use BLEND_* constants)")
	}
	o.blendMode = mode
	return value.Nil, nil
}

func (m *Module) modelSetDepth(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("MODEL.SETDEPTH expects (model, depth)")
	}
	o, err := m.getModel(args, 0, "MODEL.SETDEPTH")
	if err != nil {
		return value.Nil, err
	}
	d, ok := argInt(args[1])
	if !ok || d < 0 {
		return value.Nil, fmt.Errorf("MODEL.SETDEPTH: depth must be non-negative int (bitmask: 1=no depth test, 2=no depth write)")
	}
	o.depthBits = d
	return value.Nil, nil
}

func (m *Module) modelSetDiffuse(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("MODEL.SETDIFFUSE expects (model, r, g, b)")
	}
	o, err := m.getModel(args, 0, "MODEL.SETDIFFUSE")
	if err != nil {
		return value.Nil, err
	}
	col, err := rgbaFromArgs(args[1], args[2], args[3], value.FromInt(255))
	if err != nil {
		return value.Nil, fmt.Errorf("MODEL.SETDIFFUSE: %w", err)
	}
	mats := o.model.GetMaterials()
	for i := range mats {
		mp := mats[i].GetMap(rl.MapAlbedo)
		c := mp.Color
		c.R, c.G, c.B = col.R, col.G, col.B
		mp.Color = c
	}
	return value.Nil, nil
}

func (m *Module) modelSetSpecular(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("MODEL.SETSPECULAR expects (model, r, g, b)")
	}
	o, err := m.getModel(args, 0, "MODEL.SETSPECULAR")
	if err != nil {
		return value.Nil, err
	}
	col, err := rgbaFromArgs(args[1], args[2], args[3], value.FromInt(255))
	if err != nil {
		return value.Nil, fmt.Errorf("MODEL.SETSPECULAR: %w", err)
	}
	mats := o.model.GetMaterials()
	for i := range mats {
		mp := mats[i].GetMap(rl.MapSpecular)
		c := mp.Color
		c.R, c.G, c.B = col.R, col.G, col.B
		mp.Color = c
	}
	return value.Nil, nil
}

func (m *Module) modelSetSpecularPow(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("MODEL.SETSPECULARPOW expects (model, pow)")
	}
	o, err := m.getModel(args, 0, "MODEL.SETSPECULARPOW")
	if err != nil {
		return value.Nil, err
	}
	pow, ok := argFloat(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("MODEL.SETSPECULARPOW: pow must be numeric")
	}
	mats := o.model.GetMaterials()
	for i := range mats {
		mat := &mats[i]
		mat.Params[0] = pow
	}
	return value.Nil, nil
}

func (m *Module) modelSetEmissive(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("MODEL.SETEMISSIVE expects (model, r, g, b)")
	}
	o, err := m.getModel(args, 0, "MODEL.SETEMISSIVE")
	if err != nil {
		return value.Nil, err
	}
	ri, ok1 := argInt(args[1])
	gi, ok2 := argInt(args[2])
	bi, ok3 := argInt(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("MODEL.SETEMISSIVE: r, g, b must be numeric")
	}
	c := convert.NewColor4(ri, gi, bi, 255)
	em := color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A}
	mats := o.model.GetMaterials()
	for i := range mats {
		mp := mats[i].GetMap(rl.MapEmission)
		mp.Color = em
	}
	return value.Nil, nil
}

func (m *Module) modelSetAmbientColor(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("MODEL.SETAMBIENTCOLOR expects (model, r, g, b)")
	}
	o, err := m.getModel(args, 0, "MODEL.SETAMBIENTCOLOR")
	if err != nil {
		return value.Nil, err
	}
	ri, ok1 := argInt(args[1])
	gi, ok2 := argInt(args[2])
	bi, ok3 := argInt(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("MODEL.SETAMBIENTCOLOR: r, g, b must be numeric")
	}
	o.ambientR, o.ambientG, o.ambientB = int32(ri), int32(gi), int32(bi)
	return value.Nil, nil
}

func (m *Module) modelClone(args []value.Value) (value.Value, error) {
	return m.modelReloadCopy(args, "MODEL.CLONE")
}

func (m *Module) modelInstance(args []value.Value) (value.Value, error) {
	// Raylib does not expose shared-GPU-mesh instances in these bindings; reload the same asset (cheap vs full vertex copy).
	return m.modelReloadCopy(args, "MODEL.INSTANCE")
}

func (m *Module) modelReloadCopy(args []value.Value, op string) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("%s expects (model)", op)
	}
	src, err := m.getModel(args, 0, op)
	if err != nil {
		return value.Nil, err
	}
	if src.loadedPath == "" {
		return value.Nil, fmt.Errorf("%s: source model has no file path (only models from MODEL.LOAD can be cloned/instanced)", op)
	}
	mod := rl.LoadModel(src.loadedPath)
	id, err := m.h.Alloc(&modelObj{model: mod, loadedPath: src.loadedPath})
	if err != nil {
		rl.UnloadModel(mod)
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) modelAttachTo(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("MODEL.ATTACHTO expects (model, parent)")
	}
	child, err := m.getModel(args, 0, "MODEL.ATTACHTO")
	if err != nil {
		return value.Nil, err
	}
	if args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("MODEL.ATTACHTO: parent must be a handle")
	}
	ph := heap.Handle(args[1].IVal)
	ch := heap.Handle(args[0].IVal)
	if ph == ch {
		return value.Nil, fmt.Errorf("MODEL.ATTACHTO: cannot attach to self")
	}
	if _, err := heap.Cast[*modelObj](m.h, ph); err != nil {
		return value.Nil, fmt.Errorf("MODEL.ATTACHTO: parent must be a model handle")
	}
	child.parent = ph
	return value.Nil, nil
}

func (m *Module) modelDetach(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("MODEL.DETACH expects (model)")
	}
	o, err := m.getModel(args, 0, "MODEL.DETACH")
	if err != nil {
		return value.Nil, err
	}
	o.parent = 0
	return value.Nil, nil
}

func (m *Module) modelExists(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("MODEL.EXISTS expects (model handle)")
	}
	obj, ok := m.h.Get(heap.Handle(args[0].IVal))
	if !ok {
		return value.FromBool(false), nil
	}
	return value.FromBool(obj.TypeTag() == heap.TagModel), nil
}

func (m *Module) modelSetGPUSkinning(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("MODEL.SETGPUSKINNING expects (model, enable?)")
	}
	_, err := m.getModel(args, 0, "MODEL.SETGPUSKINNING")
	if err != nil {
		return value.Nil, err
	}
	_, ok := argBool(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("MODEL.SETGPUSKINNING: enable must be bool or numeric")
	}
	// raylib 5.6-dev in current raylib-go removed the per-model GPU skinning flag; skin path uses
	// UpdateModelAnimation / UpdateModelAnimationBones. Kept as a documented no-op for scripts.
	return value.Nil, nil
}
