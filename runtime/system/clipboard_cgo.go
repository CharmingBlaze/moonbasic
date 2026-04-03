//go:build cgo

package mbsystem

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerClipboard(reg runtime.Registrar) {
	reg.Register("SYSTEM.GETCLIPBOARD", "system", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, runtime.Errorf("SYSTEM.GETCLIPBOARD expects 0 arguments")
		}
		return rt.RetString(rl.GetClipboardText()), nil
	})
	reg.Register("SYSTEM.SETCLIPBOARD", "system", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if len(args) != 1 || args[0].Kind != value.KindString {
			return value.Nil, runtime.Errorf("SYSTEM.SETCLIPBOARD expects (text$)")
		}
		text, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		rl.SetClipboardText(text)
		return value.Nil, nil
	})
}
