package heap

import "sync/atomic"

// ReleaseOnce runs a cleanup function at most once (RULE 1 — idempotent Free).
// It is safe under concurrent Free calls: only the first caller runs fn.
//
// Embed by value in HeapObject implementations that wrap Raylib, Jolt, Box2D, ENet, or OS handles:
//
//	type texObj struct {
//	    tex rl.Texture2D
//	    release ReleaseOnce
//	}
//	func (o *texObj) Free() {
//	    o.release.Do(func() { rl.UnloadTexture(o.tex) })
//	}
type ReleaseOnce struct {
	done atomic.Bool
}

// Do invokes fn exactly once across all calls to Do for this ReleaseOnce.
func (r *ReleaseOnce) Do(fn func()) {
	if r.done.Swap(true) {
		return
	}
	fn()
}
