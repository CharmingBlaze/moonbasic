//go:build cgo || (windows && !cgo)

package mbgrid

import "moonbasic/runtime"

func (m *Module) Register(reg runtime.Registrar) {
	reg.Register("GRID.CREATE", "grid", m.gridCreate)
	reg.Register("GRID.MAKE", "grid", m.gridCreate)
	reg.Register("GRID.FREE", "grid", m.gridFree)
	reg.Register("GRID.SETCELL", "grid", m.gridSetCell)
	reg.Register("GRID.GETCELL", "grid", m.gridGetCell)
	reg.Register("GRID.WORLDTOCELL", "grid", m.gridWorldToCell)
	reg.Register("GRID.DRAW", "grid", m.gridDraw)
	reg.Register("GRID.SNAP", "grid", m.gridSnap)
	reg.Register("GRID.GETPATH", "grid", m.gridGetPath)
	reg.Register("GRID.FOLLOWTERRAIN", "grid", m.gridFollowTerrain)
	reg.Register("GRID.PLACEENTITY", "grid", m.gridPlaceEntity)
	reg.Register("GRID.RAYCAST", "grid", m.gridRaycast)
	reg.Register("GRID.GETNEIGHBORS", "grid", m.gridGetNeighbors)
}
