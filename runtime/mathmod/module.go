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
)

// Module holds RNG state for RND (not global).
type Module struct {
	rng *rand.Rand
}

// NewModule creates the math module.
func NewModule() *Module {
	return &Module{rng: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

// Register implements runtime.Module.
func (m *Module) Register(r runtime.Registrar) {
	m.registerBasic(r)
	m.registerTrig(r)
	m.registerVector(r)
	m.registerAngleInterp(r)
	m.registerRandom(r)
	m.registerLogic(r)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
