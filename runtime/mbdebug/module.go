// Package mbdebug implements DEBUG.* diagnostics, profiling hooks, and optional
// on-screen watch overlay (via window frame-end hook).
package mbdebug

import (
	"sync"
	"time"
)

type watchEntry struct {
	label string
	text  string
}

type profileFrame struct {
	label string
	start time.Time
}

// Module registers DEBUG.* builtins and holds watch / profiler state.
type Module struct {
	mu sync.Mutex

	watches []watchEntry

	profStack []profileFrame
	profSum   map[string]time.Duration
	profN     map[string]int64
}

// NewModule creates the debug module.
func NewModule() *Module { return &Module{} }
