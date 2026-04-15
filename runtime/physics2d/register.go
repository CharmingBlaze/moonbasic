package mbphysics2d

import (
	"moonbasic/runtime"
)

func (m *Module) Register(r runtime.Registrar) {
	r.Register("PHYSICS2D.START", "physics2d", runtime.AdaptLegacy(m.phStart))
	r.Register("PHYSICS2D.STOP", "physics2d", runtime.AdaptLegacy(m.phStop))
	r.Register("PHYSICS2D.SETGRAVITY", "physics2d", runtime.AdaptLegacy(m.phSetGravity))
	r.Register("PHYSICS2D.SETSTEP", "physics2d", runtime.AdaptLegacy(m.phSetStep))
	r.Register("PHYSICS2D.SETITERATIONS", "physics2d", runtime.AdaptLegacy(m.phSetIterations))
	r.Register("PHYSICS2D.STEP", "physics2d", runtime.AdaptLegacy(m.phStep))

	r.Register("BODY2D.CREATE", "physics2d", runtime.AdaptLegacy(m.bdMake))
	r.Register("BODY2D.MAKE", "physics2d", runtime.AdaptLegacy(m.bdMake))
	r.Register("BODY2D.ADDRECT", "physics2d", runtime.AdaptLegacy(m.bdAddRect))
	r.Register("BODY2D.ADDCIRCLE", "physics2d", runtime.AdaptLegacy(m.bdAddCircle))
	r.Register("BODY2D.COMMIT", "physics2d", runtime.AdaptLegacy(m.bdCommit))
	r.Register("BODY2D.X", "physics2d", runtime.AdaptLegacy(m.bdX))
	r.Register("BODY2D.Y", "physics2d", runtime.AdaptLegacy(m.bdY))
	r.Register("BODY2D.ROT", "physics2d", runtime.AdaptLegacy(m.bdRot) )
	r.Register("BODY2D.FREE", "physics2d", runtime.AdaptLegacy(m.bdFree))

	r.Register("BODY2D.SETPOS", "physics2d", runtime.AdaptLegacy(m.bdSetPos))
	r.Register("BODY2D.GETPOS", "physics2d", runtime.AdaptLegacy(m.bdGetPos))
	r.Register("BODY2D.SETROT", "physics2d", runtime.AdaptLegacy(m.bdSetRot))
	r.Register("BODY2D.GETROT", "physics2d", runtime.AdaptLegacy(m.bdGetRot))
	r.Register("BODY2D.SETMASS", "physics2d", runtime.AdaptLegacy(m.bdSetMass))
	r.Register("BODY2D.SETFRICTION", "physics2d", runtime.AdaptLegacy(m.bdSetFriction))
	r.Register("BODY2D.SETRESTITUTION", "physics2d", runtime.AdaptLegacy(m.bdSetRestitution))
	r.Register("BODY2D.APPLYFORCE", "physics2d", runtime.AdaptLegacy(m.bdApplyForce))
	r.Register("BODY2D.APPLYIMPULSE", "physics2d", runtime.AdaptLegacy(m.bdApplyImpulse))
	r.Register("BODY2D.ADDPOLYGON", "physics2d", runtime.AdaptLegacy(m.bdAddPolygon))
	r.Register("BODY2D.SETLINEARVELOCITY", "physics2d", runtime.AdaptLegacy(m.bdSetLinearVel))
	r.Register("BODY2D.SETANGULARVELOCITY", "physics2d", runtime.AdaptLegacy(m.bdSetAngularVel))
	r.Register("BODY2D.COLLIDED", "physics2d", runtime.AdaptLegacy(m.bdCollided))
	r.Register("BODY2D.COLLISIONOTHER", "physics2d", runtime.AdaptLegacy(m.bdCollisionOther))
	r.Register("BODY2D.COLLISIONNORMAL", "physics2d", runtime.AdaptLegacy(m.bdCollisionNormal))
	r.Register("BODY2D.COLLISIONPOINT", "physics2d", runtime.AdaptLegacy(m.bdCollisionPoint))
	r.Register("PHYSICS2D.DEBUGDRAW", "physics2d", runtime.AdaptLegacy(m.phDebugDraw))
	r.Register("PHYSICS2D.GETDEBUGSEGMENTS", "physics2d", runtime.AdaptLegacy(m.phGetDebugSegments))
	r.Register("JOINT2D.DISTANCE", "physics2d", runtime.AdaptLegacy(m.jtDistance))
	r.Register("JOINT2D.REVOLUTE", "physics2d", runtime.AdaptLegacy(m.jtRevolute))
	r.Register("JOINT2D.PRISMATIC", "physics2d", runtime.AdaptLegacy(m.jtPrismatic))
	r.Register("JOINT2D.FREE", "physics2d", runtime.AdaptLegacy(m.jtFree))

	// BOX2D aliases (legacy compatible names)
	r.Register("BOX2D.WORLDCREATE", "physics2d", runtime.AdaptLegacy(m.phStart))
	r.Register("BOX2D.BODYCREATE", "physics2d", runtime.AdaptLegacy(m.bdMake))
	r.Register("BOX2D.FIXTUREBOX", "physics2d", runtime.AdaptLegacy(m.bdAddRect))
	r.Register("BOX2D.FIXTURECIRCLE", "physics2d", runtime.AdaptLegacy(m.bdAddCircle))
}
