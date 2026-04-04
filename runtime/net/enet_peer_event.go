//go:build cgo

package mbnet

import (
	"fmt"

	"github.com/codecat/go-enet"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func truthy(v value.Value) bool {
	switch v.Kind {
	case value.KindBool:
		return v.IVal != 0
	case value.KindInt:
		return v.IVal != 0
	case value.KindFloat:
		return v.FVal != 0
	default:
		return false
	}
}

func peerSend(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("PEER.SEND: heap not bound")
	}
	if len(args) != 4 || args[0].Kind != value.KindHandle || args[2].Kind != value.KindString {
		return value.Nil, fmt.Errorf("PEER.SEND expects (peer, channel, data$, reliable)")
	}
	po, err := heap.Cast[*peerObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if po.peer == nil {
		return value.Nil, runtime.Errorf("PEER.SEND: peer closed")
	}
	ch, ok := args[1].ToFloat()
	if !ok || ch < 0 || ch > 255 {
		return value.Nil, fmt.Errorf("PEER.SEND: bad channel")
	}
	flags := enet.PacketFlags(0)
	if truthy(args[3]) {
		flags = enet.PacketFlagReliable
	}
	data, err := rt.ArgString(args, 2)
	if err != nil {
		return value.Nil, err
	}
	if err := po.peer.SendString(data, uint8(ch), flags); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func peerDisconnect(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("PEER.DISCONNECT: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("PEER.DISCONNECT expects peer handle")
	}
	po, err := heap.Cast[*peerObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if po.peer != nil {
		po.peer.Disconnect(0)
	}
	return value.Nil, nil
}

func peerIP(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("PEER.IP: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("PEER.IP expects peer handle")
	}
	po, err := heap.Cast[*peerObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if po.peer == nil {
		return rt.RetString(""), nil
	}
	return rt.RetString(po.peer.GetAddress().String()), nil
}

func peerPing(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("PEER.PING: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("PEER.PING expects peer handle")
	}
	po, err := heap.Cast[*peerObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if po.peer == nil {
		return value.FromInt(0), nil
	}
	return value.FromInt(int64(po.peer.GetRoundTripTime())), nil
}

func eventChannel(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("EVENT.CHANNEL: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("EVENT.CHANNEL expects event handle")
	}
	eo, err := heap.Cast[*eventObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	return value.FromInt(int64(eo.ch)), nil
}

func eventType(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("EVENT.TYPE: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("EVENT.TYPE expects event handle")
	}
	eo, err := heap.Cast[*eventObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	return value.FromInt(int64(eo.typ)), nil
}

func eventPeer(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("EVENT.PEER: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("EVENT.PEER expects event handle")
	}
	eo, err := heap.Cast[*eventObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(eo.peerH), nil
}

func eventData(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("EVENT.DATA: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("EVENT.DATA expects event handle")
	}
	eo, err := heap.Cast[*eventObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	return rt.RetString(eo.data), nil
}

func eventFree(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("EVENT.FREE: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("EVENT.FREE expects event handle")
	}
	_ = m.h.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}
