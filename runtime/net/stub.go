//go:build !cgo

package mbnet

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

const stubHint = "NET/PEER/EVENT require CGO and ENet (github.com/codecat/go-enet + libenet). Build with CGO_ENABLED=1."

func registerNetCommands(m *Module, reg runtime.Registrar) {
	_ = m
	stub := func(name string) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			_ = rt
			return value.Nil, fmt.Errorf("%s: %s", name, stubHint)
		}
	}
	keys := []string{
		"ENET.INITIALIZE", "ENET.DEINITIALIZE", "ENET.CREATEHOST", "ENET.MAKEHOST",
		"ENET.HOSTSERVICE", "ENET.HOSTBROADCAST", "ENET.PEERSEND", "ENET.PEERPING",
		"NET.START", "NET.STOP", "NET.CREATESERVER", "NET.CREATECLIENT", "NET.CONNECT",
		"NET.UPDATE", "NET.RECEIVE", "NET.CLOSE", "NET.BROADCAST", "NET.PEERCOUNT",
		"NET.SETTIMEOUT", "NET.GETPING", "NET.SETBANDWIDTH",
		"NET.SETCHANNELS", "NET.SERVICE", "NET.FLUSH",
		"NETSENDSTRING", "NETSENDINT", "NETSENDFLOAT", "NETREADSTRING", "NETREADINT", "NETREADFLOAT",
		"PACKET.CREATE", "PACKET.MAKE", "PACKET.DATA", "PACKET.FREE",
		"PEER.SEND", "PEER.SENDPACKET", "PEER.DISCONNECT", "PEER.IP", "PEER.PING",
		"EVENT.TYPE", "EVENT.PEER", "EVENT.DATA", "EVENT.FREE", "EVENT.CHANNEL",
		"SERVER.START", "SERVER.STOP", "SERVER.ONCONNECT", "SERVER.ONDISCONNECT", "SERVER.ONMESSAGE",
		"SERVER.SYNCENTITY", "SERVER.SETTICKRATE", "SERVER.TICK",
		"CLIENT.CONNECT", "CLIENT.STOP", "CLIENT.ONCONNECT", "CLIENT.ONMESSAGE", "CLIENT.ONSYNC", "CLIENT.TICK",
		"RPC.CALL", "RPC.CALLTO", "RPC.CALLSERVER",
		"LOBBY.CREATE", "LOBBY.MAKE", "LOBBY.FREE", "LOBBY.SETPROPERTY", "LOBBY.SETHOST", "LOBBY.START",
		"LOBBY.FIND", "LOBBY.GETNAME", "LOBBY.JOIN",
	}
	for _, k := range keys {
		reg.Register(k, "net", stub(k))
	}
}

func shutdownNet(m *Module) {
	resetMultiplayerState()
	_ = m
}
