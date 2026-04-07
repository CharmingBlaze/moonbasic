// Package mbnet implements NET.*, PEER.*, and EVENT.* using ENet (CGO) or stubs.
package mbnet

import (
	"fmt"
	"sync"

	mbruntime "moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// Module registers networking commands.
type Module struct {
	h *heap.Store

	invokeMu sync.Mutex
	invoke   func(string, []value.Value) (value.Value, error)

	// NETREAD* cursor: last payload copied when EVENT.DATA runs on a RECEIVE event.
	readMu    sync.Mutex
	readBuf   []byte
	readOff   int
	channels  int // ENet channel count for the next NET.CREATESERVER / NET.CREATECLIENT (default 1).
}

// NewModule creates the net module.
func NewModule() *Module { return &Module{channels: 1} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }

// channelLimit returns ENet channel count for NewHost / Connect (clamped 1–32).
func (m *Module) channelLimit() uint64 {
	if m == nil || m.channels < 1 {
		return 1
	}
	if m.channels > 32 {
		return 32
	}
	return uint64(m.channels)
}

// SetUserInvoker wires VM.CallUserFunction for SERVER/CLIENT/RPC/LOBBY callbacks.
func (m *Module) SetUserInvoker(fn func(string, []value.Value) (value.Value, error)) {
	m.invokeMu.Lock()
	defer m.invokeMu.Unlock()
	m.invoke = fn
}

func (m *Module) callUser(name string, args []value.Value) (value.Value, error) {
	m.invokeMu.Lock()
	fn := m.invoke
	m.invokeMu.Unlock()
	if fn == nil {
		return value.Nil, fmt.Errorf("multiplayer callbacks require SetUserInvoker (host must wire VM)")
	}
	return fn(name, args)
}

// Register implements runtime.Module.
func (m *Module) Register(reg mbruntime.Registrar) {
	registerNetCommands(m, reg)
}

// Shutdown releases ENet and any open hosts (CGO path).
func (m *Module) Shutdown() {
	shutdownNet(m)
}
