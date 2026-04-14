package worldmgr

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) Register(r runtime.Registrar) {
	r.Register("WORLD.GRAVITY", "world", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("WORLD.GRAVITY expects 3 floats (gx, gy, gz)")
		}
		return rt.Call("PHYSICS3D.SETGRAVITY", []value.Value{args[0], args[1], args[2]})
	})

	r.Register("WORLD.SETCENTER", "world", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return worldSetCenter(m, rt, args...) })
	r.Register("WORLD.SETCENTERENTITY", "world", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return worldSetCenterEntity(m, rt, args...) })
	r.Register("WORLD.UPDATE", "world", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return worldUpdate(m, rt, args...) })
	r.Register("WORLD.STREAMENABLE", "world", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return worldStreamEnable(m, rt, args...) })
	r.Register("WORLD.PRELOAD", "world", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return worldPreload(m, rt, args...) })
	r.Register("WORLD.STATUS", "world", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return worldStatus(m, rt, args...) })
	r.Register("WORLD.ISREADY", "world", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return worldIsReady(m, rt, args...) })
	r.Register("WORLD.SETVEGETATION", "world", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return worldSetVegetation(m, rt, args...) })

	r.Register("WORLD.SETREFLECTION", "world", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return m.worldSetReflection(rt, args...)
	})

	r.Register("WORLD.FOGMODE", "world", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return m.worldFogMode(rt, args...)
	})
	r.Register("FOGMODE", "world", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return m.worldFogMode(rt, args...)
	})
	r.Register("FOGCOLOR", "world", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return m.worldFogColor(rt, args...)
	})
	r.Register("WORLD.FOGCOLOR", "world", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return m.worldFogColor(rt, args...)
	})
	r.Register("WORLD.FOGDENSITY", "world", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return m.worldFogDensity(rt, args...)
	})
	r.Register("FOGDENSITY", "world", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return m.worldFogDensity(rt, args...)
	})
	r.Register("SKYCOLOR", "world", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return value.Nil, nil
	})
}
