// Package mbdebug implements DEBUG.* diagnostics, profiling hooks, and optional
// on-screen watch overlay (via window frame-end hook).
package mbdebug

import (
	"sync"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type logLine struct {
	msg   string
	color rl.Color
}

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

	// overlayUser: when true, DEBUG.WATCH overlay may draw without Registry.DebugMode (see overlay_cgo.go).
	overlayUser bool

	watches []watchEntry

	profStack []profileFrame
	profSum   map[string]time.Duration
	profN     map[string]int64

	// FPS Graph
	fpsHistory    []float32
	fpsIdx        int
	showFPSGraph  bool
	lastFrameTime time.Time

	// Professional QoL
	monitorOn bool
	logLines  []logLine

	inspectID     int64
	drawPhysicsOn bool
}

// NewModule creates the debug module.
func NewModule() *Module {
	return &Module{
		fpsHistory: make([]float32, 120),
	}
}

func (m *Module) Reset() {}


