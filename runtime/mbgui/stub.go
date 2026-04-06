//go:build !cgo && !windows

package mbgui

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

const hint = "GUI.* requires CGO (raygui-go wraps C): set CGO_ENABLED=1 and a C toolchain, then rebuild (on Windows with CGO off, use the built-in minimal Raylib GUI instead)"

var guiStubNames = []string{
	"GUI.ENABLE", "GUI.DISABLE", "GUI.LOCK", "GUI.UNLOCK", "GUI.ISLOCKED",
	"GUI.SETALPHA", "GUI.SETSTATE", "GUI.GETSTATE", "GUI.SETFONT",
	"GUI.SETSTYLE", "GUI.GETSTYLE", "GUI.GETCOLOR",
	"GUI.SETCOLOR", "GUI.SETTEXTSIZE", "GUI.SETTEXTSPACING", "GUI.SETTEXTLINEHEIGHT",
	"GUI.SETTEXTWRAP", "GUI.SETTEXTALIGN", "GUI.SETTEXTALIGNVERT", "GUI.GETTEXTSIZE", "GUI.THEMEAPPLY", "GUI.THEMENAMES$",
	"GUI.LOADSTYLE", "GUI.LOADDEFAULTSTYLE", "GUI.LOADSTYLEMEM", "GUI.LOADICONS", "GUI.LOADICONSMEM",
	"GUI.WINDOWBOX", "GUI.GROUPBOX", "GUI.LINE", "GUI.PANEL", "GUI.TABBAR", "GUI.SCROLLPANEL",
	"GUI.LABEL", "GUI.BUTTON", "GUI.LABELBUTTON", "GUI.TOGGLE", "GUI.TOGGLEGROUP", "GUI.TOGGLEGROUPAT",
	"GUI.TOGGLESLIDER", "GUI.CHECKBOX", "GUI.COMBOBOX", "GUI.DROPDOWNBOX",
	"GUI.TEXTBOX", "GUI.SPINNER", "GUI.VALUEBOX", "GUI.VALUEBOXFLOAT", "GUI.VALUEBOXFLOATTEXT$",
	"GUI.SLIDER", "GUI.SLIDERBAR", "GUI.PROGRESSBAR", "GUI.SCROLLBAR", "GUI.STATUSBAR", "GUI.DUMMYREC",
	"GUI.LISTVIEW", "GUI.LISTVIEWEX",
	"GUI.COLORPANEL", "GUI.COLORBARALPHA", "GUI.COLORBARHUE", "GUI.COLORPICKER",
	"GUI.COLORPICKERHSV", "GUI.COLORPANELHSV",
	"GUI.MESSAGEBOX", "GUI.TEXTINPUTBOX", "GUI.TEXTINPUTLAST$", "GUI.GRID",
	"GUI.ENABLETOOLTIP", "GUI.DISABLETOOLTIP", "GUI.SETTOOLTIP",
	"GUI.ICONTEXT", "GUI.DRAWICON", "GUI.SETICONSCALE", "GUI.GETTEXTWIDTH",
	"GUI.FADE", "GUI.DRAWRECTANGLE", "GUI.DRAWTEXT", "GUI.GETTEXTBOUNDS",
}

// Register implements runtime.Module.
func (m *Module) Register(reg runtime.Registrar) {
	for _, name := range guiStubNames {
		n := name
		reg.Register(n, "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			_ = rt
			_ = args
			return value.Nil, fmt.Errorf("%s: %s", n, hint)
		})
	}
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
