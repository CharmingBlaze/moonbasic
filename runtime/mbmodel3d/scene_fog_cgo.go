//go:build cgo || (windows && !cgo)

package mbmodel3d

import "sync"

var sceneFogMu sync.Mutex

// World-side fog (WORLD.FOGMODE / FOGCOLOR / FOGDENSITY).
var (
	sceneWorldMode    int
	sceneWorldR       uint8
	sceneWorldG       uint8
	sceneWorldB       uint8
	sceneWorldDensity float32
)

// Snapshot of weathermod fog after FOG.* builtins.
var (
	sceneWxOn   bool
	sceneWxNear float32
	sceneWxFar  float32
	sceneWxR    int
	sceneWxG    int
	sceneWxB    int
)

// SyncSceneFogWorld is called from WORLD.FOG* handlers in worldmgr.
func SyncSceneFogWorld(mode int, r, g, b uint8, density float32) {
	sceneFogMu.Lock()
	defer sceneFogMu.Unlock()
	sceneWorldMode = mode
	sceneWorldR, sceneWorldG, sceneWorldB = r, g, b
	sceneWorldDensity = density
}

// SyncSceneFogWeather updates FOG.* (weathermod) snapshot used when FogOn overrides WORLD fog.
func SyncSceneFogWeather(on bool, near, far float32, r, g, b int) {
	sceneFogMu.Lock()
	defer sceneFogMu.Unlock()
	sceneWxOn = on
	sceneWxNear = near
	sceneWxFar = far
	sceneWxR, sceneWxG, sceneWxB = r, g, b
}

// sceneFogParams: mode 0 = off, 1 = linear, 2 = exponential (WORLD.FOGMODE(2)).
func sceneFogParams() (mode int, fr, fg, fb, near, far, density float32) {
	sceneFogMu.Lock()
	defer sceneFogMu.Unlock()

	if sceneWxOn {
		mode = 1
		fr = float32(sceneWxR) / 255
		fg = float32(sceneWxG) / 255
		fb = float32(sceneWxB) / 255
		near = sceneWxNear
		far = sceneWxFar
		if far <= near {
			far = near + 1
		}
		return
	}
	if sceneWorldMode == 0 {
		return 0, 0, 0, 0, 0, 0, 0
	}
	fr = float32(sceneWorldR) / 255
	fg = float32(sceneWorldG) / 255
	fb = float32(sceneWorldB) / 255
	d := sceneWorldDensity
	if d < 1e-6 {
		d = 0.02
	}
	switch sceneWorldMode {
	case 2:
		mode = 2
		density = d * 0.22
	default:
		mode = 1
		near = 0.5 + 14*d
		far = 6 + 200*d
	}
	return
}
