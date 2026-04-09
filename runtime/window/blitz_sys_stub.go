//go:build !cgo && !windows

package window

import "moonbasic/runtime"

func (m *Module) registerBlitzSysCommands(reg runtime.Registrar) {
	stub := stubFn(stubHint)
	reg.Register("AppTitle", "window", stub("AppTitle"))
	reg.Register("Graphics3D", "window", stub("Graphics3D"))
	reg.Register("ActiveShader", "window", stub("ActiveShader"))
	reg.Register("SetMSAA", "window", stub("SetMSAA"))
	reg.Register("SetSSAO", "window", stub("SetSSAO"))
	reg.Register("SetPostProcess", "window", stub("SetPostProcess"))
	reg.Register("SetBloom", "window", stub("SetBloom"))
}
