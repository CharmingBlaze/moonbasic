//go:build !cgo && windows

package mbgui

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

const puregoNeedRaygui = "requires CGO raygui (rebuild with CGO_ENABLED=1)"

func registerPuregoUnimplemented(reg runtime.Registrar) {
	names := []string{
		"GUI.STATUSBAR", "GUI.DUMMYREC",
		"GUI.TEXTBOX", "GUI.SPINNER", "GUI.VALUEBOX", "GUI.VALUEBOXFLOAT", "GUI.VALUEBOXFLOATTEXT",
		"GUI.LISTVIEW", "GUI.LISTVIEWEX",
		"GUI.COLORPANEL", "GUI.COLORBARALPHA", "GUI.COLORBARHUE", "GUI.COLORPICKER",
		"GUI.COLORPICKERHSV", "GUI.COLORPANELHSV",
		"GUI.MESSAGEBOX", "GUI.TEXTINPUTBOX", "GUI.TEXTINPUTLAST", "GUI.GRID",
	}
	for _, name := range names {
		n := name
		reg.Register(n, "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			_ = rt
			_ = args
			return value.Nil, fmt.Errorf("%s: %s", n, puregoNeedRaygui)
		})
	}
}
