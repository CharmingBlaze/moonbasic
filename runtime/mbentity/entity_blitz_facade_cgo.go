//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

// registerBlitzFacadeCommands wires documented Blitz-style global names to ENTITY.* / physics without duplicating logic.
func registerBlitzFacadeCommands(m *Module, r runtime.Registrar) {
	r.Register("ENTITY.CREATECONE", "entity", runtime.AdaptLegacy(m.entCreateCone))
	r.Register("CreateCone", "entity", runtime.AdaptLegacy(m.entCreateConeEasy))
	r.Register("GetParent", "entity", runtime.AdaptLegacy(m.entGetParent))
	r.Register("EntityParent", "entity", runtime.AdaptLegacy(m.entParent))
	r.Register("FindChild", "entity", m.entFindChild)
	r.Register("GetChild", "entity", runtime.AdaptLegacy(m.entGetChild))
	r.Register("CountChildren", "entity", runtime.AdaptLegacy(m.entCountChildren))
	r.Register("CopyEntity", "entity", runtime.AdaptLegacy(m.entCopyExtended))
	r.Register("FreeEntity", "entity", runtime.AdaptLegacy(m.entFree))
	r.Register("EntityPick", "entity", runtime.AdaptLegacy(m.entPick))
	r.Register("Collisions", "entity", runtime.AdaptLegacy(m.entCollisions))
	r.Register("CollisionEntity", "entity", runtime.AdaptLegacy(m.entGetCollisionEntity))
	r.Register("EntityColor", "entity", runtime.AdaptLegacy(m.entColor))
	r.Register("EntityAlpha", "entity", runtime.AdaptLegacy(m.entAlpha))
	r.Register("EntityShininess", "entity", runtime.AdaptLegacy(m.entShininess))
	r.Register("EntityTexture", "entity", runtime.AdaptLegacy(m.entTexture))
	r.Register("EntityName", "entity", m.entEntityNameStr)
	r.Register("NameEntity", "entity", m.entSetName)
	r.Register("PositionEntity", "entity", runtime.AdaptLegacy(m.entSetPosition))
	r.Register("RotateEntity", "entity", runtime.AdaptLegacy(m.entRotateEntityAbs))
	r.Register("TurnEntity", "entity", runtime.AdaptLegacy(m.entRotate))
	r.Register("ScaleEntity", "entity", runtime.AdaptLegacy(m.entScale))
	r.Register("EntityX", "entity", runtime.AdaptLegacy(m.entEntityX))
	r.Register("EntityY", "entity", runtime.AdaptLegacy(m.entEntityY))
	r.Register("EntityZ", "entity", runtime.AdaptLegacy(m.entEntityZ))
	r.Register("EntityPitch", "entity", runtime.AdaptLegacy(m.entEntityPitch))
	r.Register("EntityYaw", "entity", runtime.AdaptLegacy(m.entEntityYaw))
	r.Register("EntityRoll", "entity", runtime.AdaptLegacy(m.entEntityRoll))
	r.Register("EntityVisible", "entity", runtime.AdaptLegacy(m.entVisible))
	r.Register("EntityDistance", "entity", runtime.AdaptLegacy(m.entDistance))
	r.Register("Entity.GetDistance", "entity", runtime.AdaptLegacy(m.entDistance))
	r.Register("Entity.IsType", "entity", m.entIsType)
	r.Register("Entity.SendMessage", "entity", m.entSendMessage)
	r.Register("Entity.PollMessage", "entity", m.entPollMessage)
	r.Register("Entity.FindByProperty", "entity", m.entFindByProperty)
	r.Register("EntityInView", "entity", runtime.AdaptLegacy(m.entInView))
	r.Register("EntityType", "entity", runtime.AdaptLegacy(m.entType))
	r.Register("EntityRadius", "entity", runtime.AdaptLegacy(m.entRadius))
	r.Register("EntityBox", "entity", runtime.AdaptLegacy(m.entBox))
	r.Register("ResetEntity", "entity", runtime.AdaptLegacy(m.entReset))
	r.Register("ApplyEntityForce", "entity", runtime.AdaptLegacy(m.entAddForce))
	r.Register("ApplyEntityTorque", "entity", runtime.AdaptLegacy(m.entApplyTorqueStub))
	r.Register("CollisionX", "entity", runtime.AdaptLegacy(m.entCollisionX))
	r.Register("CollisionY", "entity", runtime.AdaptLegacy(m.entCollisionY))
	r.Register("CollisionZ", "entity", runtime.AdaptLegacy(m.entCollisionZ))
	r.Register("CollisionNX", "entity", runtime.AdaptLegacy(m.entCollisionNX))
	r.Register("CollisionNY", "entity", runtime.AdaptLegacy(m.entCollisionNY))
	r.Register("CollisionNZ", "entity", runtime.AdaptLegacy(m.entCollisionNZ))
	r.Register("Animate", "entity", runtime.AdaptLegacy(m.entAnimate))
	r.Register("SetAnimTime", "entity", runtime.AdaptLegacy(m.entSetAnimTime))
	r.Register("EntityAnimTime", "entity", runtime.AdaptLegacy(m.entAnimTime))
	r.Register("ExtractAnimSeq", "entity", runtime.AdaptLegacy(m.entExtractAnimSeq))
	r.Register("AnimLength", "entity", runtime.AdaptLegacy(m.entAnimLength))
	r.Register("CreateTerrain", "entity", runtime.AdaptLegacy(m.entCreateTerrainStub))
	r.Register("CreateMirror", "entity", runtime.AdaptLegacy(m.entCreateMirrorStub))
	r.Register("EntityAutoFade", "entity", runtime.AdaptLegacy(m.entAutoFadeStub))
	r.Register("UpdateNormals", "entity", runtime.AdaptLegacy(m.entUpdateNormalsStub))
	r.Register("FlipMesh", "entity", runtime.AdaptLegacy(m.entFlipMeshStub))
	r.Register("FitMesh", "entity", runtime.AdaptLegacy(m.entFitMeshStub))
	r.Register("VertexNX", "entity", runtime.AdaptLegacy(m.entVertexNormalStub))
	r.Register("VertexNY", "entity", runtime.AdaptLegacy(m.entVertexNormalStub))
	r.Register("VertexNZ", "entity", runtime.AdaptLegacy(m.entVertexNormalStub))
	r.Register("VertexU", "entity", runtime.AdaptLegacy(m.entVertexUVNotImpl))
	r.Register("VertexV", "entity", runtime.AdaptLegacy(m.entVertexUVNotImpl))
	r.Register("CountVertices", "entity", runtime.AdaptLegacy(m.entCountVerticesStub))
	r.Register("CountTriangles", "entity", runtime.AdaptLegacy(m.entCountTrianglesStub))
}

func (m *Module) entCreateCone(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("ENTITY.CREATECONE expects (radius#, height#, segments#)")
	}
	rad, ok1 := argF32(args[0])
	h, ok2 := argF32(args[1])
	seg, ok3 := args[2].ToInt()
	if !ok1 || !ok2 || !ok3 || seg < 3 {
		return value.Nil, fmt.Errorf("ENTITY.CREATECONE: radius/height numeric, segments int >= 3")
	}
	st := m.store()
	id := st.nextID
	st.nextID++
	e := newDefaultEnt(id)
	e.kind = entKindCone
	e.radius = rad
	e.cylH = h
	e.segV = int32(seg)
	e.w, e.h, e.d = rad*2, h, rad*2
	e.static = true
	st.ents[id] = e
	return value.FromInt(id), nil
}

// CreateCone: () default; (parent#); (r,h,seg); (parent#, r,h,seg)
func (m *Module) entCreateConeEasy(args []value.Value) (value.Value, error) {
	switch len(args) {
	case 0:
		return m.entCreateCone([]value.Value{value.FromFloat(0.5), value.FromFloat(1), value.FromInt(16)})
	case 1:
		pid, ok := m.entID(args[0])
		if !ok || pid < 1 || m.store().ents[pid] == nil {
			return value.Nil, fmt.Errorf("CreateCone: invalid parent")
		}
		v, err := m.entCreateCone([]value.Value{value.FromFloat(0.5), value.FromFloat(1), value.FromInt(16)})
		if err != nil {
			return v, err
		}
		cid, _ := v.ToInt()
		_, err = m.entParent([]value.Value{value.FromInt(cid), value.FromInt(pid)})
		return v, err
	case 3:
		return m.entCreateCone(args)
	case 4:
		pid, ok := m.entID(args[0])
		if !ok || pid < 1 || m.store().ents[pid] == nil {
			return value.Nil, fmt.Errorf("CreateCone: invalid parent")
		}
		v, err := m.entCreateCone([]value.Value{args[1], args[2], args[3]})
		if err != nil {
			return v, err
		}
		cid, _ := v.ToInt()
		_, err = m.entParent([]value.Value{value.FromInt(cid), value.FromInt(pid)})
		return v, err
	default:
		return value.Nil, fmt.Errorf("CreateCone expects 0, 1 (parent#), 3 (r,h,seg), or 4 (parent,r,h,seg)")
	}
}

func (m *Module) entGetParent(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("GetParent expects entity#")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("GetParent: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("GetParent: unknown entity")
	}
	return value.FromInt(e.parentID), nil
}

func (m *Module) entEntityNameStr(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("EntityName expects entity#")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("EntityName: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("EntityName: unknown entity")
	}
	return rt.RetString(e.name), nil
}

func (m *Module) entCopyExtended(args []value.Value) (value.Value, error) {
	if len(args) != 1 && len(args) != 2 {
		return value.Nil, fmt.Errorf("CopyEntity expects entity# [, parent#]")
	}
	v, err := m.entCopy([]value.Value{args[0]})
	if err != nil {
		return v, err
	}
	if len(args) == 2 {
		pid, ok := m.entID(args[1])
		if !ok || pid < 1 {
			return value.Nil, fmt.Errorf("CopyEntity: invalid parent")
		}
		if m.store().ents[pid] == nil {
			return value.Nil, fmt.Errorf("CopyEntity: unknown parent")
		}
		_, err := m.entParent([]value.Value{v, value.FromInt(pid)})
		if err != nil {
			return v, err
		}
	}
	return v, nil
}

func (m *Module) entApplyTorqueStub(args []value.Value) (value.Value, error) {
	_ = args
	return value.Nil, fmt.Errorf("ApplyEntityTorque: not implemented (no entity torque integrator in this runtime)")
}

func (m *Module) entCreateTerrainStub(args []value.Value) (value.Value, error) {
	_ = args
	return value.Nil, fmt.Errorf("CreateTerrain: use TERRAIN.MAKE(worldW, worldH [, cellSize]) (returns terrain handle, not entity#); entity parenting for terrain is not wired")
}

func (m *Module) entCreateMirrorStub(args []value.Value) (value.Value, error) {
	_ = args
	return value.Nil, fmt.Errorf("CreateMirror: planar reflections are not implemented (deferred)")
}

func (m *Module) entAutoFadeStub(args []value.Value) (value.Value, error) {
	_ = args
	return value.Nil, fmt.Errorf("EntityAutoFade: distance-based alpha is not implemented")
}

func (m *Module) entUpdateNormalsStub(args []value.Value) (value.Value, error) {
	_ = args
	return value.Nil, fmt.Errorf("UpdateNormals: not implemented for procedural meshes in this runtime")
}

func (m *Module) entFlipMeshStub(args []value.Value) (value.Value, error) {
	_ = args
	return value.Nil, fmt.Errorf("FlipMesh: not implemented for procedural meshes in this runtime")
}

func (m *Module) entFitMeshStub(args []value.Value) (value.Value, error) {
	_ = args
	return value.Nil, fmt.Errorf("FitMesh: not implemented for procedural meshes in this runtime")
}

func (m *Module) entVertexNormalStub(args []value.Value) (value.Value, error) {
	_ = args
	return value.FromFloat(0), nil
}

func (m *Module) entVertexUVNotImpl(args []value.Value) (value.Value, error) {
	_ = args
	return value.Nil, fmt.Errorf("VertexU/VertexV: UV reads are not implemented for procedural surfaces in this runtime")
}

func (m *Module) entCountVerticesStub(args []value.Value) (value.Value, error) {
	_ = args
	return value.FromInt(0), nil
}

func (m *Module) entCountTrianglesStub(args []value.Value) (value.Value, error) {
	_ = args
	return value.FromInt(0), nil
}

// entCreatePivot: optional parent# (0 = world root, i.e. none)
func (m *Module) entCreatePivot(args []value.Value) (value.Value, error) {
	if len(args) > 1 {
		return value.Nil, fmt.Errorf("CreatePivot expects 0 or 1 (parent#) arguments")
	}
	st := m.store()
	id := st.nextID
	st.nextID++
	e := newDefaultEnt(id)
	e.kind = entKindEmpty
	e.hidden = true
	st.ents[id] = e
	v := value.FromInt(id)
	if len(args) == 1 {
		pid, ok := m.entID(args[0])
		if !ok || pid < 1 {
			return value.Nil, fmt.Errorf("CreatePivot: invalid parent entity")
		}
		if st.ents[pid] == nil {
			return value.Nil, fmt.Errorf("CreatePivot: unknown parent")
		}
		_, err := m.entParent([]value.Value{v, value.FromInt(pid)})
		if err != nil {
			return value.Nil, err
		}
	}
	return v, nil
}
