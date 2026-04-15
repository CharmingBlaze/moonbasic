//go:build !cgo

package jolt

// CharacterContactEvent mirrors character.go for builds without CGO (physics3d stub, CI headless).
type CharacterContactEvent struct {
	BodyB    uint32
	Position Vec3
	Normal   Vec3
	Distance float32
}
