//go:build linux && cgo

package mbphysics3d

import (
	"fmt"

	mbruntime "moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerPhysics3DCommands(m *Module, reg mbruntime.Registrar) {
	reg.Register("PHYSICS3D.START", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return phStart(m, a) }))
	reg.Register("PHYSICS3D.STOP", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return phStop(m, a) }))
	reg.Register("PHYSICS3D.SETGRAVITY", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return phSetGravity(m, a) }))
	reg.Register("WORLD.SETGRAVITY", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return phSetGravity(m, a) }))
	reg.Register("PHYSICS3D.STEP", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return phStep(m, a) }))
	reg.Register("PHYSICS3D.SETTIMESTEP", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return phSetTimeStep(m, a) }))
	reg.Register("PHYSICS3D.GETMATRIXBUFFER", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return phGetMatrixBuffer(m, a) }))
	reg.Register("PHYSICS3D.SETSUBSTEPS", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return phSetSubsteps(m, a) }))
	reg.Register("PHYSICS3D.ONCOLLISION", "physics3d", func(rt *mbruntime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 3 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindHandle || args[2].Kind != value.KindString {
			return value.Nil, fmt.Errorf("PHYSICS3D.ONCOLLISION expects (handle, handle, string)")
		}
		cb, err := rt.ArgString(args, 2)
		if err != nil {
			return value.Nil, err
		}
		joltMu.Lock()
		collRules = append(collRules, collRule{ha: heap.Handle(args[0].IVal), hb: heap.Handle(args[1].IVal), cb: cb})
		joltMu.Unlock()
		return value.Nil, nil
	})
	reg.Register("PHYSICS3D.PROCESSCOLLISIONS", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return phProcessCollisions(m, a) }))
	reg.Register("PHYSICS3D.RAYCAST", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return phRaycast(m, a) }))
	reg.Register("PHYSICS3D.SYNCWASMTOPHYSREGS", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return phSyncWasmToPhysRegs(m, a) }))
	reg.Register("PHYSICS3D.GETSCRATCHFLOAT", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return phGetScratchFloat(m, a) }))
	reg.Register("BODY3D.MAKE", "physics3d", func(rt *mbruntime.Runtime, args ...value.Value) (value.Value, error) {
		if m.h == nil {
			return value.Nil, mbruntime.Errorf("BODY3D.MAKE: heap not bound")
		}
		motion := "dynamic"
		if len(args) == 0 {
			// default motion type
		} else if len(args) == 1 && args[0].Kind == value.KindString {
			var err error
			motion, err = rt.ArgString(args, 0)
			if err != nil {
				return value.Nil, err
			}
		} else {
			return value.Nil, fmt.Errorf("BODY3D.MAKE expects 0 arguments (default DYNAMIC) or 1 motion string (STATIC, KINEMATIC, DYNAMIC)")
		}
		b := &BuilderObj{Motion: parseMotion(motion)}
		bid, err := m.h.Alloc(b)
		if err != nil {
			return value.Nil, err
		}
		return value.FromHandle(bid), nil
	})
	reg.Register("BODY3D.ADDBOX", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return BDAddBox(m.h, a) }))
	reg.Register("BODY3D.ADDSPHERE", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return BDAddSphere(m.h, a) }))
	reg.Register("BODY3D.ADDCAPSULE", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return BDAddCapsule(m.h, a) }))
	reg.Register("BODY3D.ADDMESH", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdAddMesh(m, a) }))
	reg.Register("BODY3D.COMMIT", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return BDCommit(m.h, a) }))
	reg.Register("BODY3D.SETPOS", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdSetPos(m, a) }))
	reg.Register("BODY3D.SETPOSITION", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdSetPos(m, a) }))
	reg.Register("BODY3D.GETPOS", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdGetPos(m, a) }))
	reg.Register("BODY3D.ACTIVATE", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdActivate(m, a) }))
	reg.Register("BODY3D.DEACTIVATE", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdDeactivate(m, a) }))
	reg.Register("BODY3D.SETROT", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdSetRotation(m, a) }))
	reg.Register("BODY3D.GETROT", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdGetRotZero(m, a) }))
	reg.Register("BODY3D.SETMASS", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdNoOp(m, a) }))
	reg.Register("BODY3D.SETFRICTION", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdSetFriction(m, a) }))
	reg.Register("BODY3D.SETRESTITUTION", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdSetRestitution(m, a) }))
	reg.Register("BODY3D.APPLYFORCE", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdApplyForce(m, a) }))
	reg.Register("BODY3D.APPLYIMPULSE", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdApplyImpulse(m, a) }))
	reg.Register("BODY3D.SETLINEARVEL", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdSetLinearVel(m, a) }))
	reg.Register("BODY3D.SETANGULARVEL", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdNoOp(m, a) }))
	reg.Register("BODY3D.X", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdAxis(m, a, 0) }))
	reg.Register("BODY3D.Y", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdAxis(m, a, 1) }))
	reg.Register("BODY3D.Z", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdAxis(m, a, 2) }))
	reg.Register("BODY3D.BUFFERINDEX", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return BDBufferIndex(m.h, a) }))
	reg.Register("BODY3D.FREE", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdFree(m, a) }))

	// Blitz-style PHYSICS.* aliases (same implementation as PHYSICS3D.*).
	reg.Register("PHYSICS.START", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return phStart(m, a) }))
	reg.Register("PHYSICS.STOP", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return phStop(m, a) }))
	reg.Register("PHYSICS.SETGRAVITY", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return phSetGravity(m, a) }))
	reg.Register("PHYSICS.STEP", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return phStep(m, a) }))
	reg.Register("PHYSICS.SETSUBSTEPS", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return phSetSubsteps(m, a) }))
	reg.Register("PHYSICS.RAYCAST", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return phRaycast(m, a) }))
	reg.Register("PHYSICS.SPHERECAST", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) {
		return value.Nil, fmt.Errorf("PHYSICS.SPHERECAST: not implemented; use PHYSICS.RAYCAST or BODY3D queries")
	}))
	reg.Register("PHYSICS.BOXCAST", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) {
		return value.Nil, fmt.Errorf("PHYSICS.BOXCAST: not implemented; use PHYSICS.RAYCAST or BODY3D queries")
	}))
	reg.Register("PHYSICS.ENABLE", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) {
		return value.Nil, fmt.Errorf("PHYSICS.ENABLE: use BODY3D.ACTIVATE on a physics body handle")
	}))
	reg.Register("PHYSICS.DISABLE", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) {
		return value.Nil, fmt.Errorf("PHYSICS.DISABLE: use BODY3D.DEACTIVATE on a physics body handle")
	}))
	reg.Register("PHYSICS3D.DEBUGDRAW", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) {
		return value.Nil, fmt.Errorf("PHYSICS3D.DEBUGDRAW: not implemented for Jolt in this runtime")
	}))
	reg.Register("BODY3D.COLLIDED", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdCollided3D(m, a) }))
	reg.Register("BODY3D.COLLISIONOTHER", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdCollisionOther3D(m, a) }))
	reg.Register("BODY3D.COLLISIONPOINT", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdCollisionPoint3D(m, a) }))
	reg.Register("BODY3D.COLLISIONNORMAL", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdCollisionNormal3D(m, a) }))
	reg.Register("JOINT3D.FIXED", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) {
		return value.Nil, fmt.Errorf("JOINT3D.FIXED: Jolt constraints not exposed in jolt-go wrapper")
	}))
	reg.Register("JOINT3D.HINGE", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) {
		return value.Nil, fmt.Errorf("JOINT3D.HINGE: Jolt constraints not exposed in jolt-go wrapper")
	}))
	reg.Register("JOINT3D.SLIDER", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) {
		return value.Nil, fmt.Errorf("JOINT3D.SLIDER: Jolt constraints not exposed in jolt-go wrapper")
	}))
	reg.Register("JOINT3D.CONE", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) {
		return value.Nil, fmt.Errorf("JOINT3D.CONE: Jolt constraints not exposed in jolt-go wrapper")
	}))
	reg.Register("JOINT3D.DELETE", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) {
		return value.Nil, fmt.Errorf("JOINT3D.DELETE: no joint handles in this runtime")
	}))
	registerPickCommands(m, reg)
}

func shutdownPhysics3D(m *Module) {
	_, _ = phStop(m, nil)
}
