//go:build cgo || (windows && !cgo)

package mbmatrix

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) registerMat4(reg runtime.Registrar) {
	// Transform.* — 4×4 object/world matrices (same heap tag "mat4" as Mat4.*).
	reg.Register("TRANSFORM.IDENTITY", "mat4", runtime.AdaptLegacy(m.mat4Identity))
	reg.Register("TRANSFORM.TRANSLATION", "mat4", runtime.AdaptLegacy(m.mat4FromTranslation))
	reg.Register("TRANSFORM.ROTATION", "mat4", runtime.AdaptLegacy(m.mat4FromRotation))
	reg.Register("TRANSFORM.SCALE", "mat4", runtime.AdaptLegacy(m.mat4FromScale))
	reg.Register("TRANSFORM.SETROTATION", "mat4", runtime.AdaptLegacy(m.mat4SetRotation))
	reg.Register("TRANSFORM.LOOKAT", "mat4", runtime.AdaptLegacy(m.mat4LookAt))
	reg.Register("TRANSFORM.PERSPECTIVE", "mat4", runtime.AdaptLegacy(m.mat4Perspective))
	reg.Register("TRANSFORM.ORTHO", "mat4", runtime.AdaptLegacy(m.mat4Ortho))
	reg.Register("TRANSFORM.MULTIPLY", "mat4", runtime.AdaptLegacy(m.mat4Multiply))
	reg.Register("TRANSFORM.INVERSE", "mat4", runtime.AdaptLegacy(m.mat4Inverse))
	reg.Register("TRANSFORM.TRANSPOSE", "mat4", runtime.AdaptLegacy(m.mat4Transpose))
	reg.Register("TRANSFORM.GETELEMENT", "mat4", runtime.AdaptLegacy(m.mat4GetElement))
	reg.Register("TRANSFORM.APPLYX", "mat4", runtime.AdaptLegacy(m.transformApplyX))
	reg.Register("TRANSFORM.APPLYY", "mat4", runtime.AdaptLegacy(m.transformApplyY))
	reg.Register("TRANSFORM.APPLYZ", "mat4", runtime.AdaptLegacy(m.transformApplyZ))
	reg.Register("TRANSFORM.FREE", "mat4", runtime.AdaptLegacy(m.mat4Free))

	// Mat4.* — legacy names (same implementations).
	reg.Register("MAT4.IDENTITY", "mat4", runtime.AdaptLegacy(m.mat4Identity))
	reg.Register("MAT4.FROMROTATION", "mat4", runtime.AdaptLegacy(m.mat4FromRotation))
	reg.Register("MAT4.ROTATION", "mat4", runtime.AdaptLegacy(m.mat4FromRotation))
	reg.Register("MAT4.SETROTATION", "mat4", runtime.AdaptLegacy(m.mat4SetRotation))
	reg.Register("MAT4.FROMSCALE", "mat4", runtime.AdaptLegacy(m.mat4FromScale))
	reg.Register("MAT4.FROMTRANSLATION", "mat4", runtime.AdaptLegacy(m.mat4FromTranslation))
	reg.Register("MAT4.LOOKAT", "mat4", runtime.AdaptLegacy(m.mat4LookAt))
	reg.Register("MAT4.PERSPECTIVE", "mat4", runtime.AdaptLegacy(m.mat4Perspective))
	reg.Register("MAT4.ORTHO", "mat4", runtime.AdaptLegacy(m.mat4Ortho))
	reg.Register("MAT4.MULTIPLY", "mat4", runtime.AdaptLegacy(m.mat4Multiply))
	reg.Register("MAT4.INVERSE", "mat4", runtime.AdaptLegacy(m.mat4Inverse))
	reg.Register("MAT4.TRANSPOSE", "mat4", runtime.AdaptLegacy(m.mat4Transpose))
	reg.Register("MAT4.GETELEMENT", "mat4", runtime.AdaptLegacy(m.mat4GetElement))
	reg.Register("MAT4.TRANSFORMX", "mat4", runtime.AdaptLegacy(m.mat4TransformX))
	reg.Register("MAT4.TRANSFORMY", "mat4", runtime.AdaptLegacy(m.mat4TransformY))
	reg.Register("MAT4.TRANSFORMZ", "mat4", runtime.AdaptLegacy(m.mat4TransformZ))
	reg.Register("MAT4.FREE", "mat4", runtime.AdaptLegacy(m.mat4Free))
}

func (m *Module) mat4Identity(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("MAT4.IDENTITY expects 0 arguments")
	}
	id, err := AllocMatrix(m.h, rl.MatrixIdentity())
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) mat4FromRotation(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("MAT4.FROMROTATION expects 3 arguments (rx, ry, rz radians)")
	}
	x, ok1 := argF(args[0])
	y, ok2 := argF(args[1])
	z, ok3 := argF(args[2])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("MAT4.FROMROTATION: angles must be numeric")
	}
	mat := rl.MatrixRotateXYZ(rl.Vector3{X: x, Y: y, Z: z})
	id, err := AllocMatrix(m.h, mat)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) mat4SetRotation(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("MAT4.SETROTATION expects 4 arguments (handle, rx, ry, rz radians)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("MAT4.SETROTATION: first argument must be matrix handle")
	}
	o, err := heap.Cast[*matObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argF(args[1])
	y, ok2 := argF(args[2])
	z, ok3 := argF(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("MAT4.SETROTATION: angles must be numeric")
	}
	o.m = rl.MatrixRotateXYZ(rl.Vector3{X: x, Y: y, Z: z})
	return value.Nil, nil
}

func (m *Module) mat4FromScale(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("MAT4.FROMSCALE expects 3 arguments (sx, sy, sz)")
	}
	sx, ok1 := argF(args[0])
	sy, ok2 := argF(args[1])
	sz, ok3 := argF(args[2])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("MAT4.FROMSCALE: scale must be numeric")
	}
	mat := rl.MatrixScale(sx, sy, sz)
	id, err := AllocMatrix(m.h, mat)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) mat4FromTranslation(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("MAT4.FROMTRANSLATION expects 3 arguments (x, y, z)")
	}
	x, ok1 := argF(args[0])
	y, ok2 := argF(args[1])
	z, ok3 := argF(args[2])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("MAT4.FROMTRANSLATION: translation must be numeric")
	}
	mat := rl.MatrixTranslate(x, y, z)
	id, err := AllocMatrix(m.h, mat)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) mat4LookAt(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 9 {
		return value.Nil, fmt.Errorf("MAT4.LOOKAT expects 9 arguments (ex,ey,ez, tx,ty,tz, ux,uy,uz)")
	}
	var eye, target, up rl.Vector3
	for i, tgt := range []*rl.Vector3{&eye, &target, &up} {
		base := i * 3
		x, okx := argF(args[base])
		y, oky := argF(args[base+1])
		z, okz := argF(args[base+2])
		if !okx || !oky || !okz {
			return value.Nil, fmt.Errorf("MAT4.LOOKAT: arguments must be numeric")
		}
		tgt.X, tgt.Y, tgt.Z = x, y, z
	}
	mat := rl.MatrixLookAt(eye, target, up)
	id, err := AllocMatrix(m.h, mat)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) mat4Perspective(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("MAT4.PERSPECTIVE expects 4 arguments (fovY, aspect, near, far)")
	}
	fov, ok1 := argF(args[0])
	aspect, ok2 := argF(args[1])
	near, ok3 := argF(args[2])
	far, ok4 := argF(args[3])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("MAT4.PERSPECTIVE: arguments must be numeric")
	}
	mat := rl.MatrixPerspective(fov, aspect, near, far)
	id, err := AllocMatrix(m.h, mat)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) mat4Ortho(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 6 {
		return value.Nil, fmt.Errorf("MAT4.ORTHO expects 6 arguments (left, right, bottom, top, near, far)")
	}
	var f [6]float32
	var ok bool
	for i := 0; i < 6; i++ {
		f[i], ok = argF(args[i])
		if !ok {
			return value.Nil, fmt.Errorf("MAT4.ORTHO: arguments must be numeric")
		}
	}
	mat := rl.MatrixOrtho(f[0], f[1], f[2], f[3], f[4], f[5])
	id, err := AllocMatrix(m.h, mat)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) mat4Multiply(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("MAT4.MULTIPLY expects 2 matrix handles")
	}
	a, err := m.matrixFromArgs(args, 0, "MAT4.MULTIPLY")
	if err != nil {
		return value.Nil, err
	}
	b, err := m.matrixFromArgs(args, 1, "MAT4.MULTIPLY")
	if err != nil {
		return value.Nil, err
	}
	mat := rl.MatrixMultiply(a, b)
	id, err := AllocMatrix(m.h, mat)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) mat4Inverse(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("MAT4.INVERSE expects matrix handle")
	}
	mat, err := m.matrixFromArgs(args, 0, "MAT4.INVERSE")
	if err != nil {
		return value.Nil, err
	}
	out := rl.MatrixInvert(mat)
	id, err := AllocMatrix(m.h, out)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) mat4Transpose(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("MAT4.TRANSPOSE expects matrix handle")
	}
	mat, err := m.matrixFromArgs(args, 0, "MAT4.TRANSPOSE")
	if err != nil {
		return value.Nil, err
	}
	out := rl.MatrixTranspose(mat)
	id, err := AllocMatrix(m.h, out)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) mat4GetElement(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("MAT4.GETELEMENT expects (matrix, row, col) with row,col in 0..3")
	}
	mat, err := m.matrixFromArgs(args, 0, "MAT4.GETELEMENT")
	if err != nil {
		return value.Nil, err
	}
	ri, ok1 := args[1].ToInt()
	ci, ok2 := args[2].ToInt()
	if !ok1 {
		if f, ok := argF(args[1]); ok {
			ri = int64(f)
			ok1 = true
		}
	}
	if !ok2 {
		if f, ok := argF(args[2]); ok {
			ci = int64(f)
			ok2 = true
		}
	}
	if !ok1 || !ok2 || ri < 0 || ri > 3 || ci < 0 || ci > 3 {
		return value.Nil, fmt.Errorf("MAT4.GETELEMENT: row and col must be integers 0..3")
	}
	v := matElement(mat, int32(ri), int32(ci))
	return value.FromFloat(float64(v)), nil
}

func (m *Module) mat4TransformComponent(args []value.Value, op string, axis int) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("%s expects (matrix, x, y, z)", op)
	}
	mat, err := m.matrixFromArgs(args, 0, op)
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argF(args[1])
	y, ok2 := argF(args[2])
	z, ok3 := argF(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("%s: position must be numeric", op)
	}
	out := rl.Vector3Transform(rl.Vector3{X: x, Y: y, Z: z}, mat)
	var comp float32
	switch axis {
	case 0:
		comp = out.X
	case 1:
		comp = out.Y
	default:
		comp = out.Z
	}
	return value.FromFloat(float64(comp)), nil
}

func (m *Module) mat4TransformX(args []value.Value) (value.Value, error) {
	return m.mat4TransformComponent(args, "MAT4.TRANSFORMX", 0)
}

func (m *Module) mat4TransformY(args []value.Value) (value.Value, error) {
	return m.mat4TransformComponent(args, "MAT4.TRANSFORMY", 1)
}

func (m *Module) mat4TransformZ(args []value.Value) (value.Value, error) {
	return m.mat4TransformComponent(args, "MAT4.TRANSFORMZ", 2)
}

func (m *Module) transformApplyX(args []value.Value) (value.Value, error) {
	return m.mat4TransformComponent(args, "TRANSFORM.APPLYX", 0)
}

func (m *Module) transformApplyY(args []value.Value) (value.Value, error) {
	return m.mat4TransformComponent(args, "TRANSFORM.APPLYY", 1)
}

func (m *Module) transformApplyZ(args []value.Value) (value.Value, error) {
	return m.mat4TransformComponent(args, "TRANSFORM.APPLYZ", 2)
}

func (m *Module) mat4Free(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("MAT4.FREE expects matrix handle")
	}
	if err := m.h.Free(heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}
