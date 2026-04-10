//go:build !cgo && !windows

package mbnav

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

const hint = "NAV/PATH/NAVAGENT/STEER/BTREE natives require CGO: set CGO_ENABLED=1 and install a C compiler, then rebuild"

func (m *Module) Register(reg runtime.Registrar) {
	stub := func(name string) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			_ = rt
			_ = args
			return value.Nil, fmt.Errorf("%s: %s", name, hint)
		}
	}
	for _, n := range []string{
		"NAV.MAKE", "NAV.FREE", "NAV.SETGRID", "NAV.ADDTERRAIN", "NAV.ADDOBSTACLE", "NAV.BUILD", "NAV.FINDPATH",
		"NAV.BAKE", "NAV.GETPATH", "NAV.ISREACHABLE",
	} {
		reg.Register(n, "nav", stub(n))
	}
	reg.Register("ENEMY.FOLLOWPATH", "enemy", stub("ENEMY.FOLLOWPATH"))
	for _, n := range []string{
		"PATH.ISVALID", "PATH.NODECOUNT", "PATH.NODEX", "PATH.NODEY", "PATH.NODEZ", "PATH.FREE",
	} {
		reg.Register(n, "path", stub(n))
	}
	for _, n := range []string{
		"NAVAGENT.MAKE", "NAVAGENT.FREE", "NAVAGENT.SETPOS", "NAVAGENT.SETSPEED", "NAVAGENT.SETMAXFORCE", "NAVAGENT.APPLYFORCE",
		"NAVAGENT.MOVETO", "NAVAGENT.UPDATE", "NAVAGENT.ISATDESTINATION", "NAVAGENT.X", "NAVAGENT.Y", "NAVAGENT.Z",
	} {
		reg.Register(n, "navagent", stub(n))
	}
	for _, n := range []string{
		"STEER.GROUPMAKE", "STEER.GROUPADD", "STEER.GROUPCLEAR",
		"STEER.SEEK", "STEER.FLEE", "STEER.ARRIVE", "STEER.WANDER", "STEER.FLOCK", "STEER.AVOIDOBSTACLES", "STEER.FOLLOWPATH",
	} {
		reg.Register(n, "steer", stub(n))
	}
	for _, n := range []string{
		"BTREE.MAKE", "BTREE.FREE", "BTREE.SEQUENCE", "BTREE.ADDCONDITION", "BTREE.ADDACTION", "BTREE.RUN",
	} {
		reg.Register(n, "btree", stub(n))
	}
}

func (m *Module) Shutdown() {}
