//go:build cgo || (windows && !cgo)

package window

import (
	"fmt"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

type computeShaderObj struct {
	programID uint32
	release   heap.ReleaseOnce
}

func (c *computeShaderObj) TypeName() string { return "ComputeShader" }

func (c *computeShaderObj) TypeTag() uint16 { return heap.TagComputeShader }

func (c *computeShaderObj) Free() {
	c.release.Do(func() {
		if c.programID != 0 {
			rl.UnloadShaderProgram(c.programID)
			c.programID = 0
		}
	})
}

type shaderBufferObj struct {
	id      uint32
	release heap.ReleaseOnce
}

func (s *shaderBufferObj) TypeName() string { return "ShaderBuffer" }

func (s *shaderBufferObj) TypeTag() uint16 { return heap.TagShaderBuffer }

func (s *shaderBufferObj) Free() {
	s.release.Do(func() {
		if s.id != 0 {
			rl.UnloadShaderBuffer(s.id)
			s.id = 0
		}
	})
}

func (m *Module) registerComputeShaderCommands(r runtime.Registrar) {
	r.Register("COMPUTESHADER.LOAD", "compute", m.csLoad)
	r.Register("COMPUTESHADER.FREE", "compute", m.csFree)
	r.Register("COMPUTESHADER.BUFFERMAKE", "compute", m.csBufferMake)
	r.Register("COMPUTESHADER.BUFFERFREE", "compute", m.csBufferFree)
	r.Register("COMPUTESHADER.SETBUFFER", "compute", m.csSetBuffer)
	r.Register("COMPUTESHADER.SETINT", "compute", m.csSetInt)
	r.Register("COMPUTESHADER.SETFLOAT", "compute", m.csSetFloat)
	r.Register("COMPUTESHADER.DISPATCH", "compute", m.csDispatch)
}

func (m *Module) requireHeapCS(rt *runtime.Runtime) (*heap.Store, error) {
	if rt != nil && rt.Heap != nil {
		return rt.Heap, nil
	}
	if m.h != nil {
		return m.h, nil
	}
	return nil, fmt.Errorf("COMPUTESHADER.*: heap not bound")
}

func (m *Module) csLoad(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeapCS(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("COMPUTESHADER.LOAD expects 1 string path")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	srcBytes, err := os.ReadFile(path)
	if err != nil {
		return value.Nil, fmt.Errorf("COMPUTESHADER.LOAD: %w", err)
	}
	src := string(srcBytes)
	csID := rl.CompileShader(src, rl.ComputeShader)
	if csID == 0 {
		return value.Nil, fmt.Errorf("COMPUTESHADER.LOAD: compile failed for %q", path)
	}
	prog := rl.LoadComputeShaderProgram(csID)
	if prog == 0 {
		return value.Nil, fmt.Errorf("COMPUTESHADER.LOAD: link failed for %q", path)
	}
	id, err := h.Alloc(&computeShaderObj{programID: prog})
	if err != nil {
		rl.UnloadShaderProgram(prog)
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) csFree(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeapCS(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("COMPUTESHADER.FREE expects compute shader handle")
	}
	_ = h.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}

func (m *Module) csBufferMake(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeapCS(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("COMPUTESHADER.BUFFERMAKE expects 1 argument (size in bytes)")
	}
	var sz uint32
	if i, ok := args[0].ToInt(); ok && i > 0 {
		sz = uint32(i)
	} else if f, ok := args[0].ToFloat(); ok && f > 0 {
		sz = uint32(f)
	} else {
		return value.Nil, fmt.Errorf("COMPUTESHADER.BUFFERMAKE: size must be a positive number")
	}
	bufID := rl.LoadShaderBuffer(sz, nil, rl.DynamicCopy)
	if bufID == 0 {
		return value.Nil, fmt.Errorf("COMPUTESHADER.BUFFERMAKE: allocation failed")
	}
	hid, err := h.Alloc(&shaderBufferObj{id: bufID})
	if err != nil {
		rl.UnloadShaderBuffer(bufID)
		return value.Nil, err
	}
	return value.FromHandle(hid), nil
}

func (m *Module) csBufferFree(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeapCS(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("COMPUTESHADER.BUFFERFREE expects buffer handle")
	}
	_ = h.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}

func (m *Module) csSetBuffer(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if rt == nil || rt.Heap == nil {
		return value.Nil, fmt.Errorf("COMPUTESHADER.SETBUFFER: heap not available")
	}
	if len(args) != 3 || args[0].Kind != value.KindHandle || args[2].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("COMPUTESHADER.SETBUFFER expects (compute$, bindingIndex, bufferHandle)")
	}
	var bind uint32
	if i, ok := args[1].ToInt(); ok {
		bind = uint32(i)
	} else if f, ok := args[1].ToFloat(); ok {
		bind = uint32(f)
	} else {
		return value.Nil, fmt.Errorf("COMPUTESHADER.SETBUFFER: binding must be numeric (GL binding point)")
	}
	store := rt.Heap
	if _, err := heap.Cast[*computeShaderObj](store, heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, fmt.Errorf("COMPUTESHADER.SETBUFFER: %w", err)
	}
	bufO, err := heap.Cast[*shaderBufferObj](store, heap.Handle(args[2].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("COMPUTESHADER.SETBUFFER: %w", err)
	}
	rl.BindShaderBuffer(bufO.id, bind)
	return value.Nil, nil
}

func (m *Module) csSetInt(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if rt == nil || rt.Heap == nil {
		return value.Nil, fmt.Errorf("COMPUTESHADER.SETINT: heap not available")
	}
	if len(args) != 3 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("COMPUTESHADER.SETINT expects (compute$, uniformName$, value)")
	}
	co, err := heap.Cast[*computeShaderObj](rt.Heap, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	name, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	var iv int32
	if i, ok := args[2].ToInt(); ok {
		iv = int32(i)
	} else if f, ok := args[2].ToFloat(); ok {
		iv = int32(f)
	} else {
		return value.Nil, fmt.Errorf("COMPUTESHADER.SETINT: value must be numeric")
	}
	rl.EnableShader(co.programID)
	loc := rl.GetLocationUniform(co.programID, name)
	if loc < 0 {
		rl.DisableShader()
		return value.Nil, fmt.Errorf("COMPUTESHADER.SETINT: uniform %q not found", name)
	}
	rl.SetUniform(loc, []int32{iv}, int32(rl.ShaderUniformInt), 1)
	rl.DisableShader()
	return value.Nil, nil
}

func (m *Module) csSetFloat(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if rt == nil || rt.Heap == nil {
		return value.Nil, fmt.Errorf("COMPUTESHADER.SETFLOAT: heap not available")
	}
	if len(args) != 3 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("COMPUTESHADER.SETFLOAT expects (compute$, uniformName$, value)")
	}
	co, err := heap.Cast[*computeShaderObj](rt.Heap, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	name, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	var fv float32
	if f, ok := args[2].ToFloat(); ok {
		fv = float32(f)
	} else if i, ok := args[2].ToInt(); ok {
		fv = float32(i)
	} else {
		return value.Nil, fmt.Errorf("COMPUTESHADER.SETFLOAT: value must be numeric")
	}
	rl.EnableShader(co.programID)
	loc := rl.GetLocationUniform(co.programID, name)
	if loc < 0 {
		rl.DisableShader()
		return value.Nil, fmt.Errorf("COMPUTESHADER.SETFLOAT: uniform %q not found", name)
	}
	rl.SetUniform(loc, []float32{fv}, int32(rl.ShaderUniformFloat), 1)
	rl.DisableShader()
	return value.Nil, nil
}

func (m *Module) csDispatch(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if rt == nil || rt.Heap == nil {
		return value.Nil, fmt.Errorf("COMPUTESHADER.DISPATCH: heap not available")
	}
	if len(args) != 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("COMPUTESHADER.DISPATCH expects (compute$, gx, gy, gz)")
	}
	co, err := heap.Cast[*computeShaderObj](rt.Heap, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	var gx, gy, gz uint32
	for i, a := range args[1:4] {
		if v, ok := a.ToInt(); ok && v > 0 {
			switch i {
			case 0:
				gx = uint32(v)
			case 1:
				gy = uint32(v)
			case 2:
				gz = uint32(v)
			}
		} else if f, ok := a.ToFloat(); ok && f > 0 {
			switch i {
			case 0:
				gx = uint32(f)
			case 1:
				gy = uint32(f)
			case 2:
				gz = uint32(f)
			}
		} else {
			return value.Nil, fmt.Errorf("COMPUTESHADER.DISPATCH: group sizes must be positive numbers")
		}
	}
	rl.EnableShader(co.programID)
	rl.ComputeShaderDispatch(gx, gy, gz)
	rl.DisableShader()
	return value.Nil, nil
}
