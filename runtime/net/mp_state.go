package mbnet

import (
	"sync"

	"moonbasic/vm/heap"
)

type syncReg struct {
	h     heap.Handle
	flags int32
	id    uint32
}

type mpState struct {
	mu sync.Mutex

	serverH    heap.Handle
	clientH    heap.Handle
	serverPeer heap.Handle

	tickRate float64
	tickAcc  float64

	onSrvConn, onSrvDisc, onSrvMsg string
	onCliConn, onCliMsg, onCliSync string

	syncRegs []syncReg
	nextSID  uint32
	sidMap   map[heap.Handle]uint32
}

var gMP = mpState{sidMap: make(map[heap.Handle]uint32), tickRate: 20}

type lobbyObj struct {
	name    string
	maxP    int
	props   map[string]string
	started bool
	hostStr string
	port    int
}

func (o *lobbyObj) TypeName() string { return "Lobby" }

func (o *lobbyObj) TypeTag() uint16 { return heap.TagLobby }

func (o *lobbyObj) Free() {}

var (
	lobbyMu      sync.Mutex
	lobbyHandles []heap.Handle
)

func resetMultiplayerState() {
	gMP.mu.Lock()
	gMP.serverH = 0
	gMP.clientH = 0
	gMP.serverPeer = 0
	gMP.syncRegs = nil
	gMP.sidMap = make(map[heap.Handle]uint32)
	gMP.nextSID = 0
	gMP.tickAcc = 0
	gMP.tickRate = 20
	gMP.onSrvConn, gMP.onSrvDisc, gMP.onSrvMsg = "", "", ""
	gMP.onCliConn, gMP.onCliMsg, gMP.onCliSync = "", "", ""
	gMP.mu.Unlock()

	lobbyMu.Lock()
	lobbyHandles = nil
	lobbyMu.Unlock()
}
