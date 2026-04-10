package texture

import "moonbasic/vm/heap"

type Module struct {
	h              *heap.Store
	enqueueCleanup func(func())
}

var globalEnqueuer func(func())

func setGlobalCleanupEnqueuer(fn func(func())) {
	globalEnqueuer = fn
}

func enqueueOnMainThread(fn func()) {
	if globalEnqueuer != nil {
		globalEnqueuer(fn)
	}
}

func (m *Module) BindCleanup(enqueuer func(func())) {
	m.enqueueCleanup = enqueuer
	setGlobalCleanupEnqueuer(enqueuer)
}

func NewModule() *Module { return &Module{} }

func (m *Module) BindHeap(h *heap.Store) { m.h = h }

func (m *Module) Reset() {}

