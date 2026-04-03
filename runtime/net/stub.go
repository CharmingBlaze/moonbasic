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
		"NET.START", "NET.STOP", "NET.CREATESERVER", "NET.CREATECLIENT", "NET.CONNECT",
		"NET.UPDATE", "NET.RECEIVE", "NET.CLOSE", "NET.BROADCAST", "NET.PEERCOUNT",
		"NET.SETTIMEOUT", "NET.GETPING", "NET.SETBANDWIDTH",
		"PEER.SEND", "PEER.DISCONNECT", "PEER.IP", "PEER.PING",
		"EVENT.TYPE", "EVENT.PEER", "EVENT.DATA", "EVENT.FREE",
	}
	for _, k := range keys {
		reg.Register(k, "net", stub(k))
	}
}

func shutdownNet(m *Module) { _ = m }
