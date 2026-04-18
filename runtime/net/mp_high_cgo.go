//go:build cgo

package mbnet

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/codecat/go-enet"

	"moonbasic/runtime"
	"moonbasic/runtime/mbmodel3d"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

const (
	chUser = uint8(0)
	chSync = uint8(1)
	chRPC  = uint8(2)
)

const (
	syncTransform = int32(1)
	syncAnimation = int32(2)
)

const (
	rpcWirePrefix  = "MBRPC1:"
	syncWirePrefix = "MBSYNC1:"
)

func registerHighLevelNet(m *Module, reg runtime.Registrar) {
	reg.Register("SERVER.START", "server", m.srvStart)
	reg.Register("SERVER.STOP", "server", m.srvStop)
	reg.Register("SERVER.ONCONNECT", "server", m.srvOnConnect)
	reg.Register("SERVER.ONDISCONNECT", "server", m.srvOnDisconnect)
	reg.Register("SERVER.ONMESSAGE", "server", m.srvOnMessage)
	reg.Register("SERVER.SYNCENTITY", "server", m.srvSyncEntity)
	reg.Register("SERVER.SETTICKRATE", "server", m.srvSetTickRate)
	reg.Register("SERVER.TICK", "server", m.srvTick)

	reg.Register("CLIENT.CONNECT", "client", m.cliConnect)
	reg.Register("CLIENT.STOP", "client", m.cliStop)
	reg.Register("CLIENT.ONCONNECT", "client", m.cliOnConnect)
	reg.Register("CLIENT.ONMESSAGE", "client", m.cliOnMessage)
	reg.Register("CLIENT.ONSYNC", "client", m.cliOnSync)
	reg.Register("CLIENT.TICK", "client", m.cliTick)

	reg.Register("RPC.CALL", "rpc", m.rpcCall)
	reg.Register("RPC.CALLTO", "rpc", m.rpcCallTo)
	reg.Register("RPC.CALLSERVER", "rpc", m.rpcCallServer)

	reg.Register("LOBBY.CREATE", "lobby", m.lobbyCreate)
	reg.Register("LOBBY.MAKE", "lobby", m.lobbyCreate)
	reg.Register("LOBBY.FREE", "lobby", m.lobbyFree)
	reg.Register("LOBBY.SETPROPERTY", "lobby", m.lobbySetProperty)
	reg.Register("LOBBY.SETHOST", "lobby", m.lobbySetHost)
	reg.Register("LOBBY.START", "lobby", m.lobbyStart)
	reg.Register("LOBBY.FIND", "lobby", m.lobbyFind)
	reg.Register("LOBBY.GETNAME", "lobby", m.lobbyGetName)
	reg.Register("LOBBY.JOIN", "lobby", m.lobbyJoin)

	// Professional Shorthands
	reg.Register("NET.HOST", "net", m.netHost)
	reg.Register("NET.CONNECT", "net", m.cliConnect)
	reg.Register("NET.SEND", "net", m.rpcCall)
	reg.Register("NET.SYNC", "net", m.srvSyncEntity)
}

func argF64MP(v value.Value) (float64, bool) {
	if f, ok := v.ToFloat(); ok {
		return f, true
	}
	if i, ok := v.ToInt(); ok {
		return float64(i), true
	}
	return 0, false
}

func valueToJSONArg(rt *runtime.Runtime, a value.Value) (interface{}, error) {
	switch a.Kind {
	case value.KindString:
		s, err := rt.ArgString([]value.Value{a}, 0)
		if err != nil {
			return nil, err
		}
		return s, nil
	case value.KindBool:
		return a.IVal != 0, nil
	case value.KindHandle:
		return float64(a.IVal), nil
	default:
		if f, ok := a.ToFloat(); ok {
			return f, nil
		}
		if i, ok := a.ToInt(); ok {
			return float64(i), nil
		}
	}
	return nil, fmt.Errorf("unsupported value kind for RPC JSON")
}

func jsonArgsToValues(rt *runtime.Runtime, raw []json.RawMessage) ([]value.Value, error) {
	out := make([]value.Value, 0, len(raw))
	for _, r := range raw {
		var v interface{}
		if err := json.Unmarshal(r, &v); err != nil {
			return nil, err
		}
		switch x := v.(type) {
		case float64:
			out = append(out, value.FromFloat(x))
		case string:
			out = append(out, rt.RetString(x))
		case bool:
			out = append(out, value.FromBool(x))
		case nil:
			out = append(out, value.Nil)
		default:
			return nil, fmt.Errorf("RPC: unsupported decoded type %T", x)
		}
	}
	return out, nil
}

func buildRPCWire(rt *runtime.Runtime, fn string, args []value.Value) (string, error) {
	arr := make([]interface{}, 0, len(args))
	for _, a := range args {
		x, err := valueToJSONArg(rt, a)
		if err != nil {
			return "", err
		}
		arr = append(arr, x)
	}
	payload := struct {
		F string        `json:"f"`
		A []interface{} `json:"a"`
	}{F: strings.ToUpper(strings.TrimSpace(fn)), A: arr}
	b, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return rpcWirePrefix + string(b), nil
}

func parseRPCWire(s string) (fn string, args []json.RawMessage, err error) {
	if !strings.HasPrefix(s, rpcWirePrefix) {
		return "", nil, fmt.Errorf("not an RPC packet")
	}
	js := strings.TrimPrefix(s, rpcWirePrefix)
	var wire struct {
		F string            `json:"f"`
		A []json.RawMessage `json:"a"`
	}
	if err := json.Unmarshal([]byte(js), &wire); err != nil {
		return "", nil, err
	}
	return wire.F, wire.A, nil
}

func broadcastHost(ho *hostObj, ch uint8, data string, reliable bool) error {
	if ho == nil || ho.host == nil {
		return fmt.Errorf("host closed")
	}
	flags := enet.PacketFlags(0)
	if reliable {
		flags = enet.PacketFlagReliable
	}
	return ho.host.BroadcastString(data, ch, flags)
}

func (m *Module) netHost(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("NET.HOST expects (port)")
	}
	return m.srvStart(rt, args[0], value.FromInt(32))
}

func (m *Module) srvStart(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("SERVER.START: heap not bound")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("SERVER.START expects (port, maxClients)")
	}
	port, ok1 := argF64MP(args[0])
	mc, ok2 := argF64MP(args[1])
	if !ok1 || !ok2 || port < 0 || port > 65535 || mc < 1 {
		return value.Nil, fmt.Errorf("SERVER.START: invalid port or maxClients")
	}
	gMP.mu.Lock()
	defer gMP.mu.Unlock()
	if gMP.serverH != 0 {
		return value.Nil, fmt.Errorf("SERVER.START: server already running")
	}
	// NET.START + CREATESERVER
	if _, err := netStart(nil); err != nil {
		return value.Nil, err
	}
	hv, err := netCreateServer(m, []value.Value{value.FromFloat(port), value.FromFloat(mc)})
	if err != nil {
		return value.Nil, err
	}
	gMP.serverH = heap.Handle(hv.IVal)
	return value.Nil, nil
}

func (m *Module) srvStop(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("SERVER.STOP expects 0 arguments")
	}
	gMP.mu.Lock()
	hid := gMP.serverH
	gMP.serverH = 0
	gMP.syncRegs = nil
	gMP.sidMap = make(map[heap.Handle]uint32)
	gMP.mu.Unlock()
	if hid != 0 && m.h != nil {
		_, _ = netClose(m, []value.Value{value.FromHandle(hid)})
	}
	return value.Nil, nil
}

func (m *Module) srvOnConnect(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("SERVER.ONCONNECT expects functionName")
	}
	s, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	gMP.mu.Lock()
	gMP.onSrvConn = strings.ToUpper(strings.TrimSpace(s))
	gMP.mu.Unlock()
	return value.Nil, nil
}

func (m *Module) srvOnDisconnect(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("SERVER.ONDISCONNECT expects functionName")
	}
	s, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	gMP.mu.Lock()
	gMP.onSrvDisc = strings.ToUpper(strings.TrimSpace(s))
	gMP.mu.Unlock()
	return value.Nil, nil
}

func (m *Module) srvOnMessage(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("SERVER.ONMESSAGE expects functionName")
	}
	s, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	gMP.mu.Lock()
	gMP.onSrvMsg = strings.ToUpper(strings.TrimSpace(s))
	gMP.mu.Unlock()
	return value.Nil, nil
}

func (m *Module) srvSyncEntity(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("SERVER.SYNCENTITY: heap not bound")
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("SERVER.SYNCENTITY expects (entityHandle, flags)")
	}
	eh := heap.Handle(args[0].IVal)
	flg, ok := argF64MP(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("SERVER.SYNCENTITY: flags must be numeric")
	}
	flags := int32(flg)
	gMP.mu.Lock()
	defer gMP.mu.Unlock()
	id, ok := gMP.sidMap[eh]
	if !ok {
		gMP.nextSID++
		id = gMP.nextSID
		gMP.sidMap[eh] = id
	}
	for i := range gMP.syncRegs {
		if gMP.syncRegs[i].h == eh {
			gMP.syncRegs = append(gMP.syncRegs[:i], gMP.syncRegs[i+1:]...)
			break
		}
	}
	gMP.syncRegs = append(gMP.syncRegs, syncReg{h: eh, flags: flags, id: id})
	return value.Nil, nil
}

func (m *Module) srvSetTickRate(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("SERVER.SETTICKRATE expects (hz)")
	}
	hz, ok := argF64MP(args[0])
	if !ok || hz <= 0 {
		return value.Nil, fmt.Errorf("SERVER.SETTICKRATE: hz must be > 0")
	}
	gMP.mu.Lock()
	gMP.tickRate = hz
	gMP.mu.Unlock()
	return value.Nil, nil
}

func (m *Module) flushSyncBroadcast(_ *runtime.Runtime) {
	gMP.mu.Lock()
	regs := append([]syncReg(nil), gMP.syncRegs...)
	srv := gMP.serverH
	gMP.mu.Unlock()
	if srv == 0 || len(regs) == 0 {
		return
	}
	ho, err := heap.Cast[*hostObj](m.h, srv)
	if err != nil || ho.host == nil {
		return
	}
	for _, sr := range regs {
		if sr.flags&syncTransform == 0 {
			continue
		}
		x, y, z, err := mbmodel3d.ModelTranslationForSync(m.h, sr.h)
		if err != nil {
			continue
		}
		m := map[string]float64{"id": float64(sr.id), "x": float64(x), "y": float64(y), "z": float64(z)}
		if sr.flags&syncAnimation != 0 {
			m["anim"] = 1
		}
		b, err := json.Marshal(m)
		if err != nil {
			continue
		}
		_ = broadcastHost(ho, chSync, syncWirePrefix+string(b), true)
	}
}

func (m *Module) processEventObj(rt *runtime.Runtime, eid heap.Handle, isServer bool) {
	if eid == 0 || m.h == nil {
		return
	}
	eo, err := heap.Cast[*eventObj](m.h, eid)
	if err != nil {
		return
	}
	defer func() { _ = m.h.Free(eid) }()

	switch eo.typ {
	case 1: // connect
		if isServer {
			fn := ""
			gMP.mu.Lock()
			fn = gMP.onSrvConn
			gMP.mu.Unlock()
			if fn != "" {
				_, _ = m.callUser(fn, []value.Value{value.FromHandle(eo.peerH)})
			}
		} else {
			fn := ""
			gMP.mu.Lock()
			fn = gMP.onCliConn
			gMP.mu.Unlock()
			if fn != "" {
				_, _ = m.callUser(fn, nil)
			}
		}
	case 2: // disconnect
		if isServer {
			fn := ""
			gMP.mu.Lock()
			fn = gMP.onSrvDisc
			gMP.mu.Unlock()
			if fn != "" {
				_, _ = m.callUser(fn, []value.Value{value.FromHandle(eo.peerH)})
			}
		}
	case 3: // receive
		switch eo.ch {
		case chUser:
			if isServer {
				fn := ""
				gMP.mu.Lock()
				fn = gMP.onSrvMsg
				gMP.mu.Unlock()
				if fn != "" {
					_, _ = m.callUser(fn, []value.Value{value.FromHandle(eo.peerH), rt.RetString(eo.data)})
				}
			} else {
				fn := ""
				gMP.mu.Lock()
				fn = gMP.onCliMsg
				gMP.mu.Unlock()
				if fn != "" {
					_, _ = m.callUser(fn, []value.Value{rt.RetString(eo.data)})
				}
			}
		case chSync:
			if !isServer {
				m.handleClientSync(rt, eo.data)
			}
		case chRPC:
			m.handleRPCPacket(rt, isServer, eo.peerH, eo.data)
		}
	}
}

func (m *Module) handleClientSync(rt *runtime.Runtime, data string) {
	if !strings.HasPrefix(data, syncWirePrefix) {
		return
	}
	js := strings.TrimPrefix(data, syncWirePrefix)
	var w map[string]float64
	if err := json.Unmarshal([]byte(js), &w); err != nil {
		return
	}
	fn := ""
	gMP.mu.Lock()
	fn = gMP.onCliSync
	gMP.mu.Unlock()
	if fn == "" {
		return
	}
	id := int64(w["id"])
	x := w["x"]
	y := w["y"]
	z := w["z"]
	_, _ = m.callUser(fn, []value.Value{value.FromInt(id), value.FromFloat(x), value.FromFloat(y), value.FromFloat(z)})
}

func (m *Module) handleRPCPacket(rt *runtime.Runtime, isServer bool, peerH heap.Handle, data string) {
	fn, rawArgs, err := parseRPCWire(data)
	if err != nil {
		return
	}
	args, err := jsonArgsToValues(rt, rawArgs)
	if err != nil {
		return
	}
	if isServer {
		args = append(args, value.FromHandle(peerH))
	}
	_, _ = m.callUser(fn, args)
}

func (m *Module) drainHost(rt *runtime.Runtime, hid heap.Handle, isServer bool) {
	for {
		eid, err := netTryPopEvent(m, hid)
		if err != nil || eid == 0 {
			break
		}
		m.processEventObj(rt, eid, isServer)
	}
}

func (m *Module) srvTick(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("SERVER.TICK expects (dt)")
	}
	dt, ok := argF64MP(args[0])
	if !ok || dt < 0 {
		return value.Nil, fmt.Errorf("SERVER.TICK: dt invalid")
	}
	gMP.mu.Lock()
	srv := gMP.serverH
	rate := gMP.tickRate
	gMP.mu.Unlock()
	if srv == 0 {
		return value.Nil, fmt.Errorf("SERVER.TICK: server not started")
	}
	_, _ = netUpdate(m, []value.Value{value.FromHandle(srv)})
	m.drainHost(rt, srv, true)

	gMP.mu.Lock()
	gMP.tickAcc += dt
	acc := gMP.tickAcc
	interval := 1.0 / rate
	if rate <= 0 {
		interval = 1e9
	}
	gMP.mu.Unlock()
	if acc >= interval {
		gMP.mu.Lock()
		gMP.tickAcc = math.Mod(acc, interval)
		gMP.mu.Unlock()
		m.flushSyncBroadcast(rt)
	}
	return value.Nil, nil
}

func (m *Module) cliConnect(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CLIENT.CONNECT: heap not bound")
	}
	if len(args) != 2 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("CLIENT.CONNECT expects (host, port)")
	}
	host, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	port, ok := argF64MP(args[1])
	if !ok || port < 0 || port > 65535 {
		return value.Nil, fmt.Errorf("CLIENT.CONNECT: bad port")
	}
	if _, err := netStart(nil); err != nil {
		return value.Nil, err
	}
	hv, err := netCreateClient(m, nil)
	if err != nil {
		return value.Nil, err
	}
	chid := heap.Handle(hv.IVal)
	peerV, err := m.innerConnectClient(rt, chid, host, uint16(port))
	if err != nil {
		_, _ = netClose(m, []value.Value{value.FromHandle(chid)})
		return value.Nil, err
	}
	gMP.mu.Lock()
	gMP.clientH = chid
	gMP.serverPeer = heap.Handle(peerV.IVal)
	gMP.mu.Unlock()
	return value.Nil, nil
}

func (m *Module) innerConnectClient(rt *runtime.Runtime, clientH heap.Handle, host string, port uint16) (value.Value, error) {
	// Reuse NET.CONNECT registration logic
	if err := requireInit(); err != nil {
		return value.Nil, err
	}
	ho, err := heap.Cast[*hostObj](m.h, clientH)
	if err != nil {
		return value.Nil, err
	}
	if ho.host == nil {
		return value.Nil, runtime.Errorf("CLIENT.CONNECT: invalid host")
	}
	peer, err := ho.host.Connect(enet.NewAddress(host, port), 1, 0)
	if err != nil {
		return value.Nil, err
	}
	pid := registerPeer(m, ho, peer)
	return value.FromHandle(pid), nil
}

func (m *Module) cliStop(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("CLIENT.STOP expects 0 arguments")
	}
	gMP.mu.Lock()
	cid := gMP.clientH
	gMP.clientH = 0
	gMP.serverPeer = 0
	gMP.mu.Unlock()
	if cid != 0 && m.h != nil {
		_, _ = netClose(m, []value.Value{value.FromHandle(cid)})
	}
	return value.Nil, nil
}

func (m *Module) cliOnConnect(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("CLIENT.ONCONNECT expects functionName")
	}
	s, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	gMP.mu.Lock()
	gMP.onCliConn = strings.ToUpper(strings.TrimSpace(s))
	gMP.mu.Unlock()
	return value.Nil, nil
}

func (m *Module) cliOnMessage(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("CLIENT.ONMESSAGE expects functionName")
	}
	s, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	gMP.mu.Lock()
	gMP.onCliMsg = strings.ToUpper(strings.TrimSpace(s))
	gMP.mu.Unlock()
	return value.Nil, nil
}

func (m *Module) cliOnSync(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("CLIENT.ONSYNC expects functionName")
	}
	s, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	gMP.mu.Lock()
	gMP.onCliSync = strings.ToUpper(strings.TrimSpace(s))
	gMP.mu.Unlock()
	return value.Nil, nil
}

func (m *Module) cliTick(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CLIENT.TICK expects (dt)")
	}
	dt, ok := argF64MP(args[0])
	if !ok || dt < 0 {
		return value.Nil, fmt.Errorf("CLIENT.TICK: dt invalid")
	}
	_ = dt
	gMP.mu.Lock()
	cid := gMP.clientH
	gMP.mu.Unlock()
	if cid == 0 {
		return value.Nil, fmt.Errorf("CLIENT.TICK: client not connected")
	}
	_, _ = netUpdate(m, []value.Value{value.FromHandle(cid)})
	m.drainHost(rt, cid, false)
	return value.Nil, nil
}

func (m *Module) rpcCall(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) < 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("RPC.CALL expects (functionName, ...)")
	}
	fn, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	gMP.mu.Lock()
	srv := gMP.serverH
	gMP.mu.Unlock()
	if srv == 0 {
		return value.Nil, fmt.Errorf("RPC.CALL: server not running")
	}
	ho, err := heap.Cast[*hostObj](m.h, srv)
	if err != nil {
		return value.Nil, err
	}
	msg, err := buildRPCWire(rt, fn, args[1:])
	if err != nil {
		return value.Nil, err
	}
	if err := broadcastHost(ho, chRPC, msg, true); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) rpcCallTo(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) < 2 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("RPC.CALLTO expects (peer, functionName, ...)")
	}
	fn, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	po, err := heap.Cast[*peerObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if po.peer == nil {
		return value.Nil, runtime.Errorf("RPC.CALLTO: peer closed")
	}
	msg, err := buildRPCWire(rt, fn, args[2:])
	if err != nil {
		return value.Nil, err
	}
	if err := po.peer.SendString(msg, chRPC, enet.PacketFlagReliable); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) rpcCallServer(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) < 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("RPC.CALLSERVER expects (functionName, ...)")
	}
	fn, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	gMP.mu.Lock()
	ph := gMP.serverPeer
	gMP.mu.Unlock()
	if ph == 0 {
		return value.Nil, fmt.Errorf("RPC.CALLSERVER: not connected")
	}
	po, err := heap.Cast[*peerObj](m.h, ph)
	if err != nil {
		return value.Nil, err
	}
	if po.peer == nil {
		return value.Nil, runtime.Errorf("RPC.CALLSERVER: peer closed")
	}
	msg, err := buildRPCWire(rt, fn, args[1:])
	if err != nil {
		return value.Nil, err
	}
	if err := po.peer.SendString(msg, chRPC, enet.PacketFlagReliable); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

// --- lobby ---

func (m *Module) lobbyCreate(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("LOBBY.CREATE: heap not bound")
	}
	if len(args) != 2 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("LOBBY.CREATE expects (name, maxPlayers)")
	}
	name, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	mp, ok := args[1].ToInt()
	if !ok {
		if f, okf := args[1].ToFloat(); okf {
			mp = int64(f)
			ok = true
		}
	}
	if !ok || mp < 1 {
		return value.Nil, fmt.Errorf("LOBBY.CREATE: maxPlayers must be >= 1")
	}
	o := &lobbyObj{name: name, maxP: int(mp), props: make(map[string]string)}
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	lobbyMu.Lock()
	lobbyHandles = append(lobbyHandles, id)
	lobbyMu.Unlock()
	return value.FromHandle(id), nil
}

func (m *Module) lobbyFree(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("LOBBY.FREE: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("LOBBY.FREE expects lobby handle")
	}
	h := heap.Handle(args[0].IVal)
	lobbyMu.Lock()
	for i, x := range lobbyHandles {
		if x == h {
			lobbyHandles = append(lobbyHandles[:i], lobbyHandles[i+1:]...)
			break
		}
	}
	lobbyMu.Unlock()
	_ = m.h.Free(h)
	return value.Nil, nil
}

func (m *Module) lobbySetProperty(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("LOBBY.SETPROPERTY: heap not bound")
	}
	if len(args) != 3 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString || args[2].Kind != value.KindString {
		return value.Nil, fmt.Errorf("LOBBY.SETPROPERTY expects (lobby, key, value)")
	}
	o, err := heap.Cast[*lobbyObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	k, e1 := rt.ArgString(args, 1)
	v, e2 := rt.ArgString(args, 2)
	if e1 != nil || e2 != nil {
		return value.Nil, fmt.Errorf("LOBBY.SETPROPERTY: bad strings")
	}
	o.props[strings.ToLower(strings.TrimSpace(k))] = v
	return value.Nil, nil
}

func (m *Module) lobbySetHost(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("LOBBY.SETHOST: heap not bound")
	}
	if len(args) != 3 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("LOBBY.SETHOST expects (lobby, host, port)")
	}
	o, err := heap.Cast[*lobbyObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	host, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	port, ok := argF64MP(args[2])
	if !ok || port < 0 || port > 65535 {
		return value.Nil, fmt.Errorf("LOBBY.SETHOST: bad port")
	}
	o.hostStr = host
	o.port = int(port)
	return value.Nil, nil
}

func (m *Module) lobbyStart(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("LOBBY.START: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("LOBBY.START expects lobby handle")
	}
	o, err := heap.Cast[*lobbyObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	o.started = true
	return value.Nil, nil
}

func (m *Module) lobbyFind(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("LOBBY.FIND: heap not bound")
	}
	if len(args) != 2 || args[0].Kind != value.KindString || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("LOBBY.FIND expects (key$, value$)")
	}
	key, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	val, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	key = strings.ToLower(strings.TrimSpace(key))
	lobbyMu.Lock()
	handles := append([]heap.Handle(nil), lobbyHandles...)
	lobbyMu.Unlock()
	var match []heap.Handle
	for _, hid := range handles {
		o, err := heap.Cast[*lobbyObj](m.h, hid)
		if err != nil {
			continue
		}
		if o.props[key] == val {
			match = append(match, hid)
		}
	}
	if len(match) == 0 {
		dims := []int64{1}
		a, err := heap.NewArray(dims)
		if err != nil {
			return value.Nil, err
		}
		a.Floats[0] = 0
		aid, err := m.h.Alloc(a)
		if err != nil {
			return value.Nil, err
		}
		return value.FromHandle(aid), nil
	}
	dims := []int64{int64(len(match))}
	a, err := heap.NewArray(dims)
	if err != nil {
		return value.Nil, err
	}
	for i, hid := range match {
		a.Floats[i] = float64(hid)
	}
	aid, err := m.h.Alloc(a)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(aid), nil
}

func (m *Module) lobbyGetName(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("LOBBY.GETNAME: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("LOBBY.GETNAME expects lobby handle")
	}
	o, err := heap.Cast[*lobbyObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	return rt.RetString(o.name), nil
}

func (m *Module) lobbyJoin(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("LOBBY.JOIN expects lobby handle")
	}
	o, err := heap.Cast[*lobbyObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if o.hostStr == "" || o.port <= 0 {
		return value.Nil, fmt.Errorf("LOBBY.JOIN: call LOBBY.SETHOST first")
	}
	return m.cliConnect(rt, rt.RetString(o.hostStr), value.FromFloat(float64(o.port)))
}
