//go:build (linux || windows) && cgo

package mbphysics3d

import (
	"fmt"
	"os"
)

// LogJoltPhysicsBackendHint prints one line to stderr so users know native Jolt is linked.
func LogJoltPhysicsBackendHint() {
	fmt.Fprintf(os.Stderr, "moonBASIC: INFO [Jolt Physics] Native backend linked (CGO). ENTITY.PHYSICS / dynamic bodies are active.\n")
}
