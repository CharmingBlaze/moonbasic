//go:build linux && cgo

package mbphysics3d

import (
	"fmt"
	"os"
)

// LogJoltPhysicsBackendHint prints one line to stderr so users know native Jolt is linked.
func LogJoltPhysicsBackendHint() {
	fmt.Fprintf(os.Stderr, "moonBASIC: INFO [Jolt Physics] Native backend linked (Linux + CGO). ENTITY.PHYSICS / dynamic bodies are active.\n")
}
