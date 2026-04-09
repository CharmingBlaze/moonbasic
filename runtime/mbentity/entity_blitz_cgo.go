//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"
	"math"
	"strings"

	"moonbasic/runtime"
	mbcamera "moonbasic/runtime/camera"
	"moonbasic/runtime/mbmatrix"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func registerEntityBlitzAPI(m *Module, r runtime.Registrar) {
	registerModernEntityFX(m, r)
	registerPickBlitz(m, r)
	registerBlitzFacadeCommands(m, r)
	// Creation
	r.Register("ENTITY.CREATESPHERE", "entity", runtime.AdaptLegacy(m.entCreateSphere))
	r.Register("ENTITY.CREATECYLINDER", "entity", runtime.AdaptLegacy(m.entCreateCylinder))
	r.Register("ENTITY.CREATEPLANE", "entity", runtime.AdaptLegacy(m.entCreatePlane))
	r.Register("ENTITY.CREATEMESH", "entity", runtime.AdaptLegacy(m.entCreateMesh))
	r.Register("ENTITY.LOADMESH", "entity", m.entLoadMesh)
	r.Register("LoadMesh", "entity", m.entLoadMesh)
	r.Register("ENTITY.LOADANIMATEDMESH", "entity", m.entLoadAnimatedMesh)

	// Blitz names (aliases)
	r.Register("ENTITY.POSITIONENTITY", "entity", runtime.AdaptLegacy(m.entSetPosition))
	r.Register("ENTITY.ROTATEENTITY", "entity", runtime.AdaptLegacy(m.entRotateEntityAbs))
	r.Register("ENTITY.TURNENTITY", "entity", runtime.AdaptLegacy(m.entRotate))
	r.Register("ENTITY.MOVEENTITY", "entity", runtime.AdaptLegacy(m.entMove))
	r.Register("ENTITY.TRANSLATEENTITY", "entity", runtime.AdaptLegacy(m.entTranslate))
	r.Register("ENTITY.GRAVITY", "entity", runtime.AdaptLegacy(m.entSetGravity))

	r.Register("ENTITY.ENTITYX", "entity", runtime.AdaptLegacy(m.entEntityX))
	r.Register("ENTITY.ENTITYY", "entity", runtime.AdaptLegacy(m.entEntityY))
	r.Register("ENTITY.ENTITYZ", "entity", runtime.AdaptLegacy(m.entEntityZ))
	r.Register("ENTITY.ENTITYPITCH", "entity", runtime.AdaptLegacy(m.entEntityPitch))
	r.Register("ENTITY.ENTITYYAW", "entity", runtime.AdaptLegacy(m.entEntityYaw))
	r.Register("ENTITY.ENTITYROLL", "entity", runtime.AdaptLegacy(m.entEntityRoll))

	// Global shorthands (Easy Mode)
	r.Register("ENTITYX", "entity", runtime.AdaptLegacy(m.entEntityX))
	r.Register("ENTITYY", "entity", runtime.AdaptLegacy(m.entEntityY))
	r.Register("ENTITYZ", "entity", runtime.AdaptLegacy(m.entEntityZ))
	r.Register("ENTITYPITCH", "entity", runtime.AdaptLegacy(m.entEntityPitch))
	r.Register("ENTITYYAW", "entity", runtime.AdaptLegacy(m.entEntityYaw))
	r.Register("ENTITYROLL", "entity", runtime.AdaptLegacy(m.entEntityRoll))

	r.Register("HIDEENTITY", "entity", runtime.AdaptLegacy(m.entHide))
	r.Register("SHOWENTITY", "entity", runtime.AdaptLegacy(m.entShow))
	r.Register("FREEENTITY", "entity", runtime.AdaptLegacy(m.entFree))
	r.Register("ENTITYTEXTURE", "entity", runtime.AdaptLegacy(m.entTexture))

	r.Register("MOVEENTITY", "entity", runtime.AdaptLegacy(m.entMove))
	r.Register("TURNENTITY", "entity", runtime.AdaptLegacy(m.entRotate))
	r.Register("POINTENTITY", "entity", runtime.AdaptLegacy(m.entPointEntity))

	r.Register("ENTITY.PARENT", "entity", runtime.AdaptLegacy(m.entParent))
	r.Register("ENTITY.PARENTCLEAR", "entity", runtime.AdaptLegacy(m.entParentClear))

	r.Register("ENTITY.VISIBLE", "entity", runtime.AdaptLegacy(m.entVisible))
	r.Register("EntityVisible", "entity", runtime.AdaptLegacy(m.entVisible))
	r.Register("ENTITY.COUNTCHILDREN", "entity", runtime.AdaptLegacy(m.entCountChildren))
	r.Register("ENTITY.GETCHILD", "entity", runtime.AdaptLegacy(m.entGetChild))
	r.Register("ENTITY.FINDCHILD", "entity", m.entFindChild)

	r.Register("ENTITY.TFORMPOINT", "entity", runtime.AdaptLegacy(m.entTFormPoint))
	r.Register("ENTITY.TFORMVECTOR", "entity", runtime.AdaptLegacy(m.entTFormVector))
	r.Register("TFormVector", "entity", runtime.AdaptLegacy(m.entTFormVector))
	r.Register("ENTITY.DELTAX", "entity", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return m.entDeltaAxis(a, 0) }))
	r.Register("ENTITY.DELTAY", "entity", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return m.entDeltaAxis(a, 1) }))
	r.Register("ENTITY.DELTAZ", "entity", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return m.entDeltaAxis(a, 2) }))
	r.Register("ENTITY.MATRIXELEMENT", "entity", runtime.AdaptLegacy(m.entMatrixElement))
	r.Register("ENTITY.INVIEW", "entity", runtime.AdaptLegacy(m.entInView))

	r.Register("LOADSPRITE", "entity", m.entLoadSprite)
	r.Register("ENTITY.LOADSPRITE", "entity", m.entLoadSprite)
	r.Register("ENTITY.CREATESPRITE", "entity", m.entCreateSpriteEntity)
	r.Register("SCALESPRITE", "entity", runtime.AdaptLegacy(m.entScaleSprite))
	r.Register("SPRITEMODE", "entity", runtime.AdaptLegacy(m.entSpriteMode))
	r.Register("ENTITY.SPRITEVIEWMODE", "entity", runtime.AdaptLegacy(m.entSpriteMode))
	r.Register("SPRITEVIEWMODE", "entity", runtime.AdaptLegacy(m.entSpriteMode))

	r.Register("ENTITY.ALPHA", "entity", runtime.AdaptLegacy(m.entAlpha))
	r.Register("ENTITY.SHININESS", "entity", runtime.AdaptLegacy(m.entShininess))
	r.Register("ENTITY.TEXTURE", "entity", runtime.AdaptLegacy(m.entTexture))
	r.Register("ENTITY.FX", "entity", runtime.AdaptLegacy(m.entFX))
	r.Register("ENTITY.BLEND", "entity", runtime.AdaptLegacy(m.entBlend))
	r.Register("ENTITY.ORDER", "entity", runtime.AdaptLegacy(m.entOrder))

	r.Register("ENTITYALPHA", "entity", runtime.AdaptLegacy(m.entAlpha))
	r.Register("ENTITYSHININESS", "entity", runtime.AdaptLegacy(m.entShininess))
	r.Register("ENTITYBLEND", "entity", runtime.AdaptLegacy(m.entBlend))

	r.Register("ENTITY.TYPE", "entity", runtime.AdaptLegacy(m.entType))
	r.Register("ENTITY.COLLIDE", "entity", runtime.AdaptLegacy(m.entCollide))
	r.Register("ENTITY.COLLISIONX", "entity", runtime.AdaptLegacy(m.entCollisionX))
	r.Register("ENTITY.COLLISIONY", "entity", runtime.AdaptLegacy(m.entCollisionY))
	r.Register("ENTITY.COLLISIONZ", "entity", runtime.AdaptLegacy(m.entCollisionZ))
	r.Register("ENTITY.COLLISIONNX", "entity", runtime.AdaptLegacy(m.entCollisionNX))
	r.Register("ENTITY.COLLISIONNY", "entity", runtime.AdaptLegacy(m.entCollisionNY))
	r.Register("ENTITY.COLLISIONNZ", "entity", runtime.AdaptLegacy(m.entCollisionNZ))
	r.Register("ENTITY.DISTANCE", "entity", runtime.AdaptLegacy(m.entDistance))

	r.Register("ENTITY.VELOCITY", "entity", runtime.AdaptLegacy(m.entVelocity))
	r.Register("ENTITY.ADDFORCE", "entity", runtime.AdaptLegacy(m.entAddForce))
	r.Register("ENTITY.SLIDE", "entity", runtime.AdaptLegacy(m.entSetSlide))
	r.Register("ENTITY.PICK", "entity", runtime.AdaptLegacy(m.entPick))
	r.Register("ENTITY.PICKMODE", "entity", runtime.AdaptLegacy(m.entPickMode))

	r.Register("ENTITY.POINTENTITY", "entity", runtime.AdaptLegacy(m.entPointEntity))
	r.Register("ENTITY.ALIGNTOVECTOR", "entity", runtime.AdaptLegacy(m.entAlignToVector))

	r.Register("ENTITY.ANIMATE", "entity", runtime.AdaptLegacy(m.entAnimate))
	r.Register("ENTITY.SETANIMTIME", "entity", runtime.AdaptLegacy(m.entSetAnimTime))
	r.Register("ENTITY.ANIMTIME", "entity", runtime.AdaptLegacy(m.entAnimTime))
	r.Register("ENTITY.ANIMLENGTH", "entity", runtime.AdaptLegacy(m.entAnimLength))
	r.Register("ENTITY.EXTRACTANIMSEQ", "entity", runtime.AdaptLegacy(m.entExtractAnimSeq))
	r.Register("ENTITY.SETANIMINDEX", "entity", runtime.AdaptLegacy(m.entSetAnimIndex))
	r.Register("ENTITY.FINDBONE", "entity", m.entFindBone)

	r.Register("LoadAnimMesh", "entity", m.entLoadAnimatedMesh)
	r.Register("Animate", "entity", runtime.AdaptLegacy(m.entAnimate))
	r.Register("SetAnimTime", "entity", runtime.AdaptLegacy(m.entSetAnimTime))
	r.Register("EntityAnimTime", "entity", runtime.AdaptLegacy(m.entAnimTime))
	r.Register("FindBone", "entity", m.entFindBone)
	r.Register("ExtractAnimSeq", "entity", runtime.AdaptLegacy(m.entExtractAnimSeq))

	r.Register("CreateBrush", "entity", runtime.AdaptLegacy(m.entCreateBrush))
	r.Register("BrushTexture", "entity", runtime.AdaptLegacy(m.entBrushTexture))
	r.Register("BrushFX", "entity", runtime.AdaptLegacy(m.entBrushFX))
	r.Register("BrushShininess", "entity", runtime.AdaptLegacy(m.entBrushShininess))
	r.Register("PaintEntity", "entity", runtime.AdaptLegacy(m.entPaintEntity))
	r.Register("EntityShadow", "entity", runtime.AdaptLegacy(m.entEntityShadow))

	r.Register("EmitSound", "entity", runtime.AdaptLegacy(m.entEmitSound))
	r.Register("CreateSurface", "entity", runtime.AdaptLegacy(m.entCreateSurface))
	r.Register("AddVertex", "entity", runtime.AdaptLegacy(m.entAddVertex))
	r.Register("AddTriangle", "entity", runtime.AdaptLegacy(m.entAddTriangle))
	r.Register("UpdateMesh", "entity", runtime.AdaptLegacy(m.entUpdateMesh))
	r.Register("VertexX", "entity", runtime.AdaptLegacy(m.entVertexX))
	r.Register("VertexY", "entity", runtime.AdaptLegacy(m.entVertexY))
	r.Register("VertexZ", "entity", runtime.AdaptLegacy(m.entVertexZ))
	r.Register("ENTITY.CREATESURFACE", "entity", runtime.AdaptLegacy(m.entCreateSurface))
	r.Register("ENTITY.ADDVERTEX", "entity", runtime.AdaptLegacy(m.entAddVertex))
	r.Register("ENTITY.ADDTRIANGLE", "entity", runtime.AdaptLegacy(m.entAddTriangle))
	r.Register("ENTITY.UPDATEMESH", "entity", runtime.AdaptLegacy(m.entUpdateMesh))
	r.Register("ENTITY.VERTEXX", "entity", runtime.AdaptLegacy(m.entVertexX))
	r.Register("ENTITY.VERTEXY", "entity", runtime.AdaptLegacy(m.entVertexY))
	r.Register("ENTITY.VERTEXZ", "entity", runtime.AdaptLegacy(m.entVertexZ))

	r.Register("ENTITY.HIDE", "entity", runtime.AdaptLegacy(m.entHide))
	r.Register("ENTITY.SHOW", "entity", runtime.AdaptLegacy(m.entShow))
	r.Register("ENTITY.FREE", "entity", runtime.AdaptLegacy(m.entFree))
	r.Register("ENTITY.COPY", "entity", runtime.AdaptLegacy(m.entCopy))
	r.Register("ENTITY.SETNAME", "entity", m.entSetName)
	r.Register("ENTITY.FIND", "entity", m.entFind)

	r.Register("ENTITY.MOVERELATIVE", "entity", runtime.AdaptLegacy(m.entMoveRelative))
	r.Register("ENTITY.APPLYGRAVITY", "entity", runtime.AdaptLegacy(m.entApplyGravity))
	r.Register("ENTITY.GROUNDED", "entity", runtime.AdaptLegacy(m.entGrounded))
	r.Register("EntityGrounded", "entity", runtime.AdaptLegacy(m.entGrounded))
	r.Register("ENTITY.MOVECAMERARELATIVE", "entity", runtime.AdaptLegacy(m.entMoveCameraRelative))
	r.Register("EntityMoveCameraRelative", "entity", runtime.AdaptLegacy(m.entMoveCameraRelative))
	r.Register("ENTITY.SETMASS", "entity", runtime.AdaptLegacy(m.entSetMass))
	r.Register("ENTITY.SETFRICTION", "entity", runtime.AdaptLegacy(m.entSetFriction))
	r.Register("ENTITY.SETBOUNCE", "entity", runtime.AdaptLegacy(m.entSetBounce))

	r.Register("CAMERA.ORBITENTITY", "entity", m.camOrbitEntity)

	// Modern Blitz-style scene graph API (aliases)
	r.Register("DrawEntities", "entity", runtime.AdaptLegacy(m.entDrawAll))
	// MoveEntity: Blitz local-axis nudge (forward, right, up) — same as MOVEENTITY / ENTITY.MOVEENTITY.
	// Use ENTITY.TRANSLATEENTITY / TranslateEntity for world-space (dx, dy, dz).
	r.Register("MoveEntity", "entity", runtime.AdaptLegacy(m.entMove))
	r.Register("TranslateEntity", "entity", runtime.AdaptLegacy(m.entTranslate))
	r.Register("CreatePivot", "entity", runtime.AdaptLegacy(m.entCreatePivot))

	// Global Easy Mode Shorthands
	r.Register("CreateCube", "entity", runtime.AdaptLegacy(m.entCreateCubeEasy))
	r.Register("CreateSphere", "entity", runtime.AdaptLegacy(m.entCreateSphereEasy))
	r.Register("CreateCylinder", "entity", runtime.AdaptLegacy(m.entCreateCylinderEasy))
	r.Register("CreateCamera", "entity", runtime.AdaptLegacy(m.entCreateCameraEasy))
	r.Register("SCALENT", "entity", runtime.AdaptLegacy(m.entScaleEntity))
	r.Register("POSENT", "entity", runtime.AdaptLegacy(m.entSetPosition))
	r.Register("ENTITYCOLOR", "entity", runtime.AdaptLegacy(m.entColor))
	r.Register("ENTRAD", "entity", runtime.AdaptLegacy(m.entRadius))
	r.Register("ENTTYPE", "entity", runtime.AdaptLegacy(m.entType))
	r.Register("SetGravity", "entity", runtime.AdaptLegacy(m.entSetGravity))
	r.Register("ENTITYFLOOR", "entity", runtime.AdaptLegacy(m.entGrounded))
	r.Register("ENTITYJUMP", "entity", runtime.AdaptLegacy(m.entJump))
}

func (m *Module) entCreateCubeEasy(args []value.Value) (value.Value, error) {
	// CreateCube() → 1×1×1; CreateCube(parent#) → 1×1×1 child of parent;
	// CreateCube(w,h,d); CreateCube(parent#, w, h, d)
	switch len(args) {
	case 0:
		return m.entCreateBox([]value.Value{value.FromFloat(1), value.FromFloat(1), value.FromFloat(1)})
	case 1:
		pid, ok := m.entID(args[0])
		if !ok || pid < 1 {
			return value.Nil, fmt.Errorf("CreateCube(parent): parent must be entity#")
		}
		if m.store().ents[pid] == nil {
			return value.Nil, fmt.Errorf("CreateCube: unknown parent entity %d", pid)
		}
		v, err := m.entCreateBox([]value.Value{value.FromFloat(1), value.FromFloat(1), value.FromFloat(1)})
		if err != nil {
			return v, err
		}
		cid, _ := v.ToInt()
		_, err = m.entParent([]value.Value{value.FromInt(cid), value.FromInt(pid)})
		return v, err
	case 3:
		return m.entCreateBox(args)
	case 4:
		pid, ok := m.entID(args[0])
		if !ok || pid < 1 {
			return value.Nil, fmt.Errorf("CreateCube(parent,w,h,d): parent must be entity#")
		}
		if m.store().ents[pid] == nil {
			return value.Nil, fmt.Errorf("CreateCube: unknown parent entity %d", pid)
		}
		v, err := m.entCreateBox([]value.Value{args[1], args[2], args[3]})
		if err != nil {
			return v, err
		}
		cid, _ := v.ToInt()
		_, err = m.entParent([]value.Value{value.FromInt(cid), value.FromInt(pid)})
		return v, err
	default:
		return value.Nil, fmt.Errorf("CreateCube expects 0, 1 (parent#), 3 (w,h,d), or 4 (parent#,w,h,d) arguments")
	}
}

func (m *Module) entCreateSphereEasy(args []value.Value) (value.Value, error) {
	switch len(args) {
	case 0:
		return m.entCreateSphere([]value.Value{value.FromFloat(0.5), value.FromInt(16)})
	case 1:
		pid, ok := m.entID(args[0])
		if !ok || pid < 1 || m.store().ents[pid] == nil {
			return value.Nil, fmt.Errorf("CreateSphere: invalid parent")
		}
		v, err := m.entCreateSphere([]value.Value{value.FromFloat(0.5), value.FromInt(16)})
		if err != nil {
			return v, err
		}
		cid, _ := v.ToInt()
		_, err = m.entParent([]value.Value{value.FromInt(cid), value.FromInt(pid)})
		return v, err
	case 2:
		return m.entCreateSphere(args)
	case 3:
		pid, ok := m.entID(args[0])
		if !ok || pid < 1 || m.store().ents[pid] == nil {
			return value.Nil, fmt.Errorf("CreateSphere: invalid parent")
		}
		v, err := m.entCreateSphere([]value.Value{args[1], args[2]})
		if err != nil {
			return v, err
		}
		cid, _ := v.ToInt()
		_, err = m.entParent([]value.Value{value.FromInt(cid), value.FromInt(pid)})
		return v, err
	default:
		return value.Nil, fmt.Errorf("CreateSphere expects 0, 1 (parent#), 2 (radius#, seg#), or 3 (parent#, radius#, seg#)")
	}
}

func (m *Module) entCreateCameraEasy(args []value.Value) (value.Value, error) {
	if len(args) > 1 {
		return value.Nil, fmt.Errorf("CreateCamera expects 0 or 1 (reserved parent#) arguments")
	}
	reg := runtime.ActiveRegistry()
	if reg == nil {
		return value.Nil, fmt.Errorf("registry not active")
	}
	// Parent is accepted for API parity with Blitz; cameras are heap objects and are not parented to entities yet.
	_ = args
	return reg.Call("CAMERA.MAKE", nil)
}

func (m *Module) entScaleEntity(args []value.Value) (value.Value, error) {
	return m.entScale(args)
}

func (m *Module) entCreateSphere(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.CREATESPHERE expects 2 arguments (radius#, segments)")
	}
	rad, ok1 := argF32(args[0])
	seg, ok2 := args[1].ToInt()
	if !ok1 || !ok2 || seg < 3 {
		return value.Nil, fmt.Errorf("ENTITY.CREATESPHERE: radius numeric, segments int >= 3")
	}
	st := m.store()
	id := st.nextID
	st.nextID++
	e := newDefaultEnt(id)
	e.kind = entKindSphere
	e.radius = rad
	e.segH, e.segV = int32(seg), int32(seg)
	e.static = true
	e.w, e.h, e.d = rad*2, rad*2, rad*2
	st.ents[id] = e
	return value.FromInt(id), nil
}

func (m *Module) entCreateCylinderEasy(args []value.Value) (value.Value, error) {
	switch len(args) {
	case 3:
		return m.entCreateCylinder(args)
	case 4:
		pid, ok := m.entID(args[0])
		if !ok || pid < 1 || m.store().ents[pid] == nil {
			return value.Nil, fmt.Errorf("CreateCylinder: invalid parent")
		}
		v, err := m.entCreateCylinder([]value.Value{args[1], args[2], args[3]})
		if err != nil {
			return v, err
		}
		cid, _ := v.ToInt()
		_, err = m.entParent([]value.Value{value.FromInt(cid), value.FromInt(pid)})
		return v, err
	default:
		return value.Nil, fmt.Errorf("CreateCylinder expects 3 (radius#, height#, segments#) or 4 (parent#, radius#, height#, segments#)")
	}
}

func (m *Module) entCreateCylinder(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("ENTITY.CREATECYLINDER expects 3 arguments (radius#, height#, segments)")
	}
	rad, ok1 := argF32(args[0])
	h, ok2 := argF32(args[1])
	seg, ok3 := args[2].ToInt()
	if !ok1 || !ok2 || !ok3 || seg < 3 {
		return value.Nil, fmt.Errorf("ENTITY.CREATECYLINDER: invalid arguments")
	}
	st := m.store()
	id := st.nextID
	st.nextID++
	e := newDefaultEnt(id)
	e.kind = entKindCylinder
	e.radius = rad
	e.cylH = h
	e.segV = int32(seg)
	e.w, e.h, e.d = rad*2, h, rad*2
	e.static = true
	st.ents[id] = e
	return value.FromInt(id), nil
}

func (m *Module) entCreatePlane(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.CREATEPLANE expects 1 argument (size#)")
	}
	sz, ok := argF32(args[0])
	if !ok || sz <= 0 {
		return value.Nil, fmt.Errorf("ENTITY.CREATEPLANE: size must be positive")
	}
	st := m.store()
	id := st.nextID
	st.nextID++
	e := newDefaultEnt(id)
	e.kind = entKindPlane
	e.w, e.h, e.d = sz, 0.01, sz
	e.static = true
	st.ents[id] = e
	return value.FromInt(id), nil
}

func (m *Module) entCreateMesh(args []value.Value) (value.Value, error) {
	if len(args) != 0 && len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.CREATEMESH expects 0 arguments or 1 (parentEntity#)")
	}
	b := &meshBuilderObj{}
	bid, err := m.h.Alloc(b)
	if err != nil {
		return value.Nil, err
	}
	st := m.store()
	id := st.nextID
	st.nextID++
	e := newDefaultEnt(id)
	e.kind = entKindMesh
	e.hasRLModel = false
	e.hidden = true
	e.static = true
	e.procMeshH = bid
	st.ents[id] = e
	v := value.FromInt(id)
	if len(args) == 1 {
		pid, okp := m.entID(args[0])
		if !okp || pid < 1 {
			return value.Nil, fmt.Errorf("ENTITY.CREATEMESH: invalid parent entity")
		}
		if st.ents[pid] == nil {
			return value.Nil, fmt.Errorf("ENTITY.CREATEMESH: unknown parent %d", pid)
		}
		_, err := m.entParent([]value.Value{v, value.FromInt(pid)})
		if err != nil {
			return value.Nil, err
		}
	}
	return v, nil
}

func (m *Module) entLoadMesh(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 && len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.LOADMESH expects 1 argument (path$) or 2 (path$, parentEntity#)")
	}
	if args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("ENTITY.LOADMESH: path must be string")
	}
	path, ok := rt.Heap.GetString(int32(args[0].IVal))
	if !ok || path == "" {
		return value.Nil, fmt.Errorf("ENTITY.LOADMESH: invalid string")
	}
	mod := rl.LoadModel(path)
	if mod.MeshCount <= 0 {
		rl.UnloadModel(mod)
		return value.Nil, fmt.Errorf("ENTITY.LOADMESH: failed to load %q", path)
	}
	st := m.store()
	id := st.nextID
	st.nextID++
	e := newDefaultEnt(id)
	e.kind = entKindModel
	e.rlModel = mod
	e.hasRLModel = true
	e.loadPath = path
	e.static = true
	st.ents[id] = e
	v := value.FromInt(id)
	if len(args) == 2 {
		pid, okp := m.entID(args[1])
		if !okp || pid < 1 {
			return value.Nil, fmt.Errorf("ENTITY.LOADMESH: invalid parent entity")
		}
		if st.ents[pid] == nil {
			return value.Nil, fmt.Errorf("ENTITY.LOADMESH: unknown parent %d", pid)
		}
		_, err := m.entParent([]value.Value{v, value.FromInt(pid)})
		if err != nil {
			return value.Nil, err
		}
	}
	return v, nil
}

func (m *Module) entLoadAnimatedMesh(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	v, err := m.entLoadMesh(rt, args...)
	if err != nil {
		return v, err
	}
	id, ok := v.ToInt()
	if !ok || id < 1 {
		return v, nil
	}
	e := m.store().ents[id]
	if e == nil || !e.hasRLModel || e.loadPath == "" {
		return v, nil
	}
	anims := rl.LoadModelAnimations(e.loadPath)
	if len(anims) > 0 {
		e.modelAnims = anims
		e.animLen = float32(anims[0].FrameCount)
		rl.UpdateModelAnimation(e.rlModel, anims[0], 0)
		rl.UpdateModelAnimationBones(e.rlModel, anims[0], 0)
	}
	return v, nil
}

func (m *Module) entRotateEntityAbs(args []value.Value) (value.Value, error) {
	if len(args) != 4 && len(args) != 5 {
		return value.Nil, fmt.Errorf("ENTITY.ROTATEENTITY expects 4–5 arguments (entity#, pitch#, yaw#, roll# [, global])")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.ROTATEENTITY: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.ROTATEENTITY: unknown entity")
	}
	p, ok1 := argF32(args[1])
	y, ok2 := argF32(args[2])
	r, ok3 := argF32(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("ENTITY.ROTATEENTITY: angles must be numeric")
	}
	_ = args // global reserved for future local-vs-world rotation composition
	e.pitch, e.yaw, e.roll = p, y, r
	return value.Nil, nil
}

func (m *Module) entEntityX(args []value.Value) (value.Value, error) {
	return m.getCoord(args, func(e *ent) float32 { return m.worldPos(e).X }, func(e *ent) float32 { return e.pos.X })
}
func (m *Module) entEntityY(args []value.Value) (value.Value, error) {
	return m.getCoord(args, func(e *ent) float32 { return m.worldPos(e).Y }, func(e *ent) float32 { return e.pos.Y })
}
func (m *Module) entEntityZ(args []value.Value) (value.Value, error) {
	return m.getCoord(args, func(e *ent) float32 { return m.worldPos(e).Z }, func(e *ent) float32 { return e.pos.Z })
}
func (m *Module) entEntityPitch(args []value.Value) (value.Value, error) {
	return m.getCoord(args, func(e *ent) float32 {
		pp, _, _ := m.worldEuler(e)
		return pp
	}, func(e *ent) float32 { return e.pitch })
}
func (m *Module) entEntityYaw(args []value.Value) (value.Value, error) {
	return m.getCoord(args, func(e *ent) float32 {
		_, yy, _ := m.worldEuler(e)
		return yy
	}, func(e *ent) float32 { return e.yaw })
}
func (m *Module) entEntityRoll(args []value.Value) (value.Value, error) {
	return m.getCoord(args, func(e *ent) float32 {
		_, _, rr := m.worldEuler(e)
		return rr
	}, func(e *ent) float32 { return e.roll })
}

func (m *Module) getCoord(args []value.Value, world, local func(*ent) float32) (value.Value, error) {
	if len(args) < 1 || len(args) > 2 {
		return value.Nil, fmt.Errorf("expected 1–2 arguments (entity# [, global])")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	global := false
	if len(args) == 2 {
		switch args[1].Kind {
		case value.KindBool:
			global = args[1].IVal != 0
		case value.KindInt:
			global = args[1].IVal != 0
		default:
			return value.Nil, fmt.Errorf("global must be bool or 0/1")
		}
	}
	if global {
		return value.FromFloat(float64(world(e))), nil
	}
	return value.FromFloat(float64(local(e))), nil
}

func (m *Module) entParent(args []value.Value) (value.Value, error) {
	if len(args) != 2 && len(args) != 3 {
		return value.Nil, fmt.Errorf("ENTITY.PARENT expects 2–3 arguments (entity#, parentEntity# [, global])")
	}
	cid, ok := m.entID(args[0])
	pid, ok2 := m.entID(args[1])
	if !ok || !ok2 || cid < 1 || pid < 1 {
		return value.Nil, fmt.Errorf("invalid entity ids")
	}
	st := m.store()
	child := st.ents[cid]
	parent := st.ents[pid]
	if child == nil || parent == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	// preserve world position when reparenting (global default true for simplicity)
	global := true
	if len(args) == 3 {
		switch args[2].Kind {
		case value.KindBool:
			global = args[2].IVal != 0
		case value.KindInt:
			global = args[2].IVal != 0
		default:
			return value.Nil, fmt.Errorf("global must be bool or 0/1")
		}
	}
	oldPID := child.parentID
	if oldPID == pid {
		return value.Nil, nil
	}
	if oldPID >= 1 {
		childLinkRemove(st, oldPID, cid)
	}
	if global {
		wp := m.worldPos(child)
		child.parentID = pid
		m.setLocalFromWorld(child, wp.X, wp.Y, wp.Z)
	} else {
		child.parentID = pid
	}
	childLinkAdd(st, pid, cid)
	return value.Nil, nil
}

func (m *Module) entParentClear(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.PARENTCLEAR expects entity#")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	st := m.store()
	e := st.ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	if e.parentID >= 1 {
		childLinkRemove(st, e.parentID, id)
	}
	wp := m.worldPos(e)
	e.parentID = 0
	e.pos = wp
	return value.Nil, nil
}

func (m *Module) entAlpha(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.ALPHA expects (entity#, alpha#)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	a, ok1 := argF32(args[1])
	if !ok1 {
		return value.Nil, fmt.Errorf("alpha must be numeric")
	}
	if a < 0 {
		a = 0
	}
	if a > 1 {
		a = 1
	}
	e.alpha = a
	return value.Nil, nil
}

func (m *Module) entShininess(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.SHININESS expects (entity#, amount#)")
	}
	id, _ := m.entID(args[0])
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	s, _ := argF32(args[1])
	e.shininess = s
	return value.Nil, nil
}

func (m *Module) entTexture(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.TEXTURE expects (entity#, textureHandle)")
	}
	id, _ := m.entID(args[0])
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	h, ok := argHandle(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("texture must be handle")
	}
	e.texHandle = h
	return value.Nil, nil
}

func (m *Module) entFX(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.FX expects (entity#, flags)")
	}
	id, _ := m.entID(args[0])
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	f, _ := args[1].ToInt()
	e.fxFlags = int32(f)
	return value.Nil, nil
}

func (m *Module) entBlend(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.BLEND expects (entity#, mode)")
	}
	id, _ := m.entID(args[0])
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	b, _ := args[1].ToInt()
	e.blendMode = int32(b)
	return value.Nil, nil
}

func (m *Module) entOrder(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.ORDER expects (entity#, order)")
	}
	id, _ := m.entID(args[0])
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	o, _ := args[1].ToInt()
	e.drawOrder = int32(o)
	return value.Nil, nil
}

func (m *Module) entType(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.TYPE expects (entity#, typeID)")
	}
	id, _ := m.entID(args[0])
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	t, _ := args[1].ToInt()
	e.collType = int32(t)
	return value.Nil, nil
}

func (m *Module) entCollide(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.COLLIDE expects (entity#, otherTypeID)")
	}
	id, ok := m.entID(args[0])
	tid, ok2 := m.entID(args[1])
	if !ok || !ok2 || id < 1 {
		return value.Nil, fmt.Errorf("invalid arguments")
	}
	a := m.store().ents[id]
	if a == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	for _, b := range m.store().ents {
		if b.id == id || b.static {
			continue
		}
		if int64(b.collType) != tid {
			continue
		}
		if !a.useSphere || !b.useSphere {
			continue
		}
		pa := m.worldPos(a)
		pb := m.worldPos(b)
		if rl.Vector3Distance(pa, pb) < a.radius+b.radius {
			return value.FromInt(b.id), nil
		}
	}
	return value.FromInt(0), nil
}

func (m *Module) entCollisionX(args []value.Value) (value.Value, error) {
	return m.hitComp(args,
		func(e *ent) float64 { return float64(e.hitX) },
		func(e *ent, i int) float64 {
			if i < 0 || i >= len(e.hitPos) {
				return 0
			}
			return float64(e.hitPos[i].X)
		})
}
func (m *Module) entCollisionY(args []value.Value) (value.Value, error) {
	return m.hitComp(args,
		func(e *ent) float64 { return float64(e.hitY) },
		func(e *ent, i int) float64 {
			if i < 0 || i >= len(e.hitPos) {
				return 0
			}
			return float64(e.hitPos[i].Y)
		})
}
func (m *Module) entCollisionZ(args []value.Value) (value.Value, error) {
	return m.hitComp(args,
		func(e *ent) float64 { return float64(e.hitZ) },
		func(e *ent, i int) float64 {
			if i < 0 || i >= len(e.hitPos) {
				return 0
			}
			return float64(e.hitPos[i].Z)
		})
}
func (m *Module) entCollisionNX(args []value.Value) (value.Value, error) {
	return m.hitComp(args,
		func(e *ent) float64 { return float64(e.hitNX) },
		func(e *ent, i int) float64 {
			if i < 0 || i >= len(e.hitN) {
				return 0
			}
			return float64(e.hitN[i].X)
		})
}
func (m *Module) entCollisionNY(args []value.Value) (value.Value, error) {
	return m.hitComp(args,
		func(e *ent) float64 { return float64(e.hitNY) },
		func(e *ent, i int) float64 {
			if i < 0 || i >= len(e.hitN) {
				return 0
			}
			return float64(e.hitN[i].Y)
		})
}
func (m *Module) entCollisionNZ(args []value.Value) (value.Value, error) {
	return m.hitComp(args,
		func(e *ent) float64 { return float64(e.hitNZ) },
		func(e *ent, i int) float64 {
			if i < 0 || i >= len(e.hitN) {
				return 0
			}
			return float64(e.hitN[i].Z)
		})
}

func (m *Module) hitComp(args []value.Value, last func(*ent) float64, atIndex func(*ent, int) float64) (value.Value, error) {
	if len(args) != 1 && len(args) != 2 {
		return value.Nil, fmt.Errorf("expects entity# [, index]")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	if len(args) == 2 {
		idx, ok := args[1].ToInt()
		if !ok || idx < 0 {
			return value.Nil, fmt.Errorf("index must be non-negative int")
		}
		return value.FromFloat(atIndex(e, int(idx))), nil
	}
	if !e.hasHit {
		return value.FromFloat(0), nil
	}
	return value.FromFloat(last(e)), nil
}

func (m *Module) entDistance(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.DISTANCE expects (entityA#, entityB#)")
	}
	ia, ok1 := m.entID(args[0])
	ib, ok2 := m.entID(args[1])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("invalid ids")
	}
	a := m.store().ents[ia]
	b := m.store().ents[ib]
	if a == nil || b == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	d := rl.Vector3Distance(m.worldPos(a), m.worldPos(b))
	return value.FromFloat(float64(d)), nil
}

func (m *Module) entVelocity(args []value.Value) (value.Value, error) {
	if len(args) != 1 && len(args) != 4 {
		return value.Nil, fmt.Errorf("ENTITY.VELOCITY: (entity#) get or (entity#, vx, vy, vz) set")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	if len(args) == 1 {
		if m.h == nil {
			return value.Nil, runtime.Errorf("heap not bound")
		}
		return mbmatrix.AllocVec3Value(m.h, e.vel.X, e.vel.Y, e.vel.Z)
	}
	vx, ok1 := argF32(args[1])
	vy, ok2 := argF32(args[2])
	vz, ok3 := argF32(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("velocity must be numeric")
	}
	e.vel = rl.Vector3{X: vx, Y: vy, Z: vz}
	e.static = false
	return value.Nil, nil
}

func (m *Module) entAddForce(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("ENTITY.ADDFORCE expects (entity#, fx#, fy#, fz#)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	fx, ok1 := argF32(args[1])
	fy, ok2 := argF32(args[2])
	fz, ok3 := argF32(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("force must be numeric")
	}
	invM := float32(1)
	if e.mass > 1e-6 {
		invM = 1 / e.mass
	}
	e.vel.X += fx * invM
	e.vel.Y += fy * invM
	e.vel.Z += fz * invM
	e.static = false
	if fy > 0.5 {
		e.groundCoyoteLeft = 0
	}
	return value.Nil, nil
}

func (m *Module) entSetSlide(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.SLIDE expects (entity#, enable)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	e.slide = args[1].Kind == value.KindBool && args[1].IVal != 0
	if args[1].Kind == value.KindInt {
		e.slide = args[1].IVal != 0
	}
	return value.Nil, nil
}

func (m *Module) entPick(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.PICK expects (entity#, range#)")
	}
	id, ok := m.entID(args[0])
	rng, ok2 := argF32(args[1])
	if !ok || !ok2 || id < 1 {
		return value.Nil, fmt.Errorf("invalid arguments")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	fwd := forwardFromYawPitch(e.yaw, e.pitch)
	origin := m.worldPos(e)
	end := rl.Vector3Add(origin, rl.Vector3Scale(fwd, rng))
	bestID := int64(0)
	bestT := float32(1e20)
	for _, s := range m.store().ents {
		if !s.static || s.id == e.id {
			continue
		}
		smn, smx := m.aabbWorldMinMax(s)
		t := rayAABB(origin, end, smn, smx)
		if t >= 0 && t < bestT {
			bestT = t
			bestID = s.id
		}
	}
	return value.FromInt(bestID), nil
}

func rayAABB(origin, end rl.Vector3, mn, mx rl.Vector3) float32 {
	dir := rl.Vector3Subtract(end, origin)
	tmax := rl.Vector3Length(dir)
	if tmax < 1e-8 {
		return -1
	}
	dir = rl.Vector3Normalize(dir)
	tmin := float32(0)
	tmaxf := tmax
	for a := 0; a < 3; a++ {
		var invD, o, n, x float32
		switch a {
		case 0:
			invD = 1 / dir.X
			o, n, x = origin.X, mn.X, mx.X
		case 1:
			invD = 1 / dir.Y
			o, n, x = origin.Y, mn.Y, mx.Y
		default:
			invD = 1 / dir.Z
			o, n, x = origin.Z, mn.Z, mx.Z
		}
		t0 := (n - o) * invD
		t1 := (x - o) * invD
		if t0 > t1 {
			t0, t1 = t1, t0
		}
		tmin = maxFloat32(tmin, t0)
		tmaxf = minFloat32(tmaxf, t1)
		if tmin > tmaxf {
			return -1
		}
	}
	if tmin >= 0 {
		return tmin
	}
	if tmaxf >= 0 {
		return tmaxf
	}
	return -1
}

func maxFloat32(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func (m *Module) entPickMode(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.PICKMODE expects (entity#, mode)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	md, _ := args[1].ToInt()
	e.pickMode = int32(md)
	return value.Nil, nil
}

func (m *Module) entPointEntity(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.POINTENTITY expects (entity#, targetEntity#)")
	}
	id, ok := m.entID(args[0])
	tid, ok2 := m.entID(args[1])
	if !ok || !ok2 {
		return value.Nil, fmt.Errorf("invalid entity ids")
	}
	e := m.store().ents[id]
	t := m.store().ents[tid]
	if e == nil || t == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	wp := m.worldPos(e)
	wt := m.worldPos(t)
	d := rl.Vector3Subtract(wt, wp)
	d.Y = 0
	if rl.Vector3Length(d) < 1e-6 {
		return value.Nil, nil
	}
	d = rl.Vector3Normalize(d)
	e.yaw = float32(math.Atan2(float64(d.X), float64(d.Z)))
	return value.Nil, nil
}

func (m *Module) entAlignToVector(args []value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("ENTITY.ALIGNTOVECTOR expects (entity#, vx#, vy#, vz#, axis)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	vx, ok1 := argF32(args[1])
	vy, ok2 := argF32(args[2])
	vz, ok3 := argF32(args[3])
	ax, ok4 := args[4].ToInt()
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("numeric vector and axis int required")
	}
	v := rl.Vector3Normalize(rl.Vector3{X: vx, Y: vy, Z: vz})
	_ = ax
	// Align local +Z to v (yaw/pitch)
	e.yaw = float32(math.Atan2(float64(v.X), float64(v.Z)))
	vyClamped := float64(v.Y)
	if vyClamped > 1 {
		vyClamped = 1
	}
	if vyClamped < -1 {
		vyClamped = -1
	}
	e.pitch = float32(math.Asin(vyClamped))
	e.roll = 0
	return value.Nil, nil
}

func (m *Module) entAnimate(args []value.Value) (value.Value, error) {
	// Mode: 0–1 = loop, 2 = ping-pong, 3+ = clamp at clip end. Dual-pose cross-fade is not implemented yet (Raylib single-clip pose).
	if len(args) < 1 || len(args) > 3 {
		return value.Nil, fmt.Errorf("ENTITY.ANIMATE expects (entity# [, mode, speed#])")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	if len(args) >= 2 {
		md, _ := args[1].ToInt()
		e.animMode = int32(md)
	}
	if len(args) >= 3 {
		s, _ := argF32(args[2])
		e.animSpeed = s
	}
	return value.Nil, nil
}

func (m *Module) entSetAnimTime(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.SETANIMTIME expects (entity#, time#)")
	}
	id, _ := m.entID(args[0])
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	t, _ := argF32(args[1])
	e.animTime = t
	return value.Nil, nil
}

func (m *Module) entAnimTime(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.ANIMTIME expects entity#")
	}
	id, _ := m.entID(args[0])
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	return value.FromFloat(float64(e.animTime)), nil
}

func (m *Module) entAnimLength(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.ANIMLENGTH expects entity#")
	}
	id, _ := m.entID(args[0])
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	if len(e.modelAnims) > 0 {
		ai := e.animIndex
		if ai < 0 || int(ai) >= len(e.modelAnims) {
			ai = 0
		}
		return value.FromFloat(float64(e.modelAnims[ai].FrameCount)), nil
	}
	return value.FromFloat(float64(e.animLen)), nil
}

func (m *Module) entHide(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.HIDE expects entity#")
	}
	id, _ := m.entID(args[0])
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	e.hidden = true
	return value.Nil, nil
}

func (m *Module) entShow(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.SHOW expects entity#")
	}
	id, _ := m.entID(args[0])
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	e.hidden = false
	return value.Nil, nil
}

func (m *Module) entFree(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.FREE expects entity#")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	st := m.store()
	if st.ents[id] == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	m.purgeEntityByID(id)
	return value.Nil, nil
}

func (m *Module) entCopy(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.COPY expects entity#")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	src := m.store().ents[id]
	if src == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	cp := *src
	cp.id = 0
	cp.parentID = 0
	cp.name = ""
	cp.boneHostID = 0
	cp.boneIndex = -1
	cp.boneWorldValid = false
	cp.brushH = 0
	if src.procMeshH != 0 {
		return value.Nil, fmt.Errorf("ENTITY.COPY: procedural mesh entities cannot be copied yet")
	}
	if cp.hasRLModel {
		if cp.loadPath == "" {
			return value.Nil, fmt.Errorf("ENTITY.COPY: model without load path (e.g. CREATEMESH) cannot be duplicated yet")
		}
		mod := rl.LoadModel(cp.loadPath)
		if mod.MeshCount <= 0 {
			rl.UnloadModel(mod)
			return value.Nil, fmt.Errorf("ENTITY.COPY: failed to load model %q", cp.loadPath)
		}
		cp.rlModel = mod
		cp.modelAnims = nil
		if anims := rl.LoadModelAnimations(cp.loadPath); len(anims) > 0 {
			cp.modelAnims = anims
		}
	}
	st := m.store()
	nid := st.nextID
	st.nextID++
	cp.id = nid
	st.ents[nid] = &cp
	return value.FromInt(nid), nil
}

func (m *Module) entSetName(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.SETNAME expects (entity#, name$)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	if args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("name must be string")
	}
	name, ok2 := rt.Heap.GetString(int32(args[1].IVal))
	if !ok2 {
		return value.Nil, fmt.Errorf("invalid string")
	}
	if e.name != "" {
		delete(m.store().byName, strings.ToUpper(e.name))
	}
	e.name = name
	m.store().byName[strings.ToUpper(name)] = id
	return value.Nil, nil
}

func (m *Module) entFind(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.FIND expects name$")
	}
	if args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("name must be string")
	}
	name, ok := rt.Heap.GetString(int32(args[0].IVal))
	if !ok {
		return value.Nil, fmt.Errorf("invalid string")
	}
	id, ok2 := m.store().byName[strings.ToUpper(name)]
	if !ok2 {
		return value.FromInt(0), nil
	}
	return value.FromInt(id), nil
}

func (m *Module) entMoveRelative(args []value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("ENTITY.MOVERELATIVE expects (entity#, forward#, right#, speed#, dt#)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	f, ok1 := argF32(args[1])
	rg, ok2 := argF32(args[2])
	sp, ok3 := argF32(args[3])
	dt, ok4 := argF32(args[4])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("numeric args required")
	}
	fwd, right, up := localAxes(e.yaw, e.pitch)
	delta := rl.Vector3Add(rl.Vector3Add(rl.Vector3Scale(fwd, f*sp*dt), rl.Vector3Scale(right, rg*sp*dt)), rl.Vector3Scale(up, 0))
	wp := m.worldPos(e)
	nw := rl.Vector3Add(wp, delta)
	m.setLocalFromWorld(e, nw.X, nw.Y, nw.Z)
	return value.Nil, nil
}

func (m *Module) entApplyGravity(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("ENTITY.APPLYGRAVITY expects (entity#, gravity#, dt#)")
	}
	id, ok := m.entID(args[0])
	g, ok1 := argF32(args[1])
	dt, ok2 := argF32(args[2])
	if !ok || !ok1 || !ok2 || id < 1 {
		return value.Nil, fmt.Errorf("invalid arguments")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	e.vel.Y += g * dt
	e.static = false
	return value.Nil, nil
}

func (m *Module) entGrounded(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.GROUNDED expects entity#")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	return value.FromBool(e.jumpGrounded), nil
}

func (m *Module) entMoveCameraRelative(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("EntityMoveCameraRelative expects (entity#, forward#, strafe#, camera)")
	}
	if m.h == nil {
		return value.Nil, fmt.Errorf("EntityMoveCameraRelative: heap not bound")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	f, ok1 := argF32(args[1])
	s, ok2 := argF32(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("forward/strafe must be numeric")
	}
	if args[3].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("camera handle required")
	}
	ch := heap.Handle(args[3].IVal)
	fwd, right, err := mbcamera.CameraXZStrafeBasis(m.h, ch)
	if err != nil {
		return value.Nil, err
	}
	dx := fwd.X*f + right.X*s
	dz := fwd.Z*f + right.Z*s
	wp := m.worldPos(e)
	nw := rl.Vector3Add(wp, rl.Vector3{X: dx, Y: 0, Z: dz})
	m.setLocalFromWorld(e, nw.X, nw.Y, nw.Z)
	return value.Nil, nil
}

func (m *Module) entSetMass(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.SETMASS expects (entity#, mass#)")
	}
	id, _ := m.entID(args[0])
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	mass, _ := argF32(args[1])
	e.mass = mass
	return value.Nil, nil
}

func (m *Module) entSetFriction(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.SETFRICTION expects (entity#, amount#)")
	}
	id, _ := m.entID(args[0])
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	f, _ := argF32(args[1])
	e.friction = f
	return value.Nil, nil
}

func (m *Module) entSetBounce(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.SETBOUNCE expects (entity#, amount#)")
	}
	id, _ := m.entID(args[0])
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	b, _ := argF32(args[1])
	e.bounce = b
	return value.Nil, nil
}

func (m *Module) camOrbitEntity(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("CAMERA.ORBITENTITY expects (camera, entity#, yaw#, pitch#, dist#)")
	}
	ch, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("invalid camera")
	}
	eid, ok2 := m.entID(args[1])
	if !ok2 || eid < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[eid]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	yaw, ok3 := argF32(args[2])
	pitch, ok4 := argF32(args[3])
	dist, ok5 := argF32(args[4])
	if !ok3 || !ok4 || !ok5 {
		return value.Nil, fmt.Errorf("numeric yaw/pitch/dist required")
	}
	wp := m.worldPos(e)
	if err := mbcamera.ApplySetOrbit(m.h, ch, wp.X, wp.Y, wp.Z, yaw, pitch, dist); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}
func (m *Module) entCreateSpriteEntity(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 && len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.CREATESPRITE expects (path$) or (path$, parentEntity#)")
	}
	return m.entLoadSprite(rt, args...)
}

func (m *Module) entLoadSprite(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if (len(args) != 1 && len(args) != 2) || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("LOADSPRITE expects (path$) or (path$, parentEntity#)")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	tex := rl.LoadTexture(path)
	if tex.ID <= 0 {
		return value.Nil, fmt.Errorf("LOADSPRITE: failed to load %q", path)
	}
	// Store texture in heap for management
	th, _ := rt.Heap.Alloc(&textureObj{tex: tex})

	st := m.store()
	id := st.nextID
	st.nextID++
	e := newDefaultEnt(id)
	e.kind = entKindMesh
	e.isSprite = true
	e.spriteMode = 1 // default Y-billboard
	e.texHandle = th
	e.scale = rl.Vector3{X: 1, Y: 1, Z: 1}
	e.w = float32(tex.Width) / 100.0 // Reasonable default size
	e.h = float32(tex.Height) / 100.0
	st.ents[id] = e
	if len(args) == 2 {
		pid, ok := m.entID(args[1])
		if !ok || pid < 1 || st.ents[pid] == nil {
			return value.Nil, fmt.Errorf("LOADSPRITE: invalid parent entity")
		}
		e.parentID = pid
		childLinkAdd(st, pid, id)
	}
	return value.FromInt(id), nil
}

func (m *Module) entScaleSprite(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("SCALESPRITE expects (sprite, x#, y#)")
	}
	id, _ := m.entID(args[0])
	e := m.store().ents[id]
	if e == nil || !e.isSprite {
		return value.Nil, fmt.Errorf("invalid sprite")
	}
	sx, _ := args[1].ToFloat()
	sy, _ := args[2].ToFloat()
	e.scale.X = float32(sx)
	e.scale.Y = float32(sy)
	return value.Nil, nil
}

func (m *Module) entSpriteMode(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("SPRITEMODE expects (sprite, mode)")
	}
	id, _ := m.entID(args[0])
	e := m.store().ents[id]
	if e == nil || !e.isSprite {
		return value.Nil, fmt.Errorf("invalid sprite")
	}
	mode, _ := args[1].ToInt()
	e.spriteMode = int32(mode)
	return value.Nil, nil
}

type textureObj struct {
	tex rl.Texture2D
}

func (o *textureObj) TypeName() string { return "Texture" }
func (o *textureObj) TypeTag() uint16  { return heap.TagTexture }
func (o *textureObj) Free()            { rl.UnloadTexture(o.tex) }
