//go:build cgo || (windows && !cgo)

package mbmodel3d

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerShaderCmds(m *Module, reg runtime.Registrar) {
	reg.Register("SHADER.LOAD", "shader", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 2 || args[0].Kind != value.KindString || args[1].Kind != value.KindString {
			return value.Nil, fmt.Errorf("SHADER.LOAD expects vertexShaderPath$, fragmentShaderPath$")
		}
		vsPath, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		fsPath, err := rt.ArgString(args, 1)
		if err != nil {
			return value.Nil, err
		}
		sh := rl.LoadShader(vsPath, fsPath)
		obj := &shaderObj{sh: sh}
		obj.setFinalizer()
		id, err := m.h.Alloc(obj)
		if err != nil {
			return value.Nil, err
		}
		// Also cache to manager if available
		return value.FromHandle(id), nil
	})

	reg.Register("SHADER.SETFLOAT", "shader", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 3 { return value.Nil, fmt.Errorf("SHADER.SETFLOAT expects shaderID, name$, value#") }
		return value.Nil, nil
	})
	
	reg.Register("SHADER.SETVECTOR", "shader", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 5 { return value.Nil, fmt.Errorf("SHADER.SETVECTOR expects shaderID, name$, x, y, z") }
		return value.Nil, nil
	})

	reg.Register("SHADER.SETTEXTURE", "shader", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 3 { return value.Nil, fmt.Errorf("SHADER.SETTEXTURE expects shaderID, name$, texHandle") }
		return value.Nil, nil
	})

	reg.Register("SHADER.FREE", "shader", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return value.Nil, nil
	})

	bindConstants(reg)
}

func bindConstants(reg runtime.Registrar) {
	reg.Register("SHADER_PBR_LIT", "shader", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return rt.RetInt(1), nil })
	reg.Register("SHADER_PS1_RETRO", "shader", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return rt.RetInt(2), nil })
	reg.Register("SHADER_CEL_STYLED", "shader", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return rt.RetInt(3), nil })
	reg.Register("SHADER_WATER_PROCEDURAL", "shader", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return rt.RetInt(4), nil })
	reg.Register("PP_BLOOM", "shader", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return rt.RetInt(101), nil })
	reg.Register("PP_CRT_SCANLINES", "shader", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return rt.RetInt(102), nil })
	reg.Register("PP_PIXELATE", "shader", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return rt.RetInt(103), nil })
}
