//go:build cgo

package mbnet

import (
	"fmt"

	"github.com/codecat/go-enet"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// registerENETCommands wires legacy ENET.* names to the same ENet stack as NET.* / PEER.* / PACKET.*.
func registerENETCommands(m *Module, reg runtime.Registrar) {
	reg.Register("ENET.INITIALIZE", "enet", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		return enetInitialize(args)
	})
	reg.Register("ENET.DEINITIALIZE", "enet", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		return enetDeinitialize(m, args)
	})
	reg.Register("ENET.CREATEHOST", "enet", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return enetCreateHost(m, rt, args...)
	})
	reg.Register("ENET.MAKEHOST", "enet", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return enetCreateHost(m, rt, args...)
	})
	reg.Register("ENET.HOSTSERVICE", "enet", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		return netService(m, args)
	})
	reg.Register("ENET.HOSTBROADCAST", "enet", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return enetHostBroadcast(m, rt, args...)
	})
	reg.Register("ENET.PEERSEND", "enet", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		// Manifest order: (peer, channel, packet); PEER.SENDPACKET uses (peer, packet, channel).
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("ENET.PEERSEND expects (peer, channel, packet)")
		}
		return peerSendPacket(m, rt, args[0], args[2], args[1])
	})
	reg.Register("ENET.PEERPING", "enet", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		return peerPing(m, args)
	})
}

func enetInitialize(args []value.Value) (value.Value, error) {
	return netStart(args)
}

func enetDeinitialize(m *Module, args []value.Value) (value.Value, error) {
	return netStop(m, args)
}

func enetCreateHost(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if err := requireInit(); err != nil {
		return value.Nil, err
	}
	if m.h == nil {
		return value.Nil, runtime.Errorf("ENET.CREATEHOST: heap not bound")
	}
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("ENET.CREATEHOST expects (address$, port, maxPeers, channels, bandwidth)")
	}
	addrStr, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	pf, ok := args[1].ToFloat()
	if !ok || pf < 0 || pf > 65535 {
		return value.Nil, fmt.Errorf("ENET.CREATEHOST: invalid port")
	}
	mc, ok := args[2].ToFloat()
	if !ok || mc < 1 {
		return value.Nil, fmt.Errorf("ENET.CREATEHOST: maxPeers must be >= 1")
	}
	chN, ok := args[3].ToFloat()
	if !ok || chN < 1 || chN > 32 {
		return value.Nil, fmt.Errorf("ENET.CREATEHOST: channels must be 1–32")
	}
	m.channels = int(chN)
	bwf, ok := args[4].ToFloat()
	if !ok || bwf < 0 {
		return value.Nil, fmt.Errorf("ENET.CREATEHOST: bandwidth must be >= 0")
	}
	bw := uint32(bwf)
	if bwf > float64(^uint32(0)) {
		bw = ^uint32(0)
	}

	var addr enet.Address
	if addrStr == "" {
		addr = enet.NewListenAddress(uint16(pf))
	} else {
		addr = enet.NewAddress(addrStr, uint16(pf))
	}

	h, err := enet.NewHost(addr, uint64(mc), m.channelLimit(), bw, bw)
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

func enetHostBroadcast(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("ENET.HOSTBROADCAST: heap not bound")
	}
	if len(args) != 4 || args[0].Kind != value.KindHandle || args[3].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("ENET.HOSTBROADCAST expects (host, channel, flags, packet)")
	}
	_ = args[2] // flags reserved; packet carries ENet flags (see PACKET.CREATE)
	ho, err := heap.Cast[*hostObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if ho.host == nil {
		return value.Nil, runtime.Errorf("ENET.HOSTBROADCAST: host closed")
	}
	ch, ok := args[1].ToFloat()
	if !ok || ch < 0 || ch > 255 {
		return value.Nil, fmt.Errorf("ENET.HOSTBROADCAST: bad channel")
	}
	pkto, err := heap.Cast[*packetObj](m.h, heap.Handle(args[3].IVal))
	if err != nil {
		return value.Nil, err
	}
	if pkto.pkt == nil {
		return value.Nil, fmt.Errorf("ENET.HOSTBROADCAST: packet freed")
	}
	if err := ho.host.BroadcastPacket(pkto.pkt, uint8(ch)); err != nil {
		return value.Nil, err
	}
	pkto.pkt = nil
	_ = m.h.Free(heap.Handle(args[3].IVal))
	return value.Nil, nil
}
