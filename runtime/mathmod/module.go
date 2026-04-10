// Package mathmod registers numeric built-ins (flat and MATH.* namespace).
//
// CLAMP, LERP, and WRAP use Raylib 5.5 raymath formulas (raymath_rl55.go, matching raylib-go rlClamp/rlLerp/rlWrap).
// SIN, COS, SQRT, POW, etc. use Go's math package (IEEE 754), same family as typical libm behind Raylib 5.5.
// The raylib Go package is not imported here: its init() would lock the OS thread for every process.
package mathmod

import (
	"math/rand"
	"time"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
)

// Module holds RNG state for RND (not global).
type Module struct {
	rng *rand.Rand
	h   *heap.Store
}

// NewModule creates the math module.
func NewModule() *Module {
	return &Module{rng: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

// BindHeap implements runtime.HeapAware (MATH.CIRCLEPOINT allocates tuple arrays).
func (m *Module) BindHeap(h *heap.Store) { m.h = h }

// Register implements runtime.Module.
func (m *Module) Register(r runtime.Registrar) {
	m.registerBasic(r)
	m.registerTrig(r)
	m.registerBlitzSurface(r)
	m.registerVector(r)
	m.registerAngleInterp(r)
	m.registerRemap(r)
	m.registerCirclePoint(r)
	m.registerRandom(r)
	m.registerLogic(r)
	m.registerMovement(r)
	m.registerGamePlaneHelpers(r)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}

func (m *Module) Reset() {}

