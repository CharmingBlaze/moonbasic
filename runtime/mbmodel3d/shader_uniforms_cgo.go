//go:build cgo || (windows && !cgo)

package mbmodel3d

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/runtime/texture"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func heapHandle(v value.Value) heap.Handle { return heap.Handle(v.IVal) }

func registerShaderUniformCmds(m *Module, reg runtime.Registrar) {
	reg.Register("SHADER.FREE", "shader", runtime.AdaptLegacy(m.shaderFree))
	reg.Register("SHADER.GETLOC", "shader", m.shaderGetLoc)
	reg.Register("SHADER.SETFLOAT", "shader", m.shaderSetFloat)
	reg.Register("SHADER.SETVEC2", "shader", m.shaderSetVec2)
	reg.Register("SHADER.SETVEC3", "shader", m.shaderSetVec3)
	reg.Register("SHADER.SETVEC4", "shader", m.shaderSetVec4)
	reg.Register("SHADER.SETINT", "shader", m.shaderSetInt)
	reg.Register("SHADER.SETTEXTURE", "shader", m.shaderSetTexture)
}

func (m *Module) shaderFree(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("SHADER.FREE expects shader handle")
	}
	if err := m.h.Free(heapHandle(args[0])); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) shaderGetLoc(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("SHADER.GETLOC expects (shader, uniformName$)")
	}
	sh, err := m.getShader(args, 0, "SHADER.GETLOC")
	if err != nil {
		return value.Nil, err
	}
	name, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	loc := rl.GetShaderLocation(sh.sh, name)
	return value.FromInt(int64(loc)), nil
}

func (m *Module) shaderSetFloat(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	sh, loc, err := m.resolveUniform(rt, args, "SHADER.SETFLOAT", 3)
	if err != nil {
		return value.Nil, err
	}
	v, ok := argF32(args[2])
	if !ok {
		return value.Nil, fmt.Errorf("SHADER.SETFLOAT: value must be numeric")
	}
	m.u1[0] = v
	rl.SetShaderValue(sh.sh, loc, m.u1, rl.ShaderUniformFloat)
	return value.Nil, nil
}

func (m *Module) shaderSetVec2(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	sh, loc, err := m.resolveUniform(rt, args, "SHADER.SETVEC2", 4)
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argF32(args[2])
	y, ok2 := argF32(args[3])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("SHADER.SETVEC2: x,y must be numeric")
	}
	m.u2[0] = x
	m.u2[1] = y
	rl.SetShaderValue(sh.sh, loc, m.u2, rl.ShaderUniformVec2)
	return value.Nil, nil
}

func (m *Module) shaderSetVec3(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	sh, loc, err := m.resolveUniform(rt, args, "SHADER.SETVEC3", 5)
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argF32(args[2])
	y, ok2 := argF32(args[3])
	z, ok3 := argF32(args[4])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("SHADER.SETVEC3: x,y,z must be numeric")
	}
	m.u3[0] = x
	m.u3[1] = y
	m.u3[2] = z
	rl.SetShaderValue(sh.sh, loc, m.u3, rl.ShaderUniformVec3)
	return value.Nil, nil
}

func (m *Module) shaderSetVec4(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	sh, loc, err := m.resolveUniform(rt, args, "SHADER.SETVEC4", 6)
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argF32(args[2])
	y, ok2 := argF32(args[3])
	z, ok3 := argF32(args[4])
	w, ok4 := argF32(args[5])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("SHADER.SETVEC4: components must be numeric")
	}
	m.u4[0] = x
	m.u4[1] = y
	m.u4[2] = z
	m.u4[3] = w
	rl.SetShaderValue(sh.sh, loc, m.u4, rl.ShaderUniformVec4)
	return value.Nil, nil
}

func (m *Module) shaderSetInt(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	sh, loc, err := m.resolveUniform(rt, args, "SHADER.SETINT", 3)
	if err != nil {
		return value.Nil, err
	}
	var iv int32
	if i, ok := args[2].ToInt(); ok {
		iv = int32(i)
	} else if f, ok := args[2].ToFloat(); ok {
		iv = int32(f)
	} else {
		return value.Nil, fmt.Errorf("SHADER.SETINT: value must be numeric")
	}
	m.u1[0] = float32(iv)
	rl.SetShaderValue(sh.sh, loc, m.u1, rl.ShaderUniformInt)
	return value.Nil, nil
}

func (m *Module) shaderSetTexture(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 3 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString || args[2].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("SHADER.SETTEXTURE expects (shader, uniformName$, texture)")
	}
	sh, err := m.getShader(args, 0, "SHADER.SETTEXTURE")
	if err != nil {
		return value.Nil, err
	}
	name, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	loc := rl.GetShaderLocation(sh.sh, name)
	if loc < 0 {
		return value.Nil, fmt.Errorf("SHADER.SETTEXTURE: uniform %q not found", name)
	}
	tex, err := texture.ForBinding(m.h, heapHandle(args[2]))
	if err != nil {
		return value.Nil, err
	}
	rl.SetShaderValueTexture(sh.sh, loc, tex)
	return value.Nil, nil
}

func (m *Module) resolveUniform(rt *runtime.Runtime, args []value.Value, op string, want int) (*shaderObj, int32, error) {
	if err := m.requireHeap(); err != nil {
		return nil, 0, err
	}
	if len(args) != want || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return nil, 0, fmt.Errorf("%s: wrong arguments (see docs)", op)
	}
	sh, err := m.getShader(args, 0, op)
	if err != nil {
		return nil, 0, err
	}
	name, err := rt.ArgString(args, 1)
	if err != nil {
		return nil, 0, err
	}
	loc := rl.GetShaderLocation(sh.sh, name)
	if loc < 0 {
		return nil, 0, fmt.Errorf("%s: uniform %q not found", op, name)
	}
	return sh, loc, nil
}

func argF32(v value.Value) (float32, bool) {
	if f, ok := v.ToFloat(); ok {
		return float32(f), true
	}
	if i, ok := v.ToInt(); ok {
		return float32(i), true
	}
	return 0, false
}
