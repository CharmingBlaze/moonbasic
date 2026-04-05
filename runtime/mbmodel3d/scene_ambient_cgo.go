//go:build cgo

package mbmodel3d

import "sync"

var (
	sceneAmbientMu            sync.Mutex
	sceneAmbientR             float32 = 0.06
	sceneAmbientG             float32 = 0.06
	sceneAmbientB             float32 = 0.06
	sceneAmbientScale         float32 = 1.0 // multiplies RGB (fourth component of RENDER.SETAMBIENT)
)

// SetSceneAmbient sets PBR hemispheric ambient tint (multiplier on albedo, per channel).
// Typical values are 0.02–0.15; default matches the built-in shader constant 0.06.
func SetSceneAmbient(r, g, b float32) {
	sceneAmbientMu.Lock()
	defer sceneAmbientMu.Unlock()
	sceneAmbientR, sceneAmbientG, sceneAmbientB = r, g, b
	sceneAmbientScale = 1
}

// SetSceneAmbientScaled sets ambient RGB and an overall scale (fourth parameter).
// Scale is 0.0–1.0, or 0–255 (values greater than 1 are normalized as 8-bit).
func SetSceneAmbientScaled(r, g, b, scale float32) {
	sceneAmbientMu.Lock()
	defer sceneAmbientMu.Unlock()
	sceneAmbientR, sceneAmbientG, sceneAmbientB = r, g, b
	if scale > 1 {
		scale /= 255
	}
	if scale < 0 {
		scale = 0
	}
	if scale > 1 {
		scale = 1
	}
	sceneAmbientScale = scale
}

func sceneAmbientRGB() (r, g, b float32) {
	sceneAmbientMu.Lock()
	defer sceneAmbientMu.Unlock()
	s := sceneAmbientScale
	return sceneAmbientR * s, sceneAmbientG * s, sceneAmbientB * s
}
