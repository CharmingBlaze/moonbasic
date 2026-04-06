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
		id, err := m.h.Alloc(&shaderObj{sh: sh})
		if err != nil {
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	})
}
