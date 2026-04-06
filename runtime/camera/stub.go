//go:build !cgo && !windows

package mbcamera

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

const hint = "CAMERA.* requires CGO: set CGO_ENABLED=1 and install a C compiler, then rebuild"

// Register implements runtime.Module.
func (m *Module) Register(r runtime.Registrar) {
	stub := func(name string) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			return value.Nil, fmt.Errorf("%s: %s", name, hint)
		}
	}
	r.Register("CAMERA.MAKE", "camera", stub("CAMERA.MAKE"))
	r.Register("CAMERA.SETPOS", "camera", stub("CAMERA.SETPOS"))
	r.Register("CAMERA.SETPOSITION", "camera", stub("CAMERA.SETPOSITION"))
	r.Register("CAMERA.SETTARGET", "camera", stub("CAMERA.SETTARGET"))
	r.Register("CAMERA.LOOKAT", "camera", stub("CAMERA.LOOKAT"))
	r.Register("CAMERA.SETPROJECTION", "camera", stub("CAMERA.SETPROJECTION"))
	r.Register("CAMERA.SETFOV", "camera", stub("CAMERA.SETFOV"))
	r.Register("CAMERA.BEGIN", "camera", stub("CAMERA.BEGIN"))
	r.Register("CAMERA.END", "camera", stub("CAMERA.END"))
	r.Register("CAMERA.MOVE", "camera", stub("CAMERA.MOVE"))
	r.Register("CAMERA.GETRAY", "camera", stub("CAMERA.GETRAY"))
	r.Register("CAMERA.PICK", "camera", stub("CAMERA.PICK"))
	r.Register("CAMERA.SHAKE", "camera", stub("CAMERA.SHAKE"))
	r.Register("CAMERA.GETVIEWRAY", "camera", stub("CAMERA.GETVIEWRAY"))
	r.Register("CAMERA.GETMATRIX", "camera", stub("CAMERA.GETMATRIX"))
	r.Register("CAMERA.GETPOS", "camera", stub("CAMERA.GETPOS"))
	r.Register("CAMERA.GETTARGET", "camera", stub("CAMERA.GETTARGET"))
	r.Register("CAMERA.SETUP", "camera", stub("CAMERA.SETUP"))
	r.Register("CAMERA.SETORBIT", "camera", stub("CAMERA.SETORBIT"))
	r.Register("CAMERA.ORBITAROUND", "camera", stub("CAMERA.ORBITAROUND"))
	r.Register("CAMERA.ORBITAROUNDEG", "camera", stub("CAMERA.ORBITAROUNDEG"))
	r.Register("CAMERA.TURN", "camera", stub("CAMERA.TURN"))
	r.Register("CAMERA.ROTATE", "camera", stub("CAMERA.ROTATE"))
	r.Register("CAMERA.ORBIT", "camera", stub("CAMERA.ORBIT"))
	r.Register("CAMERA.ZOOM", "camera", stub("CAMERA.ZOOM"))
	r.Register("CAMERA.FOLLOW", "camera", stub("CAMERA.FOLLOW"))
	r.Register("CAMERA.FREE", "camera", stub("CAMERA.FREE"))
	r.Register("CAMERA.WORLDTOSCREEN", "camera", stub("CAMERA.WORLDTOSCREEN"))
	r.Register("CAMERA.ISONSCREEN", "camera", stub("CAMERA.ISONSCREEN"))
	r.Register("CAMERA.MOUSERAY", "camera", stub("CAMERA.MOUSERAY"))
	r.Register("MATRIX.FREE", "camera", stub("MATRIX.FREE"))
	r.Register("CAMERA2D.MAKE", "camera", stub("CAMERA2D.MAKE"))
	r.Register("CAMERA2D.SETTARGET", "camera", stub("CAMERA2D.SETTARGET"))
	r.Register("CAMERA2D.SETOFFSET", "camera", stub("CAMERA2D.SETOFFSET"))
	r.Register("CAMERA2D.SETZOOM", "camera", stub("CAMERA2D.SETZOOM"))
	r.Register("CAMERA2D.SETROTATION", "camera", stub("CAMERA2D.SETROTATION"))
	r.Register("CAMERA2D.BEGIN", "camera", stub("CAMERA2D.BEGIN"))
	r.Register("CAMERA2D.END", "camera", stub("CAMERA2D.END"))
	r.Register("CAMERA2D.GETMATRIX", "camera", stub("CAMERA2D.GETMATRIX"))
	r.Register("CAMERA2D.WORLDTOSCREEN", "camera", stub("CAMERA2D.WORLDTOSCREEN"))
	r.Register("CAMERA2D.SCREENTOWORLD", "camera", stub("CAMERA2D.SCREENTOWORLD"))
	r.Register("CAMERA2D.FREE", "camera", stub("CAMERA2D.FREE"))
	r.Register("CAMERA.SETRANGE", "camera", stub("CAMERA.SETRANGE"))
	r.Register("CAMERA.SETACTIVE", "camera", stub("CAMERA.SETACTIVE"))
	r.Register("CAMERA.GETACTIVE", "camera", stub("CAMERA.GETACTIVE"))
	r.Register("CAMERA.WORLDTOSCREEN2D", "camera", stub("CAMERA.WORLDTOSCREEN2D"))
	r.Register("CAMERA.SETFPSMODE", "camera", stub("CAMERA.SETFPSMODE"))
	r.Register("CAMERA.CLEARFPSMODE", "camera", stub("CAMERA.CLEARFPSMODE"))
	r.Register("CAMERA.UPDATEFPS", "camera", stub("CAMERA.UPDATEFPS"))
	r.Register("CAMERA2D.FOLLOW", "camera", stub("CAMERA2D.FOLLOW"))
	r.Register("CAMERA2D.ZOOMTOMOUSE", "camera", stub("CAMERA2D.ZOOMTOMOUSE"))
	r.Register("CAMERA2D.ZOOMIN", "camera", stub("CAMERA2D.ZOOMIN"))
	r.Register("CAMERA2D.ZOOMOUT", "camera", stub("CAMERA2D.ZOOMOUT"))
	r.Register("CAMERA2D.ROTATION", "camera", stub("CAMERA2D.ROTATION"))
	r.Register("CAMERA2D.TARGETX", "camera", stub("CAMERA2D.TARGETX"))
	r.Register("CAMERA2D.TARGETY", "camera", stub("CAMERA2D.TARGETY"))

	r.Register("CULL.SPHEREVISIBLE", "camera", stub("CULL.SPHEREVISIBLE"))
	r.Register("CULL.AABBVISIBLE", "camera", stub("CULL.AABBVISIBLE"))
	r.Register("CULL.POINTVISIBLE", "camera", stub("CULL.POINTVISIBLE"))
	r.Register("CULL.INRANGE", "camera", stub("CULL.INRANGE"))
	r.Register("CULL.DISTANCE", "camera", stub("CULL.DISTANCE"))
	r.Register("CULL.DISTANCESQ", "camera", stub("CULL.DISTANCESQ"))
	r.Register("CULL.BEHINDHORIZON", "camera", stub("CULL.BEHINDHORIZON"))
	r.Register("CULL.BATCHSPHERE", "camera", stub("CULL.BATCHSPHERE"))
	r.Register("CULL.OCCLUSIONENABLE", "camera", stub("CULL.OCCLUSIONENABLE"))
	r.Register("CULL.OCCLUDERADD", "camera", stub("CULL.OCCLUDERADD"))
	r.Register("CULL.OCCLUDERCLEAR", "camera", stub("CULL.OCCLUDERCLEAR"))
	r.Register("CULL.ISOCCLUDED", "camera", stub("CULL.ISOCCLUDED"))
	r.Register("CULL.SETMAXDISTANCE", "camera", stub("CULL.SETMAXDISTANCE"))
	r.Register("CULL.GETMAXDISTANCE", "camera", stub("CULL.GETMAXDISTANCE"))
	r.Register("CULL.STATSRESET", "camera", stub("CULL.STATSRESET"))
	r.Register("CULL.STATSTOTAL", "camera", stub("CULL.STATSTOTAL"))
	r.Register("CULL.STATSCULLED", "camera", stub("CULL.STATSCULLED"))
	r.Register("CULL.STATSVISIBLE", "camera", stub("CULL.STATSVISIBLE"))
	r.Register("CULL.STATSFRUSTUMCULLED", "camera", stub("CULL.STATSFRUSTUMCULLED"))
	r.Register("CULL.STATSDISTANCECULLED", "camera", stub("CULL.STATSDISTANCECULLED"))
	r.Register("CULL.STATSHORIZONCULLED", "camera", stub("CULL.STATSHORIZONCULLED"))
	r.Register("CULL.STATSOCCLUSIONCULLED", "camera", stub("CULL.STATSOCCLUSIONCULLED"))
	r.Register("CULL.SETBACKFACECULLING", "camera", stub("CULL.SETBACKFACECULLING"))
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
