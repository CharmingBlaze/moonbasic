//go:build !cgo && !windows

package mbmodel3d

func SyncSceneFogWorld(mode int, r, g, b uint8, density float32) {}

func SyncSceneFogWeather(on bool, near, far float32, r, g, b int) {}
