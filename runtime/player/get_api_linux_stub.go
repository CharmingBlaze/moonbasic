//go:build linux && !cgo

package player

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerPlayerCharGetAPI(m *Module, reg runtime.Registrar) {
	zf := func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		_ = args
		return value.FromFloat(0), nil
	}
	zb := func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		_ = args
		return value.FromBool(false), nil
	}
	zi := func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		_ = args
		return value.FromInt(0), nil
	}
	zCollOn := func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		_ = args
		return value.FromBool(true), nil
	}
	shapeStub := func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if _, ok := m.kccSubjectID(args); !ok {
			if len(args) < 1 {
				return value.Nil, fmt.Errorf("PLAYER.GETSHAPETYPE: %s", kccErrNoSubject)
			}
		}
		if m.h == nil {
			return value.Nil, fmt.Errorf("PLAYER.GETSHAPETYPE: heap not bound")
		}
		return value.FromStringIndex(m.h.Intern("capsule")), nil
	}

	reg.Register("PLAYER.GETPOSITIONX", "player", zf)
	reg.Register("PLAYER.GETPOSITIONY", "player", zf)
	reg.Register("PLAYER.GETPOSITIONZ", "player", zf)
	reg.Register("PLAYER.GETROTATIONPITCH", "player", zf)
	reg.Register("PLAYER.GETROTATIONYAW", "player", zf)
	reg.Register("PLAYER.GETROTATIONROLL", "player", zf)
	reg.Register("PLAYER.GETVELOCITYX", "player", zf)
	reg.Register("PLAYER.GETVELOCITYY", "player", zf)
	reg.Register("PLAYER.GETVELOCITYZ", "player", zf)
	reg.Register("PLAYER.GETSPEED", "player", zf)
	reg.Register("PLAYER.GETONSLOPE", "player", zb)
	reg.Register("PLAYER.GETONWALL", "player", zb)
	reg.Register("PLAYER.GETSLOPEANGLE", "player", zf)
	reg.Register("PLAYER.GETISJUMPING", "player", zb)
	reg.Register("PLAYER.GETISFALLING", "player", zb)
	reg.Register("PLAYER.GETMAXSLOPE", "player", zf)
	reg.Register("PLAYER.GETSTEPHEIGHT", "player", zf)
	reg.Register("PLAYER.GETGRAVITYSCALE", "player", zf)
	reg.Register("PLAYER.GETFRICTION", "player", zf)
	reg.Register("PLAYER.GETSNAPDISTANCE", "player", zf)
	reg.Register("PLAYER.GETHEIGHT", "player", zf)
	reg.Register("PLAYER.GETRADIUS", "player", zf)
	reg.Register("PLAYER.GETLAYER", "player", zi)
	reg.Register("PLAYER.GETMASK", "player", zi)
	reg.Register("PLAYER.GETCOLLISIONENABLED", "player", zCollOn)

	reg.Register("CHAR.GETPOSITIONX", "player", zf)
	reg.Register("CHAR.GETPOSITIONY", "player", zf)
	reg.Register("CHAR.GETPOSITIONZ", "player", zf)
	reg.Register("CHAR.GETROTATIONPITCH", "player", zf)
	reg.Register("CHAR.GETROTATIONYAW", "player", zf)
	reg.Register("CHAR.GETROTATIONROLL", "player", zf)
	reg.Register("CHAR.GETVELOCITYX", "player", zf)
	reg.Register("CHAR.GETVELOCITYY", "player", zf)
	reg.Register("CHAR.GETVELOCITYZ", "player", zf)
	reg.Register("CHAR.GETSPEED", "player", zf)
	reg.Register("CHAR.GETONSLOPE", "player", zb)
	reg.Register("CHAR.GETONWALL", "player", zb)
	reg.Register("CHAR.GETSLOPEANGLE", "player", zf)
	reg.Register("CHAR.GETISJUMPING", "player", zb)
	reg.Register("CHAR.GETISFALLING", "player", zb)
	reg.Register("CHAR.GETMAXSLOPE", "player", zf)
	reg.Register("CHAR.GETSTEPHEIGHT", "player", zf)
	reg.Register("CHAR.GETGRAVITYSCALE", "player", zf)
	reg.Register("CHAR.GETFRICTION", "player", zf)
	reg.Register("CHAR.GETSNAPDISTANCE", "player", zf)
	reg.Register("CHAR.GETHEIGHT", "player", zf)
	reg.Register("CHAR.GETRADIUS", "player", zf)
	reg.Register("CHAR.GETLAYER", "player", zi)
	reg.Register("CHAR.GETMASK", "player", zi)
	reg.Register("CHAR.GETCOLLISIONENABLED", "player", zCollOn)

	reg.Register("PLAYER.GETX", "player", zf)
	reg.Register("PLAYER.GETY", "player", zf)
	reg.Register("PLAYER.GETZ", "player", zf)
	reg.Register("PLAYER.GETPITCH", "player", zf)
	reg.Register("PLAYER.GETYAW", "player", zf)
	reg.Register("PLAYER.GETROLL", "player", zf)
	reg.Register("PLAYER.GETGROUNDED", "player", zb)
	reg.Register("PLAYER.GETGRAVITY", "player", zf)
	reg.Register("PLAYER.GETCAPSULERADIUS", "player", zf)
	reg.Register("PLAYER.GETCAPSULEHEIGHT", "player", zf)
	reg.Register("PLAYER.GETSHAPETYPE", "player", shapeStub)

	reg.Register("CHAR.GETX", "player", zf)
	reg.Register("CHAR.GETY", "player", zf)
	reg.Register("CHAR.GETZ", "player", zf)
	reg.Register("CHAR.GETPITCH", "player", zf)
	reg.Register("CHAR.GETYAW", "player", zf)
	reg.Register("CHAR.GETROLL", "player", zf)
	reg.Register("CHAR.GETGROUNDED", "player", zb)
	reg.Register("CHAR.GETGRAVITY", "player", zf)
	reg.Register("CHAR.GETCAPSULERADIUS", "player", zf)
	reg.Register("CHAR.GETCAPSULEHEIGHT", "player", zf)
	reg.Register("CHAR.GETSHAPETYPE", "player", shapeStub)
}
