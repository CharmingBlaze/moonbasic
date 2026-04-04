//go:build cgo

package mbnet

import (
	"encoding/binary"
	"fmt"
	"sync"

	"github.com/codecat/go-enet"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

type queuedEv struct {
	typ   int32
	peerH int32
	ch    uint8
	data  string
}

// hostObj owns an ENet host; frees peer handles (children) before destroying the host.
type hostObj struct {
	host    enet.Host
	store   *heap.Store
	peerIDs []int32
	q       []queuedEv
	release heap.ReleaseOnce
}

func (h *hostObj) TypeName() string { return "NetHost" }

func (h *hostObj) TypeTag() uint16 { return heap.TagHost }

func (h *hostObj) Free() {
	h.release.Do(func() {
		if h.store != nil {
			ids := append([]int32(nil), h.peerIDs...)
			h.peerIDs = nil
			for _, pid := range ids {
				h.store.Free(pid)
			}
			h.store = nil
		}
		if h.host != nil {
			h.host.Destroy()
			h.host = nil
		}
	})
}

type peerObj struct {
	peer    enet.Peer
	release heap.ReleaseOnce
}

func (p *peerObj) TypeName() string { return "NetPeer" }

func (p *peerObj) TypeTag() uint16 { return heap.TagPeer }

func (p *peerObj) Free() {
	p.release.Do(func() {
		if p.peer != nil {
			p.peer.SetData(nil)
			p.peer = nil
		}
	})
}

type eventObj struct {
	typ   int32
	peerH int32
	ch    uint8
	data  string
}

func (e *eventObj) TypeName() string { return "NetEvent" }

func (e *eventObj) TypeTag() uint16 { return heap.TagEvent }

func (e *eventObj) Free() {}

type enetGlobal struct {
	mu    sync.Mutex
	ready bool
	hosts map[int32]struct{}
}

var g = &enetGlobal{hosts: make(map[int32]struct{})}

func registerNetCommands(m *Module, reg runtime.Registrar) {
	reg.Register("NET.START", "net", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		return netStart(args)
	})
	reg.Register("NET.STOP", "net", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		return netStop(m, args)
	})
	reg.Register("NET.CREATESERVER", "net", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		return netCreateServer(m, args)
	})
	reg.Register("NET.CREATECLIENT", "net", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		return netCreateClient(m, args)
	})
	reg.Register("NET.CONNECT", "net", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if err := requireInit(); err != nil {
			return value.Nil, err
		}
		if m.h == nil {
			return value.Nil, runtime.Errorf("NET.CONNECT: heap not bound")
		}
		if len(args) != 3 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
			return value.Nil, fmt.Errorf("NET.CONNECT expects (clientHost, host$, port)")
		}
		ho, err := heap.Cast[*hostObj](m.h, heap.Handle(args[0].IVal))
		if err != nil {
			return value.Nil, err
		}
		if ho.host == nil {
			return value.Nil, runtime.Errorf("NET.CONNECT: invalid host")
		}
		pf, ok := args[2].ToFloat()
		if !ok || pf < 0 || pf > 65535 {
			return value.Nil, fmt.Errorf("NET.CONNECT: invalid port")
		}
		hostName, err := rt.ArgString(args, 1)
		if err != nil {
			return value.Nil, err
		}
		peer, err := ho.host.Connect(enet.NewAddress(hostName, uint16(pf)), 1, 0)
		if err != nil {
			return value.Nil, err
		}
		pid := registerPeer(m, ho, peer)
		return value.FromHandle(pid), nil
	})
	reg.Register("NET.UPDATE", "net", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		return netUpdate(m, args)
	})
	reg.Register("NET.RECEIVE", "net", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		return netReceive(m, args)
	})
	reg.Register("NET.CLOSE", "net", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		return netClose(m, args)
	})
	reg.Register("NET.BROADCAST", "net", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if m.h == nil {
			return value.Nil, runtime.Errorf("NET.BROADCAST: heap not bound")
		}
		if len(args) != 4 || args[0].Kind != value.KindHandle || args[2].Kind != value.KindString {
			return value.Nil, fmt.Errorf("NET.BROADCAST expects (server, channel, data$, reliable)")
		}
		ho, err := heap.Cast[*hostObj](m.h, heap.Handle(args[0].IVal))
		if err != nil {
			return value.Nil, err
		}
		if ho.host == nil {
			return value.Nil, runtime.Errorf("NET.BROADCAST: host closed")
		}
		ch, ok := args[1].ToFloat()
		if !ok || ch < 0 || ch > 255 {
			return value.Nil, fmt.Errorf("NET.BROADCAST: bad channel")
		}
		data, err := rt.ArgString(args, 2)
		if err != nil {
			return value.Nil, err
		}
		flags := enet.PacketFlags(0)
		if truthy(args[3]) {
			flags = enet.PacketFlagReliable
		}
		if err := ho.host.BroadcastString(data, uint8(ch), flags); err != nil {
			return value.Nil, err
		}
		return value.Nil, nil
	})
	reg.Register("NET.PEERCOUNT", "net", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		return netPeerCount(m, args)
	})
	reg.Register("NET.SETTIMEOUT", "net", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		return netSetTimeout(m, args)
	})
	reg.Register("NET.GETPING", "net", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		return netGetPing(m, args)
	})
	reg.Register("NET.SETBANDWIDTH", "net", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		return netSetBandwidth(m, args)
	})
	reg.Register("PEER.SEND", "peer", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return peerSend(m, rt, args...)
	})
	reg.Register("PEER.DISCONNECT", "peer", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		return peerDisconnect(m, args)
	})
	reg.Register("PEER.IP", "peer", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return peerIP(m, rt, args...)
	})
	reg.Register("PEER.PING", "peer", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		return peerPing(m, args)
	})
	reg.Register("EVENT.TYPE", "event", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		return eventType(m, args)
	})
	reg.Register("EVENT.PEER", "event", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		return eventPeer(m, args)
	})
	reg.Register("EVENT.DATA", "event", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return eventData(m, rt, args...)
	})
	reg.Register("EVENT.FREE", "event", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		return eventFree(m, args)
	})
	reg.Register("EVENT.CHANNEL", "event", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		return eventChannel(m, args)
	})
	registerHighLevelNet(m, reg)
}

func shutdownNet(m *Module) {
	netFullStop(m)
}

func netFullStop(m *Module) {
	resetMultiplayerState()
	if m == nil || m.h == nil {
		g.mu.Lock()
		defer g.mu.Unlock()
		if g.ready {
			enet.Deinitialize()
			g.ready = false
		}
		g.hosts = make(map[int32]struct{})
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	for hid := range g.hosts {
		_ = m.h.Free(hid)
	}
	g.hosts = make(map[int32]struct{})
	if g.ready {
		enet.Deinitialize()
		g.ready = false
	}
}

func requireInit() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if !g.ready {
		return runtime.Errorf("NET: call NET.START first")
	}
	return nil
}

func netStart(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("NET.START expects 0 arguments")
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.hosts == nil {
		g.hosts = make(map[int32]struct{})
	}
	if !g.ready {
		enet.Initialize()
		g.ready = true
	}
	return value.Nil, nil
}

func netStop(m *Module, args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("NET.STOP expects 0 arguments")
	}
	netFullStop(m)
	return value.Nil, nil
}

func netCreateServer(m *Module, args []value.Value) (value.Value, error) {
	if err := requireInit(); err != nil {
		return value.Nil, err
	}
	if m.h == nil {
		return value.Nil, runtime.Errorf("NET.CREATESERVER: heap not bound")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("NET.CREATESERVER expects (port, maxclients)")
	}
	port, ok := args[0].ToFloat()
	if !ok || port < 0 || port > 65535 {
		return value.Nil, fmt.Errorf("NET.CREATESERVER: invalid port")
	}
	mc, ok := args[1].ToFloat()
	if !ok || mc < 1 {
		return value.Nil, fmt.Errorf("NET.CREATESERVER: maxclients must be >= 1")
	}
	addr := enet.NewListenAddress(uint16(port))
	h, err := enet.NewHost(addr, uint64(mc), 1, 0, 0)
	if err != nil {
		return value.Nil, err
	}
	ho := &hostObj{host: h, store: m.h}
	hid, err := m.h.Alloc(ho)
	if err != nil {
		return value.Nil, err
	}
	g.mu.Lock()
	g.hosts[hid] = struct{}{}
	g.mu.Unlock()
	return value.FromHandle(hid), nil
}

func netCreateClient(m *Module, args []value.Value) (value.Value, error) {
	if err := requireInit(); err != nil {
		return value.Nil, err
	}
	if m.h == nil {
		return value.Nil, runtime.Errorf("NET.CREATECLIENT: heap not bound")
	}
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("NET.CREATECLIENT expects 0 arguments")
	}
	h, err := enet.NewHost(nil, 32, 1, 0, 0)
	if err != nil {
		return value.Nil, err
	}
	ho := &hostObj{host: h, store: m.h}
	hid, err := m.h.Alloc(ho)
	if err != nil {
		return value.Nil, err
	}
	g.mu.Lock()
	g.hosts[hid] = struct{}{}
	g.mu.Unlock()
	return value.FromHandle(hid), nil
}

func registerPeer(m *Module, ho *hostObj, p enet.Peer) heap.Handle {
	if d := p.GetData(); len(d) >= 4 {
		return heap.Handle(binary.LittleEndian.Uint32(d))
	}
	pid, err := m.h.Alloc(&peerObj{peer: p})
	if err != nil {
		// Log error or handle? Networking can't easily recover if heap is full.
		return 0
	}
	ho.peerIDs = append(ho.peerIDs, pid)
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(pid))
	p.SetData(buf)
	return pid
}

func lookupPeerID(ho *hostObj, p enet.Peer, m *Module) heap.Handle {
	if d := p.GetData(); len(d) >= 4 {
		return heap.Handle(binary.LittleEndian.Uint32(d))
	}
	return registerPeer(m, ho, p)
}

func pumpHost(m *Module, ho *hostObj) {
	if ho.host == nil {
		return
	}
	for {
		ev := ho.host.Service(0)
		switch ev.GetType() {
		case enet.EventNone:
			return
		case enet.EventConnect:
			pid := lookupPeerID(ho, ev.GetPeer(), m)
			ho.q = append(ho.q, queuedEv{typ: 1, peerH: pid, ch: 0})
		case enet.EventDisconnect:
			pid := lookupPeerID(ho, ev.GetPeer(), m)
			ho.q = append(ho.q, queuedEv{typ: 2, peerH: pid, ch: 0})
			removePeerHandle(m, ho, pid)
		case enet.EventReceive:
			pkt := ev.GetPacket()
			ch := ev.GetChannelID()
			data := string(pkt.GetData())
			pkt.Destroy()
			pid := lookupPeerID(ho, ev.GetPeer(), m)
			ho.q = append(ho.q, queuedEv{typ: 3, peerH: pid, ch: ch, data: data})
		default:
			return
		}
	}
}

func netUpdate(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("NET.UPDATE: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("NET.UPDATE expects host handle")
	}
	ho, err := heap.Cast[*hostObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if ho.host == nil {
		return value.Nil, runtime.Errorf("NET.UPDATE: host closed")
	}
	pumpHost(m, ho)
	return value.Nil, nil
}

// netTryPopEvent dequeues the next queued network event for a host, pumping the host first if needed.
// Returns event handle 0 if none.
func netTryPopEvent(m *Module, hid heap.Handle) (heap.Handle, error) {
	if m.h == nil {
		return 0, runtime.Errorf("NET.RECEIVE: heap not bound")
	}
	ho, err := heap.Cast[*hostObj](m.h, hid)
	if err != nil {
		return 0, err
	}
	if ho.host == nil {
		return 0, runtime.Errorf("NET.RECEIVE: host closed")
	}
	if len(ho.q) == 0 {
		pumpHost(m, ho)
	}
	if len(ho.q) == 0 {
		return 0, nil
	}
	qe := ho.q[0]
	ho.q = ho.q[1:]
	eid, err := m.h.Alloc(&eventObj{typ: qe.typ, peerH: qe.peerH, ch: qe.ch, data: qe.data})
	if err != nil {
		return 0, err
	}
	return eid, nil
}

func netReceive(m *Module, args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("NET.RECEIVE expects host handle")
	}
	eid, err := netTryPopEvent(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(eid), nil
}

func removePeerHandle(m *Module, ho *hostObj, pid heap.Handle) {
	for i, id := range ho.peerIDs {
		if id == pid {
			ho.peerIDs = append(ho.peerIDs[:i], ho.peerIDs[i+1:]...)
			break
		}
	}
	_ = m.h.Free(pid)
}

func netClose(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("NET.CLOSE: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("NET.CLOSE expects host handle")
	}
	hid := heap.Handle(args[0].IVal)
	g.mu.Lock()
	delete(g.hosts, hid)
	g.mu.Unlock()
	_ = m.h.Free(hid)
	return value.Nil, nil
}

func netPeerCount(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("NET.PEERCOUNT: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("NET.PEERCOUNT expects host handle")
	}
	ho, err := heap.Cast[*hostObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	return value.FromInt(int64(len(ho.peerIDs))), nil
}

func netSetTimeout(m *Module, args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("NET.SETTIMEOUT expects (peer, ms)")
	}
	// go-enet does not expose enet_peer_timeout; reserved for future native wrapper.
	return value.Nil, nil
}

func netGetPing(m *Module, args []value.Value) (value.Value, error) {
	return peerPing(m, args)
}

func netSetBandwidth(m *Module, args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("NET.SETBANDWIDTH expects (host, inbps, outbps)")
	}
	// Bandwidth is fixed at host creation (NewHost); dynamic limit not exposed by go-enet.
	return value.Nil, nil
}
