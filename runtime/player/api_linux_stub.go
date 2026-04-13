//go:build linux && !cgo

package player

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerPlayerCommands(m *Module, reg runtime.Registrar) {
	stub := func(name string) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			_ = rt
			_ = args
			return value.Nil, fmt.Errorf("%s requires CGO_ENABLED=1 (Linux Jolt fullruntime)", name)
		}
	}
	type pair struct {
		name string
		ns   string
	}
	cmds := []pair{
		{"PLAYER.CREATE", "player"},
		{"PLAYER.MOVE", "player"},
		{"PLAYER.JUMP", "player"},
		{"PLAYER.ISGROUNDED", "player"},
		{"PLAYER.GETLOOKTARGET", "player"},
		{"PLAYER.GETNEARBY", "player"},
		{"ENT.GET_NEAREST", "player"},
		{"ENT.GETNEAREST", "player"},
		{"PLAYER.ONTRIGGER", "player"},
		{"PLAYER.SETSTATE", "player"},
		{"PLAYER.SYNCANIM", "player"},
		{"PLAYER.SETSTEPHEIGHT", "player"},
		{"PLAYER.SETSLOPELIMIT", "player"},
		{"PLAYER.GETVELOCITY", "player"},
		{"PLAYER.TELEPORT", "player"},
		{"PLAYER.SETGRAVITYSCALE", "player"},
		{"PLAYER.GETCROUCH", "player"},
		{"PLAYER.SETCROUCH", "player"},
		{"PLAYER.SWIM", "player"},
		{"PLAYER.SETSTEPOFFSET", "player"},
		{"PLAYER.SETSTICKFLOOR", "player"},
		{"PLAYER.NAVTO", "player"},
		{"PLAYER.NAVUPDATE", "player"},
		{"PLAYER.SETPADDING", "player"},
		{"PLAYER.MOVEWITHCAMERA", "player"},
		{"NAV.GOTO", "player"},
		{"NAV.UPDATE", "player"},
		{"NAV.CHASE", "player"},
		{"NAV.PATROL", "player"},
		{"CHAR.MAKE", "player"},
		{"CHAR.SETSTEP", "player"},
		{"CHAR.SETSLOPE", "player"},
		{"CHAR.SETPADDING", "player"},
		{"CHAR.MOVE", "player"},
		{"CHAR.MOVEWITHCAMERA", "player"},
		{"CHAR.MOVEWITHCAM", "player"},
		{"CHAR.NAVTO", "player"},
		{"CHAR.NAVUPDATE", "player"},
		{"CHAR.STICK", "player"},
		{"CHAR.ISGROUNDED", "player"},
		{"CHAR.JUMP", "player"},
		{"ENTITYREF.NAVUPDATE", "entity"},
		{"ENTITYREF.ISGROUNDED", "entity"},
		{"ENTITYREF.JUMP", "entity"},
		{"PLAYER.GETSTANDNORMAL", "player"},
		{"PLAYER.PUSH", "player"},
		{"PLAYER.GRAB", "player"},
		{"PLAYER.SETMASS", "player"},
		{"PLAYER.GETSURFACETYPE", "player"},
		{"PLAYER.SETFOVKICK", "player"},
		{"PLAYER.GETFOVKICK", "player"},
		{"PLAYER.ISMOVING", "player"},
		{"PLAYER.GETGROUNDSTATE", "player"},
		{"PLAYER.ISONSTEEPSLOPE", "player"},
		{"CHAR.GETGROUNDSTATE", "player"},
		{"CHAR.ISONSTEEPSLOPE", "player"},
		{"CHARACTER.CREATE", "player"},
		{"CHARACTERREF.ADDVELOCITY", "player"},
		{"CHARACTERREF.SETLINEARVELOCITY", "player"},
		{"CHARACTERREF.SETVELOCITY", "player"},
		{"CHARACTERREF.SETSNAPDISTANCE", "player"},
		{"CHARACTERREF.SETSTICKDOWN", "player"},
		{"CHARACTERREF.UPDATE", "player"},
		{"CHARACTERREF.UPDATEMOVE", "player"},
		{"CHARACTERREF.JUMP", "player"},
		{"CHARACTERREF.MOVEWITHCAMERA", "player"},
		{"CHARACTERREF.SETMAXSLOPE", "player"},
		{"CHARACTERREF.SETSTEPHEIGHT", "player"},
		{"CHARACTERREF.ISGROUNDED", "player"},
		{"CHARACTERREF.SETPOSITION", "player"},
		{"CHARACTERREF.GETPOSITION", "player"},
		{"CHARACTERREF.FREE", "player"},
		{"CHARACTERREF.GETGROUNDSTATE", "player"},
		{"CHARACTERREF.SETGRAVITY", "player"},
		{"CHARACTERREF.SETGRAVITYSCALE", "player"},
		{"CHARACTERREF.SETFRICTION", "player"},
		{"CHARACTERREF.SETPADDING", "player"},
		{"CHARACTERREF.SETBOUNCE", "player"},
		{"CHARACTERREF.GETSPEED", "player"},
		{"CHARACTERREF.ISMOVING", "player"},
	}
	for _, c := range cmds {
		reg.Register(c.name, c.ns, stub(c.name))
	}
	registerPlayerCharGetAPI(m, reg)
	registerPlayerTerrainCommands(m, reg)
}
