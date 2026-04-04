//go:build !cgo

package window

import "moonbasic/runtime"

func (m *Module) registerPostCommands(r runtime.Registrar) {
	stub := stubFn(stubHint)
	for _, name := range []string{"POST.ADD", "POST.SETPARAM", "POST.ADDSHADER"} {
		r.Register(name, "post", stub(name))
	}
}
