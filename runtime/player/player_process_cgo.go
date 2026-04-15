//go:build (linux || windows) && cgo

package player

// Process is invoked from UPDATEPHYSICS (after ENTITY.UPDATE). Jolt CharacterVirtual is advanced
// via CHARACTERREF.UPDATE / PHYSICS3D.STEP; there is no parallel host-side KCC solver.
func (m *Module) Process(dt float64) {
	if dt <= 0 || dt > 0.5 {
		return
	}
	_ = m
}
