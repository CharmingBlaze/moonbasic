//go:build cgo

package mbnet

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/codecat/go-enet"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

type packetObj struct {
	pkt     enet.Packet
	release heap.ReleaseOnce
}

func (p *packetObj) TypeName() string { return "NetPacket" }
func (p *packetObj) TypeTag() uint16  { return heap.TagNetPacket }
func (p *packetObj) Free() {
	p.release.Do(func() {
		if p.pkt != nil {
			p.pkt.Destroy()
			p.pkt = nil
		}
	})
}

func netSetChannels(m *Module, args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("NET.SETCHANNELS expects (count)")
	}
	n, ok := args[0].ToFloat()
	if !ok || n < 1 || n > 32 {
		return value.Nil, fmt.Errorf("NET.SETCHANNELS: count must be 1–32 (applies to the next NET.CREATESERVER / NET.CREATECLIENT)")
	}
	m.channels = int(n)
	return value.Nil, nil
}

func netFlush(m *Module, args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("NET.FLUSH expects (host)")
	}
	// go-enet Host does not expose enet_host_flush; pumping with NET.UPDATE / NET.SERVICE sends queued packets.
	return value.Nil, runtime.Errorf("NET.FLUSH: not available with upstream go-enet (use NET.UPDATE / NET.SERVICE to pump the host)")
}

func netService(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("NET.SERVICE: heap not bound")
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("NET.SERVICE expects (host, timeout_ms)")
	}
	ho, err := heap.Cast[*hostObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if ho.host == nil {
		return value.Nil, runtime.Errorf("NET.SERVICE: host closed")
	}
	tf, _ := args[1].ToFloat()
	if tf < 0 || tf > float64(^uint32(0)) {
		return value.Nil, fmt.Errorf("NET.SERVICE: bad timeout")
	}
	pumpHost(m, ho, uint32(tf))
	return value.Nil, nil
}

func packetCreate(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("PACKET.CREATE: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("PACKET.CREATE expects (data$)")
	}
	s, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	pkt, err := enet.NewPacket([]byte(s), 0)
	if err != nil {
		return value.Nil, err
	}
	id, err := m.h.Alloc(&packetObj{pkt: pkt})
	if err != nil {
		pkt.Destroy()
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func packetData(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("PACKET.DATA: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("PACKET.DATA expects (packet)")
	}
	po, err := heap.Cast[*packetObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if po.pkt == nil {
		return rt.RetString(""), nil
	}
	return rt.RetString(string(po.pkt.GetData())), nil
}

func packetFree(m *Module, args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("PACKET.FREE expects packet handle")
	}
	return value.Nil, m.h.Free(heap.Handle(args[0].IVal))
}

func peerSendPacket(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("PEER.SENDPACKET: heap not bound")
	}
	if len(args) != 3 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("PEER.SENDPACKET expects (peer, packet, channel)")
	}
	po, err := heap.Cast[*peerObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if po.peer == nil {
		return value.Nil, runtime.Errorf("PEER.SENDPACKET: peer closed")
	}
	pkto, err := heap.Cast[*packetObj](m.h, heap.Handle(args[1].IVal))
	if err != nil {
		return value.Nil, err
	}
	if pkto.pkt == nil {
		return value.Nil, fmt.Errorf("PEER.SENDPACKET: packet freed")
	}
	ch, ok := args[2].ToFloat()
	if !ok || ch < 0 || ch > 255 {
		return value.Nil, fmt.Errorf("PEER.SENDPACKET: bad channel")
	}
	if err := po.peer.SendPacket(pkto.pkt, uint8(ch)); err != nil {
		return value.Nil, err
	}
	// Packet ownership transfers to ENet; clear handle without double-destroy.
	pkto.pkt = nil
	_ = m.h.Free(heap.Handle(args[1].IVal))
	return value.Nil, nil
}

func netSendStringHelper(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("NETSENDSTRING: heap not bound")
	}
	if len(args) < 2 || len(args) > 4 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("NETSENDSTRING expects (peer, text$, [channel], [reliable])")
	}
	ch := uint8(0)
	if len(args) >= 3 {
		if cv, ok := args[2].ToFloat(); ok {
			ch = uint8(cv)
		}
	}
	flags := enet.PacketFlagReliable // Default to reliable for simple helpers
	if len(args) >= 4 {
		if !truthy(args[3]) {
			flags = 0
		}
	}
	s, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	b := []byte(s)
	buf := make([]byte, 4+len(b))
	binary.LittleEndian.PutUint32(buf, uint32(len(b)))
	copy(buf[4:], b)
	po, err := heap.Cast[*peerObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if po.peer == nil {
		return value.Nil, runtime.Errorf("NETSENDSTRING: peer closed")
	}
	pkt, err := enet.NewPacket(buf, flags)
	if err != nil {
		return value.Nil, err
	}
	err = po.peer.SendPacket(pkt, ch)
	pkt.Destroy()
	if err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func netSendIntHelper(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) < 2 || len(args) > 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("NETSENDINT expects (peer, value, [channel], [reliable])")
	}
	ch := uint8(0)
	if len(args) >= 3 {
		if cv, ok := args[2].ToFloat(); ok {
			ch = uint8(cv)
		}
	}
	flags := enet.PacketFlagReliable
	if len(args) >= 4 {
		if !truthy(args[3]) {
			flags = 0
		}
	}
	v, ok := args[1].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("NETSENDINT: bad value")
	}
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(int32(v)))
	pkt, err := enet.NewPacket(b, flags)
	if err != nil {
		return value.Nil, err
	}
	po, err := heap.Cast[*peerObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		pkt.Destroy()
		return value.Nil, err
	}
	if po.peer == nil {
		pkt.Destroy()
		return value.Nil, runtime.Errorf("NETSENDINT: peer closed")
	}
	err = po.peer.SendPacket(pkt, ch)
	pkt.Destroy()
	if err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func netSendFloatHelper(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) < 2 || len(args) > 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("NETSENDFLOAT expects (peer, value#, [channel], [reliable])")
	}
	ch := uint8(0)
	if len(args) >= 3 {
		if cv, ok := args[2].ToFloat(); ok {
			ch = uint8(cv)
		}
	}
	flags := enet.PacketFlagReliable
	if len(args) >= 4 {
		if !truthy(args[3]) {
			flags = 0
		}
	}
	f, ok := args[1].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("NETSENDFLOAT: bad value")
	}
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, math.Float64bits(f))
	pkt, err := enet.NewPacket(b, flags)
	if err != nil {
		return value.Nil, err
	}
	po, err := heap.Cast[*peerObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		pkt.Destroy()
		return value.Nil, err
	}
	if po.peer == nil {
		pkt.Destroy()
		return value.Nil, runtime.Errorf("NETSENDFLOAT: peer closed")
	}
	err = po.peer.SendPacket(pkt, ch)
	pkt.Destroy()
	if err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func netReadString(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	m.readMu.Lock()
	defer m.readMu.Unlock()
	if len(m.readBuf)-m.readOff < 4 {
		return rt.RetString(""), nil
	}
	n := int(binary.LittleEndian.Uint32(m.readBuf[m.readOff:]))
	m.readOff += 4
	if n < 0 || m.readOff+n > len(m.readBuf) {
		return value.Nil, fmt.Errorf("NETREADSTRING: truncated buffer (use NETSENDSTRING on sender)")
	}
	s := string(m.readBuf[m.readOff : m.readOff+n])
	m.readOff += n
	return rt.RetString(s), nil
}

func netReadInt(m *Module) (value.Value, error) {
	m.readMu.Lock()
	defer m.readMu.Unlock()
	if len(m.readBuf)-m.readOff < 4 {
		return value.FromInt(0), nil
	}
	v := int64(int32(binary.LittleEndian.Uint32(m.readBuf[m.readOff:])))
	m.readOff += 4
	return value.FromInt(v), nil
}

func netReadFloat(m *Module) (value.Value, error) {
	m.readMu.Lock()
	defer m.readMu.Unlock()
	if len(m.readBuf)-m.readOff < 8 {
		return value.FromFloat(0), nil
	}
	u := binary.LittleEndian.Uint64(m.readBuf[m.readOff:])
	m.readOff += 8
	return value.FromFloat(math.Float64frombits(u)), nil
}

func registerHelperNet(m *Module, reg runtime.Registrar) {
	reg.Register("NET.SETCHANNELS", "net", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		return netSetChannels(m, args)
	})
	reg.Register("NET.SERVICE", "net", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		return netService(m, args)
	})
	reg.Register("NET.FLUSH", "net", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		return netFlush(m, args)
	})
	reg.Register("PACKET.CREATE", "packet", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return packetCreate(m, rt, args...)
	})
	reg.Register("PACKET.MAKE", "packet", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return packetCreate(m, rt, args...)
	})
	reg.Register("PACKET.DATA", "packet", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return packetData(m, rt, args...)
	})
	reg.Register("PACKET.FREE", "packet", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		return packetFree(m, args)
	})
	reg.Register("PEER.SENDPACKET", "peer", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return peerSendPacket(m, rt, args...)
	})
	reg.Register("NETSENDSTRING", "net", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return netSendStringHelper(m, rt, args...)
	})
	reg.Register("NETSENDINT", "net", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return netSendIntHelper(m, rt, args...)
	})
	reg.Register("NETSENDFLOAT", "net", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return netSendFloatHelper(m, rt, args...)
	})
	reg.Register("NETREADSTRING", "net", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return netReadString(m, rt, args...)
	})
	reg.Register("NETREADINT", "net", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		return netReadInt(m)
	})
	reg.Register("NETREADFLOAT", "net", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		return netReadFloat(m)
	})
}
