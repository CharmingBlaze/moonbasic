//go:build !cgo && !windows

package window

import "moonbasic/runtime"

func (m *Module) registerComputeShaderCommands(r runtime.Registrar) {
	stub := stubFn(stubHint)
	for _, name := range []string{
		"COMPUTESHADER.LOAD", "COMPUTESHADER.FREE", "COMPUTESHADER.BUFFERMAKE",
		"COMPUTESHADER.BUFFERFREE", "COMPUTESHADER.SETBUFFER", "COMPUTESHADER.SETINT",
		"COMPUTESHADER.SETFLOAT", "COMPUTESHADER.DISPATCH",
	} {
		r.Register(name, "compute", stub(name))
	}
}
