//go:build !cgo && !windows

package mbnav

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

const hint = "NAV/PATH/NAVAGENT/STEER/BTREE natives require CGO: set CGO_ENABLED=1 and install a C compiler, then rebuild"

func navStub(name string) func([]value.Value) (value.Value, error) {
	return func(args []value.Value) (value.Value, error) {
		return value.Nil, fmt.Errorf("%s: %s", name, hint)
	}
}

func (m *Module) navMake(args []value.Value) (value.Value, error) {
	return navStub("NAV.MAKE")(args)
}
func (m *Module) navFree(args []value.Value) (value.Value, error) {
	return navStub("NAV.FREE")(args)
}
func (m *Module) navSetGrid(args []value.Value) (value.Value, error) {
	return navStub("NAV.SETGRID")(args)
}
func (m *Module) navAddTerrain(args []value.Value) (value.Value, error) {
	return navStub("NAV.ADDTERRAIN")(args)
}
func (m *Module) navAddObstacle(args []value.Value) (value.Value, error) {
	return navStub("NAV.ADDOBSTACLE")(args)
}
func (m *Module) navBuild(args []value.Value) (value.Value, error) {
	return navStub("NAV.BUILD")(args)
}
func (m *Module) navDebugDraw(args []value.Value) (value.Value, error) {
	return navStub("NAV.DEBUGDRAW")(args)
}
func (m *Module) navFindPath(args []value.Value) (value.Value, error) {
	return navStub("NAV.FINDPATH")(args)
}
func (m *Module) navBakeTerrain(args []value.Value) (value.Value, error) {
	return navStub("NAV.BAKE")(args)
}
func (m *Module) navGetPathTerrain(args []value.Value) (value.Value, error) {
	return navStub("NAV.GETPATH")(args)
}
func (m *Module) navIsReachableTerrain(args []value.Value) (value.Value, error) {
	return navStub("NAV.ISREACHABLE")(args)
}
func (m *Module) enemyFollowPath(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	return navStub("ENEMY.FOLLOWPATH")(args)
}
func (m *Module) pathIsValid(args []value.Value) (value.Value, error) {
	return navStub("PATH.ISVALID")(args)
}
func (m *Module) pathNodeCount(args []value.Value) (value.Value, error) {
	return navStub("PATH.NODECOUNT")(args)
}
func (m *Module) pathNodeX(args []value.Value) (value.Value, error) {
	return navStub("PATH.NODEX")(args)
}
func (m *Module) pathNodeY(args []value.Value) (value.Value, error) {
	return navStub("PATH.NODEY")(args)
}
func (m *Module) pathNodeZ(args []value.Value) (value.Value, error) {
	return navStub("PATH.NODEZ")(args)
}
func (m *Module) pathFree(args []value.Value) (value.Value, error) {
	return navStub("PATH.FREE")(args)
}
func (m *Module) agentMake(args []value.Value) (value.Value, error) {
	return navStub("NAVAGENT.MAKE")(args)
}
func (m *Module) agentFree(args []value.Value) (value.Value, error) {
	return navStub("NAVAGENT.FREE")(args)
}
func (m *Module) agentSetPos(args []value.Value) (value.Value, error) {
	return navStub("NAVAGENT.SETPOS")(args)
}
func (m *Module) agentSetSpeed(args []value.Value) (value.Value, error) {
	return navStub("NAVAGENT.SETSPEED")(args)
}
func (m *Module) agentSetMaxForce(args []value.Value) (value.Value, error) {
	return navStub("NAVAGENT.SETMAXFORCE")(args)
}
func (m *Module) agentApplyForce(args []value.Value) (value.Value, error) {
	return navStub("NAVAGENT.APPLYFORCE")(args)
}
func (m *Module) agentMoveTo(args []value.Value) (value.Value, error) {
	return navStub("NAVAGENT.MOVETO")(args)
}
func (m *Module) agentUpdate(args []value.Value) (value.Value, error) {
	return navStub("NAVAGENT.UPDATE")(args)
}
func (m *Module) agentIsAtDestination(args []value.Value) (value.Value, error) {
	return navStub("NAVAGENT.ISATDESTINATION")(args)
}
func (m *Module) agentX(args []value.Value) (value.Value, error) {
	return navStub("NAVAGENT.X")(args)
}
func (m *Module) agentY(args []value.Value) (value.Value, error) {
	return navStub("NAVAGENT.Y")(args)
}
func (m *Module) agentZ(args []value.Value) (value.Value, error) {
	return navStub("NAVAGENT.Z")(args)
}
func (m *Module) agentGetPos(args []value.Value) (value.Value, error) {
	return navStub("NAVAGENT.GETPOS")(args)
}
func (m *Module) agentGetRot(args []value.Value) (value.Value, error) {
	return navStub("NAVAGENT.GETROT")(args)
}
func (m *Module) agentGetSpeed(args []value.Value) (value.Value, error) {
	return navStub("NAVAGENT.GETSPEED")(args)
}
func (m *Module) agentGetMaxForce(args []value.Value) (value.Value, error) {
	return navStub("NAVAGENT.GETMAXFORCE")(args)
}
func (m *Module) steerGroupMake(args []value.Value) (value.Value, error) {
	return navStub("STEER.GROUPMAKE")(args)
}
func (m *Module) steerGroupAdd(args []value.Value) (value.Value, error) {
	return navStub("STEER.GROUPADD")(args)
}
func (m *Module) steerGroupClear(args []value.Value) (value.Value, error) {
	return navStub("STEER.GROUPCLEAR")(args)
}
func (m *Module) steerSeek(args []value.Value) (value.Value, error) {
	return navStub("STEER.SEEK")(args)
}
func (m *Module) steerFlee(args []value.Value) (value.Value, error) {
	return navStub("STEER.FLEE")(args)
}
func (m *Module) steerArrive(args []value.Value) (value.Value, error) {
	return navStub("STEER.ARRIVE")(args)
}
func (m *Module) steerWander(args []value.Value) (value.Value, error) {
	return navStub("STEER.WANDER")(args)
}
func (m *Module) steerFlock(args []value.Value) (value.Value, error) {
	return navStub("STEER.FLOCK")(args)
}
func (m *Module) steerAvoidObstacles(args []value.Value) (value.Value, error) {
	return navStub("STEER.AVOIDOBSTACLES")(args)
}
func (m *Module) steerFollowPath(args []value.Value) (value.Value, error) {
	return navStub("STEER.FOLLOWPATH")(args)
}
func (m *Module) btMake(args []value.Value) (value.Value, error) {
	return navStub("BTREE.MAKE")(args)
}
func (m *Module) btFree(args []value.Value) (value.Value, error) {
	return navStub("BTREE.FREE")(args)
}
func (m *Module) btSequence(args []value.Value) (value.Value, error) {
	return navStub("BTREE.SEQUENCE")(args)
}
func (m *Module) btAddCondition(args []value.Value) (value.Value, error) {
	return navStub("BTREE.ADDCONDITION")(args)
}
func (m *Module) btAddAction(args []value.Value) (value.Value, error) {
	return navStub("BTREE.ADDACTION")(args)
}
func (m *Module) btRun(args []value.Value) (value.Value, error) {
	return navStub("BTREE.RUN")(args)
}
