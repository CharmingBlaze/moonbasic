// Package mathmod registers numeric built-ins (flat and MATH.* namespace).
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
