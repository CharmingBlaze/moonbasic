//go:build (linux || windows) && cgo

package mbphysics3d

// MatrixBufferForEntitySync returns the shared float buffer filled by syncSharedBuffers.
// Indices are body bufferIndex*16 .. +15 (column-major 4×4: rotation in 0–10, translation in 12–14).
// Nil when physics has not started or buffer not allocated.
func MatrixBufferForEntitySync() []float32 {
	joltMu.Lock()
	defer joltMu.Unlock()
	if len(matrixBuffer) == 0 {
		return nil
	}
	return matrixBuffer
}

// MatrixBufferPrevForEntitySync returns the snapshot taken at the start of PHYSICS3D.STEP (before integration).
// Same layout as MatrixBufferForEntitySync. Used with PhysicsMatrixInterpAlpha for optional translation blending.
func MatrixBufferPrevForEntitySync() []float32 {
	joltMu.Lock()
	defer joltMu.Unlock()
	if len(prevMatrixBuffer) == 0 {
		return nil
	}
	return prevMatrixBuffer
}

// PhysicsMatrixInterpAlpha is accumulator/fixedStep after the last physics work in STEP (0..1).
// Typical use: lerp(prevTranslation, currTranslation, alpha) for high-refresh rendering (translation only).
func PhysicsMatrixInterpAlpha() float64 {
	joltMu.Lock()
	defer joltMu.Unlock()
	if matrixInterpFixed <= 1e-15 {
		return 1
	}
	a := matrixInterpAccum / matrixInterpFixed
	if a < 0 {
		return 0
	}
	if a > 1 {
		return 1
	}
	return a
}
