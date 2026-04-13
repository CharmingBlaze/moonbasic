//go:build !cgo && windows

package rl

// FrustumCullDistances returns rlgl cull near/far when raylib.dll is loaded.
// When the DLL was not loaded (deferred purego init for `go test` / MOONBASIC_SKIP_RAYLIB_DLL),
// returns Raylib defaults 0.01 / 1000 so CPU frustum math matches typical rlgl state.
func FrustumCullDistances() (near, far float32) {
	if raylibDll == 0 {
		return 0.01, 1000.0
	}
	return float32(rlGetCullDistanceNear()), float32(rlGetCullDistanceFar())
}
