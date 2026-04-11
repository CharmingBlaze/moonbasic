//go:build !cgo && !windows

package window

import (
	"fmt"

	"moonbasic/internal/driver"
	"moonbasic/runtime"
	"moonbasic/vm/value"
)

const stubHint = "WINDOW/RENDER: place the raylib shared library next to the executable (see MOONBASIC_DRIVER / internal/raylibpurego), or set CGO_ENABLED=1 and rebuild with a C compiler"

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
	reg.Register("WINDOW.OPEN", "window", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if m.usePuregoDLL() {
			if err := driver.CheckWindow(m.driverSel); err != nil {
				return value.Nil, err
			}
			return m.puregoWOpen(rt, args...)
		}
		return stub("WINDOW.OPEN")(rt, args...)
	})
	reg.Register("WINDOW.CANOPEN", "window", m.wCanOpenStub)
	reg.Register("WINDOW.SETFPS", "window", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if m.usePuregoDLL() {
			return m.puregoWSetFPS(rt, args...)
		}
		return stub("WINDOW.SETFPS")(rt, args...)
	})
	reg.Register("WINDOW.CLOSE", "window", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if m.usePuregoDLL() {
			return m.puregoWClose(rt, args...)
		}
		return stub("WINDOW.CLOSE")(rt, args...)
	})
	reg.Register("WINDOW.SHOULDCLOSE", "window", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if m.usePuregoDLL() {
			return m.puregoWShouldClose(rt, args...)
		}
		return stub("WINDOW.SHOULDCLOSE")(rt, args...)
	})
	reg.Register("RENDER.CLEAR", "render", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if m.usePuregoDLL() {
			return m.puregoRClear(rt, args...)
		}
		return stub("RENDER.CLEAR")(rt, args...)
	})
	reg.Register("RENDER.FRAME", "render", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if m.usePuregoDLL() {
			return m.puregoRFrame(rt, args...)
		}
		return stub("RENDER.FRAME")(rt, args...)
	})
	reg.Register("RENDER.BEGIN3D", "render", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if m.usePuregoDLL() {
			return m.rBegin3D(rt, args...)
		}
		return stub("RENDER.BEGIN3D")(rt, args...)
	})
	reg.Register("RENDER.END3D", "render", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if m.usePuregoDLL() {
			return m.rEnd3D(rt, args...)
		}
		return stub("RENDER.END3D")(rt, args...)
	})
	m.registerRenderAdvanced(reg)
	m.registerPostCommands(reg)
	m.registerEffectCommands(reg)
	m.registerComputeShaderCommands(reg)
	m.registerDecalCommands(reg)
	m.registerWindowStateCommands(reg)
	m.registerWindowMetricsCommands(reg)
	m.registerWindowPlacementCommands(reg)
	m.registerAutomationCommands(reg)
	m.registerBlitzSysCommands(reg)
}

func (m *Module) wCanOpenStub(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("WINDOW.CANOPEN expects 3 arguments (width, height, title)")
	}
	w, okw := argInt(args[0])
	h, okh := argInt(args[1])
	if !okw || !okh {
		return value.Nil, fmt.Errorf("WINDOW.CANOPEN: width and height must be numeric")
	}
	if args[2].Kind != value.KindString {
		return value.Nil, fmt.Errorf("WINDOW.CANOPEN: title must be a string")
	}
	title, err := rt.ArgString(args, 2)
	if err != nil {
		return value.Nil, err
	}
	ok := w > 0 && h > 0 && title != ""
	return value.FromBool(ok), nil
}

func (m *Module) Shutdown() {
	if m.usePuregoDLL() {
		m.puregoShutdown()
	}
}
