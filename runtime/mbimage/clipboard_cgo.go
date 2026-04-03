//go:build cgo

package mbimage

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerClipboardImage(m *Module, reg runtime.Registrar) {
	reg.Register("CLIPBOARD.GETIMAGE", "clipboard", m.clipboardGetImage)
}

func (m *Module) clipboardGetImage(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.h == nil {
		return value.Nil, runtime.Errorf("CLIPBOARD.GETIMAGE: heap not bound")
	}
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("CLIPBOARD.GETIMAGE expects 0 arguments")
	}
	img := rl.GetClipboardImage()
	if !rl.IsImageValid(&img) {
		return value.Nil, fmt.Errorf("CLIPBOARD.GETIMAGE: no image on clipboard (or unsupported on this platform)")
	}
	cp := rl.ImageCopy(&img)
	if cp == nil || !rl.IsImageValid(cp) {
		return value.Nil, fmt.Errorf("CLIPBOARD.GETIMAGE: could not copy clipboard image")
	}
	return m.allocImage(cp, "CLIPBOARD.GETIMAGE")
}
