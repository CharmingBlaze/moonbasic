//go:build !linux || !cgo

package mbphysics3d

import (
	"fmt"
	"os"
)

// LogJoltPhysicsBackendHint prints one line to stderr: stub builds cannot link native Jolt.
func LogJoltPhysicsBackendHint() {
	fmt.Fprintf(os.Stderr, "moonBASIC: WARN [Jolt Physics] Stub mode — native Jolt is not linked on this build. ENTITY.PHYSICS and ENTITY.ADDPHYSICS are no-ops; use scripted ENTITY.UPDATE gravity and velocity until you run on Linux + CGO (see AGENTS.md).\n")
}
