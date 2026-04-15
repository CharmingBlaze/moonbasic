//go:build (!linux && !windows) || !cgo

package player

// errPlayerRequiresCGOJolt is returned when PLAYER.* / CHARACTER.* KCC commands are used without native Jolt + CGO.
const errPlayerRequiresCGOJolt = "MoonBASIC: PLAYER commands require a CGO-enabled build with Jolt Physics (Windows/Linux desktop fullruntime)."
