//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"
	"math"

	"moonbasic/runtime"
	mbcamera "moonbasic/runtime/camera"
	"moonbasic/runtime/mbmodel3d"
	"moonbasic/runtime/texture"
	mbtime "moonbasic/runtime/time"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func registerModernEntityFX(m *Module, r runtime.Registrar) {
	r.Register("EntityPBR", "entity", runtime.AdaptLegacy(m.entityPBR))
	r.Register("EntityNormalMap", "entity", runtime.AdaptLegacy(m.entityNormalMap))
	r.Register("EntityEmission", "entity", runtime.AdaptLegacy(m.entityEmission))
	r.Register("EntityMass", "entity", runtime.AdaptLegacy(m.entSetMass))
	r.Register("EntityFriction", "entity", runtime.AdaptLegacy(m.entSetFriction))
	r.Register("EntityRestitution", "entity", runtime.AdaptLegacy(m.entSetBounce))
	r.Register("ApplyEntityImpulse", "entity", runtime.AdaptLegacy(m.entAddForce))
	r.Register("CameraSmoothFollow", "entity", m.cameraSmoothFollow)
	r.Register("CreateVehicle", "entity", runtime.AdaptLegacy(m.createVehicleStub))
	r.Register("AddWheel", "entity", runtime.AdaptLegacy(m.addWheelStub))
}

// EntityPBR(entity#, metal#, rough#) — shared PBR shader + metallic/roughness factors.
func (m *Module) entityPBR(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("EntityPBR expects (entity#, metal#, rough#)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("EntityPBR: invalid entity")
	}
	metal, ok1 := argF32(args[1])
	rough, ok2 := argF32(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("EntityPBR: metal and rough must be numeric")
	}
	e := m.store().ents[id]
	if e == nil || !e.hasRLModel {
		return value.Nil, fmt.Errorf("EntityPBR: entity has no loaded model")
	}
	mbmodel3d.ConvertEntityModelMaterialsPBR(&e.rlModel, metal, rough)
	return value.Nil, nil
}

func (m *Module) entityNormalMap(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("EntityNormalMap expects (entity#, textureHandle)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("EntityNormalMap: invalid entity")
	}
	th, ok := argHandle(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("EntityNormalMap: texture handle required")
	}
	e := m.store().ents[id]
	if e == nil || !e.hasRLModel {
		return value.Nil, fmt.Errorf("EntityNormalMap: entity has no loaded model")
	}
	tex, err := texture.ForBinding(m.h, th)
	if err != nil {
		return value.Nil, err
	}
	mats := e.rlModel.GetMaterials()
	for i := range mats {
		rl.SetMaterialTexture(&mats[i], rl.MapNormal, tex)
	}
	return value.Nil, nil
}

func (m *Module) entityEmission(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("EntityEmission expects (entity#, textureHandle, power#)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("EntityEmission: invalid entity")
	}
	th, ok := argHandle(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("EntityEmission: texture handle required")
	}
	pwr, ok := argF32(args[2])
	if !ok {
		return value.Nil, fmt.Errorf("EntityEmission: power must be numeric")
	}
	e := m.store().ents[id]
	if e == nil || !e.hasRLModel {
		return value.Nil, fmt.Errorf("EntityEmission: entity has no loaded model")
	}
	tex, err := texture.ForBinding(m.h, th)
	if err != nil {
		return value.Nil, err
	}
	mats := e.rlModel.GetMaterials()
	for i := range mats {
		rl.SetMaterialTexture(&mats[i], rl.MapEmission, tex)
		mats[i].GetMap(rl.MapEmission).Value = pwr
	}
	return value.Nil, nil
}

// CameraSmoothFollow(camera, entity#, lerp#) — third-person follow using current camera offset from target (XZ distance and Y height).
func (m *Module) cameraSmoothFollow(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("CameraSmoothFollow expects (camera, entity#, lerp#)")
	}
	ch, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CameraSmoothFollow: camera handle required")
	}
	eid, ok := args[1].ToInt()
	if !ok || eid < 1 {
		return value.Nil, fmt.Errorf("CameraSmoothFollow: entity# required")
	}
	lerp, ok := argF32(args[2])
	if !ok {
		return value.Nil, fmt.Errorf("CameraSmoothFollow: lerp must be numeric")
	}
	e := m.store().ents[eid]
	if e == nil {
		return value.Nil, fmt.Errorf("CameraSmoothFollow: unknown entity %d", eid)
	}
	cp, ok := mbcamera.CameraWorldPosition(m.h, ch)
	if !ok {
		return value.Nil, fmt.Errorf("CameraSmoothFollow: invalid camera")
	}
	wp := m.worldPos(e)
	dx := cp.X - wp.X
	dz := cp.Z - wp.Z
	dist := float32(math.Sqrt(float64(dx*dx + dz*dz)))
	height := cp.Y - wp.Y
	if dist < 0.05 {
		dist = 5
	}
	dt := mbtime.DeltaSeconds(rt)
	if dt <= 0 {
		dt = 1.0 / 60.0
	}
	if err := mbcamera.ThirdPersonFollowStep(m.h, ch, wp.X, wp.Y, wp.Z, e.yaw, dist, height, lerp, dt); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) createVehicleStub(args []value.Value) (value.Value, error) {
	_ = args
	return value.Nil, fmt.Errorf("CreateVehicle: Jolt vehicle constraint is not wired in this build (use ENTITY physics + mesh until vehicle module lands)")
}

func (m *Module) addWheelStub(args []value.Value) (value.Value, error) {
	_ = args
	return value.Nil, fmt.Errorf("AddWheel: Jolt vehicle constraint is not wired in this build")
}
