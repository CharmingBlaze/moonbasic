//go:build cgo || (windows && !cgo)

package mbfile

import (
	"os"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerFileExtras(r runtime.Registrar) {
	r.Register("FILE.EXISTS", "file", m.fileExists)
	r.Register("FILE.READALLTEXT", "file", m.fileReadAllText)
	r.Register("FILE.WRITEALLTEXT", "file", m.fileWriteAllText)
}

func (m *Module) fileExists(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("FILE.EXISTS expects (path)")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return value.FromBool(false), nil
		}
		return value.Nil, err
	}
	return value.FromBool(true), nil
}

func (m *Module) fileReadAllText(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("FILE.READALLTEXT expects (path)")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return value.Nil, err
	}
	return rt.RetString(string(b)), nil
}

func (m *Module) fileWriteAllText(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindString || args[1].Kind != value.KindString {
		return value.Nil, runtime.Errorf("FILE.WRITEALLTEXT expects (path, text)")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	text, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}
