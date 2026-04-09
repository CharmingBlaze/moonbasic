//go:build !cgo && !windows

package window

import "moonbasic/runtime"

func (m *Module) registerBlitzDisplayQueries(reg runtime.Registrar) {
	stub := stubFn(stubHint)
	for _, n := range []string{
		"WindowWidth", "WindowHeight", "ScreenWidth", "ScreenHeight",
		"GraphicsWidth", "GraphicsHeight", "GraphicsDepth",
		"AvailVidMem", "TotalVidMem",
	} {
		name := n
		reg.Register(n, "window", stub(name))
	}
}
