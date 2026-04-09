package mbphysics2d

// Physics2DBackend is the seam for swapping 2D physics implementations without changing
// BODY2D.* / JOINT2D.* surface area. The active implementation is selected at compile time.
//
// Today: [github.com/ByteArena/box2d] (C++ Box2D via CGO) — see box2d.go and box2d_extra.go.
//
// Candidates for a zero-CGO path:
//   - A pure-Go Box2D-class engine (API differs from ByteArena — needs an adapter layer).
//   - WASM-hosted Box2D or a 2D subset behind a small host ABI (heavier, sandbox-friendly).
//
// [github.com/ByteArena/box2d]: https://github.com/ByteArena/box2d
type Physics2DBackend interface {
	// Name returns a stable identifier for logging and tests (e.g. "box2d_bytearena").
	Name() string
}

// CurrentBackend reports which backend this binary was built with.
func CurrentBackend() Physics2DBackend { return byteArenaBackend{} }

type byteArenaBackend struct{}

func (byteArenaBackend) Name() string { return "box2d_bytearena" }
