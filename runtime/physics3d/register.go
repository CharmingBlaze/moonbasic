package mbphysics3d

import (
	"fmt"
	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerPhysics3DCommands(m *Module, reg runtime.Registrar) {
	// World / Controller API
	reg.Register("PHYSICS3D.START", "physics3d", runtime.AdaptLegacy(m.phStart))
	reg.Register("PHYSICS3D.STOP", "physics3d", runtime.AdaptLegacy(m.phStop))
	reg.Register("PHYSICS3D.SETGRAVITY", "physics3d", runtime.AdaptLegacy(m.phSetGravity))
	reg.Register("PHYSICS3D.GETGRAVITYX", "physics3d", runtime.AdaptLegacy(m.phGetGravityX))
	reg.Register("PHYSICS3D.GETGRAVITYY", "physics3d", runtime.AdaptLegacy(m.phGetGravityY))
	reg.Register("PHYSICS3D.GETGRAVITYZ", "physics3d", runtime.AdaptLegacy(m.phGetGravityZ))
	reg.Register("PHYSICS3D.STEP", "physics3d", runtime.AdaptLegacy(m.phStep))
	reg.Register("PHYSICS3D.SETTIMESTEP", "physics3d", runtime.AdaptLegacy(m.phSetTimeStep))
	reg.Register("PHYSICS3D.GETMATRIXBUFFER", "physics3d", runtime.AdaptLegacy(m.phGetMatrixBuffer))
	reg.Register("PHYSICS3D.SETSUBSTEPS", "physics3d", runtime.AdaptLegacy(m.phSetSubsteps))
	reg.Register("PHYSICS3D.ONCOLLISION", "physics3d", runtime.AdaptLegacy(m.phOnCollision))
	reg.Register("PHYSICS3D.PROCESSCOLLISIONS", "physics3d", runtime.AdaptLegacy(m.phProcessCollisions))
	reg.Register("PHYSICS3D.RAYCAST", "physics3d", runtime.AdaptLegacy(m.phRaycast))
	reg.Register("PHYSICS3D.SYNCWASMTOPHYSREGS", "physics3d", runtime.AdaptLegacy(m.phSyncWasmToPhysRegs))
	reg.Register("PHYSICS3D.GETSCRATCHFLOAT", "physics3d", runtime.AdaptLegacy(m.phGetScratchFloat))
	reg.Register("PHYSICS3D.DEBUGDRAW", "physics3d", runtime.AdaptLegacy(m.phDebugDraw))

	// World API (Aliases / Easy Mode)
	reg.Register("WORLD.SETUP", "physics3d", runtime.AdaptLegacy(m.phWorldSetup))
	reg.Register("WORLD.SETGRAVITY", "physics3d", runtime.AdaptLegacy(m.phSetGravity))
	reg.Register("PHYSICS.START", "physics3d", runtime.AdaptLegacy(m.phStart))
	reg.Register("PHYSICS.STOP", "physics3d", runtime.AdaptLegacy(m.phStop))
	reg.Register("PHYSICS.SETGRAVITY", "physics3d", runtime.AdaptLegacy(m.phSetGravity))
	reg.Register("PHYSICS.GETGRAVITYX", "physics3d", runtime.AdaptLegacy(m.phGetGravityX))
	reg.Register("PHYSICS.GETGRAVITYY", "physics3d", runtime.AdaptLegacy(m.phGetGravityY))
	reg.Register("PHYSICS.GETGRAVITYZ", "physics3d", runtime.AdaptLegacy(m.phGetGravityZ))
	reg.Register("PHYSICS.STEP", "physics3d", runtime.AdaptLegacy(m.phStep))
	reg.Register("PHYSICS.SETSUBSTEPS", "physics3d", runtime.AdaptLegacy(m.phSetSubsteps))
	reg.Register("PHYSICS.RAYCAST", "physics3d", runtime.AdaptLegacy(m.phRaycast))
	reg.Register("PHYSICS.SPHERECAST", "physics3d", runtime.AdaptLegacy(m.phSpherecast))
	reg.Register("PHYSICS.BOXCAST", "physics3d", runtime.AdaptLegacy(m.phBoxcast))
	reg.Register("PHYSICS.ENABLE", "physics3d", runtime.AdaptLegacy(m.phEnable))
	reg.Register("PHYSICS.DISABLE", "physics3d", runtime.AdaptLegacy(m.phDisable))

	// Shape API
	reg.Register("SHAPE.CREATEBOX", "physics3d", runtime.AdaptLegacy(m.shCreateBox))
	reg.Register("SHAPE.CREATESPHERE", "physics3d", runtime.AdaptLegacy(m.shCreateSphere))
	reg.Register("SHAPE.CREATECAPSULE", "physics3d", runtime.AdaptLegacy(m.shCreateCapsule))
	reg.Register("SHAPE.CREATECYLINDER", "physics3d", runtime.AdaptLegacy(m.shCreateCylinder))
	reg.Register("SHAPE.GETTYPE", "physics3d", runtime.AdaptLegacy(m.shGetType))
	reg.Register("SHAPE.GETWIDTH", "physics3d", runtime.AdaptLegacy(m.shGetDim1))
	reg.Register("SHAPE.GETHEIGHT", "physics3d", runtime.AdaptLegacy(m.shGetDim2))
	reg.Register("SHAPE.GETDEPTH", "physics3d", runtime.AdaptLegacy(m.shGetDim3))
	reg.Register("SHAPE.GETRADIUS", "physics3d", runtime.AdaptLegacy(m.shGetDim1))
	reg.Register("SHAPE.GETSIZEX", "physics3d", runtime.AdaptLegacy(m.shGetDim1))
	reg.Register("SHAPE.GETSIZEY", "physics3d", runtime.AdaptLegacy(m.shGetDim2))
	reg.Register("SHAPE.GETSIZEZ", "physics3d", runtime.AdaptLegacy(m.shGetDim3))
	reg.Register("SHAPEREF.FREE", "physics3d", runtime.AdaptLegacy(m.shFree))

	// Body API (High-level)
	reg.Register("KINEMATIC.CREATE", "physics3d", runtime.AdaptLegacy(m.knCreate))
	reg.Register("KINEMATIC.MAKE", "physics3d", runtime.AdaptLegacy(m.knCreate))
	reg.Register("STATIC.CREATE", "physics3d", runtime.AdaptLegacy(m.stCreate))
	reg.Register("STATIC.MAKE", "physics3d", runtime.AdaptLegacy(m.stCreate))
	reg.Register("TRIGGER.CREATE", "physics3d", runtime.AdaptLegacy(m.trCreate))
	reg.Register("TRIGGER.MAKE", "physics3d", runtime.AdaptLegacy(m.trCreate))
	reg.Register("BODY3D.CREATE", "physics3d", runtime.AdaptLegacy(m.phBodyMake))
	reg.Register("BODY3D.MAKE", "physics3d", runtime.AdaptLegacy(m.phBodyMake))
	reg.Register("BODY3D.ADDBOX", "physics3d", runtime.AdaptLegacy(m.bdAddBox))
	reg.Register("BODY3D.ADDSPHERE", "physics3d", runtime.AdaptLegacy(m.bdAddSphere))
	reg.Register("BODY3D.ADDCAPSULE", "physics3d", runtime.AdaptLegacy(m.bdAddCapsule))
	reg.Register("BODY3D.ADDMESH", "physics3d", runtime.AdaptLegacy(m.bdAddMesh))
	reg.Register("BODY3D.COMMIT", "physics3d", runtime.AdaptLegacy(m.bdCommit))
	reg.Register("BODY3D.SETPOS", "physics3d", runtime.AdaptLegacy(m.bdSetPos))
	reg.Register("BODY3D.SETPOSITION", "physics3d", runtime.AdaptLegacy(m.bdSetPos))
	reg.Register("BODY3D.GETPOS", "physics3d", runtime.AdaptLegacy(m.bdGetPos))
	reg.Register("BODY3D.ACTIVATE", "physics3d", runtime.AdaptLegacy(m.bdActivate))
	reg.Register("BODY3D.DEACTIVATE", "physics3d", runtime.AdaptLegacy(m.bdDeactivate))
	reg.Register("BODY3D.SETROT", "physics3d", runtime.AdaptLegacy(m.bdSetRotation))
	reg.Register("BODY3D.GETROT", "physics3d", runtime.AdaptLegacy(m.bdGetRotZero))
	reg.Register("BODY3D.GETSCALE", "physics3d", runtime.AdaptLegacy(m.bdGetScale))
	reg.Register("BODY3D.SETSCALE", "physics3d", runtime.AdaptLegacy(m.bdSetScale))
	reg.Register("BODY3D.SETMASS", "physics3d", runtime.AdaptLegacy(m.bdSetMass))
	reg.Register("BODY3D.SETFRICTION", "physics3d", runtime.AdaptLegacy(m.bdSetFriction))
	reg.Register("BODY3D.SETRESTITUTION", "physics3d", runtime.AdaptLegacy(m.bdSetRestitution))
	reg.Register("BODY3D.APPLYFORCE", "physics3d", runtime.AdaptLegacy(m.bdApplyForce))
	reg.Register("BODY3D.APPLYIMPULSE", "physics3d", runtime.AdaptLegacy(m.bdApplyImpulse))
	reg.Register("BODY3D.SETVELOCITY", "physics3d", runtime.AdaptLegacy(m.bdSetLinearVel))
	reg.Register("BODY3D.SETLINEARVEL", "physics3d", runtime.AdaptLegacy(m.bdSetLinearVel))
	reg.Register("BODY3D.GETVELOCITY", "physics3d", runtime.AdaptLegacy(m.bdGetLinearVel))
	reg.Register("BODY3D.GETLINEARVEL", "physics3d", runtime.AdaptLegacy(m.bdGetLinearVel))
	reg.Register("BODY3D.GETANGULARVEL", "physics3d", runtime.AdaptLegacy(m.bdGetAngularVel))
	reg.Register("BODY3D.SETANGULARVEL", "physics3d", runtime.AdaptLegacy(m.bdSetAngularVel))
	reg.Register("BODY3D.GETMASS", "physics3d", runtime.AdaptLegacy(m.bdGetMass))
	reg.Register("BODY3D.APPLYTORQUE", "physics3d", runtime.AdaptLegacy(m.bdApplyTorque))
	reg.Register("BODY3D.X", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return m.bdAxis(a, 0) }))
	reg.Register("BODY3D.Y", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return m.bdAxis(a, 1) }))
	reg.Register("BODY3D.Z", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return m.bdAxis(a, 2) }))
	reg.Register("BODY3D.BUFFERINDEX", "physics3d", runtime.AdaptLegacy(m.bdBufferIndex))
	reg.Register("BODY3D.FREE", "physics3d", runtime.AdaptLegacy(m.bdFree))
	reg.Register("BODY3D.COLLIDED", "physics3d", runtime.AdaptLegacy(m.bdCollided3D))
	reg.Register("BODY3D.COLLISIONOTHER", "physics3d", runtime.AdaptLegacy(m.bdCollisionOther3D))
	reg.Register("BODY3D.COLLISIONPOINT", "physics3d", runtime.AdaptLegacy(m.bdCollisionPoint3D))
	reg.Register("BODY3D.COLLISIONNORMAL", "physics3d", runtime.AdaptLegacy(m.bdCollisionNormal3D))

	// Body API (Shared Handle Methods)
	reg.Register("BODYREF.SETPOS", "physics3d", runtime.AdaptLegacy(m.brSetPos))
	reg.Register("BODYREF.SETPOSITION", "physics3d", runtime.AdaptLegacy(m.brSetPos))
	reg.Register("BODYREF.GETPOSITION", "physics3d", runtime.AdaptLegacy(m.bdGetPos))
	reg.Register("BODYREF.SETROTATION", "physics3d", runtime.AdaptLegacy(m.bdSetRotation))
	reg.Register("BODYREF.GETROTATION", "physics3d", runtime.AdaptLegacy(m.bdGetRotZero))
	reg.Register("BODYREF.GETVELOCITY", "physics3d", runtime.AdaptLegacy(m.bdGetLinearVel))
	reg.Register("BODYREF.SETVELOCITY", "physics3d", runtime.AdaptLegacy(m.bdSetLinearVel))
	reg.Register("BODYREF.SETLAYER", "physics3d", runtime.AdaptLegacy(m.brSetLayer))
	reg.Register("BODYREF.ENABLECOLLISION", "physics3d", runtime.AdaptLegacy(m.brEnableColl))
	reg.Register("BODYREF.FREE", "physics3d", runtime.AdaptLegacy(m.brFree))

	// Advanced Body Control
	reg.Register("BODY3D.LOCKAXIS", "physics3d", runtime.AdaptLegacy(m.bdLockAxis))
	reg.Register("BODY3D.SETGRAVITYFACTOR", "physics3d", runtime.AdaptLegacy(m.bdSetGravityFactor))
	reg.Register("BODY3D.SETDAMPING", "physics3d", runtime.AdaptLegacy(m.bdSetDamping))
	reg.Register("BODY3D.SETCCD", "physics3d", runtime.AdaptLegacy(m.btdSetCCD))

	// Joints
	reg.Register("JOINT.CREATEHINGE", "physics3d", runtime.AdaptLegacy(m.phCreateHingeJoint))
	reg.Register("JOINT.CREATEPOINT", "physics3d", runtime.AdaptLegacy(m.phCreatePointJoint))
	reg.Register("JOINT.FREE", "physics3d", runtime.AdaptLegacy(m.phJointDelete))
	reg.Register("JOINT3D.FIXED", "physics3d", runtime.AdaptLegacy(m.phJointFixed))
	reg.Register("JOINT3D.HINGE", "physics3d", runtime.AdaptLegacy(m.phJointHinge))
	reg.Register("JOINT3D.SLIDER", "physics3d", runtime.AdaptLegacy(m.phJointSlider))
	reg.Register("JOINT3D.CONE", "physics3d", runtime.AdaptLegacy(m.phJointCone))
	reg.Register("JOINT3D.DELETE", "physics3d", runtime.AdaptLegacy(m.phJointDelete))

	// Kinematic Specific
	reg.Register("KINEMATICREF.SETVELOCITY", "physics3d", runtime.AdaptLegacy(m.bdSetLinearVel))
	reg.Register("KINEMATICREF.UPDATE", "physics3d", runtime.AdaptLegacy(m.bdNoOp))

	// Debug / Picking
	reg.Register("DEBUG.DRAWCHARACTER", "physics3d", runtime.AdaptLegacy(m.bdNoOp))
	reg.Register("DEBUG.DRAWBODY", "physics3d", runtime.AdaptLegacy(m.bdNoOp))
	reg.Register("VEHICLE.CREATE", "physics3d", runtime.AdaptLegacy(m.VHCreate))
	reg.Register("VEHICLE.SETWHEEL", "physics3d", runtime.AdaptLegacy(m.VHSetWheel))
	reg.Register("VEHICLE.SETTUNING", "physics3d", runtime.AdaptLegacy(m.VHSetTuning))
	reg.Register("VEHICLE.SETSTEER", "physics3d", runtime.AdaptLegacy(m.VHSetSteering))
	reg.Register("VEHICLE.SETTHROTTLE", "physics3d", runtime.AdaptLegacy(m.VHSetThrottle))
	reg.Register("VEHICLE.CONTROL", "physics3d", runtime.AdaptLegacy(m.VHControl))
	reg.Register("VEHICLE.STEP", "physics3d", runtime.AdaptLegacy(m.VHStep))
	reg.Register("VEHICLE.WHEELX", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return m.VHWheelAxis(a, 0) }))
	reg.Register("VEHICLE.WHEELY", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return m.VHWheelAxis(a, 1) }))
	reg.Register("VEHICLE.WHEELZ", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return m.VHWheelAxis(a, 2) }))
	registerAeroCommands(m, reg)
	registerPickCommands(m, reg)
}

func (m *Module) phOnCollision(args []value.Value) (value.Value, error) {
	if len(args) != 3 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindHandle || args[2].Kind != value.KindString {
		return value.Nil, fmt.Errorf("PHYSICS3D.ONCOLLISION expects (handle, handle, string)")
	}
	cb := args[2].String()
	return phSetOnCollision(m, args[0], args[1], cb)
}

func (m *Module) phBodyMake(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("BODY3D.MAKE: heap not bound")
	}
	motion := "dynamic"
	if len(args) == 0 {
		// default motion type
	} else if len(args) == 1 && args[0].Kind == value.KindString {
		motion = args[0].String()
	} else {
		return value.Nil, fmt.Errorf("BODY3D.MAKE expects 0 arguments (default DYNAMIC) or 1 motion string (STATIC, KINEMATIC, DYNAMIC)")
	}
	return phCreateBody(m, motion)
}
