//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"
	"math"

	"moonbasic/runtime"
	mbcamera "moonbasic/runtime/camera"
	mbmatrix "moonbasic/runtime/mbmatrix"
	mbmodel3d "moonbasic/runtime/mbmodel3d"
	mbphysics3d "moonbasic/runtime/physics3d"
	"moonbasic/runtime/water"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// characterGroundNormalFn is set by runtime/player for entities with PLAYER.CREATE (Jolt CharacterVirtual).
var characterGroundNormalFn func(entityID int64) (nx, ny, nz float64, ok bool)

// SetCharacterGroundNormalResolver registers a Jolt ground-normal resolver (entity# → world normal). Clear with nil.
func SetCharacterGroundNormalResolver(fn func(int64) (float64, float64, float64, bool)) {
	characterGroundNormalFn = fn
}

func registerEntityGameplayIntelAPI(m *Module, r runtime.Registrar) {
	r.Register("ENTITY.SETCOLLISIONGROUP", "entity", runtime.AdaptLegacy(m.entCollisionLayer))
	r.Register("EntitySetCollisionGroup", "entity", runtime.AdaptLegacy(m.entCollisionLayer))
	r.Register("ENTITY.CHECKCOLLISION", "entity", runtime.AdaptLegacy(m.entEntityCollidedPair))
	r.Register("EntityCheckCollision", "entity", runtime.AdaptLegacy(m.entEntityCollidedPair))
	r.Register("ENTITY.RAYCAST", "entity", runtime.AdaptLegacy(m.entRaycast))
	r.Register("EntityRaycast", "entity", runtime.AdaptLegacy(m.entRaycast))
	r.Register("ENTITY.GETGROUNDNORMAL", "entity", runtime.AdaptLegacy(m.entGetGroundNormal))
	r.Register("EntityGetGroundNormal", "entity", runtime.AdaptLegacy(m.entGetGroundNormal))
	r.Register("ENTITY.APPLYIMPULSE", "entity", runtime.AdaptLegacy(m.entAddForce))
	r.Register("EntityApplyImpulse", "entity", runtime.AdaptLegacy(m.entAddForce))
	r.Register("ENTITY.CANSEE", "entity", runtime.AdaptLegacy(m.entCanSee))
	r.Register("EntityCanSee", "entity", runtime.AdaptLegacy(m.entCanSee))
	r.Register("ENTITY.GETCLOSESTWITHTAG", "entity", runtime.AdaptLegacy(m.entGetClosestWithTag))
	r.Register("EntityGetClosestWithTag", "entity", runtime.AdaptLegacy(m.entGetClosestWithTag))
	r.Register("ENTITY.PUSHOUTOFGEOMETRY", "entity", runtime.AdaptLegacy(m.entPushOutOfGeometry))
	r.Register("EntityPushOutOfGeometry", "entity", runtime.AdaptLegacy(m.entPushOutOfGeometry))
	r.Register("ENTITY.INFRUSTUM", "entity", runtime.AdaptLegacy(m.entInFrustumScreen))
	r.Register("EntityInFrustum", "entity", runtime.AdaptLegacy(m.entInFrustumScreen))
	r.Register("CHECK.INVIEW", "check", runtime.AdaptLegacy(m.entInFrustumScreen))
	r.Register("ENTITY.ISSUBMERGED", "entity", runtime.AdaptLegacy(m.entIsSubmerged))
	r.Register("ENTITY.LINEOFSIGHT", "entity", runtime.AdaptLegacy(m.entLineOfSight))
	r.Register("EntityLineOfSight", "entity", runtime.AdaptLegacy(m.entLineOfSight))
	r.Register("ENTITY.GETOVERLAPCOUNT", "entity", runtime.AdaptLegacy(m.entGetOverlapCount))
	r.Register("EntityGetOverlapCount", "entity", runtime.AdaptLegacy(m.entGetOverlapCount))
}

// entRaycast uses Jolt broad-phase ray cast (same path as PHYSICS3D / PICK); returns first hit entity# or 0.
func (m *Module) entRaycast(args []value.Value) (value.Value, error) {
	if len(args) != 7 {
		return value.Nil, fmt.Errorf("ENTITY.RAYCAST expects (ox#, oy#, oz#, dx#, dy#, dz#, maxDist#)")
	}
	ox, _ := args[0].ToFloat()
	oy, _ := args[1].ToFloat()
	oz, _ := args[2].ToFloat()
	dx, _ := args[3].ToFloat()
	dy, _ := args[4].ToFloat()
	dz, _ := args[5].ToFloat()
	maxd, _ := args[6].ToFloat()
	id := mbphysics3d.PickCastEntityID(ox, oy, oz, dx, dy, dz, maxd)
	return value.FromInt(id), nil
}

func (m *Module) entGetGroundNormal(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.GETGROUNDNORMAL expects (entity#)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.GETGROUNDNORMAL: invalid entity")
	}
	if m.store().ents[id] == nil {
		return value.Nil, fmt.Errorf("ENTITY.GETGROUNDNORMAL: unknown entity")
	}
	if m.h == nil {
		return value.Nil, fmt.Errorf("ENTITY.GETGROUNDNORMAL: heap not bound")
	}
	if characterGroundNormalFn != nil {
		if nx, ny, nz, ok := characterGroundNormalFn(id); ok {
			return mbmatrix.AllocVec3Value(m.h, float32(nx), float32(ny), float32(nz))
		}
	}
	cx, cy, cz, ok := m.PlayerBridgeWorldPos(id)
	if !ok {
		return value.Nil, fmt.Errorf("ENTITY.GETGROUNDNORMAL: internal")
	}
	const probeUp = 0.35
	const maxDown = 4.0
	nx, ny, nz, hit := mbphysics3d.RaycastDownNormal(cx, cy+probeUp, cz, maxDown)
	if !hit {
		return mbmatrix.AllocVec3Value(m.h, 0, 1, 0)
	}
	return mbmatrix.AllocVec3Value(m.h, float32(nx), float32(ny), float32(nz))
}

func (m *Module) entCanSee(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("ENTITY.CANSEE expects (observer#, target#, fovDeg#, maxDist#)")
	}
	obs, ok1 := m.entID(args[0])
	tgt, ok2 := m.entID(args[1])
	if !ok1 || !ok2 || obs < 1 || tgt < 1 {
		return value.Nil, fmt.Errorf("ENTITY.CANSEE: invalid entity")
	}
	if obs == tgt {
		return value.FromBool(false), nil
	}
	if m.store().ents[obs] == nil || m.store().ents[tgt] == nil {
		return value.Nil, fmt.Errorf("ENTITY.CANSEE: unknown entity")
	}
	fov, okf := args[2].ToFloat()
	maxDist, okd := args[3].ToFloat()
	if !okf || !okd || fov <= 0 || maxDist <= 0 {
		return value.Nil, fmt.Errorf("ENTITY.CANSEE: fov and maxDist must be > 0")
	}
	const eyeY = 1.65
	ox, oy, oz, dx, dy, dz, ok := m.PlayerBridgeEyeRay(obs, eyeY)
	if !ok {
		return value.FromBool(false), nil
	}
	flen := math.Sqrt(dx*dx + dy*dy + dz*dz)
	if flen < 1e-9 {
		return value.FromBool(false), nil
	}
	fx, fy, fz := dx/flen, dy/flen, dz/flen
	tx, ty, tz, ok := m.PlayerBridgeWorldPos(tgt)
	if !ok {
		return value.FromBool(false), nil
	}
	// Aim at approximate eye height on the target (same default offset as observer).
	vx := tx - ox
	vy := (ty + eyeY) - oy
	vz := tz - oz
	vlen := math.Sqrt(vx*vx + vy*vy + vz*vz)
	if vlen < 1e-9 || vlen > maxDist {
		return value.FromBool(false), nil
	}
	vx /= vlen
	vy /= vlen
	vz /= vlen
	dot := fx*vx + fy*vy + fz*vz
	half := (fov * 0.5) * math.Pi / 180.0
	if dot < math.Cos(half) {
		return value.FromBool(false), nil
	}
	// Line of sight: first physics hit along the segment must be the target (or clear).
	hit := mbphysics3d.PickCastEntityID(ox, oy, oz, vx, vy, vz, vlen)
	if hit == 0 {
		return value.FromBool(true), nil
	}
	return value.FromBool(hit == tgt), nil
}

func (m *Module) entGetClosestWithTag(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("ENTITY.GETCLOSESTWITHTAG expects (entity#, radius#, tag$)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.GETCLOSESTWITHTAG: invalid entity")
	}
	rad, ok := args[1].ToFloat()
	if !ok || rad < 0 {
		return value.Nil, fmt.Errorf("ENTITY.GETCLOSESTWITHTAG: radius must be >= 0")
	}
	if args[2].Kind != value.KindString {
		return value.Nil, fmt.Errorf("ENTITY.GETCLOSESTWITHTAG: tag must be string")
	}
	if m.h == nil {
		return value.Nil, fmt.Errorf("ENTITY.GETCLOSESTWITHTAG: heap not bound")
	}
	tag, ok := m.h.GetString(int32(args[2].IVal))
	if !ok {
		return value.Nil, fmt.Errorf("ENTITY.GETCLOSESTWITHTAG: invalid tag string")
	}
	cx, cy, cz, ok := m.PlayerBridgeWorldPos(id)
	if !ok {
		return value.FromInt(0), nil
	}
	best := m.PlayerBridgeClosestTagged(cx, cy, cz, rad, tag)
	return value.FromInt(best), nil
}

func (m *Module) entPushOutOfGeometry(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.PUSHOUTOFGEOMETRY expects (entity#)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.PUSHOUTOFGEOMETRY: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.PUSHOUTOFGEOMETRY: unknown entity")
	}
	wp := m.worldPos(e)
	// Best-effort depenetration hint: lift slightly; full Jolt recovery requires CharacterVirtual / dynamic body APIs.
	m.setLocalFromWorld(e, wp.X, wp.Y+0.2, wp.Z)
	return value.Nil, nil
}

// entIsSubmerged returns approximately the fraction [0..1] of the entity's vertical extent below water surfaces (see water.EntitySubmergedFraction).
func (m *Module) entIsSubmerged(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("ENTITY.ISSUBMERGED: heap not bound")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.ISSUBMERGED expects (entity#)")
	}
	eid, ok := m.entID(args[0])
	if !ok || eid < 1 {
		return value.Nil, fmt.Errorf("ENTITY.ISSUBMERGED: invalid entity")
	}
	e := m.store().ents[eid]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.ISSUBMERGED: unknown entity")
	}
	wp := m.worldPos(e)
	var mnY, mxY, cx, cz float32
	switch e.kind {
	case entKindSphere:
		smax := e.scale.X
		if e.scale.Y > smax {
			smax = e.scale.Y
		}
		if e.scale.Z > smax {
			smax = e.scale.Z
		}
		rs := e.radius * smax
		mnY = wp.Y - rs
		mxY = wp.Y + rs
		cx, cz = wp.X, wp.Z
	default:
		mn, mx := m.aabbWorldMinMax(e)
		mnY, mxY = mn.Y, mx.Y
		cx = (mn.X + mx.X) * 0.5
		cz = (mn.Z + mx.Z) * 0.5
	}
	f := water.EntitySubmergedFraction(m.h, mnY, mxY, cx, cz)
	return value.FromFloat(float64(f)), nil
}

// entInFrustumScreen is true if the entity AABB intersects the current CAMERA.BEGIN frustum (same math as ENTITY.INVIEW).
func (m *Module) entInFrustumScreen(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.INFRUSTUM expects (entity#)")
	}
	eid, ok := m.entID(args[0])
	if !ok || eid < 1 {
		return value.Nil, fmt.Errorf("ENTITY.INFRUSTUM: invalid entity")
	}
	e := m.store().ents[eid]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.INFRUSTUM: unknown entity")
	}
	cam, okc := mbmodel3d.ActiveCamera3D()
	if !okc {
		return value.FromBool(false), nil
	}
	rw := float32(rl.GetScreenWidth())
	rh := float32(rl.GetScreenHeight())
	aspect := float32(16.0 / 9.0)
	if rh > 1e-3 {
		aspect = rw / rh
	}
	f := mbcamera.ExtractFrustum(cam, aspect)
	return value.FromBool(entityInFrustum(m, e, f)), nil
}

// entLineOfSight tests an unobstructed Jolt ray from observer eye to target (no FOV cone). Trigger/sensor bodies still block until filtered.
func (m *Module) entLineOfSight(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.LINEOFSIGHT expects (observer#, target#)")
	}
	obs, ok1 := m.entID(args[0])
	tgt, ok2 := m.entID(args[1])
	if !ok1 || !ok2 || obs < 1 || tgt < 1 {
		return value.Nil, fmt.Errorf("ENTITY.LINEOFSIGHT: invalid entity")
	}
	if obs == tgt {
		return value.FromBool(true), nil
	}
	if m.store().ents[obs] == nil || m.store().ents[tgt] == nil {
		return value.Nil, fmt.Errorf("ENTITY.LINEOFSIGHT: unknown entity")
	}
	const eyeY = 1.65
	ox, oy, oz, _, _, _, ok := m.PlayerBridgeEyeRay(obs, eyeY)
	if !ok {
		return value.FromBool(false), nil
	}
	tx, ty, tz, ok := m.PlayerBridgeWorldPos(tgt)
	if !ok {
		return value.FromBool(false), nil
	}
	vx := tx - ox
	vy := (ty + eyeY) - oy
	vz := tz - oz
	vlen := math.Sqrt(vx*vx + vy*vy + vz*vz)
	if vlen < 1e-9 {
		return value.FromBool(true), nil
	}
	hit := mbphysics3d.PickCastEntityID(ox, oy, oz, vx, vy, vz, vlen)
	if hit == tgt {
		return value.FromBool(true), nil
	}
	if hit == 0 {
		return value.FromBool(true), nil
	}
	return value.FromBool(false), nil
}

// entGetOverlapCount counts tagged entities whose pivot lies inside the zone entity's world AABB (approximate volume overlap).
func (m *Module) entGetOverlapCount(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.GETOVERLAPCOUNT expects (zoneEntity#, tag$)")
	}
	zid, ok := m.entID(args[0])
	if !ok || zid < 1 {
		return value.Nil, fmt.Errorf("ENTITY.GETOVERLAPCOUNT: invalid entity")
	}
	if args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("ENTITY.GETOVERLAPCOUNT: tag must be string")
	}
	if m.h == nil {
		return value.Nil, fmt.Errorf("ENTITY.GETOVERLAPCOUNT: heap not bound")
	}
	tagPat, ok := m.h.GetString(int32(args[1].IVal))
	if !ok {
		return value.Nil, fmt.Errorf("ENTITY.GETOVERLAPCOUNT: invalid tag string")
	}
	zone := m.store().ents[zid]
	if zone == nil {
		return value.Nil, fmt.Errorf("ENTITY.GETOVERLAPCOUNT: unknown zone entity")
	}
	mn, mx := m.aabbWorldMinMax(zone)
	cx := float64(mn.X+mx.X) * 0.5
	cy := float64(mn.Y+mx.Y) * 0.5
	cz := float64(mn.Z+mx.Z) * 0.5
	rx := math.Abs(float64(mx.X-mn.X)) * 0.5
	ry := math.Abs(float64(mx.Y-mn.Y)) * 0.5
	rz := math.Abs(float64(mx.Z-mn.Z)) * 0.5
	rad := rx
	if ry > rad {
		rad = ry
	}
	if rz > rad {
		rad = rz
	}
	if rad < 0.05 {
		rad = 1.0
	}
	ids := m.PlayerBridgeNearbyTagged(cx, cy, cz, rad, tagPat)
	n := int64(0)
	for _, oid := range ids {
		if oid == zid {
			continue
		}
		oe := m.store().ents[oid]
		if oe == nil {
			continue
		}
		wp := m.worldPos(oe)
		if wp.X >= mn.X && wp.X <= mx.X && wp.Y >= mn.Y && wp.Y <= mx.Y && wp.Z >= mn.Z && wp.Z <= mx.Z {
			n++
		}
	}
	return value.FromInt(n), nil
}
