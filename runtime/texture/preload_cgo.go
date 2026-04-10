//go:build cgo || (windows && !cgo)

package texture

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// preloadHandles tracks textures loaded by LEVEL.PRELOAD for RENDER.CLEARCACHE.
var preloadHandles []heap.Handle

func registerTexturePreloadCmds(m *Module, r runtime.Registrar) {
	r.Register("LEVEL.PRELOAD", "entity", m.levelPreloadDir)
	r.Register("RENDER.CLEARCACHE", "render", runtime.AdaptLegacy(m.renderClearTexturePreload))
}

func (m *Module) levelPreloadDir(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("LEVEL.PRELOAD: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("LEVEL.PRELOAD expects (directoryPath$)")
	}
	dir, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	dir = strings.TrimSpace(dir)
	if dir == "" {
		return value.Nil, fmt.Errorf("LEVEL.PRELOAD: path required")
	}
	st, err := os.Stat(dir)
	if err != nil {
		return value.Nil, fmt.Errorf("LEVEL.PRELOAD: %w", err)
	}
	if !st.IsDir() {
		return value.Nil, fmt.Errorf("LEVEL.PRELOAD: not a directory: %s", dir)
	}
	var n int32
	walkErr := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".png", ".jpg", ".jpeg", ".bmp", ".tga", ".gif", ".hdr", ".dds", ".ktx", ".ktx2":
		default:
			return nil
		}
		h, err := m.TexLoadPath(path, 1)
		if err != nil {
			return err
		}
		preloadHandles = append(preloadHandles, h)
		n++
		return nil
	})
	if walkErr != nil {
		return value.Nil, walkErr
	}
	return value.FromInt(int64(n)), nil
}

func (m *Module) renderClearTexturePreload(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("RENDER.CLEARCACHE: heap not bound")
	}
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("RENDER.CLEARCACHE expects no arguments")
	}
	for _, h := range preloadHandles {
		m.h.Free(h)
	}
	preloadHandles = preloadHandles[:0]
	return value.Nil, nil
}
