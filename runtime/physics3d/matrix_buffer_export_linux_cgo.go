//go:build linux && cgo

package mbphysics3d

// MatrixBufferForEntitySync returns the shared float buffer filled by syncSharedBuffers.
// Indices are body bufferIndex*16 .. +15 (column-major 4x4, translation in slots 12–14).
// Nil when physics has not started or buffer not allocated.
func MatrixBufferForEntitySync() []float32 {
	joltMu.Lock()
	defer joltMu.Unlock()
	if len(matrixBuffer) == 0 {
		return nil
	}
	return matrixBuffer
}
