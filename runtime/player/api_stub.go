//go:build (!linux && !windows) || !cgo

package player

import (
	"fmt"

	"moonbasic/vm/value"
)

func stubErr(name string) func(args []value.Value) (value.Value, error) {
	return func(args []value.Value) (value.Value, error) {
		_ = args
		return value.Nil, fmt.Errorf("%s [%s]", errPlayerRequiresCGOJolt, name)
	}
}
func (m *Module) playerCharacterCreate(args []value.Value) (value.Value, error) {
	return stubErr("CHARACTER.CREATE")(args)
}
func (m *Module) charRefAddVel(args []value.Value) (value.Value, error) {
	return stubErr("CHARACTERREF.ADDVELOCITY")(args)
}
func (m *Module) charRefSetVel(args []value.Value) (value.Value, error) {
	return stubErr("CHARACTERREF.SETVELOCITY")(args)
}
func (m *Module) charRefSetSnapDist(args []value.Value) (value.Value, error) {
	return stubErr("CHARACTERREF.SETSNAPDISTANCE")(args)
}
func (m *Module) charRefUpdate(args []value.Value) (value.Value, error) {
	return stubErr("CHARACTERREF.UPDATE")(args)
}
func (m *Module) charRefJump(args []value.Value) (value.Value, error) {
	return stubErr("CHARACTERREF.JUMP")(args)
}
func (m *Module) charRefMoveWithCam(args []value.Value) (value.Value, error) {
	return stubErr("CHARACTERREF.MOVEWITHCAMERA")(args)
}
func (m *Module) charRefSetMaxSlope(args []value.Value) (value.Value, error) {
	return stubErr("CHARACTERREF.SETMAXSLOPE")(args)
}
func (m *Module) charRefSetStepHeight(args []value.Value) (value.Value, error) {
	return stubErr("CHARACTERREF.SETSTEPHEIGHT")(args)
}
func (m *Module) charRefIsGrounded(args []value.Value) (value.Value, error) {
	return stubErr("CHARACTERREF.ISGROUNDED")(args)
}
func (m *Module) charRefSetPos(args []value.Value) (value.Value, error) {
	return stubErr("CHARACTERREF.SETPOS")(args)
}
func (m *Module) charRefGetPos(args []value.Value) (value.Value, error) {
	return stubErr("CHARACTERREF.GETPOSITION")(args)
}
func (m *Module) charRefFree(args []value.Value) (value.Value, error) {
	return stubErr("CHARACTERREF.FREE")(args)
}
func (m *Module) charRefGetGroundState(args []value.Value) (value.Value, error) {
	return stubErr("CHARACTERREF.GETGROUNDSTATE")(args)
}
func (m *Module) charRefGetSpeed(args []value.Value) (value.Value, error) {
	return stubErr("CHARACTERREF.GETSPEED")(args)
}
func (m *Module) charRefIsMoving(args []value.Value) (value.Value, error) {
	return stubErr("CHARACTERREF.ISMOVING")(args)
}
func (m *Module) charRefSetFriction(args []value.Value) (value.Value, error) {
	return stubErr("CHARACTERREF.SETFRICTION")(args)
}
func (m *Module) charRefSetPadding(args []value.Value) (value.Value, error) {
	return stubErr("CHARACTERREF.SETPADDING")(args)
}
func (m *Module) charRefSetBounce(args []value.Value) (value.Value, error) {
	return stubErr("CHARACTERREF.SETBOUNCE")(args)
}
func (m *Module) charRefSetGravityScale(args []value.Value) (value.Value, error) {
	return stubErr("CHARACTERREF.SETGRAVITY")(args)
}

func (m *Module) playerCreate(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.CREATE")(args)
}
func (m *Module) playerMove(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.MOVE")(args)
}
func (m *Module) playerJump(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.JUMP")(args)
}
func (m *Module) playerIsGrounded(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.ISGROUNDED")(args)
}
func (m *Module) playerGetGroundState(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.GETGROUNDSTATE")(args)
}
func (m *Module) playerIsOnSteepSlope(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.ISONSTEEPSLOPE")(args)
}
func (m *Module) playerSetStepOffset(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.SETSTEPOFFSET")(args)
}
func (m *Module) playerSetSlopeLimit(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.SETSLOPELIMIT")(args)
}
func (m *Module) playerSetStickFloor(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.SETSTICKFLOOR")(args)
}
func (m *Module) playerNavTo(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.NAVTO")(args)
}
func (m *Module) playerNavUpdate(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.NAVUPDATE")(args)
}
func (m *Module) playerTeleport(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.TELEPORT")(args)
}
func (m *Module) playerSetGravityScale(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.SETGRAVITYSCALE")(args)
}
func (m *Module) playerMoveWithCam(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.MOVEWITHCAMERA")(args)
}
func (m *Module) playerMoveWithCamera(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.MOVEWITHCAMERA")(args)
}
func (m *Module) playerCharMoveDir(args []value.Value) (value.Value, error) {
	return stubErr("CHAR.MOVE")(args)
}
func (m *Module) playerGetLookTarget(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.GETLOOKTARGET")(args)
}
func (m *Module) playerGetNearby(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.GETNEARBY")(args)
}
func (m *Module) playerOnTrigger(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.ONTRIGGER")(args)
}
func (m *Module) playerSetState(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.SETSTATE")(args)
}
func (m *Module) playerSyncAnim(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.SYNCANIM")(args)
}
func (m *Module) playerSetStepHeight(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.SETSTEPHEIGHT")(args)
}
func (m *Module) playerGetVelocity(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.GETVELOCITY")(args)
}
func (m *Module) playerGetCrouch(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.GETCROUCH")(args)
}
func (m *Module) playerSetCrouch(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.SETCROUCH")(args)
}
func (m *Module) playerSwim(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.SWIM")(args)
}
func (m *Module) playerSetVelocity(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.SETVELOCITY")(args)
}
func (m *Module) playerAddImpulse(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.ADDIMPULSE")(args)
}
func (m *Module) playerGetSubmergedFraction(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.GETSUBMERGEDFACTOR")(args)
}
func (m *Module) playerIsSubmerged(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.ISSUBMERGED")(args)
}
func (m *Module) playerGetStandNormal(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.GETSTANDNORMAL")(args)
}
func (m *Module) playerPush(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.PUSH")(args)
}
func (m *Module) playerGrab(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.GRAB")(args)
}
func (m *Module) playerSetMass(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.SETMASS")(args)
}
func (m *Module) playerGetSurfaceType(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.GETSURFACETYPE")(args)
}
func (m *Module) playerSetFovKick(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.SETFOVKICK")(args)
}
func (m *Module) playerGetFovKick(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.GETFOVKICK")(args)
}
func (m *Module) playerSetJumpBuffer(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.SETJUMPBUFFER")(args)
}
func (m *Module) playerSetAirControl(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.SETAIRCONTROL")(args)
}
func (m *Module) playerSetGroundControl(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.SETGROUNDCONTROL")(args)
}
func (m *Module) playerIsMoving(args []value.Value) (value.Value, error) {
	return stubErr("PLAYER.ISMOVING")(args)
}
func (m *Module) playerNavChase(args []value.Value) (value.Value, error) {
	return stubErr("NAV.CHASE")(args)
}
func (m *Module) playerNavPatrol(args []value.Value) (value.Value, error) {
	return stubErr("NAV.PATROL")(args)
}
