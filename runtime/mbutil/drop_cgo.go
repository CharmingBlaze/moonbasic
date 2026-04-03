//go:build cgo

package mbutil

import (
	"encoding/json"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerDroppedFiles(r runtime.Registrar) {
	r.Register("UTIL.ISFILEDROPPED", "util", runtime.AdaptLegacy(m.utilIsFileDropped))
	r.Register("UTIL.GETDROPPEDFILES", "util", m.utilGetDroppedFiles)
	r.Register("UTIL.CLEARDROPPEDFILES", "util", runtime.AdaptLegacy(m.utilClearDroppedFiles))
}

func (m *Module) utilIsFileDropped(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, runtime.Errorf("UTIL.ISFILEDROPPED expects 0 arguments")
	}
	return value.FromBool(rl.IsFileDropped()), nil
}

func (m *Module) utilGetDroppedFiles(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, runtime.Errorf("UTIL.GETDROPPEDFILES expects 0 arguments")
	}
	paths := rl.LoadDroppedFiles()
	b, err := json.Marshal(paths)
	if err != nil {
		return value.Nil, err
	}
	rl.UnloadDroppedFiles()
	return rt.RetString(string(b)), nil
}

func (m *Module) utilClearDroppedFiles(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, runtime.Errorf("UTIL.CLEARDROPPEDFILES expects 0 arguments")
	}
	rl.UnloadDroppedFiles()
	return value.Nil, nil
}
