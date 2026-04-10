//go:build !cgo && !windows

package terrain

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerTerrain(m *Module, r runtime.Registrar) {
	hint := func(name string) func(*runtime.Runtime, ...value.Value) (value.Value, error) {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			return value.Nil, fmt.Errorf("%s requires CGO and Raylib (set CGO_ENABLED=1)", name)
		}
	}
	r.Register("TERRAIN.MAKE", "terrain", hint("TERRAIN.MAKE"))
	r.Register("TERRAIN.FREE", "terrain", hint("TERRAIN.FREE"))
	r.Register("TERRAIN.SETPOS", "terrain", hint("TERRAIN.SETPOS"))
	r.Register("TERRAIN.SETCHUNKSIZE", "terrain", hint("TERRAIN.SETCHUNKSIZE"))
	r.Register("TERRAIN.FILLPERLIN", "terrain", hint("TERRAIN.FILLPERLIN"))
	r.Register("TERRAIN.FILLFLAT", "terrain", hint("TERRAIN.FILLFLAT"))
	r.Register("TERRAIN.GETHEIGHT", "terrain", hint("TERRAIN.GETHEIGHT"))
	r.Register("TERRAIN.GETSLOPE", "terrain", hint("TERRAIN.GETSLOPE"))
	r.Register("TERRAIN.RAISE", "terrain", hint("TERRAIN.RAISE"))
	r.Register("TERRAIN.LOWER", "terrain", hint("TERRAIN.LOWER"))
	r.Register("TERRAIN.DRAW", "terrain", hint("TERRAIN.DRAW"))
	r.Register("TERRAIN.PLACE", "terrain", hint("TERRAIN.PLACE"))
	r.Register("TERRAIN.SNAPY", "terrain", hint("TERRAIN.SNAPY"))
	r.Register("CHUNK.GENERATE", "chunk", hint("CHUNK.GENERATE"))
	r.Register("CHUNK.COUNT", "chunk", hint("CHUNK.COUNT"))
	r.Register("CHUNK.SETRANGE", "chunk", hint("CHUNK.SETRANGE"))
	r.Register("CHUNK.ISLOADED", "chunk", hint("CHUNK.ISLOADED"))
	r.Register("TERRAIN.LOAD", "terrain", hint("TERRAIN.LOAD"))
	r.Register("TERRAIN.GETNORMAL", "terrain", hint("TERRAIN.GETNORMAL"))
	r.Register("TERRAIN.SETSCALE", "terrain", hint("TERRAIN.SETSCALE"))
	r.Register("TERRAIN.GETSPLAT", "terrain", hint("TERRAIN.GETSPLAT"))
	r.Register("TERRAIN.RAYCAST", "terrain", hint("TERRAIN.RAYCAST"))
	r.Register("TERRAIN.SETDETAIL", "terrain", hint("TERRAIN.SETDETAIL"))
}

// TickStreaming is a no-op without CGO.
func (m *Module) TickStreaming(rt *runtime.Runtime) {}

// Preload is a no-op without CGO.
func (m *Module) Preload(rt *runtime.Runtime, radius int) {}

// PreloadTerrain is a no-op without CGO.
func (m *Module) PreloadTerrain(_ heap.Handle, _ int) error { return nil }

// SetCenter is a no-op without CGO.
func (m *Module) SetCenter(x, z float32) {}

// SetStreamEnabled is a no-op without CGO.
func (m *Module) SetStreamEnabled(on bool) {}

// StatusString without CGO.
func (m *Module) StatusString() string { return "terrain requires CGO" }

// IsReadyTerrain without CGO.
func (m *Module) IsReadyTerrain(_ heap.Handle) bool { return false }
