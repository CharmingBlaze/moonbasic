package mbutil

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

// Register implements runtime.Module.
func (m *Module) Register(r runtime.Registrar) {
	r.Register("UTIL.FILEEXISTS", "util", m.utilFileExists)
	r.Register("UTIL.ISDIR", "util", m.utilIsDir)
	r.Register("UTIL.GETFILEEXT", "util", m.utilGetFileExt)
	r.Register("UTIL.GETFILENAME", "util", m.utilGetFileName)
	r.Register("UTIL.GETFILENAMENOEXT", "util", m.utilGetFileNameNoExt)
	r.Register("UTIL.GETFILEPATH", "util", m.utilGetFilePath)
	r.Register("UTIL.GETFILESIZE", "util", m.utilGetFileSize)
	r.Register("UTIL.GETFILEMODTIME", "util", m.utilGetFileModTime)
	r.Register("UTIL.LOADTEXT", "util", m.utilLoadText)
	r.Register("UTIL.SAVETEXT", "util", m.utilSaveText)
	r.Register("UTIL.GETDIRFILES", "util", m.utilGetDirFiles)
	r.Register("UTIL.CHANGEDIR", "util", m.utilChangeDir)
	r.Register("UTIL.MAKEDIRECTORY", "util", m.utilMakeDirectory)
	r.Register("UTIL.ISFILENAMEVALID", "util", m.utilIsFileNameValid)

	// Flat spec names (manifest) → same handlers as UTIL.*.
	r.Register("FILEEXISTS", "util", m.utilFileExists)
	r.Register("DIREXISTS", "util", m.utilIsDir)
	r.Register("READALLTEXT$", "util", m.utilLoadText)
	r.Register("WRITEALLTEXT", "util", m.utilSaveText)
	r.Register("MAKEDIR", "util", m.utilMakeDirectory)
	r.Register("SETDIR", "util", m.utilChangeDir)
	r.Register("GETFILEEXT$", "util", m.utilGetFileExt)
	r.Register("GETFILENAME$", "util", m.utilGetFileName)
	r.Register("GETFILENAMENOEXT$", "util", m.utilGetFileNameNoExt)
	r.Register("GETFILEPATH$", "util", m.utilGetFilePath)
	r.Register("GETFILESIZE", "util", m.utilGetFileSize)
	r.Register("GETFILEMODTIME", "util", m.utilGetFileModTime)
	r.Register("GETFILES$", "util", m.utilGetDirFiles)

	m.registerDroppedFiles(r)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}

func (m *Module) utilFileExists(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("UTIL.FILEEXISTS expects (path$)")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	_, err = os.Stat(path)
	return value.FromBool(err == nil), nil
}

func (m *Module) utilIsDir(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("UTIL.ISDIR expects (path$)")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	fi, err := os.Stat(path)
	if err != nil {
		return value.FromBool(false), nil
	}
	return value.FromBool(fi.IsDir()), nil
}

func (m *Module) utilGetFileExt(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("UTIL.GETFILEEXT expects (path$)")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	ext := filepath.Ext(path)
	return rt.RetString(ext), nil
}

func (m *Module) utilGetFileName(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("UTIL.GETFILENAME expects (path$)")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	return rt.RetString(filepath.Base(path)), nil
}

func (m *Module) utilGetFileNameNoExt(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("UTIL.GETFILENAMENOEXT expects (path$)")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	return rt.RetString(strings.TrimSuffix(base, ext)), nil
}

func (m *Module) utilGetFilePath(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("UTIL.GETFILEPATH expects (path$)")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	return rt.RetString(filepath.Dir(path)), nil
}

func (m *Module) utilGetFileSize(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("UTIL.GETFILESIZE expects (path$)")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	fi, err := os.Stat(path)
	if err != nil {
		return value.FromInt(0), nil
	}
	return value.FromInt(fi.Size()), nil
}

func (m *Module) utilGetFileModTime(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("UTIL.GETFILEMODTIME expects (path$)")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	fi, err := os.Stat(path)
	if err != nil {
		return value.FromInt(0), nil
	}
	return value.FromInt(fi.ModTime().Unix()), nil
}

func (m *Module) utilLoadText(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("UTIL.LOADTEXT expects (path$)")
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

func (m *Module) utilSaveText(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindString || args[1].Kind != value.KindString {
		return value.Nil, runtime.Errorf("UTIL.SAVETEXT expects (path$, text$)")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	text, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	err = os.WriteFile(path, []byte(text), 0644)
	return value.Nil, err
}

func (m *Module) utilGetDirFiles(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("UTIL.GETDIRFILES expects (path$)")
	}
	dir, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return value.Nil, err
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		names = append(names, e.Name())
	}
	b, err := json.Marshal(names)
	if err != nil {
		return value.Nil, err
	}
	return rt.RetString(string(b)), nil
}

func (m *Module) utilChangeDir(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("UTIL.CHANGEDIR expects (path$)")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	err = os.Chdir(path)
	return value.FromBool(err == nil), nil
}

func (m *Module) utilMakeDirectory(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("UTIL.MAKEDIRECTORY expects (path$)")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	err = os.MkdirAll(path, 0755)
	return value.FromBool(err == nil), nil
}

// utilIsFileNameValid mirrors the active rules in raylib's IsFileNameValid (invalid glyphs + not all '.').
func (m *Module) utilIsFileNameValid(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("UTIL.ISFILENAMEVALID expects (name$)")
	}
	name, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	if name == "" {
		return value.FromBool(false), nil
	}
	allDots := true
	for _, r := range name {
		if r < 32 {
			return value.FromBool(false), nil
		}
		switch r {
		case '<', '>', ':', '"', '/', '\\', '|', '?', '*':
			return value.FromBool(false), nil
		}
		if r != '.' {
			allDots = false
		}
	}
	return value.FromBool(!allDots), nil
}
