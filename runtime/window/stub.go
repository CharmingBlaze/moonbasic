//go:build !cgo && !windows

package window

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

const stubHint = "WINDOW/RENDER natives require CGO: set CGO_ENABLED=1 and install a C compiler (e.g. MinGW on Windows, gcc on Linux), then rebuild"

func stubFn(hint string) func(string) runtime.BuiltinFn {
	return func(name string) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			_ = rt
			return value.Nil, fmt.Errorf("%s: %s", name, hint)
		}
	}
}

func (m *Module) Register(reg runtime.Registrar) {
	stub := stubFn(stubHint)
	reg.Register("WINDOW.OPEN", "window", stub("WINDOW.OPEN"))
	reg.Register("WINDOW.CANOPEN", "window", stub("WINDOW.CANOPEN"))
	reg.Register("WINDOW.SETFPS", "window", stub("WINDOW.SETFPS"))
	reg.Register("WINDOW.CLOSE", "window", stub("WINDOW.CLOSE"))
	reg.Register("WINDOW.SHOULDCLOSE", "window", stub("WINDOW.SHOULDCLOSE"))
	reg.Register("RENDER.CLEAR", "render", stub("RENDER.CLEAR"))
	reg.Register("RENDER.FRAME", "render", stub("RENDER.FRAME"))
	m.registerRenderAdvanced(reg)
	m.registerPostCommands(reg)
	m.registerEffectCommands(reg)
	m.registerComputeShaderCommands(reg)
	m.registerDecalCommands(reg)
	m.registerWindowStateCommands(reg)
	m.registerWindowMetricsCommands(reg)
	m.registerWindowPlacementCommands(reg)
	m.registerAutomationCommands(reg)
}

func (m *Module) Shutdown() {
	// nothing to release without Raylib
}
