//go:build cgo || (windows && !cgo)

package mbmodel3d

import (
	"fmt"
	"image/color"
	"math"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/runtime/convert"
	"moonbasic/runtime/mbmatrix"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerModelComplete(m *Module, reg runtime.Registrar) {
	reg.Register("MODEL.MOVE", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 4 {
			return value.Nil, fmt.Errorf("MODEL.MOVE expects (model, dx, dy, dz)")
		}
		o, err := m.getModel(args, 0, "MODEL.MOVE")
		if err != nil {
			return value.Nil, err
		}
		dx, ok1 := argFloat(args[1])
		dy, ok2 := argFloat(args[2])
		dz, ok3 := argFloat(args[3])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("MODEL.MOVE: deltas must be numeric")
		}
		t := rl.MatrixTranslate(dx, dy, dz)
		o.model.Transform = rl.MatrixMultiply(t, o.model.Transform)
		return value.Nil, nil
	}))

	reg.Register("MODEL.X", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		o, err := m.getModelTransform(args, "MODEL.X")
		if err != nil {
			return value.Nil, err
		}
		return value.FromFloat(float64(o.model.Transform.M12)), nil
	}))
	reg.Register("MODEL.Y", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		o, err := m.getModelTransform(args, "MODEL.Y")
		if err != nil {
			return value.Nil, err
		}
		return value.FromFloat(float64(o.model.Transform.M13)), nil
	}))
	reg.Register("MODEL.Z", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		o, err := m.getModelTransform(args, "MODEL.Z")
		if err != nil {
			return value.Nil, err
		}
		return value.FromFloat(float64(o.model.Transform.M14)), nil
	}))

	reg.Register("MODEL.GETPOS", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		o, err := m.getModelTransform(args, "MODEL.GETPOS")
		if err != nil {
			return value.Nil, err
		}
		mat := o.model.Transform
		return mbmatrix.AllocVec3Value(m.h, mat.M12, mat.M13, mat.M14)
	}))

	reg.Register("MODEL.SETROT", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 4 {
			return value.Nil, fmt.Errorf("MODEL.SETROT expects (model, rx, ry, rz) radians")
		}
		o, err := m.getModel(args, 0, "MODEL.SETROT")
		if err != nil {
			return value.Nil, err
		}
		rx, ok1 := argFloat(args[1])
		ry, ok2 := argFloat(args[2])
		rz, ok3 := argFloat(args[3])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("MODEL.SETROT: angles must be numeric")
		}
		mat := o.model.Transform
		tx, ty, tz := mat.M12, mat.M13, mat.M14
		rot := rl.MatrixRotateXYZ(rl.Vector3{X: rx, Y: ry, Z: rz})
		o.model.Transform = rl.MatrixMultiply(rl.MatrixTranslate(tx, ty, tz), rot)
		return value.Nil, nil
	}))

	reg.Register("MODEL.ROTATE", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 4 {
			return value.Nil, fmt.Errorf("MODEL.ROTATE expects (model, drx, dry, drz) radians")
		}
		o, err := m.getModel(args, 0, "MODEL.ROTATE")
		if err != nil {
			return value.Nil, err
		}
		drx, ok1 := argFloat(args[1])
		dry, ok2 := argFloat(args[2])
		drz, ok3 := argFloat(args[3])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("MODEL.ROTATE: angles must be numeric")
		}
		d := rl.MatrixRotateXYZ(rl.Vector3{X: drx, Y: dry, Z: drz})
		o.model.Transform = rl.MatrixMultiply(d, o.model.Transform)
		return value.Nil, nil
	}))

	reg.Register("MODEL.GETROT", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("MODEL.GETROT expects model handle")
		}
		o, err := m.getModel(args, 0, "MODEL.GETROT")
		if err != nil {
			return value.Nil, err
		}
		eul := rl.QuaternionToEuler(rl.QuaternionFromMatrix(o.model.Transform))
		arr, err := heap.NewArray([]int64{3})
		if err != nil {
			return value.Nil, err
		}
		_ = arr.Set([]int64{0}, float64(eul.X))
		_ = arr.Set([]int64{1}, float64(eul.Y))
		_ = arr.Set([]int64{2}, float64(eul.Z))
		id, err := m.h.Alloc(arr)
		if err != nil {
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	}))

	reg.Register("MODEL.SETSCALE", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 4 {
			return value.Nil, fmt.Errorf("MODEL.SETSCALE expects (model, sx, sy, sz)")
		}
		o, err := m.getModel(args, 0, "MODEL.SETSCALE")
		if err != nil {
			return value.Nil, err
		}
		sx, ok1 := argFloat(args[1])
		sy, ok2 := argFloat(args[2])
		sz, ok3 := argFloat(args[3])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("MODEL.SETSCALE: scale must be numeric")
		}
		mat := o.model.Transform
		tx, ty, tz := mat.M12, mat.M13, mat.M14
		s := rl.MatrixScale(sx, sy, sz)
		o.model.Transform = rl.MatrixMultiply(rl.MatrixTranslate(tx, ty, tz), s)
		return value.Nil, nil
	}))

	reg.Register("MODEL.SETSCALEUNIFORM", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("MODEL.SETSCALEUNIFORM expects (model, s)")
		}
		o, err := m.getModel(args, 0, "MODEL.SETSCALEUNIFORM")
		if err != nil {
			return value.Nil, err
		}
		s, ok := argFloat(args[1])
		if !ok {
			return value.Nil, fmt.Errorf("MODEL.SETSCALEUNIFORM: s must be numeric")
		}
		return modelSetScaleUniform(o, s)
	}))

	reg.Register("MODEL.GETSCALE", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("MODEL.GETSCALE expects model handle")
		}
		o, err := m.getModel(args, 0, "MODEL.GETSCALE")
		if err != nil {
			return value.Nil, err
		}
		// Column lengths of rotation/scale part (rough).
		mat := o.model.Transform
		sx := float32(math.Sqrt(float64(mat.M0*mat.M0 + mat.M1*mat.M1 + mat.M2*mat.M2)))
		sy := float32(math.Sqrt(float64(mat.M4*mat.M4 + mat.M5*mat.M5 + mat.M6*mat.M6)))
		sz := float32(math.Sqrt(float64(mat.M8*mat.M8 + mat.M9*mat.M9 + mat.M10*mat.M10)))
		arr, err := heap.NewArray([]int64{3})
		if err != nil {
			return value.Nil, err
		}
		_ = arr.Set([]int64{0}, float64(sx))
		_ = arr.Set([]int64{1}, float64(sy))
		_ = arr.Set([]int64{2}, float64(sz))
		id, err := m.h.Alloc(arr)
		if err != nil {
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	}))

	reg.Register("MODEL.SETMATRIX", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("MODEL.SETMATRIX expects (model, mat4_handle)")
		}
		o, err := m.getModel(args, 0, "MODEL.SETMATRIX")
		if err != nil {
			return value.Nil, err
		}
		if args[1].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("MODEL.SETMATRIX: second argument must be MAT4 handle")
		}
		mat, err := mbmatrix.MatrixRaylib(m.h, heap.Handle(args[1].IVal))
		if err != nil {
			return value.Nil, err
		}
		o.model.Transform = mat
		return value.Nil, nil
	}))

	reg.Register("MODEL.DRAWAT", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 10 {
			return value.Nil, fmt.Errorf("MODEL.DRAWAT expects (model, x,y,z, rx,ry,rz, sx,sy,sz)")
		}
		o, err := m.getModel(args, 0, "MODEL.DRAWAT")
		if err != nil {
			return value.Nil, err
		}
		x, ok1 := argFloat(args[1])
		y, ok2 := argFloat(args[2])
		z, ok3 := argFloat(args[3])
		rx, ok4 := argFloat(args[4])
		ry, ok5 := argFloat(args[5])
		rz, ok6 := argFloat(args[6])
		sx, ok7 := argFloat(args[7])
		sy, ok8 := argFloat(args[8])
		sz, ok9 := argFloat(args[9])
		if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 || !ok8 || !ok9 {
			return value.Nil, fmt.Errorf("MODEL.DRAWAT: numeric arguments required")
		}
		if o.hidden {
			return value.Nil, nil
		}
		saved := o.model.Transform
		rot := rl.MatrixRotateXYZ(rl.Vector3{X: rx, Y: ry, Z: rz})
		scl := rl.MatrixScale(sx, sy, sz)
		o.model.Transform = rl.MatrixMultiply(rl.MatrixMultiply(rl.MatrixTranslate(x, y, z), rot), scl)
		rl.DrawModel(o.model, rl.Vector3{}, 1, rl.White)
		o.model.Transform = saved
		return value.Nil, nil
	}))

	reg.Register("MODEL.DRAWEX", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 15 {
			return value.Nil, fmt.Errorf("MODEL.DRAWEX expects (model, x,y,z, ax,ay,az, angle, sx,sy,sz, r,g,b,a)")
		}
		o, err := m.getModel(args, 0, "MODEL.DRAWEX")
		if err != nil {
			return value.Nil, err
		}
		x, _ := argFloat(args[1])
		y, _ := argFloat(args[2])
		z, _ := argFloat(args[3])
		ax, _ := argFloat(args[4])
		ay, _ := argFloat(args[5])
		az, _ := argFloat(args[6])
		ang, _ := argFloat(args[7])
		sx, _ := argFloat(args[8])
		sy, _ := argFloat(args[9])
		sz, _ := argFloat(args[10])
		ri, _ := argInt(args[11])
		gi, _ := argInt(args[12])
		bi, _ := argInt(args[13])
		ai, _ := argInt(args[14])
		if o.hidden {
			return value.Nil, nil
		}
		c := convert.NewColor4(ri, gi, bi, ai)
		tint := color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A}
		axis := rl.Vector3{X: ax, Y: ay, Z: az}
		if rl.Vector3Length(axis) < 1e-5 {
			axis = rl.Vector3{X: 0, Y: 1, Z: 0}
		} else {
			axis = rl.Vector3Normalize(axis)
		}
		angDeg := float32(ang * 180 / math.Pi)
		rl.DrawModelEx(o.model, rl.Vector3{X: x, Y: y, Z: z}, axis, angDeg, rl.Vector3{X: sx, Y: sy, Z: sz}, tint)
		return value.Nil, nil
	}))

	reg.Register("MODEL.DRAWWIRES", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 5 {
			return value.Nil, fmt.Errorf("MODEL.DRAWWIRES expects (model, r,g,b,a)")
		}
		o, err := m.getModel(args, 0, "MODEL.DRAWWIRES")
		if err != nil {
			return value.Nil, err
		}
		ri, _ := argInt(args[1])
		gi, _ := argInt(args[2])
		bi, _ := argInt(args[3])
		ai, _ := argInt(args[4])
		c := convert.NewColor4(ri, gi, bi, ai)
		tint := color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A}
		pos := rl.Vector3{X: o.model.Transform.M12, Y: o.model.Transform.M13, Z: o.model.Transform.M14}
		rl.DrawModelWires(o.model, pos, 1, tint)
		return value.Nil, nil
	}))

	reg.Register("MODEL.SHOW", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("MODEL.SHOW expects model handle")
		}
		o, err := m.getModel(args, 0, "MODEL.SHOW")
		if err != nil {
			return value.Nil, err
		}
		o.hidden = false
		return value.Nil, nil
	}))
	reg.Register("MODEL.HIDE", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("MODEL.HIDE expects model handle")
		}
		o, err := m.getModel(args, 0, "MODEL.HIDE")
		if err != nil {
			return value.Nil, err
		}
		o.hidden = true
		return value.Nil, nil
	}))
	reg.Register("MODEL.ISVISIBLE", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("MODEL.ISVISIBLE expects model handle")
		}
		o, err := m.getModel(args, 0, "MODEL.ISVISIBLE")
		if err != nil {
			return value.Nil, err
		}
		return value.FromBool(!o.hidden), nil
	}))

	reg.Register("MODEL.SETCOLOR", "model", runtime.AdaptLegacy(m.modelSetColorRGBA))
	reg.Register("MODEL.SETMETAL", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("MODEL.SETMETAL expects (model, amount)")
		}
		o, err := m.getModel(args, 0, "MODEL.SETMETAL")
		if err != nil {
			return value.Nil, err
		}
		v, ok := argFloat(args[1])
		if !ok {
			return value.Nil, fmt.Errorf("MODEL.SETMETAL: amount must be numeric")
		}
		mats := o.model.GetMaterials()
		for i := range mats {
			mats[i].Params[0] = v
		}
		return value.Nil, nil
	}))
	reg.Register("MODEL.SETROUGH", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("MODEL.SETROUGH expects (model, amount)")
		}
		o, err := m.getModel(args, 0, "MODEL.SETROUGH")
		if err != nil {
			return value.Nil, err
		}
		v, ok := argFloat(args[1])
		if !ok {
			return value.Nil, fmt.Errorf("MODEL.SETROUGH: amount must be numeric")
		}
		mats := o.model.GetMaterials()
		for i := range mats {
			mats[i].Params[1] = v
		}
		return value.Nil, nil
	}))

	reg.Register("MODEL.LOADANIMATIONS", "model", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 2 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
			return value.Nil, fmt.Errorf("MODEL.LOADANIMATIONS expects (model, path)")
		}
		o, err := m.getModel(args, 0, "MODEL.LOADANIMATIONS")
		if err != nil {
			return value.Nil, err
		}
		path, err := rt.ArgString(args, 1)
		if err != nil {
			return value.Nil, err
		}
		path = strings.TrimSpace(path)
		if path == "" {
			return value.Nil, fmt.Errorf("MODEL.LOADANIMATIONS: path required")
		}
		if len(o.anims) > 0 {
			rl.UnloadModelAnimations(o.anims)
			o.anims = nil
		}
		o.anims = rl.LoadModelAnimations(path)
		o.animIdx = 0
		o.animFrame = 0
		return value.Nil, nil
	})

	reg.Register("MODEL.UPDATEANIM", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("MODEL.UPDATEANIM expects (model, dt)")
		}
		o, err := m.getModel(args, 0, "MODEL.UPDATEANIM")
		if err != nil {
			return value.Nil, err
		}
		dt, ok := argFloat(args[1])
		if !ok {
			return value.Nil, fmt.Errorf("MODEL.UPDATEANIM: dt must be numeric")
		}
		if len(o.anims) == 0 || !o.animPlaying {
			return value.Nil, nil
		}
		anim := o.anims[o.animIdx]
		// No per-clip duration in this raylib struct — advance at ~60 base FPS scaled by animSpeed.
		o.animFrame += float32(dt) * 60 * o.animSpeed
		fc := float32(anim.FrameCount)
		if o.animLoop && fc > 0 {
			for o.animFrame >= fc {
				o.animFrame -= fc
			}
		} else if fc > 0 && o.animFrame >= fc {
			o.animFrame = fc - 1
			o.animPlaying = false
		}
		rl.UpdateModelAnimation(o.model, anim, int32(o.animFrame))
		return value.Nil, nil
	}))

	reg.Register("MODEL.PLAYIDX", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("MODEL.PLAYIDX expects (model, idx)")
		}
		o, err := m.getModel(args, 0, "MODEL.PLAYIDX")
		if err != nil {
			return value.Nil, err
		}
		idx, ok := argInt(args[1])
		if !ok || int(idx) < 0 || int(idx) >= len(o.anims) {
			return value.Nil, fmt.Errorf("MODEL.PLAYIDX: invalid animation index")
		}
		o.animIdx = int(idx)
		o.animPlaying = true
		o.animFrame = 0
		return value.Nil, nil
	}))

	reg.Register("MODEL.STOP", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("MODEL.STOP expects model handle")
		}
		o, err := m.getModel(args, 0, "MODEL.STOP")
		if err != nil {
			return value.Nil, err
		}
		o.animPlaying = false
		return value.Nil, nil
	}))

	reg.Register("MODEL.LOOP", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("MODEL.LOOP expects (model, enable)")
		}
		o, err := m.getModel(args, 0, "MODEL.LOOP")
		if err != nil {
			return value.Nil, err
		}
		en, ok := argBool(args[1])
		if !ok {
			return value.Nil, fmt.Errorf("MODEL.LOOP: enable must be bool or numeric")
		}
		o.animLoop = en
		return value.Nil, nil
	}))

	reg.Register("MODEL.SETSPEED", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("MODEL.SETSPEED expects (model, fps)")
		}
		o, err := m.getModel(args, 0, "MODEL.SETSPEED")
		if err != nil {
			return value.Nil, err
		}
		v, ok := argFloat(args[1])
		if !ok {
			return value.Nil, fmt.Errorf("MODEL.SETSPEED: fps must be numeric")
		}
		o.animSpeed = v
		return value.Nil, nil
	}))

	reg.Register("MODEL.GETPARENT", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("MODEL.GETPARENT expects model handle")
		}
		o, err := m.getModel(args, 0, "MODEL.GETPARENT")
		if err != nil {
			return value.Nil, err
		}
		return value.FromHandle(o.parent), nil
	}))

	reg.Register("MODEL.CHILDCOUNT", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		_ = args
		return value.FromInt(0), nil
	}))

	reg.Register("MODEL.LIMBCOUNT", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("MODEL.LIMBCOUNT expects model handle")
		}
		o, err := m.getModel(args, 0, "MODEL.LIMBCOUNT")
		if err != nil {
			return value.Nil, err
		}
		return value.FromInt(int64(o.model.BoneCount)), nil
	}))

	stubModel := func(msg string) func([]value.Value) (value.Value, error) {
		return func(args []value.Value) (value.Value, error) {
			_ = args
			return value.Nil, fmt.Errorf("%s", msg)
		}
	}
	reg.Register("MODEL.PLAY", "model", runtime.AdaptLegacy(stubModel("MODEL.PLAY: use MODEL.PLAYIDX after MODEL.LOADANIMATIONS")))
	reg.Register("MODEL.GETFRAME", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("MODEL.GETFRAME expects model handle")
		}
		o, err := m.getModel(args, 0, "MODEL.GETFRAME")
		if err != nil {
			return value.Nil, err
		}
		return value.FromInt(int64(o.animFrame)), nil
	}))
	reg.Register("MODEL.TOTALFRAMES", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("MODEL.TOTALFRAMES expects model handle")
		}
		o, err := m.getModel(args, 0, "MODEL.TOTALFRAMES")
		if err != nil {
			return value.Nil, err
		}
		if len(o.anims) == 0 {
			return value.FromInt(0), nil
		}
		return value.FromInt(int64(o.anims[o.animIdx].FrameCount)), nil
	}))
	reg.Register("MODEL.ISPLAYING", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("MODEL.ISPLAYING expects model handle")
		}
		o, err := m.getModel(args, 0, "MODEL.ISPLAYING")
		if err != nil {
			return value.Nil, err
		}
		return value.FromBool(o.animPlaying), nil
	}))
	reg.Register("MODEL.ANIMDONE", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("MODEL.ANIMDONE expects model handle")
		}
		o, err := m.getModel(args, 0, "MODEL.ANIMDONE")
		if err != nil {
			return value.Nil, err
		}
		if len(o.anims) == 0 {
			return value.FromBool(true), nil
		}
		fc := float32(o.anims[o.animIdx].FrameCount)
		done := !o.animPlaying && fc > 0 && o.animFrame >= fc-1
		return value.FromBool(done), nil
	}))
	reg.Register("MODEL.ANIMCOUNT", "model", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("MODEL.ANIMCOUNT expects model handle")
		}
		o, err := m.getModel(args, 0, "MODEL.ANIMCOUNT")
		if err != nil {
			return value.Nil, err
		}
		return value.FromInt(int64(len(o.anims))), nil
	}))
	modelAnimName := func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 2 {
			return value.Nil, fmt.Errorf("MODEL.ANIMNAME expects (model, idx)")
		}
		o, err := m.getModel(args, 0, "MODEL.ANIMNAME")
		if err != nil {
			return value.Nil, err
		}
		idx, ok := argInt(args[1])
		if !ok || int(idx) < 0 || int(idx) >= len(o.anims) {
			return value.Nil, fmt.Errorf("MODEL.ANIMNAME: invalid index")
		}
		return rt.RetString(o.anims[idx].GetName()), nil
	}
	reg.Register("MODEL.ANIMNAME", "model", modelAnimName)
	reg.Register("MODEL.ANIMNAME$", "model", modelAnimName)
	reg.Register("MODEL.ADDCHILD", "model", runtime.AdaptLegacy(stubModel("MODEL.ADDCHILD: use MODEL.ATTACHTO")))
	reg.Register("MODEL.REMOVECHILD", "model", runtime.AdaptLegacy(stubModel("MODEL.REMOVECHILD: use MODEL.DETACH")))
	reg.Register("MODEL.GETCHILD", "model", runtime.AdaptLegacy(stubModel("MODEL.GETCHILD: not implemented")))
	reg.Register("MODEL.LIMBX", "model", runtime.AdaptLegacy(stubModel("MODEL.LIMBX: bone API not wired")))
	reg.Register("MODEL.SETLIMBPOS", "model", runtime.AdaptLegacy(stubModel("MODEL.SETLIMBPOS: bone API not wired")))
	reg.Register("MODEL.SETCASTSHADOW", "model", runtime.AdaptLegacy(stubModel("MODEL.SETCASTSHADOW: use light/shadow pipeline")))
	reg.Register("MODEL.SETRECEIVESHADOW", "model", runtime.AdaptLegacy(stubModel("MODEL.SETRECEIVESHADOW: use light/shadow pipeline")))
}

func (m *Module) modelSetColorRGBA(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("MODEL.SETCOLOR expects (model, r,g,b,a)")
	}
	o, err := m.getModel(args, 0, "MODEL.SETCOLOR")
	if err != nil {
		return value.Nil, err
	}
	col, err := rgbaFromArgs(args[1], args[2], args[3], args[4])
	if err != nil {
		return value.Nil, err
	}
	mats := o.model.GetMaterials()
	for i := range mats {
		mp := mats[i].GetMap(rl.MapAlbedo)
		c := mp.Color
		c.R, c.G, c.B, c.A = col.R, col.G, col.B, col.A
		mp.Color = c
	}
	return value.Nil, nil
}

func (m *Module) getModelTransform(args []value.Value, op string) (*modelObj, error) {
	if err := m.requireHeap(); err != nil {
		return nil, err
	}
	if len(args) != 1 {
		return nil, fmt.Errorf("%s expects model handle", op)
	}
	return m.getModel(args, 0, op)
}

func modelSetScaleUniform(o *modelObj, s float32) (value.Value, error) {
	mat := o.model.Transform
	tx, ty, tz := mat.M12, mat.M13, mat.M14
	sc := rl.MatrixScale(s, s, s)
	o.model.Transform = rl.MatrixMultiply(rl.MatrixTranslate(tx, ty, tz), sc)
	return value.Nil, nil
}
