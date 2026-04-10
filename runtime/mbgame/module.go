// Package mbgame registers instant game utilities (shortcuts, collision math, timers, config, helpers).
package mbgame

import (
	"math/rand"
	"time"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
)

// Module holds RNG and optional config/timer state.
type Module struct {
	h  *heap.Store
	rng *rand.Rand
	t0 time.Time

	config *configStore

	shakeByCam map[int64]*shakeState
	flash      *screenFlashState
	crossfade  *musicCrossfadeState
	bursts     []*burstState
	vibrate    *vibrateState
	fpsCam     map[int64]*fpsCamState
	tpsCam     map[int64]*tpsCamState
}

// NewModule constructs game utility builtins.
func NewModule() *Module {
	return &Module{
		rng:        rand.New(rand.NewSource(time.Now().UnixNano())),
		t0:         time.Now(),
		shakeByCam: make(map[int64]*shakeState),
		fpsCam:     make(map[int64]*fpsCamState),
		tpsCam:     make(map[int64]*tpsCamState),
	}
}

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }

// Register implements runtime.Module.
func (m *Module) Register(r runtime.Registrar) {
	m.registerPure(r)
	m.registerTimeScale(r)
	m.registerShortcuts(r)
	m.registerDrawHelpers(r)
	m.registerWorldCamera(r)
	m.registerTimers(r)
	m.registerConfig(r)
	m.registerGamepad(r)
	m.registerAudioHelpers(r)
	m.registerBurst(r)
	m.registerPoolExtras(r)
	m.registerTileSprite(r)
	m.registerEntityStubs(r)
	m.registerDebugDraw(r)
	m.registerProcedural(r)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {
	m.shakeByCam = nil
	m.flash = nil
	m.crossfade = nil
	m.bursts = nil
	m.vibrate = nil
	m.fpsCam = nil
	m.tpsCam = nil
	m.config = nil
}

func (m *Module) Reset() {}

